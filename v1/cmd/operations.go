package cmd

import (
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/spf13/cobra"
	"log"
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
		}else if args[0]=="repositories" || args[0]=="repos"{

		}else if args[0]=="application" || args[0]=="app" {

		}else if args[0]=="applications" || args[0]=="apps"{

		}

			return nil
		},
	}
}
