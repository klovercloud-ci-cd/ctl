package service

// HttpClient HttpClient operations.
type HttpClient interface {
	Get(url string, header map[string]string, skipSsl bool) (httpCode int, body []byte, err error)
	Post(url string, header map[string]string, body []byte, skipSsl bool) (httpCode int, data []byte, err error)
	Put(url string, header map[string]string, body []byte, skipSsl bool) (httpCode int, err error)
}
