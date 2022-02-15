package cmd

import (
	"bufio"
	"fmt"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	_ "golang.org/x/term"
	"log"
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
			log.Println(ctlToken)
			os.Setenv("CTL_TOKEN", ctlToken)
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