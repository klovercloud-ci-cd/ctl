package service

import "github.com/spf13/cobra"

// Repository Repository operations
type Repository interface {
	Apply()
	Flag(flag string) Repository
	CompanyId(companyId string) Repository
	Repo(repoId string) Repository
	Cmd(cmd *cobra.Command) Repository
	Option(option string) Repository
}
