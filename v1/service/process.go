package service

import "github.com/spf13/cobra"

type Process interface {
	Apply()
	Cmd(cmd *cobra.Command) Process
	Id(id string) Process
	Step(step string) Process
	Claim(claim string) Process
	Footmark(footmark string) Process
	RepoId(repoId string) Process
	ApplicationId(appId string) Process
	Flag(flag string) Process
	Kind(kind string) Process
	ApiServerUrl(apiServerUrl string) Process
	Token(token string) Process
	SkipSsl(skipSsl bool) Process
	GetLogsByProcessIdAndStepAndClaimAndFootmark(page, limit string) (httpCode int, data interface{}, err error)
	GetProcessLifeCycleEventByProcessIdAndStep() (httpCode int, data interface{}, err error)
}
