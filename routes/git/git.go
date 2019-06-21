package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

func updateServerInfo(dir string) []byte {
	args := []string{"update-server-info"}
	return gitCommand(dir, "", args...)
}

func gitCommand(dir string, version string, args ...string) []byte {
	command := exec.Command("/usr/bin/git", args...)
	if len(version) > 0 {
		command.Env = append(os.Environ(), fmt.Sprintf("GIT_PROTOCOL=%s", version))
	}
	command.Dir = dir
	out, err := command.Output()

	if err != nil {
		shared.GetLogger().Warn(
			"Error while executing git command",
			zap.String("dir", dir),
			zap.String("version", version),
			zap.Strings("args", args),
		)
	}

	return out
}

func getConfigSetting(serviceName string, dir string) bool {
	serviceName = strings.Replace(serviceName, "-", "", -1)
	setting := getGitConfig("http."+serviceName, dir)

	if serviceName == "uploadpack" {
		return setting != "false"
	}

	return setting == "true"
}

func getGitConfig(configName string, dir string) string {
	args := []string{"config", configName}
	out := string(gitCommand(dir, "", args...))
	return out[0 : len(out)-1]
}
