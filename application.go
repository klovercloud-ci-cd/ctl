package main

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/v1/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	config.InitEnvironmentVariables()
	cli()
}
func cli() {
	commands := &cobra.Command{
		Use:          "ctl",
		Short:        "Cli to use klovercloud-ci apis!",
		Version:      "v1",
		SilenceUsage: true,
	}
	commands.AddCommand(cmd.Login())
	commands.AddCommand(cmd.GetLogs())
	commands.AddCommand(cmd.Trigger())
	commands.AddCommand(cmd.GetRepositoriesByCompanyId())
	commands.AddCommand(cmd.UpdateRepositories())
	commands.AddCommand(cmd.UpdateApplicationsByRepositoryId())
	commands.AddCommand(cmd.GetApplicationsByCompanyId())

	commands.AddCommand(cmd.Describe())
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
