package business

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"os"
	"log"
)

type repositoryService struct {
	httpClient service.HttpClient
}

func (r repositoryService) Apply(flag, repositoryId, companyId string) {
	var cmd *cobra.Command
	switch flag {
	case string(enums.GET_REPOSITORY):
		code, data, err := r.GetRepositoryById(repositoryId)
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
		code, data, err := r.GetApplicationsByCompanyId(companyId)
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
