package v1

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