package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (env Environment, err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env = make(Environment)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if strings.Contains(name, "=") {
			return nil, errors.New("filename should not contain '='")
		}

		file, err := os.Open(dir + "/" + name)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, " \t")
			line = strings.ReplaceAll(line, string(rune(0x00)), "\n")
			env[name] = EnvValue{Value: line}
		} else {
			env[name] = EnvValue{NeedRemove: true}
		}
	}

	return env, nil
}
