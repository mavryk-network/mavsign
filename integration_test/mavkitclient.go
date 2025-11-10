package integrationtest

import (
	"os/exec"
)

func MavkitClient(arg ...string) ([]byte, error) {
	var cmd = "docker"
	var args = []string{"exec", "mavkit", "mavkit-client"}
	args = append(args, arg...)
	return exec.Command(cmd, args...).CombinedOutput()
}

func clean_mavryk_folder() {
	delete_contracts_aliases()
	delete_wallet_lock()
	delete_watermark_files()
}

func delete_wallet_lock() {
	var cmd = "docker"
	var args = []string{"exec", "mavkit", "rm", "-f", "/home/mavryk/.mavryk-client/wallet_lock"}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic("Clean mavryk: Failed to delete wallet lock: " + string(out))
	}
}

func delete_contracts_aliases() {
	var cmd = "docker"
	var args = []string{"exec", "mavkit", "rm", "-f", "/home/mavryk/.mavryk-client/contracts"}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic("Clean mavryk: Failed to delete contracts: " + string(out))
	}
}

func delete_watermark_files() {
	var cmd = "docker"
	var args = []string{"exec", "mavkit", "/bin/sh", "-c", "rm -f /home/mavryk/.mavryk-client/*_highwatermarks"}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic("Clean mavryk: Failed to delete watermarks: " + string(out))
	}
}
