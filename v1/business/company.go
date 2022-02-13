package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/v1/service"
	"os"
)

type companyService struct {
	httpClient service.HttpClient
}

func (c companyService) Apply(company interface{}) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Post(config.ApiServerUrl+"companies", header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) ApplyUpdateRepositories(company interface{}, companyId string, option string) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Put(config.ApiServerUrl+"companies/"+companyId+"/repositories?companyUpdateOption="+option, header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) ApplyUpdateApplications(company interface{}, companyId string, repoId string, option string) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Post(config.ApiServerUrl+"applications?companyId="+companyId+"&repositoryId="+repoId+"&companyUpdateOption="+option, header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) GetCompanyById(companyId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies/"+companyId, header)
}

func (c companyService) GetCompanies() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies", header)
}

// NewCompanyService returns company type service
func NewCompanyService(httpClient service.HttpClient) service.Company {
	return &companyService{
		httpClient: httpClient,
	}
}