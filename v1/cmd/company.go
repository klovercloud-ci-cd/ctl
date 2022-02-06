package cmd

import (
	"encoding/json"
	v1 "github.com/klovercloud-ci-cd/integration-manager/core/v1"
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
		Short:     "Create Company",
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
			company := new(v1.Company)
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
			if err = companyService.Store(*company); err != nil {
				log.Fatalf("[ERROR]: %v", err)
			}
			return nil
		},
	}
}
