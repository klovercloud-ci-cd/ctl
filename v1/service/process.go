package service

import "github.com/spf13/cobra"

type Process interface {
	Apply()
	Cmd(cmd *cobra.Command) Process
	RepoId(repoId string) Process
	ApplicationId(appId string) Process
	Kind(kind string) Process
}