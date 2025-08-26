package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Store original args and restore them after the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with a valid command
	os.Args = []string{"gh-environments", "environments", "--help"}

	// We can't actually call main() since it calls os.Exit()
	// Instead, verify the command structure in the cmd package

	// This test just ensures the main package builds correctly
	t.Log("Main package builds successfully")
}
