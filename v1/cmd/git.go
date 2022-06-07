package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func Trigger() *cobra.Command {
	command := cobra.Command{
		Use:       "trigger",
		Short:     "Notify git",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: %v", "user is not logged in")
				return nil
			}
			var file string
			var apiServerUrl string
			var skipSsl bool
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
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
			}
			data, err := ioutil.ReadFile(file)
			if err != nil {
				cmd.Println("data.Get err   #%v ", err)
				return nil
			}
			webhook := new(v1.GitWebHookEvent)
			if strings.HasSuffix(file, ".yaml") {
				err = yaml.Unmarshal(data, webhook)
				if err != nil {
					cmd.Println("yaml Unmarshal: %v", err)
					return nil
				}
			} else {
				err = json.Unmarshal(data, webhook)
				if err != nil {
					cmd.Println("json Unmarshal: %v", err)
					return nil
				}
			}
			if webhook.Type == v1.GITHUB {
				var git service.Git
				if webhook.APIVersion == "v1" {
					git = dependency_manager.GetV1GithubService()
				}
				event := new(v1.GithubWebHookEvent)
				b, err := json.Marshal(webhook.Event)
				if err != nil {
					cmd.Println("failed to json marshal: %v", err.Error())
					return nil
				}
				err = json.Unmarshal(b, event)
				if err != nil {
					cmd.Println("failed to convert byte int any of the git: %v", err)
					return nil
				}
				err = git.Apply(event, webhook.CompanyId, cfg.ApiServerUrl, cfg.Token, skipSsl)
				if err != nil {
					return err
				}
			} else if webhook.Type == v1.BITBUCKET {
				var git service.Git
				if webhook.APIVersion == "v1" {
					git = dependency_manager.GetV1BitbucketService()
				}
				event := new(v1.BitbucketWebHookEvent)
				b, err := json.Marshal(webhook.Event)
				if err != nil {
					cmd.Println("failed to json marshal: %v", err.Error())
					return nil
				}
				err = json.Unmarshal(b, event)
				if err != nil {
					cmd.Println("failed to convert byte int any of the git: %v", err)
					return nil
				}
				err = git.Apply(event, webhook.CompanyId, cfg.ApiServerUrl, cfg.Token, skipSsl)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli trigger {file | -f}=WEBHOOK_PAYLOAD [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli help trigger \n" +
		"\nOptions: \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}
