package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/klovercloud-ci/ctl/v1/service"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

func Trigger() *cobra.Command {
	return &cobra.Command{
		Use:       "trigger",
		Short:     "Notify git",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						file = strs[1]
					}
				}

			}
			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("data.Get err   #%v ", err)
				return nil
			}
			webhook := new(v1.GitWebHookEvent)
			if strings.HasSuffix(file, ".yaml") {
				err = yaml.Unmarshal(data, webhook)
				if err != nil {
					log.Fatalf("yaml Unmarshal: %v", err)
					return nil
				}
			} else {
				err = json.Unmarshal(data, webhook)
				if err != nil {
					log.Fatalf("json Unmarshal: %v", err)
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
					log.Fatalf("failed to json marshal: %v", err.Error())
					return nil
				}
				err = json.Unmarshal(b, event)
				if err != nil {
					log.Fatalf("failed to convert byte int any of the git: %v", err)
					return nil
				}
				git.Apply(event, webhook.CompanyId)
			}else if  webhook.Type == v1.BITBUCKET {
				var git service.Git
				if webhook.APIVersion == "v1" {
					git = dependency_manager.GetV1GithubService()
				}
				event := new(v1.BitbucketWebHookEvent)
				b, err := json.Marshal(webhook.Event)
				if err != nil {
					log.Fatalf("failed to json marshal: %v", err.Error())
					return nil
				}
				err = json.Unmarshal(b, event)
				if err != nil {
					log.Fatalf("failed to convert byte int any of the git: %v", err)
					return nil
				}
				git.Apply(event, webhook.CompanyId)
			}
			return nil
		},
	}
}
