package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type agentService struct {
	httpClient   service.HttpClient
	processId    string
	name         string
	flag         string
	ownerId      string
	cmd          *cobra.Command
	apiServerUrl string
	token        string
	skipSsl      bool
	agentList    []string
}

func (a agentService) ProcessId(processId string) service.Agent {
	a.processId = processId
	return a
}

func (a agentService) Name(name string) service.Agent {
	a.name = name
	return a
}

func (a agentService) Flag(flag string) service.Agent {
	a.flag = flag
	return a
}

func (a agentService) OwnerReferenceId(ownerId string) service.Agent {
	a.ownerId = ownerId
	return a
}

func (a agentService) Cmd(cmd *cobra.Command) service.Agent {
	a.cmd = cmd
	return a
}

func (a agentService) ApiServerUrl(apiServerUrl string) service.Agent {
	a.apiServerUrl = apiServerUrl
	return a
}

func (a agentService) Token(token string) service.Agent {
	a.token = token
	return a
}

func (a agentService) SkipSsl(skipSsl bool) service.Agent {
	a.skipSsl = skipSsl
	return a
}

func (a agentService) AgentList(agentList []string) service.Agent {
	a.agentList = agentList
	return a
}

func (a agentService) Apply() {
	switch a.flag {
	case string(enums.GET_K8SOBJS):
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeader([]string{"Kind", "Name", "Namespace", "Available Replicas", "Unavailable Replicas", "Agent"})
		for _, agent := range a.agentList {
			code, data, err := a.GetK8sObjs(agent, a.processId)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else if code != 200 {
				a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
			} else if data != nil {
				var responseDTO v1.ResponseDTO
				err := json.Unmarshal(data, &responseDTO)
				if err != nil {
					a.cmd.Println("[ERROR]: ", err.Error())
				} else {
					jsonString, _ := json.Marshal(responseDTO.Data)
					var k8sobjs v1.K8sObjsInfo
					err := json.Unmarshal(jsonString, &k8sobjs)
					if err != nil {
						a.cmd.Println("[ERROR]: ", err.Error())
					} else {
						var tableData [][]string
						for _, each := range k8sobjs.Deployments {
							deployment := []string{each.Kind, each.Name, each.Namespace, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas)), agent}
							tableData = append(tableData, deployment)
						}
						for _, each := range k8sobjs.ReplicaSets {
							repplicaSet := []string{each.Kind, each.Name, each.Namespace, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas)), agent}
							tableData = append(tableData, repplicaSet)
						}
						for _, each := range k8sobjs.DaemonSets {
							daemonSet := []string{each.Kind, each.Name, each.Namespace, strconv.Itoa(int(each.NumberAvailable)), strconv.Itoa(int(each.NumberUnavailable)), agent}
							tableData = append(tableData, daemonSet)
						}
						for _, each := range k8sobjs.StatefulSets {
							statefulSet := []string{each.Kind, each.Name, each.Namespace, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas)), agent}
							tableData = append(tableData, statefulSet)
						}
						for _, each := range k8sobjs.Certificates {
							certificate := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, certificate)
						}
						for _, each := range k8sobjs.ClusterRoles {
							clusterRole := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, clusterRole)
						}
						for _, each := range k8sobjs.ClusterRoleBindings {
							clusterRoleBinding := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, clusterRoleBinding)
						}
						for _, each := range k8sobjs.ConfigMaps {
							configMap := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, configMap)
						}
						for _, each := range k8sobjs.Ingresses {
							ingress := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, ingress)
						}
						for _, each := range k8sobjs.Namespaces {
							namespace := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, namespace)
						}
						for _, each := range k8sobjs.NetworkPolicies {
							networkPolicy := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, networkPolicy)
						}
						for _, each := range k8sobjs.Nodes {
							node := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, node)
						}
						for _, each := range k8sobjs.PersistentVolumes {
							persistentVolume := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, persistentVolume)
						}
						for _, each := range k8sobjs.PersistentVolumeClaims {
							persistentVolumeClaim := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, persistentVolumeClaim)
						}
						for _, each := range k8sobjs.Roles {
							role := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, role)
						}
						for _, each := range k8sobjs.RoleBindings {
							roleBinding := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, roleBinding)
						}
						for _, each := range k8sobjs.Secrets {
							secret := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, secret)
						}
						for _, each := range k8sobjs.Services {
							service := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, service)
						}
						for _, each := range k8sobjs.ServiceAccounts {
							serviceAccount := []string{each.Kind, each.Name, each.Namespace, "N/A", "N/A", agent}
							tableData = append(tableData, serviceAccount)
						}
						table.AppendBulk(tableData)
					}
				}
			}
		}
		table.Render()
	case string(enums.GET_PODS_BY_DEPLOYMENT):
		code, data, err := a.GetPodsByDeployment(a.name, a.processId, a.ownerId)
		if err != nil {
			a.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var pods v1.K8sPods
				err := json.Unmarshal(jsonString, &pods)
				if err != nil {
					a.cmd.Println("[ERROR]: ", err.Error())
				} else {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Api Version", "Kind", "Name", "Namespace", "Status"})
					for _, eachPod := range pods {
						pod := []string{"api/v1", eachPod.Kind, eachPod.Name, eachPod.Namespace, getPodStatus(eachPod)}
						table.Append(pod)
					}
					table.Render()
				}
			}
		}
	case string(enums.GET_PODS_BY_DAEMONSET):
		code, data, err := a.GetPodsByDaemonSet(a.name, a.processId, a.ownerId)
		if err != nil {
			a.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var pods v1.K8sPods
				err := json.Unmarshal(jsonString, &pods)
				if err != nil {
					a.cmd.Println("[ERROR]: ", err.Error())
				} else {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Api Version", "Kind", "Name", "Namespace", "Status"})
					for _, eachPod := range pods {
						pod := []string{"api/v1", eachPod.Kind, eachPod.Name, eachPod.Namespace, getPodStatus(eachPod)}
						table.Append(pod)
					}
					table.Render()
				}
			}
		}
	case string(enums.GET_PODS_BY_REPLICASET):
		code, data, err := a.GetPodsByReplicaSet(a.name, a.processId, a.ownerId)
		if err != nil {
			a.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var pods v1.K8sPods
				err := json.Unmarshal(jsonString, &pods)
				if err != nil {
					a.cmd.Println("[ERROR]: ", err.Error())
				} else {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Api Version", "Kind", "Name", "Namespace", "Status"})
					for _, eachPod := range pods {
						pod := []string{"api/v1", eachPod.Kind, eachPod.Name, eachPod.Namespace, getPodStatus(eachPod)}
						table.Append(pod)
					}
					table.Render()
				}
			}
		}
	case string(enums.GET_PODS_BY_STATEFULSET):
		code, data, err := a.GetPodsByStatefulSet(a.name, a.processId, a.ownerId)
		if err != nil {
			a.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			a.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				a.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var pods v1.K8sPods
				err := json.Unmarshal(jsonString, &pods)
				if err != nil {
					a.cmd.Println("[ERROR]: ", err.Error())
				} else {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Api Version", "Kind", "Name", "Namespace", "Status"})
					for _, eachPod := range pods {
						pod := []string{"api/v1", eachPod.Kind, eachPod.Name, eachPod.Namespace, getPodStatus(eachPod)}
						table.Append(pod)
					}
					table.Render()
				}
			}
		}
	}
}

