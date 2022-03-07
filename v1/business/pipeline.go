package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/common"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type pipelineService struct {
	httpClient service.HttpClient
	token string
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
		return code, nil, err
	}

	return code, response.Data, nil
}

// NewPipelineService returns Pipeline type service
func NewPipelineService(httpClient service.HttpClient) service.Pipeline {
	return &pipelineService{
		httpClient: httpClient,
	}
}
