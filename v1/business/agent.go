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
	cmd          *cobra.Command
	apiServerUrl string
	token        string
	skipSsl      bool
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

func (a agentService) Apply() {
	switch a.flag {
	case string(enums.GET_K8SOBJS):
		code, data, err := a.GetK8sObjs(a.name, a.processId)
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
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Kind", "Name", "Namespace", "UID", "Available Replicas", "Unavailable Replicas"})
					var tableData [][]string
					for _, each := range k8sobjs.Deployments {
						deployment := []string{each.Kind, each.Name, each.Namespace, each.UID, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas))}
						tableData = append(tableData, deployment)
					}
					for _, each := range k8sobjs.ReplicaSets {
						repplicaSet := []string{each.Kind, each.Name, each.Namespace, each.UID, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas))}
						tableData = append(tableData, repplicaSet)
					}
					for _, each := range k8sobjs.DaemonSets {
						daemonSet := []string{each.Kind, each.Name, each.Namespace, each.UID, strconv.Itoa(int(each.NumberAvailable)), strconv.Itoa(int(each.NumberUnavailable))}
						tableData = append(tableData, daemonSet)
					}
					for _, each := range k8sobjs.StatefulSets {
						statefulSet := []string{each.Kind, each.Name, each.Namespace, each.UID, strconv.Itoa(int(each.ReadyReplicas)), strconv.Itoa(int(each.Replicas - each.ReadyReplicas))}
						tableData = append(tableData, statefulSet)
					}

					for _, each := range k8sobjs.Certificates {
						certificate := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, certificate)
					}
					for _, each := range k8sobjs.ClusterRoles {
						clusterRole := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, clusterRole)
					}
					for _, each := range k8sobjs.ClusterRoleBindings {
						clusterRoleBinding := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, clusterRoleBinding)
					}
					for _, each := range k8sobjs.ConfigMaps {
						configMap := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, configMap)
					}
					for _, each := range k8sobjs.Ingresses {
						ingress := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, ingress)
					}
					for _, each := range k8sobjs.Namespaces {
						namespace := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, namespace)
					}
					for _, each := range k8sobjs.NetworkPolicies {
						networkPolicy := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, networkPolicy)
					}
					for _, each := range k8sobjs.Nodes {
						node := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, node)
					}
					for _, each := range k8sobjs.PersistentVolumes {
						persistentVolume := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, persistentVolume)
					}
					for _, each := range k8sobjs.PersistentVolumeClaims {
						persistentVolumeClaim := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, persistentVolumeClaim)
					}
					for _, each := range k8sobjs.Roles {
						role := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, role)
					}
					for _, each := range k8sobjs.RoleBindings {
						roleBinding := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, roleBinding)
					}
					for _, each := range k8sobjs.Secrets {
						secret := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, secret)
					}
					for _, each := range k8sobjs.Services {
						service := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, service)
					}
					for _, each := range k8sobjs.ServiceAccounts {
						serviceAccount := []string{each.Kind, each.Name, each.Namespace, each.UID, "", ""}
						tableData = append(tableData, serviceAccount)
					}
					table.AppendBulk(tableData)
					table.Render()
				}
			}
		}
	}
}

func (a agentService) GetK8sObjs(name, processId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + a.token
	header["Content-Type"] = "application/json"
	return a.httpClient.Get(a.apiServerUrl+"agents/"+name+"/k8sobjs?processId="+processId, header, a.skipSsl)
}

// NewAgentService returns agent type service
func NewAgentService(httpClient service.HttpClient) service.Agent {
	return &agentService{
		httpClient: httpClient,
	}
}
