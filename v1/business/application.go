package business

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"os"
)

type applicationService struct {
	httpClient service.HttpClient
	flag string
	companyId string
	repoId string
	applicationId string
	option string
	cmd *cobra.Command
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
			a.cmd.Println(string(data))
		}
	}
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
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(config.ApiServerUrl+"application/"+applicationId+"?companyId="+companyId+"&repositoryId="+repoId, header)
}

// NewApplicationService returns application type service
func NewApplicationService(httpClient service.HttpClient) service.Application {
	return &applicationService{
		httpClient: httpClient,
	}
}