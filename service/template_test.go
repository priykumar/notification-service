package service

import (
	"strings"
	"testing"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

func TestCreateTemplate_Success(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := NewTemplateService(db)

	template := model.Template{
		Name:    "welcome",
		Subject: "Hi {0}",
		Message: "Hello {0}",
	}

	err := svc.CreateTemplate(template)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify it was inserted into datastore
	got := db.GetTemplate("welcome")
	if got == nil {
		t.Fatal("Expected template to be stored in datastore")
	}
	if got.Name != "welcome" {
		t.Errorf("Expected template name 'welcome', got %s", got.Name)
	}
}

func TestCreateTemplate_Duplicate(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := NewTemplateService(db)

	template := model.Template{
		Name:    "welcome",
		Subject: "Hi {0}",
		Message: "Hello {0}",
	}

	// First insert should succeed
	if err := svc.CreateTemplate(template); err != nil {
		t.Fatalf("Unexpected error on first insert: %v", err)
	}

	// Second insert should fail
	err := svc.CreateTemplate(template)
	if err == nil {
		t.Fatal("Expected error for duplicate template, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("Expected 'duplicate' error, got %v", err)
	}
}
