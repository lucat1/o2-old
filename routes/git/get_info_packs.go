package git

import "github.com/gin-gonic/gin"

func GetInfoPacks(c *gin.Context) {
	hdrCacheForever(c)
	sendFile("text/plain; charset=utf-8", c)
}
