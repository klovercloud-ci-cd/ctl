package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

type repositoryService struct {
	httpClient service.HttpClient
	flag string
	companyId string
	repo string
	cmd *cobra.Command
	option string
	kind string
	apiServerUrl string
	token string
}

func (r repositoryService) ApiServerUrl(apiServerUrl string) service.Repository {
	r.apiServerUrl = apiServerUrl
	return r
}

func (r repositoryService) Token(token string) service.Repository {
	r.token = token
	return r
}

func (r repositoryService) Kind(kind string) service.Repository {
	r.kind=kind
	return r
}

func (r repositoryService) Option(option string) service.Repository {
	r.option=option
	return r
}

func (r repositoryService) Cmd(cmd *cobra.Command) service.Repository {
	r.cmd=cmd
	return r
}

func (r repositoryService) Flag(flag string) service.Repository {
	r.flag=flag
	return r
}

func (r repositoryService) CompanyId(companyId string) service.Repository {
	r.companyId=companyId
	return r
}

func (r repositoryService) Repo(repoId string) service.Repository {
	r.repo=repoId
	return r
}

func (r repositoryService) Apply() {
	switch r.flag {
	case string(enums.GET_REPOSITORY):
		code, data, err := r.GetRepositoryById(r.repo)
		if err != nil {
			r.cmd.Println("[ERROR]: ", err.Error())
			r.cmd.Println("Status Code: ", code)
		} else if code != 200 {
			r.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				r.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var repository v1.Repository
				json.Unmarshal(jsonString, &repository)
				repositoryDto := v1.RepositoryDto{
					ApiVersion: "api/v1",
					Kind:       r.kind,
					Repository: repository,
				}
				b, _ := yaml.Marshal(repositoryDto)
				b = v1.AddRootIndent(b, 4)
				r.cmd.Println(string(b))
			}
		}
	case string(enums.GET_APPLICATIONS):
		httpCode, data, err := r.GetApplicationsByRepositoryId(r.repo)
		if err != nil {
			r.cmd.Println("[ERROR]: " + err.Error() + "Status Code: ", httpCode)
		} else if httpCode != 200 {
			r.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", httpCode)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				r.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var applications v1.Applications
				err := json.Unmarshal(jsonString, &applications)
				if err != nil {
					r.cmd.Println("[ERROR]: ", err.Error())
				} else {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Api Version", "Kind", "Id", "Name", "Labels", "Webhook", "Url"})
					for _, eachApp := range applications {
						var labels string
						for key, val := range eachApp.MetaData.Labels {
							labels += key + ": " + val + "\n"
						}
						webhook := "Disabled"
						if eachApp.MetaData.IsWebhookEnabled {
							webhook = "Enabled"
						}
						application := []string{"api/v1", r.kind, eachApp.MetaData.Id, eachApp.MetaData.Name, labels, webhook, eachApp.Url}
						table.Append(application)
					}
					table.Render()
				}
			}
		}
	}

}

func (r repositoryService) GetRepositoryById(repositoryId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + r.token
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(r.apiServerUrl+"repositories/"+repositoryId+"?"+r.option, header)
}

func (r repositoryService) GetApplicationsByRepositoryId(repositoryId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + r.token
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(r.apiServerUrl+"repositories/"+repositoryId+"/applications?status=ACTIVE", header)
}

// NewRepositoryService returns repository type service
func NewRepositoryService(httpClient service.HttpClient) service.Repository {
	return &repositoryService{
		httpClient: httpClient,
	}
}
