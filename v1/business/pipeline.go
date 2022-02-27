package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/common"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type pipelineService struct {
	httpClient service.HttpClient
}

func (p pipelineService) Logs(url, page, limit string) (httpCode int, data interface{}, err error) {
	var response common.ResponseDTO
	header := make(map[string]string)
	token, _ := v1.GetToken()
	header["Authorization"] = "Bearer " + token
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
