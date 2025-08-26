package create

import (
	"io"
	"strconv"
	"testing"

	"github.com/katiem0/gh-environments/internal/data"
	"github.com/katiem0/gh-environments/internal/utils"
)

func TestNewCmdCreate(t *testing.T) {
	cmd := NewCmdCreate()

	if cmd == nil {
		t.Fatal("NewCmdCreate() returned nil")
	}

	// Test basic properties
	if cmd.Use != "create  <target organization> [flags]" {
		t.Errorf("Expected Use to be 'create  <target organization> [flags]', got %s", cmd.Use)
	}

	// Test flags
	if cmd.Flag("from-file") == nil {
		t.Error("from-file flag not found")
	}

	if cmd.Flag("hostname") == nil {
		t.Error("hostname flag not found")
	}

	if cmd.Flag("debug") == nil {
		t.Error("debug flag not found")
	}

	// Test short description
	if cmd.Short == "" {
		t.Error("Command should have a short description")
	}
}

func TestCreateEnvironmentData(t *testing.T) {
	// Test environment data creation
	importedEnv := data.ImportedEnvironment{
		RepositoryName:  "testrepo",
		RepositoryID:    12345,
		EnvironmentName: "production",
		AdminBypass:     "false",
		WaitTimer:       5,
		Reviewers: []data.Reviewers{
			{Type: "User", Reviewer: data.Reviewer{ID: 1, Login: "user1"}},
			{Type: "Team", Reviewer: data.Reviewer{ID: 2, Login: "team1"}},
		},
		PreventSelfReview: true,
		DeploymentPolicy:  "protected",
		Branches: []data.CreateDeploymentBranch{
			{Name: "main", Type: "branch"},
		},
	}

	result := utils.CreateEnvironmentData(importedEnv)

	if result == nil {
		t.Fatal("CreateEnvironmentData() returned nil")
	}

	if result.WaitTimer != 5 {
		t.Errorf("Expected WaitTimer to be 5, got %d", result.WaitTimer)
	}

	if result.PreventSelfReview != true {
		t.Errorf("Expected PreventSelfReview to be true, got %v", result.PreventSelfReview)
	}

	if len(result.Reviewers) != 2 {
		t.Errorf("Expected 2 reviewers, got %d", len(result.Reviewers))
	}

	if result.DeploymentBranchPolicy == nil {
		t.Error("Expected DeploymentBranchPolicy to be non-nil")
	} else if !result.DeploymentBranchPolicy.ProtectedBranches {
		t.Error("Expected ProtectedBranches to be true")
	}
}

// Helper function to create a mock API getter for testing
func NewMockAPIGetter() *MockAPIGetter {
	return &MockAPIGetter{
		CreatedEnvironments: []data.ImportedEnvironment{},
	}
}

// MockAPIGetter implements the minimal set of methods needed for testing
type MockAPIGetter struct {
	CreatedEnvironments []data.ImportedEnvironment
}

// Implement the required methods from the APIGetter interface
func (m *MockAPIGetter) CreateEnvironmentList(filedata [][]string) []data.ImportedEnvironment {
	// Basic implementation that creates environment data from CSV rows
	result := []data.ImportedEnvironment{}

	// Skip header row
	for i := 1; i < len(filedata); i++ {
		row := filedata[i]
		if len(row) < 8 {
			continue
		}

		env := data.ImportedEnvironment{
			RepositoryName:    row[0],
			EnvironmentName:   row[2],
			AdminBypass:       row[3],
			PreventSelfReview: row[6] == "true",
		}

		// Convert WaitTimer from string to int
		if row[4] != "" {
			waitTimer := 0
			// In real code, you'd handle this error
			waitTimer, _ = strconv.Atoi(row[4])
			env.WaitTimer = waitTimer
		}

		m.CreatedEnvironments = append(m.CreatedEnvironments, env)
		result = append(result, env)
	}

	return result
}

func (m *MockAPIGetter) CreateEnvironment(owner string, repo string, env string, data io.Reader) error {
	// This is a mock implementation that does nothing
	return nil
}

func (m *MockAPIGetter) CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error {
	// This is a mock implementation that does nothing
	return nil
}
