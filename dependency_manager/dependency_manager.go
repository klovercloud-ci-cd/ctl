package dependency_manager

import (
	"github.com/klovercloud-ci/ctl/v1/business"
	"github.com/klovercloud-ci/ctl/v1/service"
)

func GetPipelineService() service.Pipeline {
	return business.NewPipelineService(business.NewHttpClientService())
}
func GetV1GithubService() service.Git {
	return business.NewGithubService(business.NewHttpClientService())
}
