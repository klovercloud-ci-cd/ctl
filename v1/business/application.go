package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type applicationService struct {
	httpClient service.HttpClient
	flag string
	companyId string
	repoId string
	applicationId string
	option string
	cmd *cobra.Command
	kind string
}

func (a applicationService) Apply() {
	switch a.flag {
	case string(enums.GET_APPLICATION):
		code, data, err := a.GetApplication(a.companyId, a.repoId, a.applicationId)
		if err != nil {
			a.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var application v1.Application
				json.Unmarshal(jsonString, &application)
				applicationDto := v1.ApplicationDto{
					ApiVersion:  "api/v1",
					Kind:        a.kind,
					Application: application,
				}
				b, _ := yaml.Marshal(applicationDto)
				b = v1.AddRootIndent(b, 4)
				a.cmd.Println(string(b))
			}
		}
	}
}

func (a applicationService) Kind(kind string) service.Application {
	a.kind = kind
	return a
}

func (a applicationService) Cmd(cmd *cobra.Command) service.Application {
	a.cmd = cmd
	return a
}

func (a applicationService) Flag(flag string) service.Application {
	a.flag = flag
	return a
}

func (a applicationService) CompanyId(companyId string) service.Application {
	a.companyId = companyId
	return a
}

func (a applicationService) RepoId(repoId string) service.Application {
	a.repoId = repoId
	return a
}

func (a applicationService) ApplicationId(applicationId string) service.Application {
	a.applicationId = applicationId
	return a
}

func (a applicationService) Option(option string) service.Application {
	a.option = option
	return a
}

func (a applicationService) GetApplication(companyId, repoId, applicationId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	token, _ := v1.GetToken()
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(v1.GetApiServerUrl()+"applications/"+applicationId+"?&repositoryId="+repoId, header)
}

// NewApplicationService returns application type service
func NewApplicationService(httpClient service.HttpClient) service.Application {
	return &applicationService{
		httpClient: httpClient,
	}
}