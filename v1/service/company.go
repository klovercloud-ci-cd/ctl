package service

// Company Company operations
type Company interface {
	Apply(flag string, company interface{}, companyId, repoId, option string)
}
