package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/priykumar/notification-service/model"
	"github.com/priykumar/notification-service/service"
)

type TemplateHandler struct {
	svc service.TemplateSaver
}

// validate request
func validateTemplate(t model.Template) (bool, string) {
	if strings.TrimSpace(t.Name) == "" {
		return false, "Provide a valid name for the template"
	} else if strings.TrimSpace(t.Subject) == "" {
		return false, "Provide a valid subject for the template"
	} else if strings.TrimSpace(t.Message) == "" {
		return false, "Provide a valid body for the template"
	}

	return true, ""
}

func NewTemplateHandler(s service.TemplateSaver) *TemplateHandler {
	return &TemplateHandler{svc: s}
}

func (t *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template model.Template

	// Parse request
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		fmt.Printf("Failed to decode the template body")
		GenerateResponse(w, 500, "Failed to decode the template body")
		return
	}

	// Validate request
	if isValid, msg := validateTemplate(template); !isValid {
		fmt.Println("Failed in template validation. Reason: ", msg)
		GenerateResponse(w, 400, "Failed in template validation. "+msg)
		return
	}

	// Create template
	if err := t.svc.CreateTemplate(template); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			GenerateResponse(w, 400, "Failed in template insertion. Reason: "+err.Error())
			return
		}
		fmt.Println("Failed inserting template into internal DB. Reason: " + err.Error())
	}

	GenerateResponse(w, 200, "Successfully inserted for "+template.Name)
}
