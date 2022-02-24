package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/config"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type processService struct {
	httpClient service.HttpClient
	cmd *cobra.Command
	repoId string
	appId string
	kind string
}

func (p processService) Kind(kind string) service.Process {
	p.kind = kind
	return p
}

func (p processService) RepoId(repoId string) service.Process {
	p.repoId = repoId
	return p
}

func (p processService) ApplicationId(appId string) service.Process {
	p.appId = appId
	return p
}

func (p processService) Cmd(cmd *cobra.Command) service.Process {
	p.cmd = cmd
	return p
}

func (p processService) Apply() {
	code, data, err := p.GetByCompanyIdAndRepositoryIdAndAppName()
	if err != nil {
		log.Fatalf("[ERROR]: %v", err)
	} else if code != 200 {
		p.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
	} else if data != nil {
		var responseDTO v1.ResponseDTO
		err := json.Unmarshal(data, &responseDTO)
		if err != nil {
			p.cmd.Println("[ERROR]: ", err.Error())
		} else {
			jsonString, _ := json.Marshal(responseDTO.Data)
			var processes v1.Processes
			json.Unmarshal(jsonString, &processes)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Api Version", "Kind", "Process Id", "Application Id", "Repository Id", "Data"})
			if len(processes) < 5 {
				for _, eachProcess := range processes {
					var data string
					for key, val := range eachProcess.Data {
						strVal, _ := json.Marshal(val)
						data += key + ": " + string(strVal) + "\n"
					}
					process := []string{"api/v1", p.kind, eachProcess.ProcessId, eachProcess.AppId, eachProcess.RepositoryId, data}
					table.Append(process)
				}
			} else {
				processes = processes[len(processes) - 5 :]
				for _, eachProcess := range processes {
					var data string
					for key, val := range eachProcess.Data {
						strVal, _ := json.Marshal(val)
						data += key + ": " + string(strVal) + "\n"
					}
					process := []string{"api/v1", p.kind, eachProcess.ProcessId, eachProcess.AppId, eachProcess.RepositoryId, data}
					table.Append(process)
				}
			}
			table.Render()
		}
	}
}

func (p processService) GetByCompanyIdAndRepositoryIdAndAppName() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return p.httpClient.Get(config.ApiServerUrl+"processes?repositoryId="+p.repoId+"&appId="+p.appId+"&companyId=12345", header)
}

// NewProcessService returns process type service
func NewProcessService(httpClient service.HttpClient) service.Process {
	return &processService{
		httpClient: httpClient,
	}
}