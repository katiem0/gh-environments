package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/katiem0/gh-environments/internal/data"
)

type MockAPIGetter struct {
	// Environment responses
	ReposResponse            *data.ReposQuery
	RepoResponse             *data.RepoSingleQuery
	EnvironmentsData         []byte
	BranchPoliciesData       []byte
	ProtectionRulesData      []byte
	EnvironmentSecretsData   []byte
	EnvironmentVariablesData []byte
	PublicKeyData            []byte
	EncryptedValue           string

	// Error flags for testing error scenarios
	ShouldFailGetEnvironments   bool
	ShouldFailCreateEnvironment bool
	ShouldFailCreateSecret      bool
	ShouldFailCreateVariable    bool
	ShouldFailGetPublicKey      bool
	ShouldFailEncryptSecret     bool
}

func NewMockAPIGetter() *MockAPIGetter {
	// Initialize with default test data
	m := &MockAPIGetter{}

	// Set default environment variables data
	variablesData := data.EnvVariables{
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
	m.EnvironmentVariablesData, _ = json.Marshal(variablesData)

	// Set default environment secrets data
	secretsData := data.EnvSecret{
		TotalCount: 2,
		Secrets: []data.Secret{
			{
				Name: "SECRET_1",
			},
			{
				Name: "SECRET_2",
			},
		},
	}
	m.EnvironmentSecretsData, _ = json.Marshal(secretsData)

	// Set default public key data
	publicKeyData := data.PublicKey{
		KeyID: "123456",
		Key:   "base64encodedkey==",
	}
	m.PublicKeyData, _ = json.Marshal(publicKeyData)

	// Set default encrypted value
	m.EncryptedValue = "encrypted-value"

	return m
}

func (m *MockAPIGetter) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	if m.ShouldFailGetEnvironments {
		return nil, fmt.Errorf("mock error: failed to get environments")
	}
	return m.EnvironmentsData, nil
}

func (m *MockAPIGetter) CreateEnvironment(owner string, repo string, env string, data io.Reader) error {
	if m.ShouldFailCreateEnvironment {
		return fmt.Errorf("mock error: failed to create environment")
	}
	return nil
}

func (m *MockAPIGetter) GetDeploymentBranchPolicies(owner string, repo string, env string) ([]byte, error) {
	return m.BranchPoliciesData, nil
}

func (m *MockAPIGetter) GetDeploymentProtectionRules(owner string, repo string, env string) ([]byte, error) {
	return m.ProtectionRulesData, nil
}

func (m *MockAPIGetter) GetEnvironmentSecrets(owner string, repo string, env string) ([]byte, error) {
	return m.EnvironmentSecretsData, nil
}

func (m *MockAPIGetter) GetEnvironmentVariables(owner string, repo string, env string) ([]byte, error) {
	return m.EnvironmentVariablesData, nil
}

func (m *MockAPIGetter) GetEnvironmentPublicKey(owner string, repo string, env string) ([]byte, error) {
	if m.ShouldFailGetPublicKey {
		return nil, fmt.Errorf("mock error: failed to get public key")
	}
	return m.PublicKeyData, nil
}

func (m *MockAPIGetter) EncryptSecret(publicKey string, secret string) (string, error) {
	if m.ShouldFailEncryptSecret {
		return "", fmt.Errorf("mock error: failed to encrypt secret")
	}
	return m.EncryptedValue, nil
}

// Add the missing CreateEnvironmentSecret method
func (m *MockAPIGetter) CreateEnvironmentSecret(owner string, repo string, env string, secret string, data io.Reader) error {
	if m.ShouldFailCreateSecret {
		return fmt.Errorf("mock error: failed to create environment secret")
	}
	return nil
}

// Add the missing CreateEnvironmentVariables method
func (m *MockAPIGetter) CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error {
	if m.ShouldFailCreateVariable {
		return fmt.Errorf("mock error: failed to create environment variable")
	}
	return nil
}

// Add CreateEnvironmentList method for completeness
func (m *MockAPIGetter) CreateEnvironmentList(filedata [][]string) []data.ImportedEnvironment {
	var environmentList []data.ImportedEnvironment
	// Implement basic conversion from CSV data to ImportedEnvironment
	for i, row := range filedata {
		if i == 0 { // Skip header row
			continue
		}

		if len(row) < 3 { // Need at least repo name, ID, and env name
			continue
		}

		env := data.ImportedEnvironment{
			RepositoryName:  row[0],
			EnvironmentName: row[2],
		}
		environmentList = append(environmentList, env)
	}
	return environmentList
}

// Add CreateSecretList method for completeness
func (m *MockAPIGetter) CreateSecretList(filedata [][]string) []data.ImportedSecret {
	var secretList []data.ImportedSecret
	// Implement basic conversion from CSV data to ImportedSecret
	for i, row := range filedata {
		if i == 0 { // Skip header row
			continue
		}

		if len(row) < 5 { // Need repo name, repo ID, env name, secret name, secret value
			continue
		}

		secret := data.ImportedSecret{
			RepositoryName:  row[0],
			EnvironmentName: row[2],
			Name:            row[3],
			Value:           row[4],
		}
		secretList = append(secretList, secret)
	}
	return secretList
}

// Add CreateVariableList method for completeness
func (m *MockAPIGetter) CreateVariableList(filedata [][]string) []data.ImportedVariable {
	var variableList []data.ImportedVariable
	// Implement basic conversion from CSV data to ImportedVariable
	for i, row := range filedata {
		if i == 0 { // Skip header row
			continue
		}

		if len(row) < 5 { // Need repo name, repo ID, env name, var name, var value
			continue
		}

		variable := data.ImportedVariable{
			RepositoryName:  row[0],
			EnvironmentName: row[2],
			Name:            row[3],
			Value:           row[4],
		}
		variableList = append(variableList, variable)
	}
	return variableList
}

// Add CreateDeploymentBranches method for completeness
func (m *MockAPIGetter) CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error {
	return nil
}
