package service

import "github.com/spf13/cobra"

// Agent Agent operations
type Agent interface {
	Apply()
	ProcessId(processId string) Agent
	Name(name string) Agent
	Flag(flag string) Agent
	OwnerReferenceId(ownerId string) Agent
	Cmd(cmd *cobra.Command) Agent
	ApiServerUrl(apiServerUrl string) Agent
	Token(token string) Agent
	SkipSsl(skipSsl bool) Agent
	AgentList(agentList []string) Agent
	GetK8sObjs(name, processId string) (httpCode int, data []byte, err error)
	GetPodsByDeployment(name, processId, ownerId string) (httpCode int, data []byte, err error)
	GetPodsByDaemonSet(name, processId, ownerId string) (httpCode int, data []byte, err error)
	GetPodsByReplicaSet(name, processId, ownerId string) (httpCode int, data []byte, err error)
	GetPodsByStatefulSet(name, processId, ownerId string) (httpCode int, data []byte, err error)
}
