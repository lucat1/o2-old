package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
)

// Static handles the serving of static files like
// /favicon.ico
// /shared.css
func Static(assets map[string]*assets.File) func(*gin.Context) {
	if os.Getenv("O2") == "dev" {
		// Serve fiels from the fs
		return func(c *gin.Context) {
			url := c.Request.URL.Path
			if strings.HasPrefix(url, "/static/") || url == "/favicon.ico" {
				if url == "/favicon.ico" {
					url = "/static/favicon.ico"
				}

				http.ServeFile(c.Writer, c.Request, "."+url)
				c.Abort()
			}
		}
	}

	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if strings.HasPrefix(url, "/static/") || url == "/favicon.ico" {
			if url == "/favicon.ico" {
				url = "/static/favicon.ico"
			}

			if assets[url] != nil {
				c.Writer.Write(assets[url].Data)
				c.Abort()
				return
			}

			c.Next()
		}
	}
}
