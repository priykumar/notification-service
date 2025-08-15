package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/service"
)

func TestCreateNotification(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := service.NewNotificationService(db)
	h := NewNotificationHandler(svc)

	body := `{
		"to": "to@example.com",
		"from": "from@example.com",
		"message": { "subject": "Hi", "body": "Hello" },
		"channel": "email"
	}`

	req := httptest.NewRequest(http.MethodPost, "/producer/notify", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateNotification(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateNotificationInvalidSender(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := service.NewNotificationService(db)
	h := NewNotificationHandler(svc)

	body := []string{`{ "to": "to@abc.com", "from": "", "message": {}, "channel": "" }`,
		`{ "to": "", "from": "from@abc.com", "message": {}, "channel": "" }`,
		`{ "to": "to@abc.com", "from": "from@abc.com", "message": {}, "channel": "" }`}

	for _, b := range body {
		req := httptest.NewRequest(http.MethodPost, "/producer/notify", bytes.NewBufferString(b))
		w := httptest.NewRecorder()

		h.CreateNotification(w, req)

		if w.Code == 200 {
			t.Fatal("Expected failure for invalid notification", w.Code)
		}
	}
}
