package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strings"
	"gopkg.in/yaml.v2"
)

func Describe() *cobra.Command{
	return &cobra.Command{
		Use:       "describe",
		Short:     "Describe resource[company/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {

		if len(args)<2{
			log.Fatalf("[ERROR]: %v", "please provide a resource name!")
			return nil
		}

		if args[0]=="company"{
			companyService := dependency_manager.GetCompanyService()
		// get companyId from token
			var companyId string
			companyService.Flag(string(enums.GET_COMPANY_BY_ID)).CompanyId(companyId).Apply()
		}else if args[0]=="repository" || args[0]=="repo"{
			repositoryService := dependency_manager.GetRepositoryService()
			repositoryService.Flag(string(enums.GET_REPOSITORY)).Repo(args[1]).Apply()
		}else if args[0]=="application" || args[0]=="app" {

		}

			return nil
		},
	}
}

func List() *cobra.Command{
	return &cobra.Command{
		Use:       "list",
		Short:     "Describe resource[company/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args)<2{
				log.Fatalf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			if args[0]=="repositories" || args[0]=="repos"{
				var companyId string
				companyService := dependency_manager.GetCompanyService()
				companyService.Flag(string(enums.GET_REPOSITORIES)).CompanyId(companyId).Apply()
			}else if args[0]=="applications" || args[0]=="apps"{
				if len(args)<3{
					log.Fatalf("[ERROR]: %v", "please provide a repository id!")
					return nil
				}
				strs:=strings.Split(args[1],"=")
				if len(strs)<2{
					log.Fatalf("[ERROR]: %v", "please provide sub resource like repo=123!")
					return nil
				}
				if strs[0]=="repo" || strs[0]=="repository"{
					repositoryService := dependency_manager.GetRepositoryService()
					repositoryService.Flag(string(enums.GET_APPLICATIONS)).Repo(strs[1]).Apply()
				}else{
					log.Fatalf("[ERROR]: %v", "please provide a valid sub resource[repo/repository]!")
					return nil
				}

			}

			return nil
		},
	}
}


func Update() *cobra.Command{
	return &cobra.Command{
		Use:       "update",
		Short:     "Update resource [repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var option string
			var companyId string
			var repoId string

			if args[0]=="repositories" || args[0]=="repos"{
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					}
					if strings.Contains(strings.ToLower(each), "option") {
						strs := strings.Split(strings.ToUpper(each), "=")
						if len(strs) > 0 {
							option = strs[1]
						}
					}
				}

				data, err := ioutil.ReadFile(file)
				if err != nil {
					log.Printf("data.Get err   #%v ", err)
					return nil
				}
				repos := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, repos)
					if err != nil {
						log.Fatalf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, repos)
					if err != nil {
						log.Fatalf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option(option).Apply()
				return nil
			}else if args[0]=="applications" || args[0]=="apps"{
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") {
						strs := strings.Split(strings.ToUpper(each), "=")
						if len(strs) > 0 {
							option = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "repo") {
						strs := strings.Split(each, "=")
						if len(strs) > 0 {
							repoId = strs[1]
						}
					}
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					log.Printf("data.Get err   #%v ", err)
					return nil
				}
				company := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, company)
					if err != nil {
						log.Fatalf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, company)
					if err != nil {
						log.Fatalf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).RepoId(repoId).Option(option).Apply()
				return nil
			}

			return nil
		},
	}
}
