package business

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"os"
)

type repositoryService struct {
	httpClient service.HttpClient
	flag string
	companyId string
	repo string
}

func (r repositoryService) Flag(flag string) service.Repository {
	r.flag=flag
	return r
}

func (r repositoryService) CompanyId(companyId string) service.Repository {
	r.companyId=companyId
	return r
}

func (r repositoryService) Repo(repoId string) service.Repository {
	r.repo=repoId
	return r
}

func (r repositoryService) Apply() {
	var cmd *cobra.Command
	switch r.flag {
	case string(enums.GET_REPOSITORY):
		code, data, err := r.GetRepositoryById(r.repo)
		if err != nil {
			cmd.Println("[ERROR]: ", err.Error())
		}
		if code != 200 {
			cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		}
		if data != nil {
			cmd.Println(string(data))
		}
	case string(enums.GET_APPLICATIONS):
		code, data, err := r.GetApplicationsByCompanyId(r.companyId)
		if err != nil {
			cmd.Println("[ERROR]: ", err.Error())
		}
		if code != 200 {
			cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		}
		if data != nil {
			cmd.Println(string(data))
		}
	}

}

func (r repositoryService) GetRepositoryById(repositoryId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(config.ApiServerUrl+"repositories/"+repositoryId, header)
}

func (r repositoryService) GetApplicationsByCompanyId(companyId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(config.ApiServerUrl+"repositories/"+companyId+"/applications", header)
}

// NewCompanyService returns repository type service
func NewRepositoryService(httpClient service.HttpClient) service.Repository {
	return &repositoryService{
		httpClient: httpClient,
	}
}
