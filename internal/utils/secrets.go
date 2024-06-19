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

func (g *APIGetter) CreateEnvironmentSecret(owner string, repo string, env string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets/%s", owner, repo, env, secret)

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

func (g *APIGetter) EncryptSecret(publicKey string, secret string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	var decodedKey [32]byte
	copy(decodedKey[:], bytes)

	encrypted, err := box.SealAnonymous(nil, []byte(secret), (*[32]byte)(bytes), rand.Reader)
	if err != nil {
		return "", err
	}
	// Encode the encrypted value in base64
	encryptedValue := base64.StdEncoding.EncodeToString(encrypted)

	return encryptedValue, nil
}

func (g *APIGetter) GetEnvironmentPublicKey(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets/public-key", owner, repo, env)
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

func (g *APIGetter) GetEnvironmentSecrets(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/secrets", owner, repo, env)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}
