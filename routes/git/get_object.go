package git

import "github.com/gin-gonic/gin"

func GetLooseObject(c *gin.Context) {
	hdrCacheForever(c)
	sendFile("application/x-git-loose-object", c)
}