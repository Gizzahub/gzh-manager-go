package cmd

import "github.com/spf13/cobra"

func newSetcloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "setclone",
		Short:        "setclone subcommand to manage GitHub org repositories",
		SilenceUsage: true,
	}

	cmd.AddCommand(newSetcloneGiteaCmd())
	cmd.AddCommand(newSetcloneGithubCmd())
	cmd.AddCommand(newSetcloneGitlabCmd())

	return cmd
}
