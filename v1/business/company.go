package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"os"
	"log"
)

type companyService struct {
	httpClient service.HttpClient
}

func (c companyService) Apply(flag string, company interface{}, companyId, repoId, option string) {
	var cmd *cobra.Command
	switch flag {
	case string(enums.CREATE_COMPANY):
		err := c.CreateCompany(company)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.UPDATE_REPOSITORIES):
		err :=  c.UpdateRepositoriesByCompanyId(company, companyId, option)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.UPDATE_APPLICATIONS):
		err :=  c.UpdateApplicationsByRepositoryId(company, companyId, repoId, option)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.GET_COMPANY_BY_ID):
		code, data, err := c.GetCompanyById(companyId)
		if err != nil {
			cmd.Println("[ERROR]: ", err.Error())
		}
		if code != 200 {
			cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		}
		if data != nil {
			cmd.Println(string(data))
		}
	case string(enums.GET_COMPANIES):
		code, data, err := c.GetCompanies()
		if err != nil {
			cmd.Println("[ERROR]: ", err.Error())
		}
		if code != 200 {
			cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		}
		if data != nil {
			cmd.Println(string(data))
		}
	case string(enums.GET_REPOSITORIES):
		code, data, err := c.GetRepositoriesByCompanyId(companyId)
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

func (c companyService) CreateCompany(company interface{}) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, _, err = c.httpClient.Post(config.ApiServerUrl+"companies", header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) UpdateRepositoriesByCompanyId(company interface{}, companyId string, option string) error {
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

func (c companyService) UpdateApplicationsByRepositoryId(company interface{}, companyId string, repoId string, option string) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, _, err = c.httpClient.Post(config.ApiServerUrl+"applications?companyId="+companyId+"&repositoryId="+repoId+"&companyUpdateOption="+option, header, b)
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

func (c companyService) GetRepositoriesByCompanyId(companyId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies/"+companyId+"/repositories", header)
}

// NewCompanyService returns company type service
func NewCompanyService(httpClient service.HttpClient) service.Company {
	return &companyService{
		httpClient: httpClient,
	}
}