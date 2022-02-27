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
	user v1.UserRegistrationDto
	cmd *cobra.Command
	company interface{}
	passwordResetDto interface{}
	email string
}

func (u userService) Email(email string) service.User {
	u.email = email
	return u
}

func (u userService) PasswordResetDto(passwordResetDto interface{}) service.User {
	u.passwordResetDto = passwordResetDto
	return u
}

func (u userService) User(user v1.UserRegistrationDto) service.User {
	u.user = user
	return u
}

func (u userService) Company(company interface{}) service.User {
	u.company = company
	return u
}

func (u userService) Flag(flag string) service.User {
	u.flag = flag
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
	case string(enums.ATTACH_COMPANY):
		err := u.AttachCompany(u.company)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.RESET_PASSWORD):
		err := u.ResetPassword(u.passwordResetDto)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.FORGOT_PASSWORD):
		err := u.ForgotPassword(u.email)
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

func (u userService) CreateAdmin(user interface{}) error {
	securityUrl, err := v1.AddOrGetSecurityUrl()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(u.user)
	if err != nil {
		return err
	}
	_, _, err = u.httpClient.Post(securityUrl+"users", header, b)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) AttachCompany(company interface{}) error {
	header := make(map[string]string)
	token, _ := v1.GetToken()
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = u.httpClient.Put(v1.GetSecurityUrl()+"users?action="+string(enums.ATTACH_COMPANY), header, b)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) ResetPassword(passwordResetDto interface{}) interface{} {
	securityUrl, err := v1.AddOrGetSecurityUrl()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(passwordResetDto)
	if err != nil {
		return err
	}
	_, err = u.httpClient.Put(securityUrl+"users?action="+string(enums.RESET_PASSWORD), header, b)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) ForgotPassword(email string) error {
	securityUrl, err := v1.AddOrGetSecurityUrl()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	_, err = u.httpClient.Put(securityUrl+"users?action="+string(enums.FORGOT_PASSWORD)+"&media="+email, header, nil)
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