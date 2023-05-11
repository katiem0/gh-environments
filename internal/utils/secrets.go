package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

func (g *APIGetter) CreateEnvironmentSecret(repo_id int, env string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repositories/%s/environments/%s/secrets/%s", strconv.Itoa(repo_id), env, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) GetEnvironmentPublicKey(repo_id int, env string) ([]byte, error) {
	url := fmt.Sprintf("repositories/%s/environments/%s/secrets/public-key", strconv.Itoa(repo_id), env)
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

func (g *APIGetter) GetEnvironmentSecrets(repo_id int, env string) ([]byte, error) {
	url := fmt.Sprintf("repositories/%s/environments/%s/secrets", strconv.Itoa(repo_id), env)
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
