package variableslist

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

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
	hostname   string
	token      string
	reportFile string
	debug      bool
}

func NewCmdList() *cobra.Command {
	//var repository string
	cmdFlags := cmdFlags{}
	var authToken string

	exportCmd := cobra.Command{
		Use:   "list [flags] <organization> [repo ...] ",
		Short: "Generate a report of Environment variable.",
		Long:  "Generate a report of variables for each environment per repository in an organization.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(exportCmd *cobra.Command, args []string) error {
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
			repos := args[1:]

			if _, err := os.Stat(cmdFlags.reportFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(cmdFlags.reportFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmdList(owner, repos, utils.NewAPIGetter(gqlClient, restClient), reportWriter)
		},
	}

	// Determine default report file based on current timestamp; for more info see https://pkg.go.dev/time#pkg-constants
	reportFileDefault := fmt.Sprintf("report-variables-%s.csv", time.Now().Format("20060102150405"))
	// Configure flags for command
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	exportCmd.Flags().StringVarP(&cmdFlags.reportFile, "output-file", "o", reportFileDefault, "Name of file to write CSV report")
	exportCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	//cmd.MarkPersistentFlagRequired("app")

	return &exportCmd
}

func runCmdList(owner string, repos []string, g *utils.APIGetter, reportWriter io.Writer) error {
	var reposCursor *string
	var allRepos []data.RepoInfo

	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"RepositoryID",
		"RepositoryName",
		"EnvironmentName",
		"VariableName",
		"VariableValue",
		"VariableCreatedAt",
		"VariableUpdatedAt",
	})
	if err != nil {
		zap.S().Error("Error raised in writing to csv", zap.Error(err))
	}
	if len(repos) > 0 {
		zap.S().Infof("Processing repos: %s", repos)

		for _, repo := range repos {

			zap.S().Debugf("Processing %s/%s", owner, repo)

			repoQuery, err := g.GetRepo(owner, repo)
			if err != nil {
				zap.S().Error("Error raised in getting repos", zap.Error(err))
				return err
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}

	} else {
		// Prepare writer for outputting report
		for {
			zap.S().Debugf("Processing list of repositories for %s", owner)
			reposQuery, err := g.GetReposList(owner, reposCursor)

			if err != nil {
				zap.S().Error("Error raised in processing list of repos", zap.Error(err))
				return err
			}

			allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)

			reposCursor = &reposQuery.Organization.Repositories.PageInfo.EndCursor

			if !reposQuery.Organization.Repositories.PageInfo.HasNextPage {
				break
			}
		}
	}

	// Writing to CSV environment Variables
	for _, singleRepo := range allRepos {
		// Writing to CSV repository level Actions Variables
		zap.S().Debugf("Gathering environments for %s", singleRepo.Name)
		envListResp, err := g.GetRepoEnvironments(owner, singleRepo.Name)
		if err != nil {
			zap.S().Error("Error raised in getting repo environments", zap.Error(err))
		}
		var envList data.EnvResponse
		err = json.Unmarshal(envListResp, &envList)
		if err != nil {
			return err
		}

		for _, env := range envList.Environments {
			zap.S().Debugf("Gathering environment %s variables for %s", env.Name, singleRepo.Name)
			envVarsResp, err := g.GetEnvironmentVariables(owner, singleRepo.Name, env.Name)
			if err != nil {
				zap.S().Error("Error raised in getting environment variables", zap.Error(err))
			}

			var envVars data.EnvVariables
			err = json.Unmarshal(envVarsResp, &envVars)
			if err != nil {
				return err
			}

			for _, evar := range envVars.Variables {
				err = csvWriter.Write([]string{
					strconv.Itoa(singleRepo.DatabaseId),
					singleRepo.Name,
					env.Name,
					evar.Name,
					evar.Value,
					evar.CreatedAt.Format(time.RFC3339),
					evar.UpdatedAt.Format(time.RFC3339),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
					return err
				}
			}

		}
	}

	csvWriter.Flush()
	fmt.Printf("Successfully exported variables for %s", owner)
	return nil
}
