package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type bitbucketService struct {
	httpClient service.HttpClient
}

func (b bitbucketService) Apply(git v1.Git, companyId, apiServerUrl, token string, skipSsl bool) error {
	err := git.Validate()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	body, err := json.Marshal(git)
	if err != nil {
		return err
	}
	_, _, err = b.httpClient.Post(apiServerUrl+"bitbuckets?companyId="+companyId, header, body, skipSsl)
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
