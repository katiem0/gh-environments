package cmd

import (
	"bytes"
	"testing"
)

func TestNewCmdRoot(t *testing.T) {
	cmd := NewCmdRoot()

	if cmd == nil {
		t.Fatal("NewCmdRoot() returned nil")
	}

	// Test basic properties
	if cmd.Use != "environments <command> <subcommand> [flags]" {
		t.Errorf("Expected Use to be 'environments <command> <subcommand> [flags]', got %s", cmd.Use)
	}

	// Test that subcommands are added
	subcommands := cmd.Commands()
	if len(subcommands) < 2 {
		t.Errorf("Expected at least 2 subcommands, got %d", len(subcommands))
	}

	// Test completion options
	if !cmd.CompletionOptions.DisableDefaultCmd {
		t.Error("Default completion command should be disabled")
	}

	// Test that the command has a short description
	if cmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	// Test command execution with help flag
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected help output, got empty string")
	}
}
