package git

import (
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

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

	c.Header("Content-Type", "application/x-git-" + rpc + "-result")
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

	// TODO: config
	/*
	user, password, authok := c.Request.BasicAuth()
	fmt.Println(authok)
	if authok {

		if config.AuthUserEnvVar != "" {
			env = append(env, fmt.Sprintf("%s=%s", config.AuthUserEnvVar, user))
		}
		if config.AuthPassEnvVar != "" {
			env = append(env, fmt.Sprintf("%s=%s", config.AuthPassEnvVar, password))
		}
	}*/

	args := []string{rpc, "--stateless-rpc", dir}
	cmd := exec.Command("/usr/bin/git", args...)
	version := c.GetHeader("Git-Protocol")
	if len(version) != 0 {
		cmd.Env = append(os.Environ(), "GIT_PROTOCOL=" + version)
	}
	cmd.Dir = dir
	cmd.Env = env
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Print(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Print(err)
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
		n_read, err := stdout.Read(p)
		if err == io.EOF {
			break
		}
		n_write, err := c.Writer.Write(p[:n_read])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n_read != n_write {
			fmt.Printf("failed to write data: %d read, %d written\n", n_read, n_write)
			os.Exit(1)
		}
		flusher.Flush()
	}

	cmd.Wait()
}