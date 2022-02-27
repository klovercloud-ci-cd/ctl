package cmd

import (
	"bufio"
	"fmt"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	_ "golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

type LoginDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func Login() *cobra.Command{
	return &cobra.Command{
		Use:       "login",
		Short:     "Login using email and password",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var apiServerUrl string
			var securityUrl string
			for idx, each := range args {
				if strings.Contains(strings.ToLower(each), "option") {
					if idx + 1 < len(args) {
						if strings.Contains(strings.ToLower(args[idx+1]), "apiserver") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								apiServerUrl = strs[1]
							}
						} else if strings.Contains(strings.ToLower(args[idx+1]), "security") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								securityUrl = strs[1]
							}
						}
					}
				}
			}
			_, err := v1.AddOrGetUrl()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
			}
			email, password := credentials()
			loginDto := LoginDto{
				Email:    email,
				Password: password,
			}
			oauthService := dependency_manager.GetOauthService()
			ctlToken, err := oauthService.Apply(loginDto)
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			if ctlToken == "" {
				cmd.Println("[ERROR]: Something went wrong!")
				return nil
			}
			err = v1.AddToConfigFile(ctlToken, apiServerUrl, securityUrl)
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
			}
			cmd.Println("[SUCCESS]: Successfully logged in!")
			return nil
		},
	}
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Enter Password: ")
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	return strings.TrimSpace(email), strings.TrimSpace(password)
}