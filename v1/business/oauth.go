package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type oauthService struct {
	httpClient service.HttpClient
	securityUrl   string
}

func (o oauthService) SecurityUrl(securityUrl string) service.Oauth {
	o.securityUrl = securityUrl
	return o
}

type JWTPayLoad struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
type ResponseDTO struct {
	Data     JWTPayLoad `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

func (o oauthService) Apply(loginDto interface{}) (string, error, int) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(loginDto)
	if err != nil {
		return "", err, 0
	}
	code, data, err := o.httpClient.Post(o.securityUrl+"oauth/login?grant_type=password&token_type=ctl", header, b)
	if err != nil {
		return "", err, code
	}
	var payload ResponseDTO
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return "", err, code
	}
	return payload.Data.AccessToken, nil, code
}

// NewOauthService returns oauth type service
func NewOauthService(httpClient service.HttpClient) service.Oauth {
	return &oauthService{
		httpClient: httpClient,
	}
}