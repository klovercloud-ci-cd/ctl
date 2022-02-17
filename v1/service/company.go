package service

// Company Company operations
type Company interface {
	Apply()
	Flag(flag string) Company
	Company(company interface{}) Company
	CompanyId(companyId string) Company
	RepoId(repoId string) Company
	Option(option string) Company
}
