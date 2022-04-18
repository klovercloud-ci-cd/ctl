package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/common"
	"github.com/klovercloud-ci/ctl/v1/service"
	"log"
	"net/http"
)

type pipelineService struct {
	httpClient service.HttpClient
	token      string
}

func (p pipelineService) Get(processId string, action string, url string, token string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token
	header["Content-Type"] = "application/json"
	url = url + "pipelines/" + processId + "?action=" + action
	code, b, err := p.httpClient.Get(url, header)
	if err != nil {
		return code, nil, err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		log.Println("[ERROR]: ", err.Error())
	}

	return code, response.Data, nil
}

func (p pipelineService) Token(token string) service.Pipeline {
	p.token = token
	return p
}

func (p pipelineService) Logs(url, page, limit string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + p.token
	header["Content-Type"] = "application/json"
	url = url + "?order=&page=" + page + "&limit=" + limit
	code, b, err := p.httpClient.Get(url, header)

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
