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
	"log"
	"os"
	"strconv"
)

type companyService struct {
	httpClient service.HttpClient
	flag string
    company interface{}
	companyId string
	repoId string
	option string
	cmd *cobra.Command
	kind string
}

func (c companyService) Kind(kind string) service.Company {
	c.kind=kind
	return c
}

func (c companyService) Cmd(cmd *cobra.Command) service.Company {
	c.cmd=cmd
	return c
}
func (c companyService) Flag(flag string) service.Company {
	c.flag=flag
	return c
}

func (c companyService) Company(company interface{}) service.Company {
	c.company=company
	return c
}

func (c companyService) CompanyId(companyId string) service.Company {
	c.companyId=companyId
	return c
}

func (c companyService) RepoId(repoId string) service.Company {
	c.repoId=repoId
	return c
}

func (c companyService) Option(option string) service.Company {
	c.option=option
	return c
}

func (c companyService) Apply() {
	switch c.flag {
	case string(enums.CREATE_COMPANY):
		err := c.CreateCompany(c.company)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.UPDATE_REPOSITORIES):
		err :=  c.UpdateRepositoriesByCompanyId(c.company, c.companyId, c.option)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.UPDATE_APPLICATIONS):
		err :=  c.UpdateApplicationsByRepositoryId(c.company, c.companyId, c.repoId, c.option)
		if err != nil {
			log.Fatalf("[ERROR]: %v", err)
		}
	case string(enums.GET_COMPANY_BY_ID):
		code, data, err := c.GetCompanyById(c.companyId, c.option)
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
				var company v1.Company
				json.Unmarshal(jsonString, &company)
				companyDto := v1.CompanyDto{
					ApiVersion: "api/v1",
					Kind:       c.kind,
					Company:    company,
				}
				c.cmd.Println(responseDTO.Data)
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
				json.Unmarshal(jsonString, &repositories)
				c.cmd.Println(responseDTO)
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

func (c companyService) CreateCompany(company interface{}) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, _, err = c.httpClient.Post(config.ApiServerUrl+"companies", header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) UpdateRepositoriesByCompanyId(company interface{}, companyId string, option string) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Put(config.ApiServerUrl+"companies/"+companyId+"/repositories?companyUpdateOption="+option, header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) UpdateApplicationsByRepositoryId(company interface{}, companyId string, repoId string, option string) error {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	b, err := json.Marshal(company)
	if err != nil {
		return err
	}
	_, _, err = c.httpClient.Post(config.ApiServerUrl+"applications?companyId="+companyId+"&repositoryId="+repoId+"&companyUpdateOption="+option, header, b)
	if err != nil {
		return err
	}
	return nil
}

func (c companyService) GetCompanyById(companyId string, option string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies/"+companyId+"?"+option, header)
}

func (c companyService) GetCompanies() (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies", header)
}

func (c companyService) GetRepositoriesByCompanyId(companyId string) (httpCode int, data []byte, err error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + os.Getenv("CTL_TOKEN")
	header["Content-Type"] = "application/json"
	return c.httpClient.Get(config.ApiServerUrl+"companies/"+companyId+"/repositories?"+c.option, header)
}

// NewCompanyService returns company type service
func NewCompanyService(httpClient service.HttpClient) service.Company {
	return &companyService{
		httpClient: httpClient,
	}
}