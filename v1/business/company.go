package business

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strconv"
)

type companyService struct {
	httpClient   service.HttpClient
	flag         string
	company      interface{}
	companyId    string
	repoId       string
	option       string
	cmd          *cobra.Command
	kind         string
	apiServerUrl string
	token        string
	skipSsl      bool
}

func (c companyService) SkipSsl(skipSsl bool) service.Company {
	c.skipSsl = skipSsl
	return c
}

func (c companyService) ApiServerUrl(apiServerUrl string) service.Company {
	c.apiServerUrl = apiServerUrl
	return c
}

func (c companyService) Token(token string) service.Company {
	c.token = token
	return c
}

func (c companyService) Kind(kind string) service.Company {
	c.kind = kind
	return c
}

func (c companyService) Cmd(cmd *cobra.Command) service.Company {
	c.cmd = cmd
	return c
}
func (c companyService) Flag(flag string) service.Company {
	c.flag = flag
	return c
}

func (c companyService) Company(company interface{}) service.Company {
	c.company = company
	return c
}

func (c companyService) CompanyId(companyId string) service.Company {
	c.companyId = companyId
	return c
}

func (c companyService) RepoId(repoId string) service.Company {
	c.repoId = repoId
	return c
}

func (c companyService) Option(option string) service.Company {
	c.option = option
	return c
}

func (c companyService) Apply() {
	switch c.flag {
	case string(enums.CREATE_COMPANY):
		httpCode, _, err := c.CreateCompany(c.company)
		if err != nil {
			c.cmd.Println("[ERROR]: %v", err)
			c.cmd.Println("Status Code: %v", httpCode)
		} else {
			c.cmd.Println("Successfully Created Company")
		}
	case string(enums.UPDATE_REPOSITORIES):
		if c.option != string(enums.APPEND_REPOSITORY) && c.option != string(enums.SOFT_DELETE_REPOSITORY) && c.option != string(enums.DELETE_REPOSITORY) {
			c.cmd.Println("[ERROR]: Invalid Repository Update Option")
		} else {
			httpCode, _, err := c.UpdateRepositoriesByCompanyId(c.company, c.companyId, c.option)
			if err != nil {
				c.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
			} else {
				c.cmd.Println("Successfully Updated Repositories")
			}
		}
	case string(enums.UPDATE_APPLICATIONS):
		if c.option != string(enums.APPEND_APPLICATION) && c.option != string(enums.SOFT_DELETE_APPLICATION) && c.option != string(enums.DELETE_APPLICATION) {
			c.cmd.Println("[ERROR]: Invalid Application Update Option")
		} else {
			httpCode, _, err := c.UpdateApplicationsByRepositoryId(c.company, c.companyId, c.repoId, c.option)
			if err != nil {
				c.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
			} else {
				c.cmd.Println("Successfully Updated Applications")
			}
		}
	case string(enums.GET_COMPANY_BY_ID):
		httpCode, data, err := c.GetCompanyById()
		if err != nil {
			c.cmd.Println("[ERROR]: "+err.Error()+"Status Code: ", httpCode)
		} else if httpCode != 200 {
			c.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", httpCode)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				c.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var company v1.Company
				err := json.Unmarshal(jsonString, &company)
				if err != nil {
					return
				}
				companyDto := v1.CompanyDto{
					ApiVersion: "api/v1",
					Kind:       c.kind,
					Company:    company,
				}
				b, _ := yaml.Marshal(companyDto)
				b = v1.AddRootIndent(b, 4)
				c.cmd.Println(string(b))
			}
		}
	case string(enums.GET_COMPANIES):
		code, data, err := c.GetCompanies()
		if err != nil {
			c.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			c.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		} else if data != nil {
			c.cmd.Println(string(data))
		}
	case string(enums.GET_REPOSITORIES):
		code, data, err := c.GetRepositoriesByCompanyId(c.companyId)
		if err != nil {
			c.cmd.Println("[ERROR]: ", err.Error())
		} else if code != 200 {
			c.cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		} else if data != nil {
			var responseDTO v1.ResponseDTO
			err := json.Unmarshal(data, &responseDTO)
			if err != nil {
				c.cmd.Println("[ERROR]: ", err.Error())
			} else {
				jsonString, _ := json.Marshal(responseDTO.Data)
				var repositories v1.Repositories
				err := json.Unmarshal(jsonString, &repositories)
				if err != nil {
					return
				}
				table := tablewriter.NewWriter(os.Stdout)
				if c.option == "loadApplications=false" {
					table.SetHeader([]string{"Api Version", "Kind", "Id", "Type"})
					for _, eachRepo := range repositories {
						repository := []string{"api/v1", c.kind, eachRepo.Id, eachRepo.Type}
						table.Append(repository)
					}
				} else {
					table.SetHeader([]string{"Api Version", "Kind", "Id", "Type", "Applications Count"})
					for _, eachRepo := range repositories {
						repository := []string{"api/v1", c.kind, eachRepo.Id, eachRepo.Type, strconv.Itoa(len(eachRepo.Applications))}
						table.Append(repository)
					}
				}
				table.Render()
			}
		}
	default:
		c.cmd.Println("[ERROR]: ", "Please provide valid options")
	}

}

func (c companyService) CreateCompany(company interface{}) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	httpCode, data, err = c.httpClient.Post(c.apiServerUrl+"companies", header, b, c.skipSsl)
	if err != nil {
		return httpCode, nil, err
	}
	return httpCode, data, nil
}

func (c companyService) UpdateRepositoriesByCompanyId(company interface{}, companyId string, option string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	httpCode, err = c.httpClient.Put(c.apiServerUrl+"companies/"+companyId+"/repositories?companyUpdateOption="+option, header, b, c.skipSsl)
	if err != nil {
		return httpCode, nil, err
	}
	return httpCode, data, nil
}

func (c companyService) UpdateApplicationsByRepositoryId(company interface{}, companyId, repoId, option string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	httpCode, err = c.httpClient.Put(c.apiServerUrl+"companies/"+companyId+"/repositories/"+repoId+"/applications?companyUpdateOption="+option, header, b, c.skipSsl)
	if err != nil {
		return httpCode, nil, err
	}
	return httpCode, nil, nil
}

func (c companyService) GetCompanyById() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(c.apiServerUrl+"companies/"+c.companyId+"?"+c.option, header, c.skipSsl)
}

func (c companyService) GetCompanies() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(c.apiServerUrl+"companies", header, c.skipSsl)
}

func (c companyService) GetRepositoriesByCompanyId(companyId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.token
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(c.apiServerUrl+"companies/"+companyId+"/repositories?"+c.option, header, c.skipSsl)
}

// NewCompanyService returns company type service
func NewCompanyService(httpClient service.HttpClient) service.Company {
	return &companyService{
		httpClient: httpClient,
	}
}
