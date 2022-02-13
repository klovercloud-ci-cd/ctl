package cmd

import (
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func GetRepositoryById() *cobra.Command{
	return &cobra.Command{
		Use:       "get-repository",
		Short:     "Get repository by repository ID",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println(args)
			var repositoryId string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "repoid") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						repositoryId = strs[1]
					}
				}
			}
			repositoryService := dependency_manager.GetRepositoryService()
			code, data, err := repositoryService.GetRepositoryById(repositoryId)
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			if code != 200 {
				cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
				return nil
			}
			if data != nil {
				cmd.Println(string(data))
			}
			return nil
		},
	}
}

func GetApplicationsByCompanyId() *cobra.Command{
	return &cobra.Command{
		Use:       "get-applications",
		Short:     "Get applications by company ID",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println(args)
			var companyId string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "companyid") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						companyId = strs[1]
					}
				}
			}
			repositoryService := dependency_manager.GetRepositoryService()
			code, data, err := repositoryService.GetApplicationsByCompanyId(companyId)
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			if code != 200 {
				cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
				return nil
			}
			if data != nil {
				cmd.Println(string(data))
			}
			return nil
		},
	}
}