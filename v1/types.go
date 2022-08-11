package v1

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/v1/encryption"
	"io/ioutil"
	"os"
	"time"
)

// ResponseDTO Http response dto
type ResponseDTO struct {
	Metadata interface{} `json:"_metadata"`
	Data     interface{} `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

// RepositoryDto contains repository dto
type RepositoryDto struct {
	ApiVersion string     `bson:"apiVersion" json:"apiVersion"`
	Kind       string     `bson:"kind" json:"kind"`
	Repository Repository `bson:"repository" json:"repository"`
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
	ApiVersion  string      `bson:"apiVersion" json:"apiVersion"`
	Kind        string      `bson:"kind" json:"kind"`
	Application Application `bson:"application" json:"application"`
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
	IsWebhookEnabled bool              `json:"is_webhook_enabled" yaml:"isWebhookEnabled"`
}

// Applications contains application list
type Applications []struct {
	MetaData ApplicationMetadata `bson:"_metadata" json:"_metadata"`
	Url      string              `bson:"url" json:"url"`
}

// CompanyDto contains company data
type CompanyDto struct {
	ApiVersion string  `bson:"apiVersion" json:"apiVersion"`
	Kind       string  `bson:"kind" json:"kind"`
	Company    Company `bson:"company" json:"company"`
}

// Company contains company data
type Company struct {
	MetaData     CompanyMetadata `bson:"_metadata" json:"_metadata"`
	Id           string          `bson:"id" json:"id"`
	Name         string          `bson:"name" json:"name"`
	Repositories []Repository    `bson:"repositories" json:"repositories"`
	Status       string          `bson:"status" json:"status"`
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
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
}

//ProcessesWithStatus contains process list with status
type ProcessesWithStatus struct {
	ProcessId    string                 `bson:"process_id" json:"process_id"`
	AppId        string                 `bson:"app_id" json:"app_id"`
	RepositoryId string                 `bson:"repository_id" json:"repository_id"`
	Data         map[string]interface{} `bson:"data" json:"data"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	Status       string                 `bson:"status" json:"status"`
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
	Status             string                 `json:"status" bson:"status"`
	CreatedDate        time.Time              `json:"created_date" bson:"created_date"`
	UpdatedDate        time.Time              `json:"updated_date" bson:"updated_date"`
	AuthType           string                 `json:"auth_type" bson:"auth_type"`
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

type PipelineMetadata struct {
	CompanyId       string          `json:"company_id" yaml:"company_id"`
	CompanyMetadata CompanyMetadata `json:"company_metadata" yaml:"company_metadata"`
}

type PipelineApplyOption struct {
	Purging string
}

type Pipeline struct {
	MetaData   PipelineMetadata    `json:"_metadata" yaml:"_metadata"`
	Option     PipelineApplyOption `json:"option" yaml:"option"`
	ApiVersion string              `json:"api_version" yaml:"api_version"`
	Name       string              `json:"name"  yaml:"name"`
	ProcessId  string              `json:"process_id" yaml:"process_id"`
	Label      map[string]string   `json:"label" yaml:"label"`
	Steps      []Step              `json:"steps" yaml:"steps"`
}

type Step struct {
	Name    string            `json:"name" yaml:"name"`
	Type    string            `json:"type" yaml:"type"`
	Status  string            `json:"status" yaml:"status"`
	Trigger string            `json:"trigger" yaml:"trigger"`
	Params  map[string]string `json:"params" yaml:"params"`
	Next    []string          `json:"next" yaml:"next"`
	ArgData map[string]string `json:"arg_data"  yaml:"arg_data"`
	EnvData map[string]string `json:"env_data"  yaml:"env_data"`
}

// UserMetadata holds users metadata
type UserMetadata struct {
	CompanyId string `json:"company_id" bson:"company_id"`
}

// UserResourcePermissionDto holds metadata and user
type UserResourcePermissionDto struct {
	Metadata UserMetadata `json:"metadata" bson:"-"`
	UserId   string       `json:"user_id" bson:"user_id"`
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
	Token        string `json:"token" bson:"token"`
	ApiServerUrl string `json:"api_server_url" bson:"api_server_url"`
	SecurityUrl  string `json:"security_url" bson:"security_url"`
	RepositoryId string `json:"repository_id" bson:"repository_id"`
}

type K8sObjsInfo struct {
	Deployments            []DeploymentShortInfo  `json:"deployments"`
	Services               []K8sObjShortInfo      `json:"services"`
	ConfigMaps             []K8sObjShortInfo      `json:"config_maps"`
	StatefulSets           []StatefulSetShortInfo `json:"stateful_sets"`
	ClusterRoles           []K8sObjShortInfo      `json:"cluster_roles"`
	ClusterRoleBindings    []K8sObjShortInfo      `json:"cluster_role_bindings"`
	DaemonSets             []DaemonSetShortInfo   `json:"daemon_sets"`
	Ingresses              []K8sObjShortInfo      `json:"ingresses"`
	Namespaces             []K8sObjShortInfo      `json:"namespaces"`
	NetworkPolicies        []K8sObjShortInfo      `json:"network_policies"`
	Nodes                  []K8sObjShortInfo      `json:"nodes"`
	PersistentVolumes      []K8sObjShortInfo      `json:"persistent_volumes"`
	PersistentVolumeClaims []K8sObjShortInfo      `json:"persistent_volume_claims"`
	ReplicaSets            []ReplicaSetShortInfo  `json:"replica_sets"`
	Roles                  []K8sObjShortInfo      `json:"roles"`
	RoleBindings           []K8sObjShortInfo      `json:"role_bindings"`
	Secrets                []K8sObjShortInfo      `json:"secrets"`
	ServiceAccounts        []K8sObjShortInfo      `json:"service_accounts"`
	Certificates           []K8sObjShortInfo      `json:"certificates"`
}

type K8sObjShortInfo struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

type DeploymentShortInfo struct {
	Kind                string `json:"kind"`
	Name                string `json:"name"`
	Namespace           string `json:"namespace"`
	UID                 string `json:"uid"`
	Replicas            int32  `json:"replicas"`
	AvailableReplicas   int32  `json:"available_replicas"`
	UnavailableReplicas int32  `json:"unavailable_replicas"`
	ReadyReplicas       int32  `json:"ready_replicas"`
}

type StatefulSetShortInfo struct {
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	UID           string `json:"uid"`
	Replicas      int32  `json:"replicas"`
	ReadyReplicas int32  `json:"ready_replicas"`
}

type ReplicaSetShortInfo struct {
	Kind              string `json:"kind"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	UID               string `json:"uid"`
	Replicas          int32  `json:"replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
}

type DaemonSetShortInfo struct {
	Kind              string `json:"kind"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	UID               string `json:"uid"`
	NumberReady       int32  `json:"number_ready"`
	NumberAvailable   int32  `json:"number_available"`
	NumberUnavailable int32  `json:"number_unavailable"`
}

func (cfg Config) Store() error {
	if cfg.Token != "" {
		aes := encryption.AES256()
		encryptedToken, err := aes.Encrypt(cfg.Token)
		if err != nil {
			return err
		}
		cfg.Token = encryptedToken
	}
	cfg.ApiServerUrl = FixUrl(cfg.ApiServerUrl)
	cfg.SecurityUrl = FixUrl(cfg.SecurityUrl)
	data, err := json.MarshalIndent(cfg, "", "")
	if err != nil {
		return err
	}
	path := GetCfgPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(path+"config.cfg", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigFile() Config {
	jsonFile, err := os.Open(GetCfgPath() + "config.cfg")
	if err != nil {
		return Config{}
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configFile Config
	err = json.Unmarshal(byteValue, &configFile)
	if err != nil {
		return Config{}
	}
	if configFile.Token != "" {
		aes := encryption.AES256()
		decryptedToken, err := aes.Decrypt(configFile.Token)
		if err != nil {
			return Config{}
		}
		configFile.Token = decryptedToken
	}
	return configFile
}
