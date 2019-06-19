package git

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)
func hdrNocache(c *gin.Context) {
	c.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	c.Header("Pragma", "no-cache")
	c.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func hdrCacheForever(c *gin.Context) {
	now := time.Now().Unix()
	expires := now + 31536000
	c.Header("Date", fmt.Sprintf("%d", now))
	c.Header("Expires", fmt.Sprintf("%d", expires))
	c.Header("Cache-Control", "public, max-age=31536000")
}
