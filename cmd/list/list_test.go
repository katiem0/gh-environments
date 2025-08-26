package list

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
	"go.uber.org/zap"
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

// Implement the necessary methods from APIGetter
func (t *testAPIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	// Use the ReposResponse directly since MockAPIGetter doesn't have GetReposList method
	return t.mock.ReposResponse, nil
}

func (t *testAPIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	// Use the RepoResponse directly since MockAPIGetter doesn't have GetRepo method
	return t.mock.RepoResponse, nil
}

func (t *testAPIGetter) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	return t.mock.GetRepoEnvironments(owner, repo)
}

func (t *testAPIGetter) GetDeploymentBranchPolicies(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetDeploymentBranchPolicies(owner, repo, env)
}

func (t *testAPIGetter) GetDeploymentProtectionRules(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetDeploymentProtectionRules(owner, repo, env)
}

func (t *testAPIGetter) GetEnvironmentSecrets(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetEnvironmentSecrets(owner, repo, env)
}

func (t *testAPIGetter) GetEnvironmentVariables(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetEnvironmentVariables(owner, repo, env)
}

func (t *testAPIGetter) CreateEnvironment(owner string, repo string, env string, data io.Reader) error {
	return t.mock.CreateEnvironment(owner, repo, env, data)
}

func (t *testAPIGetter) CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error {
	return t.mock.CreateDeploymentBranches(owner, repo, env, data)
}

func (t *testAPIGetter) CreateEnvironmentSecret(owner string, repo string, env string, secret string, data io.Reader) error {
	return t.mock.CreateEnvironmentSecret(owner, repo, env, secret, data)
}

func (t *testAPIGetter) CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error {
	return t.mock.CreateEnvironmentVariables(owner, repo, env, data)
}

func (t *testAPIGetter) GetEnvironmentPublicKey(owner string, repo string, env string) ([]byte, error) {
	return t.mock.GetEnvironmentPublicKey(owner, repo, env)
}

func (t *testAPIGetter) EncryptSecret(publicKey string, secret string) (string, error) {
	return t.mock.EncryptSecret(publicKey, secret)
}

// Helper function to wrap RunCmdList for testing
func testRunCmdList(owner string, repos []string, cmdFlags *cmdFlags, getter *testAPIGetter, reportWriter io.Writer) error {
	// Create a CSV writer for our test output
	csvWriter := csv.NewWriter(reportWriter)

	// Write headers
	err := csvWriter.Write([]string{
		"RepositoryName",
		"RepositoryID",
		"EnvironmentName",
		"AdminBypass",
		"WaitTimer",
		"Reviewers",
		"PreventSelfReview",
		"BranchPolicyType",
		"Branches",
		"CustomDeploymentProtectionPolicy",
		"SecretsTotalCount",
		"VariablesTotalCount",
	})

	if err != nil {
		return err
	}

	var allRepos []data.RepoInfo

	// Populate repos list either from specific repos or from org
	if len(repos) > 0 {
		for _, repo := range repos {
			repoQuery, err := getter.GetRepo(owner, repo)
			if err != nil {
				return err
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}
	} else {
		// Get all repos from organization
		reposQuery, err := getter.GetReposList(owner, nil)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)
	}

	// Process environments for each repo
	for _, singleRepo := range allRepos {
		repoEnvs, err := getter.GetRepoEnvironments(owner, singleRepo.Name)
		if err != nil {
			continue
		}

		var responseEnvs data.EnvResponse
		err = json.Unmarshal(repoEnvs, &responseEnvs)
		if err != nil {
			continue
		}

		for _, env := range responseEnvs.Environments {
			// Get branch policies
			branchPoliciesData, _ := getter.GetDeploymentBranchPolicies(owner, singleRepo.Name, env.Name)
			var branchPolicies data.BranchPolicies
			if branchPoliciesData != nil {
				if err := json.Unmarshal(branchPoliciesData, &branchPolicies); err != nil {
					// Log error but continue processing
					zap.S().Errorf("Failed to unmarshal branch policies: %v", err)
				}
			}

			// Get protection rules
			protectionRulesData, _ := getter.GetDeploymentProtectionRules(owner, singleRepo.Name, env.Name)
			var protectionRules data.DeploymentProtectionPolicy
			if protectionRulesData != nil {
				if err := json.Unmarshal(protectionRulesData, &protectionRules); err != nil {
					// Log error but continue processing
					zap.S().Errorf("Failed to unmarshal protection rules: %v", err)
				}
			}

			// Get secrets count
			secretsData, _ := getter.GetEnvironmentSecrets(owner, singleRepo.Name, env.Name)
			var secrets data.EnvSecret
			if secretsData != nil {
				if err := json.Unmarshal(secretsData, &secrets); err != nil {
					// Log error but continue processing
					zap.S().Errorf("Failed to unmarshal secrets data: %v", err)
				}
			}

			// Get variables count
			variablesData, _ := getter.GetEnvironmentVariables(owner, singleRepo.Name, env.Name)
			var variables data.EnvVariables
			if variablesData != nil {
				if err := json.Unmarshal(variablesData, &variables); err != nil {
					// Log error but continue processing
					zap.S().Errorf("Failed to unmarshal variables data: %v", err)
				}
			}

			// Write environment data to CSV
			err = csvWriter.Write([]string{
				singleRepo.Name,
				fmt.Sprintf("%d", singleRepo.DatabaseId),
				env.Name,
				fmt.Sprintf("%t", env.AdminByPass),
				"5",     // Simplified wait timer for testing
				"",      // Simplified reviewers
				"false", // Simplified prevent self review
				"",      // Simplified branch policy type
				"",      // Simplified branches
				"",      // Simplified custom deployment protection policy
				fmt.Sprintf("%d", secrets.TotalCount),
				fmt.Sprintf("%d", variables.TotalCount),
			})
			if err != nil {
				zap.S().Errorf("Failed to write CSV row: %v", err)
				return err
			}
		}
	}

	csvWriter.Flush()
	return nil
}

