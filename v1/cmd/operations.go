package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

func Create() *cobra.Command{
	return &cobra.Command{
		Use:       "create",
		Short:     "Create resource [user/repositories/applications]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				cmd.Printf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			var apiServerUrl string
			var securityUrl string
			var file string
			if args[0] == "user" {
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "security=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									securityUrl = strs[1]
								}
							}
						}
					}
				}
				if file == "" {
					cmd.Printf("[ERROR]: %v", "please provide a file!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Printf("data.Get err   #%v ", err.Error())
					return nil
				}
				var user v1.UserRegistrationDto
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, &user)
					if err != nil {
						cmd.Printf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, &user)
					if err != nil {
						cmd.Printf("json Unmarshal: %v", err)
						return nil
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cfg.ApiServerUrl = "http://localhost:8080/api/v1/"
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				if securityUrl == "" {
					if cfg.SecurityUrl == "" {
						cfg.SecurityUrl = "http://localhost:8085/api/v1/"
					}
				} else {
					cfg.SecurityUrl = securityUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
				}
				userService := dependency_manager.GetUserService()
				userService.Cmd(cmd).Flag(string(enums.CREATE_USER)).User(user).Apply()
				return nil
			} else if args[0]=="repositories" || args[0]=="repos"{
				if err := v1.IsUserLoggedIn(); err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				userMetadata, err := v1.GetUserMetadataFromBearerToken()
				if err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if file == "" {
					cmd.Printf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Printf("data.Get err   #%v ", err)
					return nil
				}
				repos := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, repos)
					if err != nil {
						cmd.Printf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, repos)
					if err != nil {
						cmd.Printf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option("APPEND_REPOSITORY").Apply()
				return nil
			} else if args[0]=="applications" || args[0]=="apps"{
				var repoId string
				if err := v1.IsUserLoggedIn(); err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "repoid=") {
						strs := strings.Split(each, "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if file == "" {
					cmd.Printf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if repoId == "" {
					cmd.Printf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Printf("data.Get err   #%v ", err)
					return nil
				}
				company := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, company)
					if err != nil {
						cmd.Printf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, company)
					if err != nil {
						cmd.Printf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).RepoId(repoId).Option("APPEND_APPLICATION").Apply()
				return nil
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
		},
	}
}

