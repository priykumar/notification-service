package service

import (
	"testing"

	"github.com/priykumar/notification-service/model"
)

func TestShouldBeSentNow(t *testing.T) {
	if !shouldBeSentNow(nil) {
		t.Error("Expected true for nil time")
	}
	val := 5
	if shouldBeSentNow(&val) {
		t.Error("Expected false for non-nil time")
	}
}

func TestGetPlaceholderCount(t *testing.T) {
	count := getPlaceholderCount("Hello {0} and {1}")
	if count != 2 {
		t.Errorf("Expected 2 placeholders, got %d", count)
	}
}

func TestPopulatePlaceholders(t *testing.T) {
	template := model.Template{
		Subject: "Hi {0}",
		Message: "Welcome {0} to {1}",
	}
	content := model.Content{
		SubPlaceHolder:  []string{"John"},
		BodyPlaceHolder: []string{"John", "Go"},
	}

	tmpl, err := populatePlaceholders(template, content)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if tmpl.Subject != "Hi John" || tmpl.Message != "Welcome John to Go" {
		t.Errorf("Placeholders not populated correctly: %+v", tmpl)
	}
}

func TestPopulatePlaceholdersMismatch(t *testing.T) {
	template := model.Template{Subject: "Hi {0}", Message: "Hello {0}"}
	content := model.Content{SubPlaceHolder: []string{}, BodyPlaceHolder: []string{"John"}}

	_, err := populatePlaceholders(template, content)
	if err == nil {
		t.Error("Expected error for mismatched placeholders")
	}
}
