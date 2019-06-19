package routes

import (
	"github.com/gin-gonic/gin"
)

// Index is the / route
// /
func Index(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"user": c.Keys["user"],
	})
}