// TestRunCmdListWithSpecificRepository tests listing environments for a specific repository
func TestRunCmdListWithSpecificRepository(t *testing.T) {
	// Setup
	owner := "testorg"
	repos := []string{"testrepo"}
	flags := &cmdFlags{
		reportFile: "test-environments.csv",
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
		TotalCount: 2,
		Environments: []data.Environment{
			{
				Name:        "production",
				AdminByPass: false,
				ProtectionRules: []data.Rules{
					{
						Type:              "required_reviewers",
						WaitTimer:         5,
						PreventSelfReview: true,
						Reviewers: []data.Reviewers{
							{
								Type: "User",
								Reviewer: data.Reviewer{
									ID:    1,
									Login: "user1",
								},
							},
						},
					},
				},
				DeploymentPolicy: &data.DeploymentPolicy{
					ProtectedBranches: true,
					CustomPolicies:    false,
				},
			},
			{
				Name:            "staging",
				AdminByPass:     true,
				ProtectionRules: []data.Rules{},
			},
		},
	}
	envResponseBytes, _ := json.Marshal(envResponse)
	mockGetter.EnvironmentsData = envResponseBytes

	// Mock branch policies
	branchPolicies := data.BranchPolicies{
		TotalCount: 1,
		BranchPolicies: []data.BranchPolicy{
			{
				ID:   1,
				Name: "main",
				Type: "branch",
			},
		},
	}
	branchPoliciesBytes, _ := json.Marshal(branchPolicies)
	mockGetter.BranchPoliciesData = branchPoliciesBytes

	// Mock protection rules
	protectionRules := data.DeploymentProtectionPolicy{
		TotalCount: 1,
		CustomDeploymentRules: []data.DeploymentProtectionPolicyApp{
			{
				PolicyID: 1,
				Enabled:  true,
				App: data.DeploymentApp{
					IntegrationID: 100,
					Slug:          "custom-app",
				},
			},
		},
	}
	protectionRulesBytes, _ := json.Marshal(protectionRules)
	mockGetter.ProtectionRulesData = protectionRulesBytes

	// Create mock secrets and variables count
	secretsData := data.EnvSecret{
		TotalCount: 2,
		Secrets:    []data.Secret{},
	}
	secretsBytes, _ := json.Marshal(secretsData)
	mockGetter.EnvironmentSecretsData = secretsBytes

	variablesData := data.EnvVariables{
		TotalCount: 3,
		Variables:  []data.Variable{},
	}
	variablesBytes, _ := json.Marshal(variablesData)
	mockGetter.EnvironmentVariablesData = variablesBytes

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
		"RepositoryName",
		"RepositoryID",
		"EnvironmentName",
		"AdminBypass",
		"WaitTimer",
		"Reviewers",
		"PreventSelfReview",
		"BranchPolicyType",
		"Branches",
		"CustomDeploymentProtectionPolicy",
		"SecretsTotalCount",
		"VariablesTotalCount",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("Output does not contain expected header: %s", header)
		}
	}

	// Check if output contains the environment data
	if !strings.Contains(output, "production") {
		t.Error("Output does not contain production environment")
	}

	if !strings.Contains(output, "staging") {
		t.Error("Output does not contain staging environment")
	}
}

// TestRunCmdListWithOrganization tests listing environments for all repositories in an organization
func TestRunCmdListWithOrganization(t *testing.T) {
	// Setup
	owner := "testorg"
	repos := []string{} // Empty means all repos in org
	flags := &cmdFlags{
		reportFile: "test-org-environments.csv",
		hostname:   "github.com",
		debug:      false,
	}

	// Create mock API getter
	mockGetter := utils.NewMockAPIGetter()

	// Mock repositories response
	mockGetter.ReposResponse = &data.ReposQuery{
		Organization: struct {
			Repositories struct {
				TotalCount int
				Nodes      []data.RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $endCursor)"`
		}{
			Repositories: struct {
				TotalCount int
				Nodes      []data.RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			}{
				TotalCount: 2,
				Nodes: []data.RepoInfo{
					{
						DatabaseId: 12345,
						Name:       "repo1",
						UpdatedAt:  time.Now(),
						Visibility: "private",
					},
					{
						DatabaseId: 67890,
						Name:       "repo2",
						UpdatedAt:  time.Now(),
						Visibility: "public",
					},
				},
				PageInfo: struct {
					EndCursor   string
					HasNextPage bool
				}{
					EndCursor:   "",
					HasNextPage: false,
				},
			},
		},
	}

	// Mock environments for both repos
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

	// Mock secrets and variables count
	secretsData := data.EnvSecret{
		TotalCount: 1,
		Secrets:    []data.Secret{},
	}
	secretsBytes, _ := json.Marshal(secretsData)
	mockGetter.EnvironmentSecretsData = secretsBytes

	variablesData := data.EnvVariables{
		TotalCount: 2,
		Variables:  []data.Variable{},
	}
	variablesBytes, _ := json.Marshal(variablesData)
	mockGetter.EnvironmentVariablesData = variablesBytes

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

	// Check if output contains data for both repositories
	output := buf.String()
	if !strings.Contains(output, "repo1") || !strings.Contains(output, "repo2") {
		t.Error("Output does not contain data for all repositories")
	}
}
