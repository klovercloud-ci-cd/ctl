package service

type Company interface {
	Apply(company interface{}) error
	ApplyUpdateRepositories(company interface{}, companyId string, option string) error
	ApplyUpdateApplications(compant interface{}, companyId string, repoId string, option string) error
}
