package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": EnvValue{Value: "123"},
		"BAR": EnvValue{Value: "value"},
	}

	var stdout, stderr bytes.Buffer

	returnCode := RunCmd([]string{"env"}, env, &stdout, &stderr)

	if returnCode != 0 {
		t.Fatalf("Expected return code 0, got %d", returnCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "FOO=123") {
		t.Errorf("Expected output to contain FOO=123, got: %s", output)
	}

	if !strings.Contains(output, "BAR=value") {
		t.Errorf("Expected output to contain BAR=value, got: %s", output)
	}
}
