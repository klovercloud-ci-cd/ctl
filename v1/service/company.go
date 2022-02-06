package service
import v1 "github.com/klovercloud-ci-cd/integration-manager/core/v1"

type Company interface {
	Store(company v1.Company) error
}
