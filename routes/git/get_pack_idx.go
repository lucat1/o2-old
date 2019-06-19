package git

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetPackFile handles the streaming of single packages
func GetPackFile(c *gin.Context) {
	hdrCacheForever(c)
	c.Keys["file"] = strings.Replace(c.Param("pack"), ".pack", "", 1)
	sendFile("application/x-git-packed-objects", c)
}

// GetIdxFile handles the streaming of single packages
func GetIdxFile(c *gin.Context) {
	hdrCacheForever(c)
	c.Keys["file"] = strings.Replace(c.Param("pack"), ".idx", "", 1)
	sendFile("application/x-git-packed-objects-toc", c)
}

// GetPackOrIdx handles the streaming of both formats above
// Providing the correnct route for each request
// /:user/:repo/objects/pack/pack-:pack
func GetPackOrIdx(c *gin.Context) {
	if strings.HasSuffix(c.Request.URL.Path, ".pack") {
		GetPackFile(c)
	} else {
		GetIdxFile(c)
	}
}
