package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/katiem0/gh-environments/internal/data"
)

func (g *APIGetter) CreateEnvironmentVariables(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()
	return err
}

func CreateVariableData(variable data.ImportedVariable) *data.CreateVariable {
	s := data.CreateVariable{
		Name:  variable.Name,
		Value: variable.Value,
	}
	return &s
}

func (g *APIGetter) CreateVariableList(filedata [][]string) []data.ImportedVariable {
	// convert csv lines to array of structs
	var variableList []data.ImportedVariable
	var vars data.ImportedVariable
	for _, each := range filedata[1:] {
		// Check if we have enough columns
		if len(each) < 5 {
			continue // Skip rows with insufficient data
		}
		vars.RepositoryName = each[0]
		repositoryID, _ := strconv.Atoi(each[1])
		vars.RepositoryID = repositoryID
		vars.EnvironmentName = each[2]
		vars.Name = each[3]
		vars.Value = each[4]

		variableList = append(variableList, vars)
	}
	return variableList
}

func (g *APIGetter) GetEnvironmentVariables(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)

	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return responseData, nil
}
