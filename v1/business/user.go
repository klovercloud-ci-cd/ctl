package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"log"
)

type userService struct {
	httpClient service.HttpClient
	flag string
	companyId string
	user v1.UserRegistrationDto
	cmd *cobra.Command
}

func (u userService) User(user v1.UserRegistrationDto) service.User {
	u.user = user
	return u
}

func (u userService) Flag(flag string) service.User {
	u.flag = flag
	return u
}

func (u userService) CompanyId(companyId string) service.User {
	u.companyId = companyId
	return u
}

func (u userService) Cmd(cmd *cobra.Command) service.User {
	u.cmd = cmd
	return u
}

func (u userService) Apply() {
	switch u.flag {
	case string(enums.CREATE_USER):
		err := u.CreateUser(u.user)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.CREATE_ADMIN):
		err := u.CreateAdmin(u.user)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	default:
		u.cmd.Println("[ERROR]: ", "Please provide valid options")
	}
}

func (u userService) CreateUser(user v1.UserRegistrationDto) error {
	header := make(map[string]string)
	token, _ := v1.GetToken()
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, _, err = u.httpClient.Post(v1.GetSecurityUrl()+"users?action=create_user", header, b)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) CreateAdmin(user interface{}) interface{} {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(u.user)
	if err != nil {
		return err
	}
	_, _, err = u.httpClient.Post(v1.GetSecurityUrl()+"users", header, b)
	if err != nil {
		return err
	}
	return nil
}

// NewUseryService returns user type service
func NewUserService(httpClient service.HttpClient) service.User {
	return &userService{
		httpClient: httpClient,
	}
}