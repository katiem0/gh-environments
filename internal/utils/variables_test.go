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

func TestGetEnvironmentVariables(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockResponse := `{"total_count": 2, "variables": [{"name": "VAR_1", "value": "value1"}, {"name": "VAR_2", "value": "value2"}]}`
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the path is correct
			expectedPath := "repos/testorg/testrepo/environments/production/variables"
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
	result, err := getter.GetEnvironmentVariables("testorg", "testrepo", "production")

	// Verify
	if err != nil {
		t.Errorf("GetEnvironmentVariables() error = %v", err)
	}

	var variables data.EnvVariables
	err = json.Unmarshal(result, &variables)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variables.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", variables.TotalCount)
	}

	if len(variables.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(variables.Variables))
	}

	if variables.Variables[0].Name != "VAR_1" {
		t.Errorf("Expected first variable name to be 'VAR_1', got %s", variables.Variables[0].Name)
	}

	if variables.Variables[0].Value != "value1" {
		t.Errorf("Expected first variable value to be 'value1', got %s", variables.Variables[0].Value)
	}
}

func TestCreateEnvironmentVariables(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the request method and path
			if method != "POST" {
				t.Errorf("Expected method POST, got %s", method)
			}

			expectedPath := "repos/testorg/testrepo/environments/production/variables"
			if !strings.Contains(path, expectedPath) {
				t.Errorf("Expected path to contain %s, got %s", expectedPath, path)
			}

			// Read and verify request body
			bodyBytes, err := io.ReadAll(body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
			}

			var varData data.CreateVariable
			err = json.Unmarshal(bodyBytes, &varData)
			if err != nil {
				t.Errorf("Failed to unmarshal request body: %v", err)
			}

			if varData.Name != "TEST_VAR" {
				t.Errorf("Expected variable name to be TEST_VAR, got %s", varData.Name)
			}

			if varData.Value != "test-value" {
				t.Errorf("Expected variable value to be test-value, got %s", varData.Value)
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
	varData := data.CreateVariable{
		Name:  "TEST_VAR",
		Value: "test-value",
	}
	jsonData, _ := json.Marshal(varData)

	// Call the method
	err := getter.CreateEnvironmentVariables("testorg", "testrepo", "production", bytes.NewReader(jsonData))

	// Verify
	if err != nil {
		t.Errorf("CreateEnvironmentVariables() error = %v", err)
	}
}

func TestCreateVariableData(t *testing.T) {
	// Create test data
	importedVar := data.ImportedVariable{
		RepositoryName:  "testrepo",
		RepositoryID:    12345,
		EnvironmentName: "production",
		Name:            "TEST_VAR",
		Value:           "test-value",
	}

	// Call the function
	result := CreateVariableData(importedVar)

	// Verify
	if result.Name != "TEST_VAR" {
		t.Errorf("Expected Name to be TEST_VAR, got %s", result.Name)
	}

	if result.Value != "test-value" {
		t.Errorf("Expected Value to be test-value, got %s", result.Value)
	}
}

func TestCreateVariableList(t *testing.T) {
	// Setup
	g := &APIGetter{}
	filedata := [][]string{
		{"RepositoryName", "RepositoryID", "EnvironmentName", "VariableName", "VariableValue", "VariableCreatedAt", "VariableUpdatedAt"},
		{"testrepo", "12345", "production", "VAR_1", "value1", "", ""},
		{"testrepo", "12345", "staging", "VAR_2", "value2", "", ""},
	}

	// Execute
	result := g.CreateVariableList(filedata)

	// Verify
	if len(result) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(result))
		return
	}

	// Check first variable
	if result[0].RepositoryName != "testrepo" {
		t.Errorf("Expected RepositoryName testrepo, got %s", result[0].RepositoryName)
	}

	if result[0].RepositoryID != 12345 {
		t.Errorf("Expected RepositoryID 12345, got %d", result[0].RepositoryID)
	}

	if result[0].EnvironmentName != "production" {
		t.Errorf("Expected EnvironmentName production, got %s", result[0].EnvironmentName)
	}

	if result[0].Name != "VAR_1" {
		t.Errorf("Expected VariableName VAR_1, got %s", result[0].Name)
	}

	if result[0].Value != "value1" {
		t.Errorf("Expected VariableValue value1, got %s", result[0].Value)
	}

	// Check second variable
	if result[1].EnvironmentName != "staging" {
		t.Errorf("Expected EnvironmentName staging, got %s", result[1].EnvironmentName)
	}

	if result[1].Name != "VAR_2" {
		t.Errorf("Expected VariableName VAR_2, got %s", result[1].Name)
	}
}