func Registration() *cobra.Command{
	return &cobra.Command{
		Use:       "register",
		Short:     "Register user",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var actionType string
			var apiServerUrl string
			var securityUrl string
			for idx, each := range args {
				if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						file = strs[1]
					}
				} else if strings.ToLower(each) == "user" {
					actionType = "user"
				} else if strings.Contains(strings.ToLower(each), "option") {
					if idx + 1 < len(args) {
						if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								apiServerUrl = strs[1]
							}
						} else if strings.Contains(strings.ToLower(args[idx+1]), "security=") {
							strs := strings.Split(strings.ToLower(args[idx+1]), "=")
							if len(strs) > 1 {
								securityUrl = strs[1]
							}
						}
					}
				}
			}
			if actionType == "user" {
				if err := v1.IsUserLoggedIn(); err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
			}
			if file == "" {
				cmd.Printf("[ERROR]: %v", "please provide a file!")
				return nil
			}
			data, err := ioutil.ReadFile(file)
			if err != nil {
				cmd.Printf("data.Get err   #%v ", err.Error())
				return nil
			}
			var user v1.UserRegistrationDto
			if strings.HasSuffix(file, ".yaml") {
				err = yaml.Unmarshal(data, &user)
				if err != nil {
					cmd.Printf("yaml Unmarshal: %v", err)
					return nil
				}
			} else {
				err = json.Unmarshal(data, &user)
				if err != nil {
					cmd.Printf("json Unmarshal: %v", err)
					return nil
				}
			}
			cfg := v1.GetConfigFile()
			if apiServerUrl == "" {
				if cfg.ApiServerUrl == "" {
					cfg.ApiServerUrl = "http://localhost:8080/api/v1/"
				}
			} else {
				cfg.ApiServerUrl = apiServerUrl
			}
			if securityUrl == "" {
				if cfg.SecurityUrl == "" {
					cfg.SecurityUrl = "http://localhost:8085/api/v1/"
				}
			} else {
				cfg.SecurityUrl = securityUrl
			}
			err = cfg.Store()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
			}
			userService := dependency_manager.GetUserService()
			if strings.ToLower(actionType) == "user" {
				userService.Cmd(cmd).Flag(string(enums.CREATE_USER)).User(user).Apply()
			} else {
				userService.Cmd(cmd).Flag(string(enums.CREATE_ADMIN)).User(user).Apply()
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
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				cmd.Printf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken()
			if err != nil {
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			var apiServerUrl string
			if args[0]=="company"{
				loadRepo := false
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadrepositories=") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadrepos=") || strings.Contains(strings.ToLower(args[idx+1]), "lr=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadRepo = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la="){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Kind("Company").Cmd(cmd).Flag(string(enums.GET_COMPANY_BY_ID)).CompanyId(companyId).Option("loadRepositories="+strconv.FormatBool(loadRepo)+"&loadApplications="+strconv.FormatBool(loadApp)).Apply()
			}else if args[0]=="repository" || args[0]=="repo"{
				if len(args)<2 {
					cmd.Printf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				var repoId string
				if strings.Contains(strings.ToLower(args[1]), "repoid=") {
					strs := strings.Split(strings.ToLower(args[1]), "=")
					if len(strs) > 1 {
						repoId = strs[1]
					}
				}
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la="){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				repositoryService := dependency_manager.GetRepositoryService()
				repositoryService.Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORY)).Repo(repoId).Option("loadApplications="+strconv.FormatBool(loadApp)).Apply()
			}else if args[0]=="application" || args[0]=="app" {
				var repoId string
				var appId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid=") || strings.Contains(strings.ToLower(each), "repoid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "applicationid=") || strings.Contains(strings.ToLower(each), "appid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				if repoId == "" {
					cmd.Printf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if appId == "" {
					cmd.Printf("[ERROR]: %v", "please provide application id!")
					return nil
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				applicationService := dependency_manager.GetApplicationService()
				applicationService.Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATION)).RepoId(repoId).ApplicationId(appId).Apply()
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
	}
}

