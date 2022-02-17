package service

// Repository Repository operations
type Repository interface {
	Apply()
	Flag(flag string) Repository
	CompanyId(companyId string) Repository
	Repo(repoId string) Repository
}
