package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
)

// NotFound is the 404 route
// 404 route
func NotFound(c *gin.Context) {
	shared.GetLogger().Info(
		"Hit 404",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)

	c.HTML(404, "notfound.tmpl", gin.H{
		"path": c.Request.URL.Path,
		"user": c.Keys["user"],
	})
	c.Abort()
}
