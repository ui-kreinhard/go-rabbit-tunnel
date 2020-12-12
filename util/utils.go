package util

import "os/exec"

func convertOutput(output []byte, err error) (string, error) {
	return string(output), err
}

func Exec(cmd string, params ...string) (string, error) {
	return convertOutput(exec.Command(cmd, params...).CombinedOutput())
}
