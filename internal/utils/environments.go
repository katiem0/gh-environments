package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/katiem0/gh-environments/internal/data"
)

func (g *APIGetter) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments", owner, repo)
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

func (g *APIGetter) GetDeploymentBranchPolicies(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment-branch-policies", owner, repo, env)
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

func (g *APIGetter) CreateEnvironmentList(fileData [][]string) []data.ImportedEnvironment {
	// convert csv lines to array of structs
	var environmentList []data.ImportedEnvironment
	var envs data.ImportedEnvironment

	for _, each := range fileData[1:] {
		var reviewers []data.Reviewers
		var branches []data.CreateDeploymentBranch

		envs.RepositoryName = each[0]
		envs.RepositoryID, _ = strconv.Atoi(each[1])
		envs.EnvironmentName = each[2]
		envs.AdminBypass = each[3]
		envs.WaitTimer, _ = strconv.Atoi(each[4])

		reviewersData := each[5]
		reviewersParts := strings.Split(reviewersData, "|")
		for _, part := range reviewersParts {
			fields := strings.Split(part, ";")

			if len(fields) == 3 {
				id, _ := strconv.Atoi(fields[2])

				reviewer := data.Reviewer{
					Login: fields[1],
					ID:    id,
				}
				reviewerStruct := data.Reviewers{
					Type:     fields[0],
					Reviewer: reviewer,
				}
				reviewers = append(reviewers, reviewerStruct)
			} else {
				log.Printf("No reviewers listed for environment %v", envs.EnvironmentName)
			}
		}
		envs.Reviewers = reviewers
		envs.PreventSelfReview, _ = strconv.ParseBool(each[6])
		envs.DeploymentPolicy = each[7]

		deploymentData := each[8]
		deploymentDataParts := strings.Split(deploymentData, "|")
		for _, policy := range deploymentDataParts {
			policyFields := strings.Split(policy, ";")

			if len(policyFields) == 2 {
				branch := data.CreateDeploymentBranch{
					Name: policyFields[0],
					Type: policyFields[1],
				}
				branches = append(branches, branch)
			} else {
				log.Printf("No branches listed for environment %v", envs.EnvironmentName)
			}
		}
		envs.Branches = branches
		environmentList = append(environmentList, envs)
	}
	return environmentList
}

func CreateEnvironmentData(environment data.ImportedEnvironment) *data.CreateEnvironment {
	var createReviewers []data.CreateReviewer
	var deploymentPolicy *data.DeploymentPolicy

	for _, reviewerStruct := range environment.Reviewers {
		cr := data.CreateReviewer{
			ID:   reviewerStruct.Reviewer.ID,
			Type: reviewerStruct.Type,
		}
		createReviewers = append(createReviewers, cr)
	}
	if environment.DeploymentPolicy == "protected" {
		deploymentPolicy = &data.DeploymentPolicy{
			ProtectedBranches: true,
			CustomPolicies:    false,
		}
	} else if environment.DeploymentPolicy == "custom" {
		deploymentPolicy = &data.DeploymentPolicy{
			ProtectedBranches: false,
			CustomPolicies:    true,
		}
	}

	s := data.CreateEnvironment{
		WaitTimer:              environment.WaitTimer,
		PreventSelfReview:      environment.PreventSelfReview,
		Reviewers:              createReviewers,
		DeploymentBranchPolicy: deploymentPolicy,
	}
	return &s
}

func (g *APIGetter) CreateEnvironment(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, env)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) CreateDeploymentBranches(owner string, repo string, env string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment-branch-policies", owner, repo, env)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) GetDeploymentProtectionRules(owner string, repo string, env string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments/%s/deployment_protection_rules", owner, repo, env)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body response data read error, %v", err)
	}
	return responseData, err
}
