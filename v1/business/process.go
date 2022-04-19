package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

type processService struct {
	httpClient   service.HttpClient
	pipeline     service.Pipeline
	cmd          *cobra.Command
	repoId       string
	appId        string
	kind         string
	apiServerUrl string
	token        string
}

func (p processService) ApiServerUrl(apiServerUrl string) service.Process {
	p.apiServerUrl = apiServerUrl
	return p
}

func (p processService) Token(token string) service.Process {
	p.token = token
	return p
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
	httpCode, data, err := p.GetByCompanyIdAndRepositoryIdAndAppName()
	if err != nil {
		p.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
	} else if httpCode != 200 {
		p.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", httpCode)
	} else if data != nil {
		var responseDTO v1.ResponseDTO
		err := json.Unmarshal(data, &responseDTO)
		if err != nil {
			p.cmd.Println("[ERROR]: ", err.Error())
		} else {
			jsonString, _ := json.Marshal(responseDTO.Data)
			var processes v1.Processes
			err := json.Unmarshal(jsonString, &processes)
			if err != nil {
				return
			}
			var processList []v1.ProcessesWithStatus
			for _, each := range processes {
				_, data, err := p.pipeline.Get(each.ProcessId, "get_pipeline", p.apiServerUrl, p.token)
				if err != nil {
					return
				}
				var pipeline v1.Pipeline
				byteBody, _ := json.Marshal(data)
				json.Unmarshal(byteBody, &pipeline)
				totalStep := len(pipeline.Steps)
				statusMap := make(map[string][]v1.Step)
				for _, eachStep := range pipeline.Steps {
					statusMap[eachStep.Status] = append(statusMap[eachStep.Status], eachStep)
				}
				if _, ok := statusMap["active"]; ok {
					processList = append(processList, v1.ProcessesWithStatus{
						ProcessId:    each.ProcessId,
						AppId:        each.AppId,
						RepositoryId: each.RepositoryId,
						Data:         each.Data,
						CreatedAt:    each.CreatedAt,
						Status:       "active",
					})
				} else if steps, ok := statusMap["completed"]; ok {
					count := len(steps)
					if count < totalStep {
						if _, ok := statusMap["failed"]; ok {
							processList = append(processList, v1.ProcessesWithStatus{
								ProcessId:    each.ProcessId,
								AppId:        each.AppId,
								RepositoryId: each.RepositoryId,
								Data:         each.Data,
								CreatedAt:    each.CreatedAt,
								Status:       "failed",
							})
						} else if steps, ok := statusMap["paused"]; ok {
							for _, eachStp := range steps {
								if eachStp.Trigger == "AUTO" {
									processList = append(processList, v1.ProcessesWithStatus{
										ProcessId:    each.ProcessId,
										AppId:        each.AppId,
										RepositoryId: each.RepositoryId,
										Data:         each.Data,
										CreatedAt:    each.CreatedAt,
										Status:       "paused",
									})
								} else {
									processList = append(processList, v1.ProcessesWithStatus{
										ProcessId:    each.ProcessId,
										AppId:        each.AppId,
										RepositoryId: each.RepositoryId,
										Data:         each.Data,
										CreatedAt:    each.CreatedAt,
										Status:       "completed",
									})
								}
							}
						} else if steps, ok := statusMap["non_initialized"]; ok {
							for _, eachStp := range steps {
								if eachStp.Trigger == "AUTO" {
									processList = append(processList, v1.ProcessesWithStatus{
										ProcessId:    each.ProcessId,
										AppId:        each.AppId,
										RepositoryId: each.RepositoryId,
										Data:         each.Data,
										CreatedAt:    each.CreatedAt,
										Status:       "non_initialized",
									})
								} else {
									processList = append(processList, v1.ProcessesWithStatus{
										ProcessId:    each.ProcessId,
										AppId:        each.AppId,
										RepositoryId: each.RepositoryId,
										Data:         each.Data,
										CreatedAt:    each.CreatedAt,
										Status:       "completed",
									})
								}
							}
						} else {
							processList = append(processList, v1.ProcessesWithStatus{
								ProcessId:    each.ProcessId,
								AppId:        each.AppId,
								RepositoryId: each.RepositoryId,
								Data:         each.Data,
								CreatedAt:    each.CreatedAt,
								Status:       "completed",
							})
						}
					}
				}
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Api Version", "Kind", "Process Id", "Application Id", "Repository Id", "Created At", "Status"})
			if len(processList) < 5 {
				for _, eachProcess := range processList {
					createdAt := strconv.Itoa(eachProcess.CreatedAt.Local().Day()) + "-" + strconv.Itoa(int(eachProcess.CreatedAt.Local().Month())) + "-" + strconv.Itoa(eachProcess.CreatedAt.Local().Year()) + " " + eachProcess.CreatedAt.Local().Format(time.Kitchen)
					process := []string{"api/v1", p.kind, eachProcess.ProcessId, eachProcess.AppId, eachProcess.RepositoryId, createdAt, eachProcess.Status}
					table.Append(process)
				}
			} else {
				processList = processList[0:5]
				for _, eachProcess := range processList {
					createdAt := strconv.Itoa(eachProcess.CreatedAt.Local().Day()) + "-" + strconv.Itoa(int(eachProcess.CreatedAt.Local().Month())) + "-" + strconv.Itoa(eachProcess.CreatedAt.Local().Year()) + " " + eachProcess.CreatedAt.Local().Format(time.Kitchen)
					process := []string{"api/v1", p.kind, eachProcess.ProcessId, eachProcess.AppId, eachProcess.RepositoryId, createdAt, eachProcess.Status}
					table.Append(process)
				}
			}
			table.Render()
		}
	}
}

func (p processService) GetByCompanyIdAndRepositoryIdAndAppName() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	return p.httpClient.Get(p.apiServerUrl+"processes?repositoryId="+p.repoId+"&appId="+p.appId, header)
}

// NewProcessService returns process type service
func NewProcessService(httpClient service.HttpClient, pipeline service.Pipeline) service.Process {
	return &processService{
		httpClient: httpClient,
		pipeline:   pipeline,
	}
}
