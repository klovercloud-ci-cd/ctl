package service

import (
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
)

// User User operations
type User interface {
	Apply()
	User(user v1.UserRegistrationDto) User
	Flag(flag string) User
	CompanyId(companyId string) User
	Cmd(cmd *cobra.Command) User
}