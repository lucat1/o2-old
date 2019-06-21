package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
)

// Static handles the serving of static files like
// /favicon.ico
// /shared.css
func Static(assets map[string]*assets.File) func(*gin.Context) {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if strings.HasPrefix(url, "/static/") || url == "/favicon.ico" {
			asset := assets[url]
			if url == "/favicon.ico" {
				asset = assets["/shared/favicon.ico"]
			}

			if asset != nil {
				c.Writer.Write(asset.Data)
				c.Abort()
				return
			}

			c.Next()
		}
	}
}
