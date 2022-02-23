package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

func UpdateRepositories() *cobra.Command{
	return &cobra.Command{
		Use:       "update repositories",
		Short:     "Update repositories by company id with option",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var option string
			var companyId string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 1 {
						file = strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each), "option") {
					strs := strings.Split(strings.ToUpper(each), "=")
					if len(strs) > 1 {
						option = strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each), "companyid") {
					strs := strings.Split(each, "=")
					if len(strs) > 1 {
						companyId = strs[1]
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
			companyService.Cmd(cmd).Flag(string(enums.UPDATE_REPOSITORIES)).Company(*company).CompanyId(companyId).Option(option).Apply()
			return nil
		},
	}
}

func UpdateApplicationsByRepositoryId() *cobra.Command{
	return &cobra.Command{
		Use:       "update applications",
		Short:     "Update applications by repository id and company id with option",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var option string
			var companyId string
			var repoId string
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
				} else if strings.Contains(strings.ToLower(each), "companyid") {
					strs := strings.Split(each, "=")
					if len(strs) > 1 {
						companyId = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "repoid") {
					strs := strings.Split(each, "=")
					if len(strs) > 1 {
						repoId = strs[1]
					}
				}
			}
			log.Println(companyId)
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
		},
	}
}