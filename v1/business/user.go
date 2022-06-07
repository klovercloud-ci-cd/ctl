package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"net/http"
)

type userService struct {
	httpClient       service.HttpClient
	flag             string
	user             v1.UserRegistrationDto
	cmd              *cobra.Command
	company          interface{}
	passwordResetDto interface{}
	email            string
	securityUrl      string
	token            string
	skipSsl          bool
}

func (u userService) SkipSsl(skipSsl bool) service.User {
	u.skipSsl = skipSsl
	return u
}

func (u userService) SecurityUrl(securityUrl string) service.User {
	u.securityUrl = securityUrl
	return u
}

func (u userService) Token(token string) service.User {
	u.token = token
	return u
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
		httpCode, _, err := u.CreateUser(u.user)
		if err != nil {
			u.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else {
			u.cmd.Println("Successfully Created User")
		}
	case string(enums.CREATE_ADMIN):
		httpCode, _, err := u.CreateAdmin()
		if err != nil {
			u.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else {
			u.cmd.Println("Successfully Created User")
		}
	case string(enums.ATTACH_COMPANY):
		httpCode, err := u.AttachCompany(u.company)
		if err != nil {
			u.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else {
			u.cmd.Println("[SUCCESS]: Successfully Attached Company")
		}
	case string(enums.RESET_PASSWORD):
		httpCode, err := u.ResetPassword(u.passwordResetDto)
		if err != nil {
			u.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else {
			u.cmd.Println("[SUCCESS]: Successfully Reset Password")
		}
	case string(enums.FORGOT_PASSWORD):
		httpCode, err := u.ForgotPassword(u.email)
		if err != nil {
			u.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else {
			u.cmd.Println("[SUCCESS]: Otp sent sucessfully")
		}
	default:
		u.cmd.Println("[ERROR]: ", "Please provide valid options")
	}
}

func (u userService) CreateUser(user v1.UserRegistrationDto) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + u.token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(user)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	httpCode, data, err = u.httpClient.Post(u.securityUrl+"users?action=create_user", header, b, u.skipSsl)
	if err != nil {
		return httpCode, nil, err
	}
	return httpCode, nil, nil
}

func (u userService) CreateAdmin() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(u.user)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	httpCode, data, err = u.httpClient.Post(u.securityUrl+"users", header, b, u.skipSsl)
	if err != nil {
		return httpCode, nil, err
	}
	return httpCode, nil, nil
}

func (u userService) AttachCompany(company interface{}) (httpCode int, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + u.token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return http.StatusBadRequest, err
	}
	httpCode, err = u.httpClient.Put(u.securityUrl+"users?action="+string(enums.ATTACH_COMPANY), header, b, u.skipSsl)
	if err != nil {
		return httpCode, err
	}
	return httpCode, nil
}

func (u userService) ResetPassword(passwordResetDto interface{}) (httpCode int, err error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(passwordResetDto)
	if err != nil {
		return http.StatusBadRequest, err
	}
	_, err = u.httpClient.Put(u.securityUrl+"users?action="+string(enums.RESET_PASSWORD), header, b, u.skipSsl)
	if err != nil {
		return httpCode, err
	}
	return httpCode, nil
}

func (u userService) ForgotPassword(email string) (httpCode int, err error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	httpCode, err = u.httpClient.Put(u.securityUrl+"users?action="+string(enums.FORGOT_PASSWORD)+"&media="+email, header, nil, u.skipSsl)
	if err != nil {
		return httpCode, err
	}
	return httpCode, nil
}

// NewUserService returns user type service
func NewUserService(httpClient service.HttpClient) service.User {
	return &userService{
		httpClient: httpClient,
	}
}
