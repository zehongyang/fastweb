package fastweb

import (
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/utils/config"
	"github.com/zehongyang/utils/logger"
	"go.uber.org/zap"
	"reflect"
	"runtime"
	"strings"
)

type WebConfig struct {
	Web []WebServer
}

type WebServer struct {
	Server      string
	Port        string
	Middlewares []string
	Routers     []WebRouter
}

type WebRouter struct {
	Group       string
	Middlewares []string
	Paths       []WebPath
}

type WebPath struct {
	Method  string
	Path    string
	Handler string
}

type Engine struct {
	engines  []*gin.Engine
	config   WebConfig
	handlers map[string]HandleFunc
	errCh    chan error
}

type HandleFunc func(ctx *WebContext)

func New() *Engine {
	var wc WebConfig
	err := config.Load(&wc)
	if err != nil {
		logger.F("", zap.Error(err))
	}
	return &Engine{config: wc, handlers: make(map[string]HandleFunc), errCh: make(chan error)}
}

func (e *Engine) getHandler(hf string) gin.HandlerFunc {
	handleFunc, ok := e.handlers[hf]
	if !ok {
		logger.E("", zap.Any("handler not found", hf))
		return nil
	}
	return func(ctx *gin.Context) {
		if handleFunc != nil {
			handleFunc(&WebContext{ctx})
		}
	}
}

func (e *Engine) Register(fns ...func() HandleFunc) {
	for _, fn := range fns {
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		names := strings.Split(name, "/")
		if len(names) < 1 {
			logger.W("", zap.Any("not found func name", name))
			continue
		}
		e.handlers[names[len(names)-1]] = fn()
	}
}

func (e *Engine) Run() error {
	if len(e.config.Web) < 1 {
		return EngineConfigErr
	}
	for _, server := range e.config.Web {
		engine := gin.Default()
		if len(server.Middlewares) > 0 {
			var fns []gin.HandlerFunc
			for _, middleware := range server.Middlewares {
				fn := e.getHandler(middleware)
				if fn != nil {
					fns = append(fns, fn)
				}
			}
			if len(fns) > 0 {
				engine.Use(fns...)
			}
		}
		if len(server.Routers) > 0 {
			for _, router := range server.Routers {
				routerGroup := &engine.RouterGroup
				if len(router.Group) > 0 {
					routerGroup = engine.Group(router.Group)
				}
				if len(router.Middlewares) > 0 {
					var fns []gin.HandlerFunc
					for _, middleware := range router.Middlewares {
						fn := e.getHandler(middleware)
						if fn != nil {
							fns = append(fns, fn)
						}
						if len(fns) > 0 {
							routerGroup.Use(fns...)
						}
					}
				}
				if len(router.Paths) > 0 {
					for _, path := range router.Paths {
						fn := e.getHandler(path.Handler)
						if fn == nil {
							continue
						}
						routerGroup.Handle(path.Method, path.Path, fn)
					}
				}
			}
		}
		go func(port string) {
			err := engine.Run(port)
			if err != nil {
				e.errCh <- err
			}
		}(server.Port)
	}
	select {
	case err := <-e.errCh:
		return err
	}
}
