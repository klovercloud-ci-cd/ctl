package service

// Pipeline Pipeline operations.
type Pipeline interface {
	Logs(url, page, limit string) (httpCode int, data interface{}, err error)
	Token(token string) Pipeline
	SkipSsl(skipSsl bool) Pipeline
	Get(processId string, action string, url string, token string) (httpCode int, data interface{}, err error)
}
