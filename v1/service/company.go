package service

import "github.com/spf13/cobra"

// Company Company operations
type Company interface {
	Apply()
	Flag(flag string) Company
	Company(company interface{}) Company
	CompanyId(companyId string) Company
	RepoId(repoId string) Company
	Option(option string) Company
	Cmd(cmd *cobra.Command) Company
	Kind(kind string) Company
}
