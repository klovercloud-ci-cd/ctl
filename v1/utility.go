package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"log"
	"os"
)

func GetUserMetadataFromBearerToken() (UserMetadata, error) {
	token, err := GetToken()
	if err != nil {
		return UserMetadata{}, nil
	}
	if token == "" {
		return UserMetadata{}, errors.New("no token found")
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	jsonbody, err := json.Marshal(claims["data"])
	if err != nil {
		log.Println(err)
	}
	userResourcePermission := UserResourcePermissionDto{}
	if err := json.Unmarshal(jsonbody, &userResourcePermission); err != nil {
		return UserMetadata{}, errors.New("no resource permissions")
	}
	return userResourcePermission.Metadata, nil
}

func GetToken() (string, error) {
	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("user is not logged in")
	}
	jsonFile, err := os.Open("config.cfg")
	if err != nil {
		return "", err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configFile Config
	err = json.Unmarshal(byteValue, &configFile)
	if err != nil {
		return "", err
	}
	if configFile.Token == "" {
		return "", errors.New("user is not logged in")
	}
	return configFile.Token, nil
}

func GetApiServerUrl() string {
	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		return ""
	}
	jsonFile, err := os.Open("config.cfg")
	if err != nil {
		return ""
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configFile Config
	err = json.Unmarshal(byteValue, &configFile)
	if err != nil {
		return ""
	}
	return configFile.ApiServerUrl
}

func GetSecurityUrl() string {
	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		return ""
	}
	jsonFile, err := os.Open("config.cfg")
	if err != nil {
		return ""
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configFile Config
	err = json.Unmarshal(byteValue, &configFile)
	if err != nil {
		return ""
	}
	return configFile.SecurityUrl
}

func IsUserLoggedIn() error {
	if _, err := GetToken(); err != nil {
		return err
	}
	return nil
}

func AddToConfigFile(token, apiServerUrl, securityUrl string) error {
	configFile := Config{
		Token:        token,
	}
	if apiServerUrl == "" {
		configFile.ApiServerUrl = "http://localhost:8080/api/v1/"
	} else {
		configFile.ApiServerUrl = apiServerUrl
	}
	if securityUrl == "" {
		configFile.SecurityUrl = "http://localhost:8085/api/v1/"
	} else {
		configFile.SecurityUrl = securityUrl
	}

	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("config.cfg")
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(configFile, "", "")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.cfg", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func AddRootIndent(b []byte, n int) []byte {
	prefix := append([]byte("\n"), bytes.Repeat([]byte(" "), n)...)
	b = append(prefix[1:], b...)
	return bytes.ReplaceAll(b, []byte("\n"), prefix)
}

func AddOrGetSecurityUrl() (string, error){
	token, _ := GetToken()
	apiServerUrl := GetApiServerUrl()
	securityUrl := GetSecurityUrl()
	if securityUrl == "" {
		err := AddToConfigFile(token, apiServerUrl, securityUrl)
		if err != nil {
			return "", err
		}
		securityUrl = GetSecurityUrl()
	}
	return securityUrl, nil
}