package service

// Oauth Oauth operations
type Oauth interface {
	Apply(loginDto interface{}) (string, error)
	SecurityUrl(security string)	Oauth
}