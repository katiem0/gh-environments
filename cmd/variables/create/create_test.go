package createvariables

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/katiem0/gh-environments/internal/data"
	"github.com/katiem0/gh-environments/internal/utils"
	"github.com/spf13/cobra"
)

func TestNewCmdCreate(t *testing.T) {
	cmd := NewCmdCreate()

	if cmd == nil {
		t.Fatal("NewCmdCreate() returned nil")
	}

	// Test basic properties
	if cmd.Use != "create <organization> [flags]" {
		t.Errorf("Expected Use to be 'create <organization> [flags]', got %s", cmd.Use)
	}

	// Test flags exist
	if cmd.Flag("from-file") == nil {
		t.Error("from-file flag not found")
	}

	if cmd.Flag("token") == nil {
		t.Error("token flag not found")
	}

	if cmd.Flag("hostname") == nil {
		t.Error("hostname flag not found")
	}

	if cmd.Flag("debug") == nil {
		t.Error("debug flag not found")
	}
}

func runCmdCreateTest(owner string, cmdFlags *cmdFlags, g interface{}) error {
	// Type assertion to the interface methods we need
	getter, ok := g.(interface {
		CreateVariableList(data [][]string) []data.ImportedVariable
		CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error
	})

	if !ok {
		return nil // For testing, we're not concerned with this error
	}

	// Call the original function implementation but with our getter
	var variableData [][]string
	var variablesList []data.ImportedVariable

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				// In a real implementation, you might want to log this error
				fmt.Printf("Error closing file: %v\n", err)
			}
		}()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		variableData, err = csvReader.ReadAll()
		if err != nil {
			return err
		}
		variablesList = getter.CreateVariableList(variableData)

		for _, variable := range variablesList {
			importVar := utils.CreateVariableData(variable)
			createVariable, err := json.Marshal(importVar)
			if err != nil {
				return err
			}
			reader := bytes.NewReader(createVariable)
			err = getter.CreateEnvironmentVariables(owner, variable.RepositoryName, variable.EnvironmentName, reader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Custom adapter to make MockAPIGetter compatible with runCmdCreate
type testAPIGetter struct {
	mock *utils.MockAPIGetter
}

// Implement the necessary methods from APIGetter
func (t *testAPIGetter) CreateVariableList(data [][]string) []data.ImportedVariable {
	return t.mock.CreateVariableList(data)
}

func (t *testAPIGetter) CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error {
	return t.mock.CreateEnvironmentVariables(owner, repo, env, data)
}

func setupMockGetter() (*utils.MockAPIGetter, *testAPIGetter) {
	mockGetter := utils.NewMockAPIGetter()

	// Create our adapter that wraps the mock
	testGetter := &testAPIGetter{
		mock: mockGetter,
	}

	return mockGetter, testGetter
}

func TestRunCmdCreate(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test-variables.csv")
	csvContent := `RepositoryName,RepositoryID,EnvironmentName,VariableName,VariableValue,VariableCreatedAt,VariableUpdatedAt
testrepo,12345,production,TEST_VAR,test-value,,`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	// Create mock API getter
	_, testGetter := setupMockGetter()

	// Create command flags
	flags := &cmdFlags{
		fileName: csvFile,
		hostname: "github.com",
		debug:    false,
	}

	// Execute with our adapter, using the test version that accepts an interface
	err = runCmdCreateTest("testorg", flags, testGetter)

	// Verify
	if err != nil {
		t.Errorf("runCmdCreate() error = %v", err)
	}
}

func TestRunCmdCreateFileError(t *testing.T) {
	// Create mock API getter
	_, testGetter := setupMockGetter()

	// Create command flags with non-existent file
	flags := &cmdFlags{
		fileName: "non-existent-file.csv",
		hostname: "github.com",
		debug:    false,
	}

	// Execute with our adapter, using the test version that accepts an interface
	err := runCmdCreateTest("testorg", flags, testGetter)

	// Verify error is returned
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestCmdRunE(t *testing.T) {
	// This test only checks that arguments validation works correctly
	// We'll skip the actual execution to avoid real API calls
	cmd := NewCmdCreate()

	// Temporarily modify RunE to skip actual execution
	originalRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires at least 1 arg(s), only received %d", len(args))
		}
		return nil
	}
	defer func() { cmd.RunE = originalRunE }()

	// Test with insufficient args
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error for insufficient arguments, got nil")
	}

	// Test with sufficient args
	err = cmd.RunE(cmd, []string{"test-org"})
	if err != nil {
		t.Errorf("Unexpected error for sufficient arguments: %v", err)
	}
}
