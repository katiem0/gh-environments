package variables

import (
	listCmd "github.com/katiem0/gh-environments/cmd/variables/list"
	"github.com/spf13/cobra"
)

func NewCmdVariables() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "variables <command>",
		Short: "List and Create Environment variables for an organization, repository, or environment.",
		Long:  "List and Create Environment variables for an organization, repository, or environment.",
	}
	cmd.Flags().Bool("help", false, "Show help for command")

	cmd.AddCommand(listCmd.NewCmdList())
	//cmd.AddCommand(createCmd.NewCmdCreate())

	return cmd
}
