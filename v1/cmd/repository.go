package cmd

import (
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
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
			repositoryService.Apply(string(enums.GET_REPOSITORY), repositoryId, "")
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
			repositoryService.Apply(string(enums.GET_APPLICATIONS), "", companyId)
			return nil
		},
	}
}