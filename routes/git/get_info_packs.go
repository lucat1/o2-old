package git

import "github.com/gin-gonic/gin"

// GetInfoPacks handles git pull
// /:user/:repo/objects/info/packs
func GetInfoPacks(c *gin.Context) {
	hdrCacheForever(c)
	sendFile("text/plain; charset=utf-8", c)
}
