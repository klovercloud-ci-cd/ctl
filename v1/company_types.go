package v1
//
//import (
//	"errors"
//	"github.com/klovercloud-ci/ctl/enums"
//	"reflect"
//)
//
//// Company contains company data
//type Company struct {
//	MetaData     CompanyMetadata      `bson:"_metadata" json:"_metadata"`
//	Id           string               `bson:"id" json:"id"`
//	Name         string               `bson:"name" json:"name"`
//	Repositories []Repository         `bson:"repositories" json:"repositories"`
//	Status       enums.COMPANY_STATUS `bson:"status" json:"status"`
//}
//
//// Validate validates company data
//func (dto Company) Validate() error {
//	err := dto.MetaData.Validate()
//	if err != nil {
//		return err
//	}
//	if dto.Id == "" {
//		return errors.New("Company id is required!")
//	}
//	if dto.Name == "" {
//		return errors.New("Company name is required!")
//	}
//	for _, each := range dto.Repositories {
//		err := each.Validate()
//		if err != nil {
//			return err
//		}
//	}
//	if dto.Status == enums.ACTIVE || dto.Status == enums.INACTIVE {
//		return nil
//	} else if dto.Status == "" {
//		return errors.New("Company status is required!")
//	}
//	return errors.New("Company status invalid!")
//}
//
//// Repository contains repository info
//type Repository struct {
//	Id           string                `bson:"id" json:"id"`
//	Type         enums.REPOSITORY_TYPE `bson:"type" json:"type"`
//	Token        string                `bson:"token" json:"token"`
//	Applications []Application         `bson:"applications" json:"applications"`
//}
//
//// Validate validates repository info
//func (repository Repository) Validate() error {
//	if repository.Id == "" {
//		return errors.New("Repository id is required!")
//	}
//	if repository.Token == "" {
//		return errors.New("Repository token is required!")
//	}
//	for _, each := range repository.Applications {
//		err := each.Validate()
//		if err != nil {
//			return err
//		}
//	}
//	if repository.Type == enums.GITHUB || repository.Type == enums.BIT_BUCKET {
//		return nil
//	} else if repository.Type == "" {
//		return errors.New("Repository type is required")
//	}
//	return errors.New("Repository type is invalid!")
//
//}
//
//// Application contains application info
//type Application struct {
//	MetaData ApplicationMetadata  `bson:"_metadata" json:"_metadata"`
//	Url      string               `bson:"url" json:"url"`
//	Webhook  GitWebhook           `bson:"webhook" json:"webhook"`
//	Status   enums.COMPANY_STATUS `bson:"status" json:"status"`
//}
//
//// Validate validates application info
//func (application Application) Validate() error {
//	if application.Url == "" {
//		return errors.New("Application url is required!")
//	}
//	err := application.MetaData.Validate()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// ApplicationMetadata contains application metadata info
//type ApplicationMetadata struct {
//	Labels           map[string]string `bson:"labels" json:"labels"`
//	Id               string            `bson:"id" json:"id"`
//	Name             string            `bson:"name" json:"name"`
//	IsWebhookEnabled bool              `bson:"is_webhook_enabled" json:"is_webhook_enabled"`
//}
//
//// Validate validates application metadata
//func (metadata ApplicationMetadata) Validate() error {
//	keys := reflect.ValueOf(metadata.Labels).MapKeys()
//	for i := 0; i < len(keys); i++ {
//		if metadata.Labels[keys[i].String()] == "" {
//			return errors.New("Application metadata label is missing!")
//		}
//	}
//	if metadata.Id == "" {
//		return errors.New("Application metadata id is required!")
//	}
//	if metadata.Name == "" {
//		return errors.New("Application metadata name is required!")
//	}
//	return nil
//}
//
//// CompanyMetadata contains company metadata info
//type CompanyMetadata struct {
//	Labels                    map[string]string `bson:"labels" json:"labels" yaml:"labels"`
//	NumberOfConcurrentProcess int64             `bson:"number_of_concurrent_process" json:"number_of_concurrent_process" yaml:"number_of_concurrent_process"`
//	TotalProcessPerDay        int64             `bson:"total_process_per_day" json:"total_process_per_day" yaml:"total_process_per_day"`
//}
//
//// Validate validates company metadata
//func (metadata CompanyMetadata) Validate() error {
//	keys := reflect.ValueOf(metadata.Labels).MapKeys()
//	for i := 0; i < len(keys); i++ {
//		if metadata.Labels[keys[i].String()] == "" {
//			return errors.New("Company metadata label is missing!")
//		}
//	}
//	return nil
//}
//
//
