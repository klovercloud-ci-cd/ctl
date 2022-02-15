package service

// Oauth Oauth operations
type Oauth interface {
	Relogin(loginDto interface{}) (string, error)
}