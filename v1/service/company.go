package service

type Company interface {
	Apply(company interface{}) error
	ApplyUpdateRepositories(company interface{}, companyId string, option string) error
	ApplyUpdateApplications(company interface{}, companyId string, repoId string, option string) error
	GetCompanyById(companyId string) (httpCode int, data []byte, err error)
	GetCompanies() (httpCode int, data []byte, err error)
}
