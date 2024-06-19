package secrets

import (
	createCmd "github.com/katiem0/gh-environments/cmd/secrets/create"
	listCmd "github.com/katiem0/gh-environments/cmd/secrets/list"
	"github.com/spf13/cobra"
)

func NewCmdSecrets() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "secrets <command> [flags]",
		Args:  cobra.MinimumNArgs(1),
		Short: "List and Create Environment secrets.",
		Long:  "List and Create Environment specific secrets in repositories under an organization.",
	}
	cmd.Flags().Bool("help", false, "Show help for command")
	cmd.AddCommand(listCmd.NewCmdList())
	cmd.AddCommand(createCmd.NewCmdCreate())

	return cmd
}