func getPodStatus(pod v1.K8sPod) string {
	if pod.Status.Phase == v1.PodSucceeded || len(pod.Status.ContainerStatuses) == 0 {
		return string(pod.Status.Phase)
	}
	running := true
	var status string
	for _, eachContainerStatus := range pod.Status.ContainerStatuses {
		if eachContainerStatus.State.Terminated != nil {
			status = eachContainerStatus.State.Terminated.Reason
			running = false
			break
		} else if eachContainerStatus.State.Waiting != nil {
			status = eachContainerStatus.State.Waiting.Reason
			running = false
			if eachContainerStatus.State.Waiting.Reason == "ImagePullBackOff" || eachContainerStatus.State.Waiting.Reason == "CrashLoopBackOff" {
				break
			}
		}
	}
	if running {
		status = string(v1.PodRunning)
	}
	return status
}

func (a agentService) GetK8sObjs(name, processId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/k8sobjs?processId="+processId, header, a.skipSsl)
}

func (a agentService) GetPodsByDeployment(name, processId, ownerId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/deployments/"+ownerId+"/pods?processId="+processId, header, a.skipSsl)
}

func (a agentService) GetPodsByDaemonSet(name, processId, ownerId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/daemonSets/"+ownerId+"/pods?processId="+processId, header, a.skipSsl)
}

func (a agentService) GetPodsByReplicaSet(name, processId, ownerId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/replicaSets/"+ownerId+"/pods?processId="+processId, header, a.skipSsl)
}

func (a agentService) GetPodsByStatefulSet(name, processId, ownerId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/statefulSets/"+ownerId+"/pods?processId="+processId, header, a.skipSsl)
}

// NewAgentService returns agent type service
func NewAgentService(httpClient service.HttpClient) service.Agent {
	return &agentService{
		httpClient: httpClient,
	}
}
