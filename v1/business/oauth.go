package business

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
)

type oauthService struct {
	httpClient service.HttpClient
}

type JWTPayLoad struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
type ResponseDTOWithPagination struct {
	Data     JWTPayLoad `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

func (o oauthService) Apply(loginDto interface{}) (string, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(loginDto)
	if err != nil {
		return "", err
	}
	_, data, err := o.httpClient.Post(v1.GetSecurityUrl()+"oauth/login?grant_type=password&token_type=ctl", header, b)
	if err != nil {
		return "", err
	}
	var payload ResponseDTOWithPagination
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return "", err
	}
	return payload.Data.AccessToken, nil
}

// NewOauthService returns oauth type service
func NewOauthService(httpClient service.HttpClient) service.Oauth {
	return &oauthService{
		httpClient: httpClient,
	}
}