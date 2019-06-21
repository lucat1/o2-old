package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

// NotFound is the 404 route
// 404 route
func NotFound(c *gin.Context) {
	if c.Keys["notfound"] != nil {
		shared.GetLogger().Warn("Called NotFound twice!!")
		return
	}

	shared.GetLogger().Info(
		"Rendering NotFound for route",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)

	c.HTML(404, "notfound.tmpl", gin.H{
		"path": c.Request.URL.Path,
		"user": c.Keys["user"],
	})
	c.Abort()
	c.Keys["notfound"] = true
}
