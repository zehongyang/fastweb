package fastweb

import (
	"fmt"
	"testing"
)

func MiddlewareTest() HandleFunc {
	return func(ctx *WebContext) {
		fmt.Println("middleware")
		ctx.Next()
	}
}

func HandlerTest() HandleFunc {
	return func(ctx *WebContext) {
		ctx.Response(map[string]string{"hello": "world"})
	}
}

func TestEngine(t *testing.T) {
	eng := New()
	eng.Register(HandlerTest, MiddlewareTest)
	eng.Run()
}
