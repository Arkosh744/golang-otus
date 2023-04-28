package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s /path/to/env/dir command arg1 arg2 ...\n", os.Args[0])
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	returnCode := RunCmd(os.Args[2:], env, os.Stdout, os.Stderr)
	os.Exit(returnCode)
}
