package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Create() *cobra.Command {
	command := cobra.Command{
		Use:       "create",
		Short:     "Create any resource [user/repositories/applications]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: %v", "user is not logged in")
				return nil
			}
			if len(args) < 1 {
				cmd.Println("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			var apiServerUrl, securityUrl, file string
			var skipSsl bool
			if strings.ToLower(args[0]) == "-u" || strings.ToLower(args[0]) == "user" {
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if file == "" {
					cmd.Println("[ERROR]: %v", "please provide a file!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Println("data.Get err   #%v ", err.Error())
					return nil
				}
				var user v1.UserRegistrationDto
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, &user)
					if err != nil {
						cmd.Println("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, &user)
					if err != nil {
						cmd.Println("json Unmarshal: %v", err)
						return nil
					}
				}
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
				userService.Cmd(cmd).SecurityUrl(cfg.SecurityUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.CREATE_USER)).User(user).Apply()
				return nil
			} else if strings.ToLower(args[0]) == "repositories" || strings.ToLower(args[0]) == "repos" || strings.ToLower(args[0]) == "-r" {
				userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
				if err != nil {
					cmd.Println("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				if companyId == "" {
					cmd.Println("[ERROR]: %v", "User got no company attached!")
					return nil
				}
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if file == "" {
					cmd.Println("[ERROR]: %v", "please provide update file!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Println("data.Get err   #%v ", err)
					return nil
				}
				repos := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, repos)
					if err != nil {
						cmd.Println("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, repos)
					if err != nil {
						cmd.Println("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option("APPEND_REPOSITORY").Apply()
				return nil
			} else if strings.ToLower(args[0]) == "applications" || strings.ToLower(args[0]) == "apps" || strings.ToLower(args[0]) == "-a" {
				var repoId string
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(each, "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
				if err != nil {
					cmd.Println("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				if companyId == "" {
					cmd.Println("[ERROR]: %v", "User got no company attached!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err := cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if file == "" {
					cmd.Println("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if repoId == "" {
					repoId = cfg.RepositoryId
					if repoId == "" {
						cmd.Println("[ERROR]: %v", "please provide repository id!")
						return nil
					}
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Println("data.Get err   #%v ", err)
					return nil
				}
				company := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, company)
					if err != nil {
						cmd.Println("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, company)
					if err != nil {
						cmd.Println("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Cmd(cmd).Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).CompanyId(companyId).RepoId(repoId).Option("APPEND_APPLICATION").Apply()
				return nil
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli create {user | -u} {file | -f}=USER_PAYLOAD [{option | -o} [apiserver=APISERVER_URL | security=SECURITY_SERVER_URL]]... [--skipssl] \n" +
		"  cli create {repositories | repos | -r} {file | -f}=REPOSITORY_PAYLOAD [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli create {applications | apps | -a} {file | -f}=APPLICATION_PAYLOAD repo=REPOSITORY_ID [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli help create \n" +
		"\nOptions: \n" +
		"  option | -o\t" + "Provide apiserver or security server url option while creating user resource. \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}

func Registration() *cobra.Command {
	command := cobra.Command{
		Use:       "register",
		Short:     "Register user",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file, apiServerUrl, securityUrl string
			var skipSsl bool
			for idx, each := range args {
				if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						file = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
					if idx+1 < len(args) {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
			}
			if file == "" {
				cmd.Println("[ERROR]: %v", "please provide a file!")
				return nil
			}
			data, err := ioutil.ReadFile(file)
			if err != nil {
				cmd.Println("data.Get err   #%v ", err.Error())
				return nil
			}
			var user v1.UserRegistrationDto
			if strings.HasSuffix(file, ".yaml") {
				err = yaml.Unmarshal(data, &user)
				if err != nil {
					cmd.Println("yaml Unmarshal: %v", err)
					return nil
				}
			} else {
				err = json.Unmarshal(data, &user)
				if err != nil {
					cmd.Println("json Unmarshal: %v", err)
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
			userService.SecurityUrl(cfg.SecurityUrl).SkipSsl(skipSsl).Cmd(cmd).Flag(string(enums.CREATE_ADMIN)).User(user).Apply()
			return nil
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli register {file | -f}=USER_REGISTRATION_PAYLOAD [{option | -o} [apiserver=APISERVER_URL | security=SECURITY_SERVER_URL]]... [--skipssl]\n" +
		"  cli help register \n" +
		"\nOptions: \n" +
		"  option | -o\t" + "Provide apiserver or security server url option \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}

func Describe() *cobra.Command {
	command := cobra.Command{
		Use:       "describe",
		Short:     "Describe resource [company/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: %v", "user is not logged in")
				return nil
			}
			if len(args) < 1 {
				cmd.Println("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
			if err != nil {
				cmd.Println("[ERROR]: %v", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			if companyId == "" {
				cmd.Println("[ERROR]: %v", "User got no company attached!")
				return nil
			}
			var apiServerUrl string
			if strings.ToLower(args[0]) == "company" || strings.ToLower(args[0]) == "-c" {
				loadRepo := false
				loadApp := false
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadrepositories=") || strings.Contains(strings.ToLower(args[idx+1]), "loadrepos=") || strings.Contains(strings.ToLower(args[idx+1]), "lr=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									if strs[1] == "true" {
										loadRepo = true
									}
								}
							} else if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") || strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la=") {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err := cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Company").Cmd(cmd).Flag(string(enums.GET_COMPANY_BY_ID)).CompanyId(companyId).Option("loadRepositories=" + strconv.FormatBool(loadRepo) + "&loadApplications=" + strconv.FormatBool(loadApp)).Apply()
			} else if strings.ToLower(args[0]) == "repository" || strings.ToLower(args[0]) == "repo" || strings.ToLower(args[0]) == "-r" {
				var repoId string
				var skipSsl bool
				loadApp := false
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "repository=") || strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") || strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la=") {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err := cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				repoId = cfg.RepositoryId
				if repoId == "" {
					cmd.Println("[ERROR]: %v", "please provide repository id!")
					return nil
				}
				repositoryService := dependency_manager.GetRepositoryService()
				repositoryService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORY)).Repo(repoId).Option("loadApplications=" + strconv.FormatBool(loadApp)).Apply()
			} else if strings.ToLower(args[0]) == "application" || strings.ToLower(args[0]) == "app" || strings.ToLower(args[0]) == "-a" {
				var repoId, appId string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "repository=") || strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "application=") || strings.Contains(strings.ToLower(each), "app=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if repoId == "" {
					repoId = cfg.RepositoryId
					if repoId == "" {
						cmd.Println("[ERROR]: %v", "please provide repository id!")
						return nil
					}
				}
				if appId == "" {
					cmd.Println("[ERROR]: %v", "please provide application id!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				applicationService := dependency_manager.GetApplicationService()
				applicationService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATION)).RepoId(repoId).ApplicationId(appId).Apply()
			} else if strings.ToLower(args[0]) == "process" || strings.ToLower(args[0]) == "-p" {
				var processId string
				var follow, skipSsl bool
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "processid=") || strings.Contains(strings.ToLower(each), "process=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							processId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					} else if strings.Contains(strings.ToLower(each), "follow") || strings.Contains(strings.ToLower(each), "-f") {
						follow = true
					}
				}
				return getPipeline(cmd, processId, "get_pipeline", cfg.ApiServerUrl, cfg.Token, follow, true, skipSsl, make(map[string]string))
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli describe {company | -c} [{option | -o} [{loadrepositories | loadrepos | lr}={true | false} | {loadapplications | loadapps | la}={true | false} | apiserver=APISERVER_URL]]... [--skipssl] \n" +
		"  cli describe {repository | repo | -r} {repository | repo}=REPOSITORY_ID [{option | -o} [{loadapplications | loadapps | la}={true | false} | apiserver=APISERVER_URL]]... [--skipssl] \n" +
		"  cli describe {application | app | -a} {repository | repo}=REPOSITORY_ID {application | app}=APPLICATION_ID [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli describe {process | -p} {processid | process}=PROCESS_ID [--skipssl] \n" +
		"  cli help describe \n" +
		"\nOptions: \n" +
		"  option | -o\t" + "Provide load repositories, load applications or apiserver url option \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}
func getPipeline(cmd *cobra.Command, processId, action, url, token string, follow, firstFollow, skipSsl bool, stepMap map[string]string) error {
	pipelineService := dependency_manager.GetPipelineService()
	code, data, err := pipelineService.SkipSsl(skipSsl).Get(processId, action, url, token)
	if err != nil {
		cmd.Println("[ERROR]: ", err.Error())
		return nil
	} else if code != 200 {
		cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		return nil
	} else if data != nil {
		byteBody, _ := json.Marshal(data)
		var pipeline v1.Pipeline
		err := json.Unmarshal(byteBody, &pipeline)
		if err != nil {
			return err
		}
		var tableData [][]string
		for _, each := range pipeline.Steps {
			if val, ok := stepMap[each.Name]; !ok || val != each.Status {
				stepMap[each.Name] = each.Status
				name := each.Name
				status := each.Status
				for i := len(each.Name); i < 18; i++ {
					name += " "
				}
				for i := len(each.Status); i < 15; i++ {
					status += " "
				}
				process := []string{"api/v1     ", "Process", name, strings.Title(status)}
				tableData = append(tableData, process)
			}
		}
		if len(tableData) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			if firstFollow {
				table.SetHeader([]string{"Api Version", "Kind   ", "Step              ", "Status         "})
				if follow {
					table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
				} else {
					table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
				}
			} else {
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			}
			table.AppendBulk(tableData)
			table.Render()
		}
		if follow {
			return getPipeline(cmd, processId, action, url, token, follow, false, skipSsl, stepMap)
		}
	}
	return nil
}
func List() *cobra.Command {
	command := cobra.Command{
		Use:       "list",
		Short:     "List resources [repository/application/process/k8sobjs]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: user is not logged in")
				return nil
			}
			if len(args) < 1 {
				cmd.Println("[ERROR]: please provide a resource name!")
				return nil
			}
			userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			companyId := userMetadata.CompanyId
			if companyId == "" {
				cmd.Println("[ERROR]: ", "User got no company attached!")
				return nil
			}
			var apiServerUrl string
			if strings.ToLower(args[0]) == "repositories" || strings.ToLower(args[0]) == "repos" || strings.ToLower(args[0]) == "-r" {
				loadApp := false
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "loadapplications=") || strings.Contains(strings.ToLower(args[idx+1]), "loadapps=") || strings.Contains(strings.ToLower(args[idx+1]), "la=") {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Repository").Cmd(cmd).Flag(string(enums.GET_REPOSITORIES)).CompanyId(companyId).Option("loadApplications=" + strconv.FormatBool(loadApp)).Apply()
			} else if strings.ToLower(args[0]) == "applications" || strings.ToLower(args[0]) == "apps" || strings.ToLower(args[0]) == "-a" {
				var repoId string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "repository=") || strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if repoId == "" {
					applicationService := dependency_manager.GetApplicationService()
					applicationService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Application").Cmd(cmd).Flag(string(enums.GET_All_APPLICATIONS)).Apply()
				} else {
					repositoryService := dependency_manager.GetRepositoryService()
					repositoryService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Application").Cmd(cmd).Flag(string(enums.GET_APPLICATIONS)).Repo(repoId).Apply()
				}
			} else if strings.ToLower(args[0]) == "process" || strings.ToLower(args[0]) == "-p" {
				var repoId, appId string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "repository=") || strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "application=") || strings.Contains(strings.ToLower(each), "app=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							appId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if repoId == "" {
					repoId = cfg.RepositoryId
					if repoId == "" {
						cmd.Println("[ERROR]: please provide repository id!")
						return nil
					}
				}
				if appId == "" {
					cmd.Println("[ERROR]: please provide application id!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				processService := dependency_manager.GetProcessService()
				processService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Kind("Process").Flag(string(enums.GET_PROCESS)).Cmd(cmd).RepoId(repoId).ApplicationId(appId).Apply()
			} else if strings.ToLower(args[0]) == "agents" {
				var processId string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "processid=") || strings.Contains(strings.ToLower(each), "process=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							processId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if processId == "" {
					cmd.Println("[ERROR]: please provide process id!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				pipelineService := dependency_manager.GetPipelineService()
				pipelineService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_AGENTS)).Cmd(cmd).ProcessId(processId).Apply()
			} else if strings.ToLower(args[0]) == "objects" || strings.ToLower(args[0]) == "objs" {
				var processId string
				var agentList []string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "processid=") || strings.Contains(strings.ToLower(each), "process=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							processId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "agent=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							agentList = append(agentList, strs[1])
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if processId == "" {
					cmd.Println("[ERROR]: please provide process id!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if len(agentList) == 0 {
					pipelineService := dependency_manager.GetPipelineService()
					_, data, err := pipelineService.Get(processId, "get_pipeline", cfg.ApiServerUrl, cfg.Token)
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
					var pipeline v1.Pipeline
					byteBody, _ := json.Marshal(data)
					err = json.Unmarshal(byteBody, &pipeline)
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
					agentsMap := make(map[string]bool)
					for _, step := range pipeline.Steps {
						agent := step.Params["agent"]
						if step.Type == "DEPLOY" {
							if _, ok := agentsMap[agent]; !ok {
								agentsMap[agent] = true
								agentList = append(agentList, agent)
							}
						}
					}
				}
				agentService := dependency_manager.GetAgentService()
				agentService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_K8SOBJS)).Cmd(cmd).AgentList(agentList).ProcessId(processId).Apply()
			} else if strings.ToLower(args[0]) == "pods" {
				var agent, processId, owner, ownerid string
				var skipSsl bool
				for idx, each := range args {
					if strings.Contains(strings.ToLower(each), "processid=") || strings.Contains(strings.ToLower(each), "process=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							processId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "agent=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							agent = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "owner=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							owner = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "ownerid=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							ownerid = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option") || strings.Contains(strings.ToLower(each), "-o") {
						if idx+1 < len(args) {
							if strings.Contains(strings.ToLower(args[idx+1]), "apiserver=") {
								strs := strings.Split(strings.ToLower(args[idx+1]), "=")
								if len(strs) > 1 {
									apiServerUrl = strs[1]
								}
							}
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if agent == "" {
					cmd.Println("[ERROR]: please provide agent name!")
					return nil
				}
				if processId == "" {
					cmd.Println("[ERROR]: please provide process id!")
					return nil
				}
				if owner != string(enums.DEPLOYMENT) && owner != string(enums.DAEMONSET) && owner != string(enums.REPLICASET) && owner != string(enums.STATEFULSET) {
					cmd.Println("[ERROR]: please provide valid owner reference type! [DEPLOYMENT/DAEMONSET/REPLICASET/STATEFULSET]")
					return nil
				}
				if ownerid == "" {
					cmd.Println("[ERROR]: please provide owner reference id!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				agentService := dependency_manager.GetAgentService()
				if owner == string(enums.DEPLOYMENT) {
					agentService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_PODS_BY_DEPLOYMENT)).OwnerReferenceId(ownerid).Cmd(cmd).Name(agent).ProcessId(processId).Apply()
				} else if owner == string(enums.DAEMONSET) {
					agentService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_PODS_BY_DAEMONSET)).OwnerReferenceId(ownerid).Cmd(cmd).Name(agent).ProcessId(processId).Apply()
				} else if owner == string(enums.REPLICASET) {
					agentService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_PODS_BY_REPLICASET)).OwnerReferenceId(ownerid).Cmd(cmd).Name(agent).ProcessId(processId).Apply()
				} else if owner == string(enums.STATEFULSET) {
					agentService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Flag(string(enums.GET_PODS_BY_STATEFULSET)).OwnerReferenceId(ownerid).Cmd(cmd).Name(agent).ProcessId(processId).Apply()
				}
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli list {repositories | repos | -r} [{option | -o} [{loadapplications | loadapps | la}={true | false} | apiserver=APISERVER_URL]]... [--skipssl] \n" +
		"  cli list {applications | apps | -a} {repository | repo}=REPOSITORY_ID [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli list {process | -p} {repository | repo}=REPOSITORY_ID {application | app}=APPLICATION_ID [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli list agents {processid | process}=PROCESS_ID [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli list {objects | objs} {processid | process}=PROCESS_ID [agent=AGENT_NAME] [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli list pods {processid | process}=PROCESS_ID agent=AGENT_NAME owner={deployment | daemonset | replicaset | statefulset} ownerid=[OWNER_ID] [{option | -o} apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli help list \n" +
		"\nOptions: \n" +
		"  option | -o\t" + "Provide load applications or apiserver url option \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}

func Set() *cobra.Command {
	command := cobra.Command{
		Use:       "set",
		Short:     "Set Default Id [repository]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: %v", "user is not logged in")
				return nil
			}
			if len(args) < 1 {
				cmd.Println("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			var repositoryId string
			if strings.Contains(strings.ToLower(args[0]), "repository=") || strings.Contains(strings.ToLower(args[0]), "repo=") {
				strs := strings.Split(strings.ToLower(args[0]), "=")
				if len(strs) > 1 {
					repositoryId = strs[1]
				}
			} else if strings.ToLower(args[0]) == "-r" {
				if len(args) > 1 {
					repositoryId = args[1]
				}
			}
			if repositoryId != "" {
				cfg.RepositoryId = repositoryId
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
				cmd.Println("[Success]: " + "successfully updated repository id")
			}
			return nil
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli set {repository | repo}=REPOSITORY_ID \n" +
		"  cli set -r REPOSITORY_ID \n" +
		"  cli help set \n" +
		"\nOptions: \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}

func Update() *cobra.Command {
	command := cobra.Command{
		Use:       "update",
		Short:     "Update resource [user/repository/application]",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Println("[ERROR]: %v", "please provide a resource name!")
				return nil
			}
			var file, option, repoId, email, apiServerUrl string
			var skipSsl bool
			if strings.ToLower(args[0]) == "user" || strings.ToLower(args[0]) == "-u" {
				for _, each := range args {
					if strings.Contains(strings.ToLower(each), "file=") || strings.Contains(strings.ToLower(each), "-f=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 0 {
							file = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "option=") || strings.Contains(strings.ToLower(each), "-o=") {
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
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
					err := cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if option == string(enums.ATTACH_COMPANY) || option == "ac" {
					if cfg.Token == "" {
						cmd.Println("[ERROR]: %v", "user is not logged in")
						return nil
					}
					if file == "" {
						cmd.Println("[ERROR]: %v", "please provide a file!")
						return nil
					}
					data, err := ioutil.ReadFile(file)
					if err != nil {
						cmd.Println("data.Get err   #%v ", err)
						return nil
					}
					company := new(v1.Company)
					if strings.HasSuffix(file, ".yaml") {
						err = yaml.Unmarshal(data, company)
						if err != nil {
							cmd.Println("yaml Unmarshal: %v", err)
							return nil
						}
					} else {
						err = json.Unmarshal(data, company)
						if err != nil {
							cmd.Println("json Unmarshal: %v", err)
							return nil
						}
					}
					userService := dependency_manager.GetUserService()
					userService.SecurityUrl(cfg.SecurityUrl).SkipSsl(skipSsl).Token(cfg.Token).Cmd(cmd).Flag(string(enums.ATTACH_COMPANY)).Company(company).Apply()
				} else if strings.ToLower(option) == string(enums.RESET_PASSWORD) || strings.ToLower(option) == "rp" {
					if file == "" {
						cmd.Println("[ERROR]: %v", "please provide a file!")
						return nil
					}
					data, err := ioutil.ReadFile(file)
					if err != nil {
						cmd.Println("data.Get err   #%v ", err)
						return nil
					}
					var passwordResetDto v1.PasswordResetDto
					if strings.HasSuffix(file, ".yaml") {
						err = yaml.Unmarshal(data, passwordResetDto)
						if err != nil {
							cmd.Println("yaml Unmarshal: %v", err)
							return nil
						}
					} else {
						err = json.Unmarshal(data, &passwordResetDto)
						if err != nil {
							cmd.Println("json Unmarshal: %v", err)
							return nil
						}
					}
					userService := dependency_manager.GetUserService()
					userService.SecurityUrl(cfg.SecurityUrl).SkipSsl(skipSsl).Cmd(cmd).Flag(string(enums.RESET_PASSWORD)).PasswordResetDto(passwordResetDto).Apply()
				} else if option == string(enums.FORGOT_PASSWORD) || option == "fp" {
					userService := dependency_manager.GetUserService()
					userService.SecurityUrl(cfg.SecurityUrl).SkipSsl(skipSsl).Cmd(cmd).Flag(string(enums.FORGOT_PASSWORD)).Email(email).Apply()
				}
			} else if strings.ToLower(args[0]) == "repositories" || strings.ToLower(args[0]) == "repos" || strings.ToLower(args[0]) == "-r" {
				cfg := v1.GetConfigFile()
				if cfg.Token == "" {
					cmd.Println("[ERROR]: %v", "user is not logged in")
					return nil
				}
				userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
				if err != nil {
					cmd.Println("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				if companyId == "" {
					cmd.Println("[ERROR]: %v", "User got no company attached!")
					return nil
				}
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
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err = cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if file == "" {
					cmd.Println("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if option == "" {
					cmd.Println("[ERROR]: %v", "please provide update option!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Println("data.Get err   #%v ", err)
					return nil
				}
				repos := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, repos)
					if err != nil {
						cmd.Println("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, repos)
					if err != nil {
						cmd.Println("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*repos).CompanyId(companyId).Option(option).Apply()
				return nil
			} else if strings.ToLower(args[0]) == "applications" || strings.ToLower(args[0]) == "apps" || strings.ToLower(args[0]) == "-a" {
				cfg := v1.GetConfigFile()
				if cfg.Token == "" {
					cmd.Println("[ERROR]: %v", "user is not logged in")
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
					} else if strings.Contains(strings.ToLower(each), "repository=") || strings.Contains(strings.ToLower(each), "repo=") {
						strs := strings.Split(each, "=")
						if len(strs) > 1 {
							repoId = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "apiserver=") {
						strs := strings.Split(strings.ToLower(each), "=")
						if len(strs) > 1 {
							apiServerUrl = strs[1]
						}
					} else if strings.Contains(strings.ToLower(each), "--skipssl") {
						skipSsl = true
					}
				}
				userMetadata, err := v1.GetUserMetadataFromBearerToken(cfg.Token)
				if err != nil {
					cmd.Println("[ERROR]: %v", err.Error())
					return nil
				}
				companyId := userMetadata.CompanyId
				if companyId == "" {
					cmd.Println("[ERROR]: %v", "User got no company attached!")
					return nil
				}
				if apiServerUrl == "" {
					if cfg.ApiServerUrl == "" {
						cmd.Println("[ERROR]: Api server url not found!")
						return nil
					}
				} else {
					cfg.ApiServerUrl = apiServerUrl
					err := cfg.Store()
					if err != nil {
						cmd.Println("[ERROR]: ", err.Error())
						return nil
					}
				}
				if file == "" {
					cmd.Println("[ERROR]: %v", "please provide update file!")
					return nil
				}
				if repoId == "" {
					repoId = cfg.RepositoryId
					if repoId == "" {
						cmd.Println("[ERROR]: %v", "please provide repository id!")
						return nil
					}
				}
				if option == "" {
					cmd.Println("[ERROR]: %v", "please provide update option!")
					return nil
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					cmd.Println("data.Get err   #%v ", err)
					return nil
				}
				company := new(interface{})
				if strings.HasSuffix(file, ".yaml") {
					err = yaml.Unmarshal(data, company)
					if err != nil {
						cmd.Println("yaml Unmarshal: %v", err)
						return nil
					}
				} else {
					err = json.Unmarshal(data, company)
					if err != nil {
						cmd.Println("json Unmarshal: %v", err)
						return nil
					}
				}
				companyService := dependency_manager.GetCompanyService()
				companyService.ApiServerUrl(cfg.ApiServerUrl).SkipSsl(skipSsl).Token(cfg.Token).Cmd(cmd).Flag(string(enums.UPDATE_APPLICATIONS)).Company(*company).RepoId(repoId).CompanyId(companyId).Option(option).Apply()
				return nil
			} else {
				cmd.Println("[ERROR]: Wrong command")
				return nil
			}
			return nil
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli update {user | -u} {option | -o}={attach_company | ac} {file | -f}=COMPANY_ATTACH_PAYLOAD [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli update {user | -u} {option | -o}={forgot_password | fp} email={USER_EMAIL} [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli update {user | -u} {option | -o}={reset_password | rp} {file | -f}=RESET_PASSWORD_PAYLOAD [apiserver=APISERVER_URL]  [--skipssl]\n" +
		"  cli update {repositories | repos | -r} {file | -f}=REPOSITORY_UPDATE_PAYLOAD option={APPEND_REPOSITORY | SOFT_DELETE_REPOSITORY | DELETE_REPOSITORY} [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli update {applications | apps | -a} {repository | repo}=REPOSITORY_ID {file | -f}=APPLICATION_UPDATE_PAYLOAD option={APPEND_APPLICATION | SOFT_DELETE_APPLICATION | DELETE_APPLICATION} [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli help update \n" +
		"\nOptions: \n" +
		"  option | -o\t" + "Provide resource update option \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}
