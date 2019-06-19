package git

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetInfoRefs(c *gin.Context) {
	dir := c.Keys["dir"].(string)
	service_name := getServiceType(c.Request)
	access := hasAccess(c.Request, dir, service_name, false)
	version := c.GetHeader("Git-Protocol")
	if access {
		args := []string{service_name, "--stateless-rpc", "--advertise-refs", "."}
		refs := gitCommand(dir, version, args...)

		hdrNocache(c)
		c.Header("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service_name))
		c.Status(http.StatusOK)
		if len(version) == 0 {
			c.Writer.Write(packetWrite("# service=git-" + service_name + "\n"))
			c.Writer.Write(packetFlush())
		}
		c.Writer.Write(refs)
	} else {
		updateServerInfo(dir)
		hdrNocache(c)
		sendFile("text/plain; charset=utf-8", c)
	}
}
