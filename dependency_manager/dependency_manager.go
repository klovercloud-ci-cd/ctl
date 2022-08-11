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

func GetV1BitbucketService() service.Git {
	return business.NewBitbucketService(business.NewHttpClientService())
}

func GetCompanyService() service.Company {
	return business.NewCompanyService(business.NewHttpClientService())
}

func GetRepositoryService() service.Repository {
	return business.NewRepositoryService(business.NewHttpClientService())
}

func GetApplicationService() service.Application {
	return business.NewApplicationService(business.NewHttpClientService())
}

func GetOauthService() service.Oauth {
	return business.NewOauthService(business.NewHttpClientService())
}

func GetProcessService() service.Process {
	return business.NewProcessService(business.NewHttpClientService(), business.NewPipelineService(business.NewHttpClientService()))
}

func GetUserService() service.User {
	return business.NewUserService(business.NewHttpClientService())
}

func GetAgentService() service.Agent {
	return business.NewAgentService(business.NewHttpClientService())
}
