package service

import "github.com/spf13/cobra"

// Application Application operations
type Application interface {
	Apply()
	Flag(flag string) Application
	CompanyId(companyId string) Application
	RepoId(repoId string) Application
	ApplicationId(applicationId string) Application
	Option(option string) Application
	Cmd(cmd *cobra.Command) Application
	Kind(kind string) Application
}