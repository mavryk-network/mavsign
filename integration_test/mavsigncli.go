package integrationtest

import (
	"os/exec"
)

func MavSignCli(arg ...string) ([]byte, error) {
	var cmd = "docker"
	var args = []string{"exec", "mavsign", "mavsign-cli"}
	args = append(args, arg...)
	return exec.Command(cmd, args...).CombinedOutput()
}
