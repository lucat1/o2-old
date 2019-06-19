package git

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetPackFile(c *gin.Context) {
	hdrCacheForever(c)
	c.Keys["file"] = strings.Replace(c.Param("pack"), ".pack", "", 1)
	sendFile("application/x-git-packed-objects", c)
}

func GetIdxFile(c *gin.Context) {
	hdrCacheForever(c)
	c.Keys["file"] = strings.Replace(c.Param("pack"), ".idx", "", 1)
	sendFile("application/x-git-packed-objects-toc", c)
}

func GetPackOrIdx(c *gin.Context) {
	if strings.HasSuffix(c.Request.URL.Path, ".pack") {
		GetPackFile(c)
	} else {
		GetIdxFile(c)
	}
}