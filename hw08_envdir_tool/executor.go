package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, stdout, stderr io.Writer) int {
	for key, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				fmt.Println(err)
			}

			continue
		}

		if err := os.Setenv(key, value.Value); err != nil {
			fmt.Println(err)
		}
	}

	executablePath, err := exec.LookPath(cmd[0])
	if err != nil {
		fmt.Printf("Executable not found: %v\n", err)

		return 1
	}

	execCmd := exec.Command(executablePath, cmd[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = stdout
	execCmd.Stderr = stderr

	var exitErr *exec.ExitError
	if err = execCmd.Run(); errors.As(err, &exitErr) {
		waitStatus := exitErr.Sys().(syscall.WaitStatus)

		return waitStatus.ExitStatus()
	}

	return 0
}
