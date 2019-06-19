package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
		log.Print(err)
	}

	return out
}

func getConfigSetting(service_name string, dir string) bool {
	service_name = strings.Replace(service_name, "-", "", -1)
	setting := getGitConfig("http."+service_name, dir)

	if service_name == "uploadpack" {
		return setting != "false"
	}

	return setting == "true"
}

func getGitConfig(config_name string, dir string) string {
	args := []string{"config", config_name}
	out := string(gitCommand(dir, "", args...))
	return out[0 : len(out)-1]
}