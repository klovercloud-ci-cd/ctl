package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type githubService struct {
	httpClient service.HttpClient
}

func (g githubService) Apply(git v1.Git, companyId string) error {
	err := git.Validate()
	if err != nil {
		return err
	}
	header := make(map[string]string)
	token, _ := v1.GetToken()
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(git)
	if err != nil {
		return err
	}
	_, _, err = g.httpClient.Post(v1.GetApiServerUrl()+"githubs?companyId="+companyId, header, b)
	if err != nil {
		return err
	}
	return nil
}

// NewGithubService returns git type service
func NewGithubService(httpClient service.HttpClient) service.Git {
	return &githubService{
		httpClient: httpClient,
	}
}
