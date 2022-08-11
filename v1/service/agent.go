package service

import "github.com/spf13/cobra"

// Agent Agent operations
type Agent interface {
	Apply()
	ProcessId(processId string) Agent
	Name(name string) Agent
	Flag(flag string) Agent
	Cmd(cmd *cobra.Command) Agent
	ApiServerUrl(apiServerUrl string) Agent
	Token(token string) Agent
	SkipSsl(skipSsl bool) Agent
	GetK8sObjs(name, processId string) (httpCode int, data []byte, err error)
}
