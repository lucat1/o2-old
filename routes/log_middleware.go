package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	"time"
)

// LogMiddleware reports the time spent in requests
func LogMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	end := time.Now().Sub(start)
	shared.GetLogger().Info(c.Request.Method + " " + c.Request.URL.Path + " -- " + end.String())
}