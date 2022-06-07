package service

// Oauth Oauth operations
type Oauth interface {
	Apply(loginDto interface{}) (string, error, int)
	SecurityUrl(securityUrl string) Oauth
	SkipSsl(skipSsl bool) Oauth
}
