package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

// NotFound is the 404 route
// 404 route
func NotFound(c *gin.Context) {
	fmt.Println("404 " + c.Request.Method + " " + c.Request.URL.Path);
	debug.PrintStack()
	c.HTML(404, "notfound.tmpl", gin.H{
		"path": c.Request.URL.Path,
		"user": c.Keys["user"],
	})
}
