package utils

import (
	"io"

	"github.com/cli/go-gh/pkg/api"
	"github.com/katiem0/gh-environments/internal/data"
	"github.com/shurcooL/graphql"
)

type Getter interface {
	CreateEnvironment(owner string, repo string, env string, data io.Reader) error
	CreateEnvironmentList(filedata [][]string) []data.ImportedEnvironment
	CreateEnvironmentVariables(repo_id int, env string, data io.Reader) error
	CreateEnvironmentSecret(repo_id int, env string, secret string, data io.Reader) error
	CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error
	CreateSecretList(filedata [][]string) []data.ImportedSecret
	EncryptSecret(publickey string, secret string) (string, error)
	GetDeploymentBranchPolicies(owner string, repo string, env string) ([]byte, error)
	GetDeploymentProtectionRules(owner string, repo string, env string) ([]byte, error)
	GetEnvironmentPublicKey(repo_id int, env string) ([]byte, error)
	GetEnvironmentVariables(repo_id int, env string) ([]byte, error)
	GetEnvironmentSecrets(repo_id int, env string) ([]byte, error)
	GetRepo(owner string, name string) ([]data.RepoSingleQuery, error)
	GetRepoEnvironments(owner string, repo string) ([]byte, error)
	GetReposList(owner string, endCursor *string) ([]data.ReposQuery, error)
}

type APIGetter struct {
	gqlClient  api.GQLClient
	restClient api.RESTClient
}

func NewAPIGetter(gqlClient api.GQLClient, restClient api.RESTClient) *APIGetter {
	return &APIGetter{
		gqlClient:  gqlClient,
		restClient: restClient,
	}
}

type sourceAPIGetter struct {
	restClient api.RESTClient
}

func NewSourceAPIGetter(restClient api.RESTClient) *sourceAPIGetter {
	return &sourceAPIGetter{
		restClient: restClient,
	}
}

func (g *APIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	query := new(data.ReposQuery)
	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(endCursor),
		"owner":     graphql.String(owner),
	}

	err := g.gqlClient.Query("getRepos", &query, variables)

	return query, err
}

func (g *APIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	query := new(data.RepoSingleQuery)
	variables := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(name),
	}

	err := g.gqlClient.Query("getRepo", &query, variables)
	return query, err
}
