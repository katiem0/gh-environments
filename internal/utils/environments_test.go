package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/katiem0/gh-environments/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestGetRepoEnvironments(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockResponse := `{"total_count": 2, "environments": [{"name": "production", "can_admins_bypass": true}, {"name": "staging", "can_admins_bypass": false}]}`
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	// Create API getter with mock client
	getter := newAPIGetterWithMockREST(mockClient)

	// Call the method
	result, err := getter.GetRepoEnvironments("testorg", "testrepo")

	// Verify
	if err != nil {
		t.Errorf("GetRepoEnvironments() error = %v", err)
	}

	// Check if response contains expected data
	var envResponse data.EnvResponse
	err = json.Unmarshal(result, &envResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if envResponse.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", envResponse.TotalCount)
	}

	if len(envResponse.Environments) != 2 {
		t.Errorf("Expected 2 environments, got %d", len(envResponse.Environments))
	}

	if envResponse.Environments[0].Name != "production" {
		t.Errorf("Expected first environment name to be 'production', got %s", envResponse.Environments[0].Name)
	}

	if !envResponse.Environments[0].AdminByPass {
		t.Errorf("Expected AdminByPass for production to be true")
	}

	if envResponse.Environments[1].Name != "staging" {
		t.Errorf("Expected second environment name to be 'staging', got %s", envResponse.Environments[1].Name)
	}

	if envResponse.Environments[1].AdminByPass {
		t.Errorf("Expected AdminByPass for staging to be false")
	}
}

func TestGetDeploymentBranchPolicies(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockResponse := `{"total_count": 1, "branch_policies": [{"id": 123, "name": "main", "type": "branch"}]}`
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the path is correct
			expectedPath := "repos/testorg/testrepo/environments/production/deployment-branch-policies"
			if !strings.Contains(path, expectedPath) {
				t.Errorf("Expected path to contain %s, got %s", expectedPath, path)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	// Create API getter with mock client
	getter := newAPIGetterWithMockREST(mockClient)

	// Call the method
	result, err := getter.GetDeploymentBranchPolicies("testorg", "testrepo", "production")

	// Verify
	if err != nil {
		t.Errorf("GetDeploymentBranchPolicies() error = %v", err)
	}

	var branchPolicies data.BranchPolicies
	err = json.Unmarshal(result, &branchPolicies)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if branchPolicies.TotalCount != 1 {
		t.Errorf("Expected TotalCount to be 1, got %d", branchPolicies.TotalCount)
	}

	if len(branchPolicies.BranchPolicies) != 1 {
		t.Errorf("Expected 1 branch policy, got %d", len(branchPolicies.BranchPolicies))
	}

	if branchPolicies.BranchPolicies[0].ID != 123 {
		t.Errorf("Expected branch policy ID to be 123, got %d", branchPolicies.BranchPolicies[0].ID)
	}

	if branchPolicies.BranchPolicies[0].Name != "main" {
		t.Errorf("Expected branch policy name to be 'main', got %s", branchPolicies.BranchPolicies[0].Name)
	}
}

func TestCreateEnvironmentList(t *testing.T) {
	// Setup
	g := &APIGetter{}
	filedata := [][]string{
		{"RepositoryName", "RepositoryID", "EnvironmentName", "AdminBypass", "WaitTimer", "Reviewers", "PreventSelfReview", "BranchPolicyType", "Branches", "CustomDeploymentProtectionPolicy", "SecretsTotalCount", "VariablesTotalCount"},
		{"testrepo", "12345", "production", "false", "5", "User;user1;1", "true", "protected", "", "", "", ""},
		{"testrepo", "12345", "staging", "true", "0", "", "false", "custom", "main;branch", "", "", ""},
	}

	// Execute
	result := g.CreateEnvironmentList(filedata)

	// Verify
	if len(result) != 2 {
		t.Errorf("Expected 2 environments, got %d", len(result))
		return
	}

	// Check first environment
	if result[0].RepositoryName != "testrepo" {
		t.Errorf("Expected RepositoryName testrepo, got %s", result[0].RepositoryName)
	}

	if result[0].RepositoryID != 12345 {
		t.Errorf("Expected RepositoryID 12345, got %d", result[0].RepositoryID)
	}

	if result[0].EnvironmentName != "production" {
		t.Errorf("Expected EnvironmentName production, got %s", result[0].EnvironmentName)
	}

	if result[0].AdminBypass != "false" {
		t.Errorf("Expected AdminBypass false, got %s", result[0].AdminBypass)
	}

	if result[0].WaitTimer != 5 {
		t.Errorf("Expected WaitTimer 5, got %d", result[0].WaitTimer)
	}

	if !result[0].PreventSelfReview {
		t.Error("Expected PreventSelfReview to be true")
	}

	if result[0].DeploymentPolicy != "protected" {
		t.Errorf("Expected DeploymentPolicy protected, got %s", result[0].DeploymentPolicy)
	}

	// Check second environment
	if result[1].EnvironmentName != "staging" {
		t.Errorf("Expected EnvironmentName staging, got %s", result[1].EnvironmentName)
	}

	if result[1].AdminBypass != "true" {
		t.Errorf("Expected AdminBypass true, got %s", result[1].AdminBypass)
	}

	if result[1].DeploymentPolicy != "custom" {
		t.Errorf("Expected DeploymentPolicy custom, got %s", result[1].DeploymentPolicy)
	}
}

func TestCreateEnvironment(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the request method and path
			if method != "PUT" {
				t.Errorf("Expected method PUT, got %s", method)
			}

			expectedPath := "repos/testorg/testrepo/environments/production"
			if !strings.Contains(path, expectedPath) {
				t.Errorf("Expected path to contain %s, got %s", expectedPath, path)
			}

			// Read and verify request body
			bodyBytes, err := io.ReadAll(body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
			}

			var envData data.CreateEnvironment
			err = json.Unmarshal(bodyBytes, &envData)
			if err != nil {
				t.Errorf("Failed to unmarshal request body: %v", err)
			}

			if envData.WaitTimer != 5 {
				t.Errorf("Expected WaitTimer to be 5, got %d", envData.WaitTimer)
			}

			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		},
	}

	// Create API getter with mock client
	getter := newAPIGetterWithMockREST(mockClient)

	// Create test data
	envData := data.CreateEnvironment{
		WaitTimer:         5,
		PreventSelfReview: true,
		Reviewers: []data.CreateReviewer{
			{Type: "User", ID: 1},
		},
		DeploymentBranchPolicy: &data.DeploymentPolicy{
			ProtectedBranches: true,
			CustomPolicies:    false,
		},
	}
	jsonData, _ := json.Marshal(envData)

	// Call the method
	err := getter.CreateEnvironment("testorg", "testrepo", "production", bytes.NewReader(jsonData))

	// Verify
	if err != nil {
		t.Errorf("CreateEnvironment() error = %v", err)
	}
}
