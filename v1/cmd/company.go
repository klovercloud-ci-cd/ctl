package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

func CreateCompany() *cobra.Command{
	return &cobra.Command{
		Use:       "create-company",
		Short:     "Create company",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println(args)
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
			if err = companyService.Apply(*company); err != nil {
				log.Fatalf("[ERROR]: %v", err)
			}
			return nil
		},
	}
}

func UpdateCompanyRepositories() *cobra.Command{
	return &cobra.Command{
		Use:       "update-company-repositories",
		Short:     "Update company repositories",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var option string
			var companyId string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						file = strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each), "option") {
					strs := strings.Split(strings.ToUpper(each), "=")
					if len(strs) > 0 {
						option = strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each), "companyid") {
					strs := strings.Split(each, "=")
					if len(strs) > 0 {
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
			if err = companyService.ApplyUpdateRepositories(*company, companyId, option); err != nil {
				log.Fatalf("[ERROR]: %v", err)
			}
			return nil
		},
	}
}

func UpdateRepositoryApplications() *cobra.Command{
	return &cobra.Command{
		Use:       "update-repository-applications",
		Short:     "Update repositories applications",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			var file string
			var option string
			var companyId string
			var repoId string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "file") || strings.Contains(strings.ToLower(each), "-f") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						file = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "option") {
					strs := strings.Split(strings.ToUpper(each), "=")
					if len(strs) > 0 {
						option = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "companyid") {
					strs := strings.Split(each, "=")
					if len(strs) > 0 {
						companyId = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "repoid") {
					strs := strings.Split(each, "=")
					if len(strs) > 0 {
						repoId = strs[1]
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
			if err = companyService.ApplyUpdateApplications(*company, companyId, repoId, option); err != nil {
				log.Fatalf("[ERROR]: %v", err)
			}
			return nil
		},
	}
}

func GetCompanyById() *cobra.Command{
	return &cobra.Command{
		Use:       "get-company",
		Short:     "Get company by company ID",
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
			companyService := dependency_manager.GetCompanyService()
			code, data, err := companyService.GetCompanyById(companyId)
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

func GetCompanies() *cobra.Command{
	return &cobra.Command{
		Use:       "get-companies",
		Short:     "Get all companies",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println(args)
			companyService := dependency_manager.GetCompanyService()
			code, data, err := companyService.GetCompanies()
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

func GetRepositoriesByCompanyId() *cobra.Command{
	return &cobra.Command{
		Use:       "get-repositories",
		Short:     "Get repositories by company id",
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
			companyService := dependency_manager.GetCompanyService()
			code, data, err := companyService.GetRepositoriesByCompanyId(companyId)
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