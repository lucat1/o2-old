package git

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/routes"
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

	// usr, passwd
	authHead := c.GetHeader("Authorization")
	if len(authHead) == 0 {
		c.Header("WWW-Authenticate", "Basic realm=\".\"")
		c.Status(http.StatusUnauthorized)
		return
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	shared.GetLogger().Info(
		"New login",
		zap.String("username", username),
		zap.String("password", password),
	)

	user := routes.FindUser(username)
	if user == nil {
		// User with the provided username doesnt exist
		c.Status(http.StatusUnauthorized)
		return
	}
	// If the user exists lets check for the password
	ok = shared.CheckPassword(user.Password, password)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	shared.GetLogger().Info("Git '"+rpc+"' rpc with authenticated user", zap.String("user", user.Username))

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
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	p := make([]byte, 1024)
	for {
		nRead, err := stdout.Read(p)
		if err == io.EOF {
			break
		}
		nWrite, err := c.Writer.Write(p[:nRead])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if nRead != nWrite {
			fmt.Printf("failed to write data: %d read, %d written\n", nRead, nWrite)
			os.Exit(1)
		}
		flusher.Flush()
	}

	cmd.Wait()
}
