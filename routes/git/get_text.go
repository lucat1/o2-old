package git

import "github.com/gin-gonic/gin"

func GetTextFile(c *gin.Context) {
	hdrNocache(c)
	sendFile("text/plain", c)
}