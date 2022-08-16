package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/common"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strconv"
	"time"
)

type processService struct {
	httpClient   service.HttpClient
	pipeline     service.Pipeline
	cmd          *cobra.Command
	id           string
	step         string
	claim        string
	footmark     string
	repoId       string
	appId        string
	flag         string
	kind         string
	apiServerUrl string
	token        string
	skipSsl      bool
}

func (p processService) Id(id string) service.Process {
	p.id = id
	return p
}

func (p processService) Step(step string) service.Process {
	p.step = step
	return p
}

func (p processService) Claim(claim string) service.Process {
	p.claim = claim
	return p
}

func (p processService) Footmark(footmark string) service.Process {
	p.footmark = footmark
	return p
}

func (p processService) Flag(flag string) service.Process {
	p.flag = flag
	return p
}

func (p processService) SkipSsl(skipSsl bool) service.Process {
	p.skipSsl = skipSsl
	return p
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
	switch p.flag {
	case string(enums.GET_PROCESS):
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
}

func (p processService) GetByCompanyIdAndRepositoryIdAndAppName() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	return p.httpClient.Get(p.apiServerUrl+"processes?repositoryId="+p.repoId+"&appId="+p.appId, header, p.skipSsl)
}

func (p processService) GetLogsByProcessIdAndStepAndClaimAndFootmark(page, limit string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	url := p.apiServerUrl + "processes/" + p.id + "/steps/" + p.step + "/footmarks/" + p.footmark + "/logs?claims=" + p.claim + "&page=" + page + "&limit=" + limit
	code, b, err := p.httpClient.Get(url, header, p.skipSsl)
	if err != nil {
		return code, nil, err
	}
	er := json.Unmarshal(b, &response)
	if er != nil {
		return http.StatusBadRequest, nil, err
	}
	return code, response.Data, nil
}

func (p processService) GetProcessLifeCycleEventByProcessIdAndStep() (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	url := p.apiServerUrl + "processes/" + p.id + "/process_life_cycle_events?step=" + p.step
	code, b, err := p.httpClient.Get(url, header, p.skipSsl)
	if err != nil {
		return code, nil, err
	}
	er := json.Unmarshal(b, &response)
	if er != nil {
		return http.StatusBadRequest, nil, err
	}
	return code, response.Data, nil
}

// NewProcessService returns process type service
func NewProcessService(httpClient service.HttpClient, pipeline service.Pipeline) service.Process {
	return &processService{
		httpClient: httpClient,
		pipeline:   pipeline,
	}
}
