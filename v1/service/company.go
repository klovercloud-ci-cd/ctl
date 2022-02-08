package service

type Company interface {
	Apply(company interface{}) error
}
