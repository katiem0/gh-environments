//go:build !cover
// +build !cover

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// mockRESTClient for testing
type mockRESTClient struct {
	RequestFunc func(method string, path string, body io.Reader) (*http.Response, error)
}

func (m *mockRESTClient) Request(method string, path string, body io.Reader) (*http.Response, error) {
	if m.RequestFunc != nil {
		return m.RequestFunc(method, path, body)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

func (m *mockRESTClient) RequestWithContext(ctx context.Context, method string, path string, body io.Reader) (*http.Response, error) {
	return m.Request(method, path, body)
}

// testAPIGetterWrapper allows us to use mock REST clients in tests
type testAPIGetterWrapper struct {
	mockClient *mockRESTClient
}

// Create a new APIGetter with a mock REST client
func newAPIGetterWithMockREST(client *mockRESTClient) *testAPIGetterWrapper {
	return &testAPIGetterWrapper{
		mockClient: client,
	}
}

func (t *testAPIGetterWrapper) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments", owner, repo)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) GetDeploymentBranchPolicies(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment-branch-policies", owner, repo, env)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) GetDeploymentProtectionRules(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment_protection_rules", owner, repo, env)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) CreateEnvironment(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, env)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

func (t *testAPIGetterWrapper) CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment-branch-policies", owner, repo, env)
	resp, err := t.mockClient.Request("POST", url, data)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

// Environment secrets methods
func (t *testAPIGetterWrapper) GetEnvironmentSecrets(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets", owner, repo, env)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) GetEnvironmentPublicKey(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets/public-key", owner, repo, env)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) CreateEnvironmentSecret(owner string, repo string, env string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets/%s", owner, repo, env, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}()
	return nil
}

// Environment variables methods
func (t *testAPIGetterWrapper) GetEnvironmentVariables(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return responseData, nil
}

func (t *testAPIGetterWrapper) CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)
	resp, err := t.mockClient.Request("POST", url, data)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error or handle it appropriately
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}()
	return nil
}
