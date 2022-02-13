package business

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/v1/service"
	"os"
)

type repositoryService struct {
	httpClient service.HttpClient
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
