package service

import (
	"fmt"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

type TemplateService interface {
	CreateTemplate(model.Template) error
}

type templateService struct {
	db datastore.DataStore
}

func NewTemplateService(d datastore.DataStore) TemplateService {
	return &templateService{db: d}
}

// func (t *templateService) GetTemplate(tname string) (*model.Template, error) {
// 	template := t.db.GetTemplate(tname)
// 	if template == nil {
// 		return nil, fmt.Errorf("could not find desired template")
// 	}

// 	return template, nil
// }

func (t *templateService) CreateTemplate(template model.Template) error {

	if t.db.GetTemplate(template.Name) == nil {
		if err := t.db.PutTemplate(template); err != nil {
			return fmt.Errorf("failed inserting template, %v", err)
		}
	} else {
		fmt.Println("Template for", template.Name, "already exist in the DB, hence couldn't make the duplicate entry")
		return fmt.Errorf("duplicate entry")
	}

	fmt.Println("New entry for template ", template.Name, "successfully done")
	return nil
}

// func (t *templateService) ModifyTemplate(template model.Template) error {
// 	if isValid, msg := validateTemplate(template); !isValid {
// 		fmt.Println("Failed in template validation. Reason: ", msg)
// 		return fmt.Errorf("failed in template validation. %s", msg)
// 	}

// 	if err := t.db.ModifyTemplate(template); err != nil {
// 		return fmt.Errorf("failed modifying template, %v", err)
// 	}

// 	return nil
// }
