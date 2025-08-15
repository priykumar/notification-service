package service

import (
	"fmt"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

type TemplateSaver interface {
	CreateTemplate(model.Template) error
}

type templateService struct {
	db datastore.DataStore
}

func NewTemplateService(d datastore.DataStore) TemplateSaver {
	return &templateService{db: d}
}

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
