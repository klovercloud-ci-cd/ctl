package cmd

import (
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/klovercloud-ci/ctl/enums"
	"github.com/spf13/cobra"
	"log"
	"strings"
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
					var companyId string
					companyService := dependency_manager.GetCompanyService()
					companyService.Flag(string(enums.GET_APPLICATIONS)).Company(companyId).RepoId(strs[1]).Apply()
				}else{
					log.Fatalf("[ERROR]: %v", "please provide a valid sub resource[repo/repository]!")
					return nil
				}

			}

			return nil
		},
	}
}