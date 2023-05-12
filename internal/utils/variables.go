package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/katiem0/gh-environments/internal/data"
)

func (g *APIGetter) CreateEnvironmentVariables(repo_id int, env string, data io.Reader) error {
	url := fmt.Sprintf("repositories/%s/environments/%s/variables", strconv.Itoa(repo_id), env)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
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
		vars.RepositoryID, _ = strconv.Atoi(each[0])
		vars.RepositoryName = each[1]
		vars.EnvironmentName = each[2]
		vars.Name = each[3]
		vars.Value = each[4]

		variableList = append(variableList, vars)
	}
	return variableList
}

func (g *APIGetter) GetEnvironmentVariables(repo_id int, env string) ([]byte, error) {
	url := fmt.Sprintf("repositories/%s/environments/%s/variables", strconv.Itoa(repo_id), env)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}
