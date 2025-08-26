package utils

import (
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
)

func TestNewAPIGetter(t *testing.T) {
	// Create mock clients
	gqlClient := &api.GraphQLClient{}
	restClient := &api.RESTClient{}

	getter := NewAPIGetter(gqlClient, restClient)

	if getter == nil {
		t.Error("NewAPIGetter() returned nil")
		return
	}

	// Since the struct fields are not exported, we can't directly verify them
	// Instead, we'll check that the getter is functional by verifying it's not nil
	if getter == nil {
		t.Error("Expected non-nil APIGetter")
	}
}

func TestMockAPIGetter(t *testing.T) {
	// Create a mock API getter
	mockGetter := NewMockAPIGetter()

	// Set mock data
	mockGetter.EnvironmentsData = []byte(`{"total_count": 1, "environments": [{"name": "production"}]}`)

	// Call a method on the mock
	data, err := mockGetter.GetRepoEnvironments("testorg", "testrepo")

	// Verify
	if err != nil {
		t.Errorf("GetRepoEnvironments() error = %v", err)
	}

	if string(data) != `{"total_count": 1, "environments": [{"name": "production"}]}` {
		t.Errorf("Expected mock data, got %s", string(data))
	}

	// Test error scenario
	mockGetter.ShouldFailGetEnvironments = true
	_, err = mockGetter.GetRepoEnvironments("testorg", "testrepo")
	if err == nil {
		t.Error("Expected error from mock, got nil")
	}
}
