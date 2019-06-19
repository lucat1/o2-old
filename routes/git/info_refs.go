package git

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetInfoRefs handles the diffing between local branches and the remote
// by providing the current refs in the server
// /:user/:repo/info/refs
func GetInfoRefs(c *gin.Context) {
	dir := c.Keys["dir"].(string)
	serviceName := getServiceType(c.Request)
	access := hasAccess(c.Request, dir, serviceName, false)
	version := c.GetHeader("Git-Protocol")
	if access {
		args := []string{serviceName, "--stateless-rpc", "--advertise-refs", "."}
		refs := gitCommand(dir, version, args...)

		hdrNocache(c)
		c.Header("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", serviceName))
		c.Status(http.StatusOK)
		if len(version) == 0 {
			c.Writer.Write(packetWrite("# service=git-" + serviceName + "\n"))
			c.Writer.Write(packetFlush())
		}
		c.Writer.Write(refs)
	} else {
		updateServerInfo(dir)
		hdrNocache(c)
		sendFile("text/plain; charset=utf-8", c)
	}
}
