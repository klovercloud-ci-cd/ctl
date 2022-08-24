package service

import "github.com/spf13/cobra"

// Pipeline Pipeline operations.
type Pipeline interface {
	Apply()
	Logs(url, page, limit string) (httpCode int, data interface{}, err error)
	ProcessId(processId string) Pipeline
	Token(token string) Pipeline
	SkipSsl(skipSsl bool) Pipeline
	Flag(flag string) Pipeline
	ApiServerUrl(apiServerUrl string) Pipeline
	Cmd(cmd *cobra.Command) Pipeline
	Get(processId string, action string, url string, token string) (httpCode int, data interface{}, err error)
}
