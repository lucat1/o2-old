package git

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/routes"
)

func sendFile(contentType string, c *gin.Context) {
	reqFile := path.Join(c.Keys["dir"].(string), c.Keys["file"].(string))

	f, err := os.Stat(reqFile)
	if os.IsNotExist(err) {
		routes.NotFound(c)
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", f.Size()))
	c.Header("Last-Modified", f.ModTime().Format(http.TimeFormat))
	http.ServeFile(c.Writer, c.Request, reqFile)
}
