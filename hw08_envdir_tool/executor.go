package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, stdout, stderr io.Writer) (returnCode int) {
	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := os.Setenv(key, value.Value)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = stdout
	execCmd.Stderr = stderr

	err := execCmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		waitStatus := exitErr.Sys().(syscall.WaitStatus)
		returnCode = waitStatus.ExitStatus()
	} else {
		returnCode = 0
	}

	return
}
