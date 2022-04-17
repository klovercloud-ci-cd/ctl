package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
)

func GetUserMetadataFromBearerToken(token string) (UserMetadata, error) {
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

func AddRootIndent(b []byte, n int) []byte {
	prefix := append([]byte("\n"), bytes.Repeat([]byte(" "), n)...)
	b = append(prefix[1:], b...)
	return bytes.ReplaceAll(b, []byte("\n"), prefix)
}

func GetCfgPath() string {
	path := os.Getenv("KCPATH")
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func FixUrl(url string) string {
	if url[len(url)-1] != '/' {
		url += "/"
	}
	return url
}
