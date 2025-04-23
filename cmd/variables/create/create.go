package createvariables

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
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
		Use:   "create <organization> [flags]",
		Short: "Create Environment variables.",
		Long:  "Create Environment variables for specified environments per repository in an organization from a file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient *api.GraphQLClient
			var restClient *api.RESTClient

			// Reinitialize logging if debugging was enabled
			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			if cmdFlags.token != "" {
				authToken = cmdFlags.token
			} else {
				t, _ := auth.TokenForHost(cmdFlags.hostname)
				authToken = t
			}

			gqlClient, err = api.NewGraphQLClient(api.ClientOptions{
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

			restClient, err = api.NewRESTClient(api.ClientOptions{
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
	createCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create variables from")
	createCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	if err := createCmd.MarkFlagRequired("from-file"); err != nil {
		zap.S().Errorf("Error marking flag 'from-file' as required: %v", err)
	}

	return &createCmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var variableData [][]string
	var variablesList []data.ImportedVariable

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening variables csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		variableData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading variables from csv file")
		}
		variablesList = g.CreateVariableList(variableData)
		zap.S().Debugf("Identifying Variable list to create under %s", owner)
		zap.S().Debugf("Determining variables to create")

		for _, variable := range variablesList {

			zap.S().Debugf("Gathering variable %s for repo %s and env %s", variable.Name, variable.RepositoryName, variable.EnvironmentName)
			importVar := utils.CreateVariableData(variable)
			createVariable, err := json.Marshal(importVar)
			if err != nil {
				return err
			}
			reader := bytes.NewReader(createVariable)
			zap.S().Debugf("Creating variable %s under %s/%s for env %s", variable.Name, owner, variable.RepositoryName, variable.EnvironmentName)
			err = g.CreateEnvironmentVariables(owner, variable.RepositoryName, variable.EnvironmentName, reader)
			if err != nil {
				zap.S().Errorf("Error arose creating variable with %s", variable.Name)
			}
		}
	} else {
		zap.S().Errorf("Error arose identifying variables")
	}

	fmt.Printf("Successfully created variables from file: %s.", cmdFlags.fileName)
	return nil
}
