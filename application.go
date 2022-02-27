package main

import (
	"github.com/klovercloud-ci/ctl/v1/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
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
	commands.AddCommand(cmd.Describe())
	commands.AddCommand(cmd.List())
	commands.AddCommand(cmd.Update())
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
