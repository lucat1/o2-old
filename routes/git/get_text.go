package git

import "github.com/gin-gonic/gin"

// GetTextFile returns a raw file in plain text
// /:user/:repo/HEAD
// /:user/:repo/objects/info/alternates
// /:user/:repo/objects/info/http-alternates
// /:user/:repo/objects/info/*path
func GetTextFile(c *gin.Context) {
	hdrNocache(c)
	sendFile("text/plain", c)
}
