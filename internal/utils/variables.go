package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

func (g *APIGetter) CreateEnvironmentVariables(repo_id int, env string, data io.Reader) error {
	url := fmt.Sprintf("repositories/%s/environments/%s/variables", strconv.Itoa(repo_id), env)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
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
