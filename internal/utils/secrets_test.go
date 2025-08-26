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

func TestGetEnvironmentSecrets(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockResponse := `{"total_count": 2, "secrets": [{"name": "SECRET_1"}, {"name": "SECRET_2"}]}`
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the path is correct
			expectedPath := "repos/testorg/testrepo/environments/production/secrets"
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
	result, err := getter.GetEnvironmentSecrets("testorg", "testrepo", "production")

	// Verify
	if err != nil {
		t.Errorf("GetEnvironmentSecrets() error = %v", err)
	}

	var secrets data.EnvSecret
	err = json.Unmarshal(result, &secrets)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if secrets.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", secrets.TotalCount)
	}

	if len(secrets.Secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(secrets.Secrets))
	}

	if secrets.Secrets[0].Name != "SECRET_1" {
		t.Errorf("Expected first secret name to be 'SECRET_1', got %s", secrets.Secrets[0].Name)
	}
}

func TestGetEnvironmentPublicKey(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockResponse := `{"key_id": "123456", "key": "base64encodedkey=="}`
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the path is correct
			expectedPath := "repos/testorg/testrepo/environments/production/secrets/public-key"
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
	result, err := getter.GetEnvironmentPublicKey("testorg", "testrepo", "production")

	// Verify
	if err != nil {
		t.Errorf("GetEnvironmentPublicKey() error = %v", err)
	}

	var publicKey data.PublicKey
	err = json.Unmarshal(result, &publicKey)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if publicKey.KeyID != "123456" {
		t.Errorf("Expected KeyID to be '123456', got %s", publicKey.KeyID)
	}

	if publicKey.Key != "base64encodedkey==" {
		t.Errorf("Expected Key to be 'base64encodedkey==', got %s", publicKey.Key)
	}
}

func TestCreateEnvironmentSecret(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Setup mock REST client
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Verify the request method and path
			if method != "PUT" {
				t.Errorf("Expected method PUT, got %s", method)
			}

			expectedPath := "repos/testorg/testrepo/environments/production/secrets/TEST_SECRET"
			if !strings.Contains(path, expectedPath) {
				t.Errorf("Expected path to contain %s, got %s", expectedPath, path)
			}

			// Read and verify request body
			bodyBytes, err := io.ReadAll(body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
			}

			var secretData struct {
				EncryptedValue string `json:"encrypted_value"`
				KeyID          string `json:"key_id"`
			}
			err = json.Unmarshal(bodyBytes, &secretData)
			if err != nil {
				t.Errorf("Failed to unmarshal request body: %v", err)
			}

			if secretData.KeyID != "123456" {
				t.Errorf("Expected KeyID to be 123456, got %s", secretData.KeyID)
			}

			if secretData.EncryptedValue == "" {
				t.Error("Expected EncryptedValue to be non-empty")
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
	secretData := struct {
		EncryptedValue string `json:"encrypted_value"`
		KeyID          string `json:"key_id"`
	}{
		EncryptedValue: "base64encryptedvalue==",
		KeyID:          "123456",
	}
	jsonData, _ := json.Marshal(secretData)

	// Call the method
	err := getter.CreateEnvironmentSecret("testorg", "testrepo", "production", "TEST_SECRET", bytes.NewReader(jsonData))

	// Verify
	if err != nil {
		t.Errorf("CreateEnvironmentSecret() error = %v", err)
	}
}

func TestCreateSecretList(t *testing.T) {
	// Setup
	g := &APIGetter{}
	filedata := [][]string{
		{"RepositoryName", "RepositoryID", "EnvironmentName", "SecretName", "SecretValue", "SecretCreatedAt", "SecretUpdatedAt"},
		{"testrepo", "12345", "production", "SECRET_1", "value1", "", ""},
		{"testrepo", "12345", "staging", "SECRET_2", "value2", "", ""},
	}

	// Execute
	result := g.CreateSecretList(filedata)

	// Verify
	if len(result) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(result))
		return
	}

	// Check first secret
	if result[0].RepositoryName != "testrepo" {
		t.Errorf("Expected RepositoryName testrepo, got %s", result[0].RepositoryName)
	}

	if result[0].RepositoryID != 12345 {
		t.Errorf("Expected RepositoryID 12345, got %d", result[0].RepositoryID)
	}

	if result[0].EnvironmentName != "production" {
		t.Errorf("Expected EnvironmentName production, got %s", result[0].EnvironmentName)
	}

	if result[0].Name != "SECRET_1" {
		t.Errorf("Expected SecretName SECRET_1, got %s", result[0].Name)
	}

	if result[0].Value != "value1" {
		t.Errorf("Expected SecretValue value1, got %s", result[0].Value)
	}

	// Check second secret
	if result[1].EnvironmentName != "staging" {
		t.Errorf("Expected EnvironmentName staging, got %s", result[1].EnvironmentName)
	}

	if result[1].Name != "SECRET_2" {
		t.Errorf("Expected SecretName SECRET_2, got %s", result[1].Name)
	}
}

func TestEncryptSecret(t *testing.T) {
	// Create API getter
	g := &APIGetter{}

	// Test data
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu9+CQrTDDHx/sJQzQ1ghX7pRAzNpf46VoLYf9XZwPZ4nJDD/UHkdCjYlkzDsxAHQxz2H8N45nPz5pJSHbJK+ppN9y9QyFx5aqP4xQhNbCRcjgaecOj4q4Z5JzQwu0EsrLxfzm9+FGJd5zLQvIzfWQwXbXXKs1+QudCnH1fIlJDLpk4pwkLvBQSDqyXeQ1HGZVyvv4+eiOVlovFJGv/Z5VKLxGy+MrJCMRP4EQXdoW0VxFRvyTxFwQp9KMkJJKTFIcwQNYgNbEcTL1KeFA1e0uvQjBpUkPzTnRQdUQZpHzASCHAJ7J56t/KMKVzGF8xQXMw6IcyH8pz7jdWMGNIeNgQIDAQAB"
	secretValue := "test-secret-value"

	// Call the encryption function
	encryptedValue, err := g.EncryptSecret(publicKey, secretValue)

	// Verify
	if err != nil {
		t.Errorf("EncryptSecret() error = %v", err)
	}

	if encryptedValue == "" {
		t.Error("Expected encrypted value to be non-empty")
	}

	// The encrypted value should be base64 encoded
	if !strings.HasSuffix(encryptedValue, "==") && !strings.HasSuffix(encryptedValue, "=") {
		t.Error("Expected encrypted value to be base64 encoded")
	}
}
