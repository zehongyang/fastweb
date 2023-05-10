package fastweb

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type WebContext struct {
	*gin.Context
}

func (s *WebContext) Response(data ...any) {
	if len(data) > 0 {
		if s.GetHeader("Content-Type") == binding.MIMEPROTOBUF {
			s.ProtoBuf(http.StatusOK, data[0])
		} else {
			s.JSON(http.StatusOK, data[0])
		}
	}
}
