package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/klovercloud-ci/ctl/config"
	"log"
	"strings"
)

// UserMetadata holds users metadata
type UserMetadata struct {
	CompanyId string `json:"company_id" bson:"company_id"`
}

// UserResourcePermissionDto holds metadata and user
type UserResourcePermissionDto struct {
	Metadata  UserMetadata           `json:"metadata" bson:"-"`
	UserId    string                 `json:"user_id" bson:"user_id"`
}

func GetUserMetadataFromBearerToken() (UserMetadata, error) {
	bearerToken := GetBearerToken()
	if bearerToken == "" {
		return UserMetadata{}, errors.New("no token found")
	}
	var token string
	if len(strings.Split(bearerToken, " ")) == 2 {
		token = strings.Split(bearerToken, " ")[1]
	} else {
		return UserMetadata{}, errors.New("no token found")
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Publickey), nil
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

func GetBearerToken() string {
	return "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7Im1ldGFkYXRhIjp7ImNvbXBhbnlfaWQiOiIxMjM0NSJ9LCJ1c2VyX2lkIjoiYjg3NmVjOGEtOTY1MC00MDhlLTg0YmItZTVmM2QzNmI0NzA0IiwicmVzb3VyY2VzIjpbeyJuYW1lIjoicGlwZWxpbmUiLCJyb2xlcyI6W3sibmFtZSI6IkFETUlOIiwicGVybWlzc2lvbnMiOlt7Im5hbWUiOiJDUkVBVEUifSx7Im5hbWUiOiJSRUFEIn0seyJuYW1lIjoiVVBEQVRFIn0seyJuYW1lIjoiREVMRVRFIn1dfV19LHsibmFtZSI6InByb2Nlc3MiLCJyb2xlcyI6W3sibmFtZSI6IkFETUlOIiwicGVybWlzc2lvbnMiOlt7Im5hbWUiOiJDUkVBVEUifSx7Im5hbWUiOiJSRUFEIn0seyJuYW1lIjoiVVBEQVRFIn0seyJuYW1lIjoiREVMRVRFIn1dfV19LHsibmFtZSI6ImNvbXBhbnkiLCJyb2xlcyI6W3sibmFtZSI6IkFETUlOIiwicGVybWlzc2lvbnMiOlt7Im5hbWUiOiJDUkVBVEUifSx7Im5hbWUiOiJSRUFEIn0seyJuYW1lIjoiVVBEQVRFIn0seyJuYW1lIjoiREVMRVRFIn1dfV19LHsibmFtZSI6InJlcG9zaXRvcnkiLCJyb2xlcyI6W3sibmFtZSI6IkFETUlOIiwicGVybWlzc2lvbnMiOlt7Im5hbWUiOiJDUkVBVEUifSx7Im5hbWUiOiJSRUFEIn0seyJuYW1lIjoiVVBEQVRFIn0seyJuYW1lIjoiREVMRVRFIn1dfV19LHsibmFtZSI6ImFwcGxpY2F0aW9uIiwicm9sZXMiOlt7Im5hbWUiOiJBRE1JTiIsInBlcm1pc3Npb25zIjpbeyJuYW1lIjoiQ1JFQVRFIn0seyJuYW1lIjoiUkVBRCJ9LHsibmFtZSI6IlVQREFURSJ9LHsibmFtZSI6IkRFTEVURSJ9XX1dfSx7Im5hbWUiOiJ1c2VyIiwicm9sZXMiOlt7Im5hbWUiOiJBRE1JTiIsInBlcm1pc3Npb25zIjpbeyJuYW1lIjoiQ1JFQVRFIn0seyJuYW1lIjoiUkVBRCJ9LHsibmFtZSI6IlVQREFURSJ9LHsibmFtZSI6IkRFTEVURSJ9XX1dfV19LCJleHAiOjE5NTQ2NDg5OTgsImlhdCI6MTY0MzYwODk5OCwic3ViIjoiYjg3NmVjOGEtOTY1MC00MDhlLTg0YmItZTVmM2QzNmI0NzA0In0.vUxIiWWcbUA3aWBrOzzNMzqGb1vInZ3nJCuUyXOGv_RFSlBUhPWStLgBt0EjNq8_DuXckSa7qM-UNHlCI-z9Ma4KIFjUIBf3Q8hIOBt5WTrfyT4_v64Vocxl0qgIbyZg1zH_WP0i3yHxMuo3Pvq7DrRXOasROhncF9yhxxngyyghq9RXGzSyRzmM8KDepfuHQ5JHoPECR2oz07MwvlD8yg1nCziTFDlwJwBMWoUFc44vuo5g84JfWiSKQw6dgQyMSOZtGkhqdCwnK871T2jZRuRPKjGXhr__XKIjPFlSQNj3Cw0TngMiYYCPt0ti1VvRtZViBzIZptSHSPNo27GVnTabzkxUHTlqE95DoYtgTcibJBdy9q1SMPmifqUV5XHEGOpEcebAvPIBHTCm5xAgEq73p-dRKFZUvI8Plrayb7UcNzlJGMPULI6fZDk9ts3jZ362-vPBh41fZPPK_up3d-jzn6j779I7kpb9_kxzYXxYfsOvBCDUkgQ572575llO3hptFwFQACBxMg5pXUiZp20SImOCS475O6j8OHFrPlWoLdyxsz-1GkpkA_jC3v93aELAR5GsEy0QfYxP2GkaZUpJwtXNVejbngT0wtXKKlXi5FKs6rB1-TEMaOsBJvYLtSwxZm3UXffIyRARumBW6AHow448_HL8rqc0P62fWFg"
}

func AddRootIndent(b []byte, n int) []byte {
	prefix := append([]byte("\n"), bytes.Repeat([]byte(" "), n)...)
	b = append(prefix[1:], b...)
	return bytes.ReplaceAll(b, []byte("\n"), prefix)
}