func List() *cobra.Command{
	return &cobra.Command{
		Use:       "list",
		Short:     "List resources [company/repository/application/process]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			if len(args) < 1{
				cmd.Printf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken()
			if err != nil {
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			var apiServerUrl string
			if args[0]=="repositories" || args[0]=="repos"{
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") {
						if idx + 1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") ||  strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la="){
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadApp = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORIES)).CompanyId(companyId).Option("loadApplications="+strconv.FormatBool(loadApp)).Apply()
			} else if args[0]=="applications" || args[0]=="apps"{
				var repoId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid=") || strings.Contains(strings.ToLower(each), "repoid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if repoId == "" {
					applicationService := dependency_manager.GetApplicationService()
					applicationService.Kind("Application").Cmd(cmd).Flag(string(enums.GET_All_APPLICATIONS)).Apply()
				} else {
					repositoryService := dependency_manager.GetRepositoryService()
					repositoryService.Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATIONS)).Repo(repoId).Apply()
				}
			} else if args[0]=="process" {
				var repoId string
				var appId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "repositoryid=") || strings.Contains(strings.ToLower(each), "repoid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "applicationid=") || strings.Contains(strings.ToLower(each), "appid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				if repoId == "" {
					cmd.Printf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if appId == "" {
					cmd.Printf("[ERROR]: %v", "please provide application id!")
					return nil
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				processService := dependency_manager.GetProcessService()
				processService.Kind("Process").Cmd(cmd).RepoId(repoId).ApplicationId(appId).Apply()
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
	}
}

func Update() *cobra.Command{
	return &cobra.Command{
		Use:       "update",
		Short:     "Update resource [user/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1{
				cmd.Printf("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			var file string
			var option string
			var repoId string
			var email string
			var apiServerUrl string
			if args[0]=="user" {
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							option = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "email=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							email = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if option == string(enums.ATTACH_COMPANY) || option == "ac" {
					if err := v1.IsUserLoggedIn(); err != nil {
						cmd.Printf("[ERROR]: %v", err.Error())
						return nil
					}
					if file == "" {
						cmd.Printf("[ERROR]: %v", "please provide a file!")
						return nil
					}
					data, err := ioutil.ReadFile(file)
					if err != nil {
						cmd.Printf("data.Get err   #%v ", err)
						return nil
					}
					company := new(v1.Company)
					if strings.HasSuffix(file, ".yaml") {
						err = yaml.Unmarshal(data, company)
						if err != nil {
							cmd.Printf("yaml Unmarshal: %v", err)
							return nil
						}
					} else {
						err = json.Unmarshal(data, company)
						if err != nil {
							cmd.Printf("json Unmarshal: %v", err)
							return nil
						}
					}
					userService := dependency_manager.GetUserService()
					userService.Cmd(cmd).Flag(string(enums.ATTACH_COMPANY)).Company(company).Apply()
				} else if option == string(enums.RESET_PASSWORD) || option == "rp" {
					if file == "" {
						cmd.Printf("[ERROR]: %v", "please provide a file!")
						return nil
					}
					data, err := ioutil.ReadFile(file)
					if err != nil {
						cmd.Printf("data.Get err   #%v ", err)
						return nil
					}
					var passwordResetDto v1.PasswordResetDto
					if strings.HasSuffix(file, ".yaml") {
						err = yaml.Unmarshal(data, passwordResetDto)
						if err != nil {
							cmd.Printf("yaml Unmarshal: %v", err)
							return nil
						}
					} else {
						err = json.Unmarshal(data, &passwordResetDto)
						if err != nil {
							cmd.Printf("json Unmarshal: %v", err)
							return nil
						}
					}
					userService := dependency_manager.GetUserService()
					userService.Cmd(cmd).Flag(string(enums.RESET_PASSWORD)).PasswordResetDto(passwordResetDto).Apply()
				} else if option == string(enums.FORGOT_PASSWORD) || option == "fp" {
					userService := dependency_manager.GetUserService()
					userService.Cmd(cmd).Flag(string(enums.FORGOT_PASSWORD)).Email(email).Apply()
				}
			} else if args[0]=="repositories" || args[0]=="repos"{
				if err := v1.IsUserLoggedIn(); err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				userMetadata, err := v1.GetUserMetadataFromBearerToken()
				if err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option=") {
						strs := strings.Split(strings.ToUpper(each), "=")
						if len(strs) > 1 {
							option = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err = cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if file == "" {
					cmd.Printf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if option == "" {
					cmd.Printf("[ERROR]: %v", "please provide update option!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Printf("data.Get err   #%v ", err)
					return nil
				}
				repos := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, repos)
					if err != nil {
						cmd.Printf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, repos)
					if err != nil {
						cmd.Printf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option(option).Apply()
				return nil
			} else if args[0]=="applications" || args[0]=="apps"{
				if err := v1.IsUserLoggedIn(); err != nil {
					cmd.Printf("[ERROR]: %v", err.Error())
					return nil
				}
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option=") {
						strs := strings.Split(strings.ToUpper(each), "=")
						if len(strs) > 1 {
							option = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "repoid=") {
						strs := strings.Split(each, "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					}
				}
				cfg := v1.GetConfigFile()
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
				}
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				if file == "" {
					cmd.Printf("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if repoId == "" {
					cmd.Printf("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				if option == "" {
					cmd.Printf("[ERROR]: %v", "please provide update option!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Printf("data.Get err   #%v ", err)
					return nil
				}
				company := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, company)
					if err != nil {
						cmd.Printf("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, company)
					if err != nil {
						cmd.Printf("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.Cmd(cmd).Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).RepoId(repoId).Option(option).Apply()
				return nil
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
	}
}
