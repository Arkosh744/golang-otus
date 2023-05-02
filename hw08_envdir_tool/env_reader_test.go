package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	testDir, err := os.MkdirTemp("", "envdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	err = os.WriteFile(filepath.Join(testDir, "FOO"), []byte("123\n"), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	err = os.WriteFile(filepath.Join(testDir, "BAR"), []byte("value"), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	env, err := ReadDir(testDir)
	if err != nil {
		t.Fatalf("Error reading directory: %v", err)
	}

	if len(env) != 2 {
		t.Fatalf("Expected 2 environment variables, got %d", len(env))
	}

	if env["FOO"].Value != "123" {
		t.Errorf("Expected FOO=123, got FOO=%s", env["FOO"].Value)
	}

	if env["BAR"].Value != "value" {
		t.Errorf("Expected BAR=value, got BAR=%s", env["BAR"].Value)
	}
}
