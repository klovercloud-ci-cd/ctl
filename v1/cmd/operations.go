package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func Create() *cobra.Command{
	return &cobra.Command{
		Use:       "create",
		Short:     "Create resource [Company]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				log.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				log.Fatalf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			if args[0]=="company"{
				var file string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					}
				}
				if file == "" {
					log.Fatalf("[ERROR]: %v", "please provide a file!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					log.Printf("data.Get err   #%v ", err)
					return nil
				}
				company := new(v1.Company)
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
				companyService.Cmd(cmd).Flag(string(enums.CREATE_COMPANY)).Company(company).Apply()
			}
			return nil
		},
	}
}

func Describe() *cobra.Command{
	return &cobra.Command{
		Use:       "describe",
		Short:     "Describe resource [company/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				log.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				log.Fatalf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken()
			if err != nil {
				log.Fatalf("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			if args[0]=="company"{
				loadRepo := false
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadrepositories") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadrepos") || strings.Contains(strings.ToLower(args[idx+1]), "lr") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadRepo = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps") || strings.Contains(strings.ToLower(args[idx+1]), "la"){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							}
						}
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Kind("Company").Cmd(cmd).Flag(string(enums.GET_COMPANY_BY_ID)).CompanyId(companyId).Option("loadRepositories="+strconv.FormatBool(loadRepo)+"&loadApplications="+strconv.FormatBool(loadApp)).Apply()
			}else if args[0]=="repository" || args[0]=="repo"{
				if len(args)<2 {
					log.Fatalf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				var repoId string
				if strings.Contains(strings.ToLower(args[1]), "repoid") {
					strs := strings.Split(strings.ToLower(args[1]), "=")
					if len(strs) > 1 {
						repoId = strs[1]
					}
				}
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps") || strings.Contains(strings.ToLower(args[idx+1]), "la"){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							}
						}
					}
				}
				repositoryService := dependency_manager.GetRepositoryService()
				repositoryService.Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORY)).Repo(repoId).Option("loadApplications="+strconv.FormatBool(loadApp)).Apply()
			}else if args[0]=="application" || args[0]=="app" {
				var repoId string
				var appId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid") || strings.Contains(strings.ToLower(each), "repoid") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "applicationid") || strings.Contains(strings.ToLower(each), "appid") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					}
				}
				if repoId == "" {
					log.Fatalf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if appId == "" {
					log.Fatalf("[ERROR]: %v", "please provide application id!")
					return nil
				}
				applicationService := dependency_manager.GetApplicationService()
				applicationService.Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATION)).CompanyId(companyId).RepoId(repoId).ApplicationId(appId).Apply()
			}
			return nil
		},
	}
}

func List() *cobra.Command{
	return &cobra.Command{
		Use:       "list",
		Short:     "Describe resource [company/repository/application/process]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				log.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				log.Fatalf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken()
			if err != nil {
				log.Fatalf("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			if args[0]=="repositories" || args[0]=="repos"{
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps") || strings.Contains(strings.ToLower(args[idx+1]), "la"){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							}
						}
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORIES)).CompanyId(companyId).Option("loadApplications="+strconv.FormatBool(loadApp)).Apply()
			} else if args[0]=="applications" || args[0]=="apps"{
				var repoId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid") || strings.Contains(strings.ToLower(each), "repoid") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					}
				}
				if repoId == "" {
					log.Fatalf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				repositoryService := dependency_manager.GetRepositoryService()
				repositoryService.Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATIONS)).Repo(repoId).Apply()
			} else if args[0]=="process" {
				var repoId string
				var appId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid") || strings.Contains(strings.ToLower(each), "repoid") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "applicationid") || strings.Contains(strings.ToLower(each), "appid") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					}
				}
				if repoId == "" {
					log.Fatalf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if appId == "" {
					log.Fatalf("[ERROR]: %v", "please provide application id!")
					return nil
				}
				processService := dependency_manager.GetProcessService()
				processService.Kind("Process").Cmd(cmd).RepoId(repoId).ApplicationId(appId).Apply()
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
			if err := v1.IsUserLoggedIn(); err != nil {
				log.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				log.Fatalf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken()
			if err != nil {
				log.Fatalf("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId

			var file string
			var option string
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
						if len(strs) > 1 {
							option = strs[1]
						}
					}
				}
				if file == "" {
					log.Fatalf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if option == "" {
					log.Fatalf("[ERROR]: %v", "please provide update option!")
					return nil
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
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option(option).Apply()
				return nil
			}else if args[0]=="applications" || args[0]=="apps"{
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") {
						strs := strings.Split(strings.ToUpper(each), "=")
						if len(strs) > 1 {
							option = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "repoid") {
						strs := strings.Split(each, "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					}
				}
				if file == "" {
					log.Fatalf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if repoId == "" {
					log.Fatalf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if option == "" {
					log.Fatalf("[ERROR]: %v", "please provide update option!")
					return nil
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
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).RepoId(repoId).Option(option).Apply()
				return nil
			}
			return nil
		},
	}
}
