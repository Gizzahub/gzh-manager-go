package cmd

import (
	"fmt"
	githubpkg "github.com/gizzahub/gzh-manager-go/pkg/github"
	"github.com/spf13/cobra"
)

type setcloneGithubOptions struct {
	targetPath string
	orgName    string
}

func defaultSetcloneGithubOptions() *setcloneGithubOptions {
	return &setcloneGithubOptions{}
}

func newSetcloneGithubCmd() *cobra.Command {
	o := defaultSetcloneGithubOptions()

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Clone repositories from a GitHub organization",
		Args:  cobra.NoArgs,
		RunE:  o.run,
	}

	cmd.Flags().StringVarP(&o.targetPath, "targetPath", "t", o.targetPath, "targetPath")
	cmd.Flags().StringVarP(&o.orgName, "orgName", "o", o.orgName, "orgName")

	return cmd
}

func (o *setcloneGithubOptions) run(_ *cobra.Command, args []string) error {
	if o.targetPath == "" || o.orgName == "" {
		return fmt.Errorf("both targetPath and orgName must be specified")
	}

	err := githubpkg.RefreshAll(o.targetPath, o.orgName)
	if err != nil {
		return err
	}

	return nil
}
