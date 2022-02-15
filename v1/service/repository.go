package service

// Repository Repository operations
type Repository interface {
	Apply(flag, repositoryId, companyId string)
	//GetRepositoryById(repositoryId string) (httpCode int, data []byte, err error)
	//GetApplicationsByCompanyId(companyId string) (httpCode int, data []byte, err error)
}
