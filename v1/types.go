package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

// ResponseDTO Http response dto
type ResponseDTO struct {
	Metadata interface{}      `json:"_metadata"`
	Data     interface{} `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

// RepositoryDto contains repository dto
type RepositoryDto struct {
	ApiVersion string `bson:"apiVersion" json:"apiVersion"`
	Kind       string `bson:"kind" json:"kind"`
	Repository Repository	`bson:"repository" json:"repository"`
}

// Repositories contains repository list
type Repositories []struct {
	Id           string        `bson:"id" json:"id"`
	Type         string        `bson:"type" json:"type"`
	Applications []Application `bson:"applications" json:"applications"`
}

// Repository contains repository info
type Repository struct {
	Id           string        `bson:"id" json:"id"`
	Type         string        `bson:"type" json:"type"`
	Applications []Application `bson:"applications" json:"applications"`
}

// ApplicationDto contains application dto
type ApplicationDto struct {
	ApiVersion string `bson:"apiVersion" json:"apiVersion"`
	Kind       string `bson:"kind" json:"kind"`
	Application Application	`bson:"application" json:"application"`
}

// Application contains application info
type Application struct {
	MetaData ApplicationMetadata `bson:"_metadata" json:"_metadata"`
	Url      string              `bson:"url" json:"url"`
}

// ApplicationMetadata contains application metadata info
type ApplicationMetadata struct {
	Labels           map[string]string `bson:"labels" json:"labels"`
	Id               string            `bson:"id" json:"id"`
	Name             string            `bson:"name" json:"name"`
	IsWebhookEnabled bool            `bson:"is_webhook_enabled" json:"is_webhook_enabled"`
}

// Applications contains application list
type Applications []struct {
	MetaData ApplicationMetadata `bson:"_metadata" json:"_metadata"`
	Url      string              `bson:"url" json:"url"`
}

// CompanyDto contains company data
type CompanyDto struct {
	ApiVersion string `bson:"apiVersion" json:"apiVersion"`
	Kind       string `bson:"kind" json:"kind"`
	Company Company	`bson:"company" json:"company"`
}

// Company contains company data
type Company struct {
	MetaData     CompanyMetadata      `bson:"_metadata" json:"_metadata"`
	Id           string               `bson:"id" json:"id"`
	Name         string               `bson:"name" json:"name"`
	Repositories []Repository         `bson:"repositories" json:"repositories"`
	Status       string `bson:"status" json:"status"`
}

// CompanyMetadata contains company metadata info
type CompanyMetadata struct {
	Labels                    map[string]string `bson:"labels" json:"labels" yaml:"labels"`
	NumberOfConcurrentProcess int64             `bson:"number_of_concurrent_process" json:"number_of_concurrent_process" yaml:"number_of_concurrent_process"`
	TotalProcessPerDay        int64             `bson:"total_process_per_day" json:"total_process_per_day" yaml:"total_process_per_day"`
}

// Processes contains process list
type Processes []struct {
	ProcessId    string                 `bson:"process_id" json:"process_id"`
	AppId        string                 `bson:"app_id" json:"app_id"`
	RepositoryId string                 `bson:"repository_id" json:"repository_id"`
	Data         map[string]interface{} `bson:"data" json:"data"`
}

// UserRegistrationDto dto that holds user registration info.
type UserRegistrationDto struct {
	Metadata           UserMetadata           `json:"metadata"`
	ID                 string                 `json:"id" bson:"id"`
	FirstName          string                 `json:"first_name" bson:"first_name" `
	LastName           string                 `json:"last_name" bson:"last_name"`
	Email              string                 `json:"email" bson:"email" `
	Phone              string                 `json:"phone" bson:"phone"`
	Password           string                 `json:"password" bson:"password" `
	Status             string           	  `json:"status" bson:"status"`
	CreatedDate        time.Time              `json:"created_date" bson:"created_date"`
	UpdatedDate        time.Time              `json:"updated_date" bson:"updated_date"`
	AuthType           string       		  `json:"auth_type" bson:"auth_type"`
	ResourcePermission UserResourcePermission `json:"resource_permission" bson:"resource_permission"`
}

// UserResourcePermission dto that holds metadata, user and resource wise roles.
type UserResourcePermission struct {
	Resources []ResourceWiseRoles `json:"resources" bson:"resources"`
}

// ResourceWiseRoles dto that holds resource wise role dtos.
type ResourceWiseRoles struct {
	Name  string `json:"name" bson:"name"`
	Roles []Role `json:"roles" bson:"roles"`
}

// Role dto that holds role name.
type Role struct {
	Name string `json:"name" bson:"name"`
}

// UserMetadata holds users metadata
type UserMetadata struct {
	CompanyId string `json:"company_id" bson:"company_id"`
}

// UserResourcePermissionDto holds metadata and user
type UserResourcePermissionDto struct {
	Metadata  UserMetadata           `json:"metadata" bson:"-"`
	UserId    string                 `json:"user_id" bson:"user_id"`
}

// PasswordResetDto contains data for password reset
type PasswordResetDto struct {
	Otp             string `json:"otp" bson:"otp"`
	Email           string `json:"email" bson:"email"`
	CurrentPassword string `json:"current_password" bson:"current_password"`
	NewPassword     string `json:"new_password" bson:"new_password"`
}

// Config contains config file struct
type Config struct {
	Token 			string `json:"token" bson:"token"`
	ApiServerUrl 	string `json:"api_server_url" bson:"api_server_url"`
	SecurityUrl 	string `json:"security_url" bson:"security_url"`
}

func (cfg Config) Store() error {
	data, err := json.MarshalIndent(cfg, "", "")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(os.Getenv("CONFIG_FILE_PATH"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigFile() Config {
	jsonFile, err := os.Open(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		return Config{}
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configFile Config
	err = json.Unmarshal(byteValue, &configFile)
	if err != nil {
		return Config{}
	}
	return configFile
}

func IsUserLoggedIn() error {
	cfg := GetConfigFile()
	if cfg.Token == "" {
		return errors.New("user is not logged in")
	}
	return nil
}