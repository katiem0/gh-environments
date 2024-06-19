package create

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-environments/internal/data"
	"github.com/katiem0/gh-environments/internal/log"
	"github.com/katiem0/gh-environments/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	fileName string
	token    string
	hostname string
	debug    bool
}

func NewCmdCreate() *cobra.Command {
	//var repository string
	cmdFlags := cmdFlags{}
	var authToken string

	createCmd := cobra.Command{
		Use:   "create  <target organization> [flags]",
		Short: "Create environments and metadata.",
		Long:  "Create environments and metadata for specified environments per repository in an organization from a file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient api.GQLClient
			var restClient api.RESTClient
			if cmdFlags.token != "" {
				authToken = cmdFlags.token
			} else {
				t, _ := auth.TokenForHost(cmdFlags.hostname)
				authToken = t
			}

			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			gqlClient, err = gh.GQLClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github.hawkgirl-preview+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving graphql client")
				return err
			}

			restClient, err = gh.RESTClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving rest client")
				return err
			}

			owner := args[0]

			return runCmdCreate(owner, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient))
		},
	}

	// Configure flags for command

	createCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub personal access token for organization to write to (default "gh auth token")`)
	createCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	createCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create environments from")
	createCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	createCmd.MarkFlagRequired("from-file")

	return &createCmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var environmentData [][]string
	var environmentList []data.ImportedEnvironment

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening environments csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		environmentData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading environments from csv file")
		}

		environmentList = g.CreateEnvironmentList(environmentData)
		zap.S().Debugf("Identifying Environments list to create under %s", owner)
		zap.S().Debugf("Determining environments to create")

		for _, environment := range environmentList {
			fmt.Printf("Gathering environment %s for repo %s", environment.EnvironmentName, environment.RepositoryName)
			importEnv := utils.CreateEnvironmentData(environment)
			createEnvironment, err := json.Marshal(importEnv)
			if err != nil {
				return err
			}
			reader := bytes.NewReader(createEnvironment)
			zap.S().Debugf("Creating Environment %s for %s/%s", environment.EnvironmentName, owner, environment.RepositoryName)
			err = g.CreateEnvironment(owner, environment.RepositoryName, environment.EnvironmentName, reader)
			if err != nil {
				zap.S().Errorf("Error arose creating environment %s", environment.EnvironmentName)
			}
			if environment.DeploymentPolicy == "custom" {
				zap.S().Debugf("Creating Branch/Tag Deployment Policy for %s/%s/%s", owner, environment.RepositoryName, environment.EnvironmentName)
				for _, branch := range environment.Branches {
					createEnvironmentBranch, err := json.Marshal(branch)
					if err != nil {
						return err
					}
					readerBranch := bytes.NewReader(createEnvironmentBranch)
					err = g.CreateDeploymentBranches(owner, environment.RepositoryName, environment.EnvironmentName, readerBranch)
					if err != nil {
						zap.S().Errorf("Error arose creating deployment policy for %s", environment.EnvironmentName)
					}
				}
			}
		}
		// Gathering Envs for each repository listed
	} else {
		zap.S().Errorf("Error arose identifying environments")
	}
	fmt.Printf("Successfully created environments from file: %s.", cmdFlags.fileName)
	return nil
}
