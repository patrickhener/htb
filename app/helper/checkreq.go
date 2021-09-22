package helper

import "os/exec"

// CheckTool will check if the required tool is there and then return true
func CheckTool(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
