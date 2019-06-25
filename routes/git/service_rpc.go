package git

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

// ServiceRPC handles the biggest part of pushing/pulling
// via the git headless rpc service
// /:user/:repo/git-upload-pack
func ServiceRPC(c *gin.Context) {
	dir := c.Keys["dir"].(string)
	var rpc string
	if strings.Contains(c.Request.URL.Path, "upload-pack") {
		rpc = "upload-pack"
	} else {
		rpc = "receive-pack"
	}
	shared.GetLogger().Info("RPC " + rpc)
	access := hasAccess(c.Request, dir, rpc, true)

	if access == false {
		c.Status(500)
		return
	}

	c.Header("Content-Type", "application/x-git-"+rpc+"-result")
	c.Header("Connection", "Keep-Alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Status(http.StatusOK)

	var env []string
	// TODO: config
	/*
		if config.DefaultEnv != "" {
			env = append(env, config.DefaultEnv)
		}*/

	args := []string{rpc, "--stateless-rpc", dir}
	cmd := exec.Command("/usr/bin/git", args...)
	version := c.GetHeader("Git-Protocol")
	if len(version) != 0 {
		cmd.Env = append(os.Environ(), "GIT_PROTOCOL="+version)
	}
	cmd.Dir = dir
	cmd.Env = env
	in, err := cmd.StdinPipe()
	if err != nil {
		shared.GetLogger().Warn("Error in service-rpc while reading stdin", zap.Error(err))
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		shared.GetLogger().Warn("Error in service-rpc while reading stdout", zap.Error(err))
	}

	err = cmd.Start()
	if err != nil {
		shared.GetLogger().Warn("Error in service-rpc while reading exit error code", zap.Error(err))
	}

	var reader io.ReadCloser
	switch c.GetHeader("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(c.Request.Body)
		defer reader.Close()
	default:
		reader = c.Request.Body
	}
	io.Copy(in, reader)
	defer in.Close()

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		shared.GetLogger().Error("expected http.ResponseWriter to be an http.Flusher")
		return
	}

	p := make([]byte, 1024)
	for {
		nRead, err := stdout.Read(p)
		if err == io.EOF {
			break
		}
		nWrite, err := c.Writer.Write(p[:nRead])
		if err != nil {
			shared.GetLogger().Error("Could not write to response in git rpc", zap.Error(err))
			return
		}
		if nRead != nWrite {
			shared.GetLogger().Error(
				"Written/Ridden data do not match in git rpc",
				zap.Int("nRead", nRead),
				zap.Int("nWrite", nWrite),
			)
			return
		}
		flusher.Flush()
	}

	cmd.Wait()
}
