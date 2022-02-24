package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/config"
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
		code, data, err := r.GetApplicationsByRepositoryId(r.repo)
		if err != nil {
			r.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			r.cmd.Println("[ERROR]: ", "Something went wrong! Status Code: ", code)
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
					table.SetHeader([]string{"Api Version", "Kind", "Id", "Name", "Labels", "IsWebhookEnabled", "Url"})
					for _, eachApp := range applications {
						var labels string
						for key, val := range eachApp.MetaData.Labels {
							labels += key + ": " + val + "\n"
						}
						application := []string{"api/v1", r.kind, eachApp.MetaData.Id, eachApp.MetaData.Name, labels, eachApp.MetaData.IsWebhookEnabled, eachApp.Url}
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
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(config.ApiServerUrl+"repositories/"+repositoryId+"?"+r.option, header)
}

func (r repositoryService) GetApplicationsByRepositoryId(repositoryId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return r.httpClient.Get(config.ApiServerUrl+"repositories/"+repositoryId+"/applications", header)
}

// NewRepositoryService returns repository type service
func NewRepositoryService(httpClient service.HttpClient) service.Repository {
	return &repositoryService{
		httpClient: httpClient,
	}
}
