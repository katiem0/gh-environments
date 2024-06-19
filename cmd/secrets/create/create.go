package createsecrets

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
		Use:   "create <organization> [flags]",
		Short: "Create Environment secrets.",
		Long:  "Create Environment secrets for specified environments per repository in an organization from a file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient api.GQLClient
			var restClient api.RESTClient

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
	createCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create secrets from")
	createCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	createCmd.MarkFlagRequired("from-file")

	return &createCmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var secretData [][]string
	var secretList []data.ImportedSecret

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening secrets csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		secretData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading secrets from csv file")
		}
		secretList = g.CreateSecretList(secretData)
		zap.S().Debugf("Identifying secrets list to create under %s", owner)
		zap.S().Debugf("Determining secrets to create")

		for _, secret := range secretList {

			zap.S().Debugf("Gathering secret %s for repo %s and env %s", secret.Name, secret.RepositoryName, secret.EnvironmentName)
			publicKey, err := g.GetEnvironmentPublicKey(owner, secret.RepositoryName, secret.EnvironmentName)
			if err != nil {
				zap.S().Errorf("Error arose reading secret from csv file")
			}
			var responsePublicKey data.PublicKey
			err = json.Unmarshal(publicKey, &responsePublicKey)
			if err != nil {
				return err
			}

			encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, secret.Value)
			if err != nil {
				return err
			}
			importSecret := utils.CreateSecretData(responsePublicKey.KeyID, encryptedSecret)
			createSecret, err := json.Marshal(importSecret)

			if err != nil {
				return err
			}

			reader := bytes.NewReader(createSecret)
			zap.S().Debugf("Creating secret %s under %s/%s for env %s", secret.Name, owner, secret.RepositoryName, secret.EnvironmentName)
			err = g.CreateEnvironmentSecret(owner, secret.RepositoryName, secret.EnvironmentName, secret.Name, reader)
			if err != nil {
				zap.S().Errorf("Error arose creating variable with %s", secret.Name)
			}
		}
	} else {
		zap.S().Errorf("Error arose identifying secrets")
	}

	fmt.Printf("Successfully created secrets from file: %s.", cmdFlags.fileName)
	return nil
}
