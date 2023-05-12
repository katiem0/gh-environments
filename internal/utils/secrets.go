package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/katiem0/gh-environments/internal/data"
	"golang.org/x/crypto/nacl/box"
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

func CreateSecretData(keyID string, encryptedValue string) *data.CreateEnvSecret {
	s := data.CreateEnvSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
	}
	return &s
}

func (g *APIGetter) CreateSecretList(filedata [][]string) []data.ImportedSecret {
	// convert csv lines to array of structs
	var secretList []data.ImportedSecret
	var secret data.ImportedSecret
	for _, each := range filedata[1:] {
		secret.RepositoryID, _ = strconv.Atoi(each[0])
		secret.RepositoryName = each[1]
		secret.EnvironmentName = each[2]
		secret.Name = each[3]
		secret.Value = each[4]

		secretList = append(secretList, secret)
	}
	return secretList
}

func (g *APIGetter) EncryptSecret(publickey string, secret string) (string, error) {
	var pkBytes [32]byte
	copy(pkBytes[:], publickey)
	secretBytes := secret

	out := make([]byte, 0,
		len(secretBytes)+
			box.Overhead+
			len(pkBytes))

	enc, err := box.SealAnonymous(
		out, []byte(secretBytes), &pkBytes, rand.Reader,
	)
	if err != nil {
		return "", err
	}

	encEnc := base64.StdEncoding.EncodeToString(enc)

	return encEnc, nil
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
