package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type bitbucketService struct {
	httpClient service.HttpClient
}

func (b bitbucketService) Apply(git v1.Git, companyId string) error {
	err := git.Validate()
	cfg := v1.GetConfigFile()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + cfg.Token
	header["Content-Type"] = "application/json"
	body, err := json.Marshal(git)
	if err != nil {
		return err
	}
	_, _, err = b.httpClient.Post(cfg.ApiServerUrl+"bitbuckets?companyId="+companyId, header, body)
	if err != nil {
		return err
	}
	return nil
}

// NewBitbucketService returns git type service
func NewBitbucketService(httpClient service.HttpClient) service.Git {
	return &bitbucketService{
		httpClient: httpClient,
	}
}
