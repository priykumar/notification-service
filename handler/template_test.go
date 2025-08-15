package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/service"
)

func TestCreateTemplate(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := service.NewTemplateService(db)
	h := NewTemplateHandler(svc)

	body := `{"name": "welcome", "subject": "Hi {0}", "message": "Hello {0}"}`
	req := httptest.NewRequest(http.MethodPost, "/producer/template", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateTemplate(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
}

func TestCreateTemplateInvalid(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := service.NewTemplateService(db)
	h := NewTemplateHandler(svc)

	body := `{"name": "", "subject": "", "message": ""}`
	req := httptest.NewRequest(http.MethodPost, "/producer/template", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateTemplate(w, req)

	if w.Code == 200 {
		t.Fatal("Expected error for invalid template")
	}
}
