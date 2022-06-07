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
	SkipSsl(skipSsl bool) User
	Cmd(cmd *cobra.Command) User
	Company(company interface{}) User
	PasswordResetDto(passwordResetDto interface{}) User
	Email(email string) User
	SecurityUrl(securityUrl string) User
	Token(token string) User
}
