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
	"strconv"
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
						if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								apiServerUrl = strs[1]
							}
						} else if strings.Contains(strings.ToLower(args[idx+1]), "security=") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								securityUrl = strs[1]
							}
						}
					}
				}
			}
			cfg := v1.GetConfigFile()
			if apiServerUrl != "" {
				cfg.ApiServerUrl = apiServerUrl
			}
			if securityUrl != "" {
				cfg.SecurityUrl = securityUrl
			} else {
				if cfg.SecurityUrl == "" {
					cmd.Println("[ERROR]: Security server url not found!")
					return nil
				}
			}
			err := cfg.Store()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			email, password := credentials()
			loginDto := LoginDto{
				Email:    email,
				Password: password,
			}
			cfg = v1.GetConfigFile()
			oauthService := dependency_manager.GetOauthService()
			ctlToken, err, code := oauthService.SecurityUrl(cfg.SecurityUrl).Apply(loginDto)
			if err != nil {
				cmd.Println("[ERROR]: " + err.Error() + " Status Code: " + strconv.Itoa(code))
				return nil
			}
			if ctlToken == "" {
				cmd.Println("[ERROR]: Something went wrong!")
				return nil
			}
			cfg.Token = ctlToken
			err = cfg.Store()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			cmd.Println("[SUCCESS]: Successfully logged in!")
			return nil
		},
	}
}

func Logout() *cobra.Command{
	return &cobra.Command{
		Use:       "logout",
		Short:     "Logout user from the ctl",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			cfg.Token = ""
			err := cfg.Store()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			cmd.Println("[SUCCESS]: Successfully logged-out!")
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
	fmt.Println()
	password := string(bytePassword)
	return strings.TrimSpace(email), strings.TrimSpace(password)
}