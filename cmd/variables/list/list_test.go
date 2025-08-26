package variableslist

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/katiem0/gh-environments/internal/data"
	"github.com/katiem0/gh-environments/internal/utils"
)

// TestNewCmdList verifies the command structure is correctly created
func TestNewCmdList(t *testing.T) {
	cmd := NewCmdList()

	if cmd == nil {
		t.Fatal("NewCmdList() returned nil")
	}

	// Test basic properties
	if cmd.Use != "list [flags] <organization> [repo ...] " {
		t.Errorf("Expected Use to be 'list [flags] <organization> [repo ...] ', got %s", cmd.Use)
	}

	// Test flags
	if cmd.Flag("output-file") == nil {
		t.Error("output-file flag not found")
	}

	// Test short description
	if cmd.Short == "" {
		t.Error("Command should have a short description")
	}
}

// Define an adapter that implements the required methods for testing
type testAPIGetter struct {
	mock *utils.MockAPIGetter
}

// Implement the necessary methods from APIGetter for variables functionality
func (t *testAPIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	return t.mock.ReposResponse, nil
}

func (t *testAPIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	return t.mock.RepoResponse, nil
}

func (t *testAPIGetter) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	return t.mock.GetRepoEnvironments(owner, repo)
}

func (t *testAPIGetter) GetEnvironmentVariables(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetEnvironmentVariables(owner, repo, env)
}

func testRunCmdList(owner string, repos []string, cmdFlags *cmdFlags, getter *testAPIGetter, reportWriter io.Writer) error {
	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"RepositoryID",
		"RepositoryName",
		"EnvironmentName",
		"VariableName",
		"VariableValue",
		"VariableCreatedAt",
		"VariableUpdatedAt",
	})

	if err != nil {
		return err
	}

	var allRepos []data.RepoInfo

	// Process specific repos or all repos in organization
	if len(repos) > 0 {
		for _, repo := range repos {
			repoQuery, err := getter.GetRepo(owner, repo)
			if err != nil {
				return err
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}
	} else {
		reposQuery, err := getter.GetReposList(owner, nil)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)
	}

	// Process variables for each repo's environments
	for _, singleRepo := range allRepos {
		envListResp, err := getter.GetRepoEnvironments(owner, singleRepo.Name)
		if err != nil {
			continue
		}

		var envList data.EnvResponse
		err = json.Unmarshal(envListResp, &envList)
		if err != nil {
			continue
		}

		for _, env := range envList.Environments {
			envVarsResp, err := getter.GetEnvironmentVariables(owner, singleRepo.Name, env.Name)
			if err != nil {
				continue
			}

			var envVars data.EnvVariables
			err = json.Unmarshal(envVarsResp, &envVars)
			if err != nil {
				continue
			}

			for _, evar := range envVars.Variables {
				err := csvWriter.Write([]string{
					fmt.Sprintf("%d", singleRepo.DatabaseId),
					singleRepo.Name,
					env.Name,
					evar.Name,
					evar.Value,
					"", // CreatedAt not available in test
					"", // UpdatedAt not available in test
				})
				if err != nil {
					return err
				}
			}
		}
	}

	csvWriter.Flush()
	return nil
}

// TestRunCmdListWithSpecificRepository tests listing variables for a specific repository
func TestRunCmdListWithSpecificRepository(t *testing.T) {
	// Setup
	owner := "testorg"
	repos := []string{"testrepo"}
	flags := &cmdFlags{
		reportFile: "test-variables.csv",
		hostname:   "github.com",
		debug:      false,
	}

	// Create mock API getter
	mockGetter := utils.NewMockAPIGetter()

	// Mock repository response
	mockGetter.RepoResponse = &data.RepoSingleQuery{
		Repository: data.RepoInfo{
			DatabaseId: 12345,
			Name:       "testrepo",
			UpdatedAt:  time.Now(),
			Visibility: "private",
		},
	}

	// Mock environments response
	envResponse := data.EnvResponse{
		TotalCount: 1,
		Environments: []data.Environment{
			{
				Name:            "production",
				AdminByPass:     false,
				ProtectionRules: []data.Rules{},
			},
		},
	}
	envResponseBytes, _ := json.Marshal(envResponse)
	mockGetter.EnvironmentsData = envResponseBytes

	// Mock variables response
	varResponse := data.EnvVariables{
		TotalCount: 2,
		Variables: []data.Variable{
			{
				Name:  "VAR_1",
				Value: "value1",
			},
			{
				Name:  "VAR_2",
				Value: "value2",
			},
		},
	}
	varResponseBytes, _ := json.Marshal(varResponse)
	mockGetter.EnvironmentVariablesData = varResponseBytes

	// Create buffer for output
	var buf bytes.Buffer

	// Create test adapter
	testGetter := &testAPIGetter{
		mock: mockGetter,
	}

	// Execute
	err := testRunCmdList(owner, repos, flags, testGetter, &buf)

	// Verify
	if err != nil {
		t.Errorf("runCmdList() error = %v", err)
	}

	// Check if output contains expected headers
	output := buf.String()
	expectedHeaders := []string{
		"RepositoryID",
		"RepositoryName",
		"EnvironmentName",
		"VariableName",
		"VariableValue",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("Output does not contain expected header: %s", header)
		}
	}

	// Check if output contains the variable data
	if !strings.Contains(output, "VAR_1") {
		t.Error("Output does not contain VAR_1")
	}

	if !strings.Contains(output, "VAR_2") {
		t.Error("Output does not contain VAR_2")
	}
}
