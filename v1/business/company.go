package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci-cd/integration-manager/core/v1"
	"github.com/klovercloud-ci/ctl/config"
	"github.com/klovercloud-ci/ctl/v1/service"
	"os"
)

type companyService struct {
	httpClient service.HttpClient
}

func (c companyService) Store(company v1.Company) error {
	err := company.Validate()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Post(config.IntegrationManagerUrl+"companies",header, b)
	if err != nil {
		return err
	}
	return nil
}

// NewCompanyService returns company type service
func NewCompanyService(httpClient service.HttpClient) service.Company {
	return &companyService{
		httpClient: httpClient,
	}
}