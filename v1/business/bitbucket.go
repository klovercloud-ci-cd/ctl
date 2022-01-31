package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/config"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"os"
)

type bitbucketService struct {
	httpClient service.HttpClient
}

func (b bitbucketService) Apply(git v1.Git, companyId string) error {
	err := git.Validate()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	body, err := json.Marshal(git)
	if err != nil {
		return err
	}
	_, err = b.httpClient.Post(config.IntegrationManagerUrl+"bitbuckets?companyId="+companyId, header, body)
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
