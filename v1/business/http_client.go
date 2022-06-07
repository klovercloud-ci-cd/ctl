package business

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/klovercloud-ci/ctl/v1/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type httpClientService struct {
}

// Put method that fires a Put request.
func (h httpClientService) Put(url string, header map[string]string, body []byte, skipSsl bool) (httpCode int, err error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err.Error())
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	var client http.Client
	if skipSsl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = http.Client{Transport: tr}
	} else {
		client = http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(resp.Body)
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, err
		} else {
			return resp.StatusCode, errors.New(string(body))
		}
	}
	return resp.StatusCode, nil
}

// Get method that fires a get request.
func (h httpClientService) Get(url string, header map[string]string, skipSsl bool) (httpCode int, body []byte, err error) {
	var client http.Client
	if skipSsl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = http.Client{Transport: tr}
	} else {
		client = http.Client{}
	}
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return res.StatusCode, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(res.Body)
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		jsonDataFromHttp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return res.StatusCode, nil, err
		}
		return res.StatusCode, jsonDataFromHttp, nil
	}
	return res.StatusCode, nil, errors.New("Status: " + res.Status + ", code: " + strconv.Itoa(res.StatusCode))
}

// Post method that fires a Post request.
func (h httpClientService) Post(url string, header map[string]string, body []byte, skipSsl bool) (httpCode int, data []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	var client http.Client
	if skipSsl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = http.Client{Transport: tr}
	} else {
		client = http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, nil, err
		} else {
			return resp.StatusCode, nil, errors.New(string(body))
		}
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}

// NewHttpClientService returns HttpClient type service
func NewHttpClientService() service.HttpClient {
	return &httpClientService{}
}
