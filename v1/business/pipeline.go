package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/common"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

type pipelineService struct {
	httpClient   service.HttpClient
	processId    string
	flag         string
	token        string
	apiServerUrl string
	skipSsl      bool
	cmd          *cobra.Command
}

func (p pipelineService) ProcessId(processId string) service.Pipeline {
	p.processId = processId
	return p
}

func (p pipelineService) Flag(flag string) service.Pipeline {
	p.flag = flag
	return p
}

func (p pipelineService) SkipSsl(skipSsl bool) service.Pipeline {
	p.skipSsl = skipSsl
	return p
}

func (p pipelineService) Token(token string) service.Pipeline {
	p.token = token
	return p
}

func (p pipelineService) ApiServerUrl(apiServerUrl string) service.Pipeline {
	p.apiServerUrl = apiServerUrl
	return p
}

func (p pipelineService) Cmd(cmd *cobra.Command) service.Pipeline {
	p.cmd = cmd
	return p
}

func (p pipelineService) Apply() {
	switch p.flag {
	case string(enums.GET_AGENTS):
		httpCode, data, err := p.Get(p.processId, "get_pipeline", p.apiServerUrl, p.token)
		if err != nil {
			p.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else if httpCode != 200 {
			p.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", httpCode)
		} else if data != nil {
			var pipeline v1.Pipeline
			byteBody, _ := json.Marshal(data)
			err := json.Unmarshal(byteBody, &pipeline)
			if err != nil {
				return
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name"})
			agentsMap := make(map[string]bool)
			for _, step := range pipeline.Steps {
				agent := step.Params["agent"]
				if step.Type == "DEPLOY" {
					if _, ok := agentsMap[agent]; !ok {
						agentsMap[agent] = true
						table.Append([]string{agent})
					}
				}
			}
			table.Render()
		}
	}
}

func (p pipelineService) Get(processId string, action string, url string, token string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	url = url + "pipelines/" + processId + "?action=" + action
	code, b, err := p.httpClient.Get(url, header, p.skipSsl)
	if err != nil {
		return code, nil, err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		log.Println("[ERROR]: ", err.Error())
	}

	return code, response.Data, nil
}

func (p pipelineService) Logs(url, page, limit string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	url = url + "?order=&page=" + page + "&limit=" + limit
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

// NewPipelineService returns Pipeline type service
func NewPipelineService(httpClient service.HttpClient) service.Pipeline {
	return &pipelineService{
		httpClient: httpClient,
	}
}
