package integrationtest

import (
	"fmt"
	"os/exec"
)

func restart_mavsign() {
	_, err := exec.Command("docker", "compose", "-f", "./docker-compose.yml", "stop", "mavsign").CombinedOutput()
	if err != nil {
		panic("failed to stop mavsign")
	}
	out, err := exec.Command("docker", "compose", "-f", "./docker-compose.yml", "up", "-d", "--wait", "mavsign").CombinedOutput()
	if err != nil {
		fmt.Println("restart mavsign: failed to start: " + string(out))
		panic("failed to start mavsign during restart")
	}
}

func backup_then_update_config(c Config) {
	_, err := exec.Command("cp", "mavsign.yaml", "mavsign.original.yaml").CombinedOutput()
	if err != nil {
		panic("failed to backup config")
	}
	err = c.Write()
	if err != nil {
		panic("failed to write new config")
	}
}

func restore_config() {
	_, err := exec.Command("mv", "mavsign.original.yaml", "mavsign.yaml").CombinedOutput()
	if err != nil {
		panic("failed to restore original config")
	}
	restart_mavsign()
}

func restart_stack() {
	_, err := exec.Command("docker", "compose", "-f", "./docker-compose.yml", "kill").CombinedOutput()
	if err != nil {
		panic("failed to kill stack")
	}
	_, err = exec.Command("docker", "compose", "-f", "./docker-compose.yml", "up", "-d", "--wait").CombinedOutput()
	if err != nil {
		panic("failed to up stack")
	}
}
