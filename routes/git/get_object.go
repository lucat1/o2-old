package git

import "github.com/gin-gonic/gin"

// GetLooseObject handles the stream of raw git objects
// /:user/:repo/objects/*path
func GetLooseObject(c *gin.Context) {
	hdrCacheForever(c)
	sendFile("application/x-git-loose-object", c)
}
