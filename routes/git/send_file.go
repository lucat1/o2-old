package git

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/routes"
	"net/http"
	"os"
	"path"
)

func sendFile(content_type string, c *gin.Context) {
	req_file := path.Join(c.Keys["dir"].(string), c.Keys["file"].(string))

	f, err := os.Stat(req_file)
	if os.IsNotExist(err) {
		routes.NotFound(c)
		return
	}

	c.Header("Content-Type", content_type)
	c.Header("Content-Length", fmt.Sprintf("%d", f.Size()))
	c.Header("Last-Modified", f.ModTime().Format(http.TimeFormat))
	http.ServeFile(c.Writer, c.Request, req_file)
}