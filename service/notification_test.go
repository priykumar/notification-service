package service

import (
	"testing"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

// SUCCESS: Test Immediate
func TestSendImmediate(t *testing.T) {
	db := datastore.InitialiseDB()
	svc := NewNotificationService(db)

	for _, ch := range []model.Channel{model.EMAIL, model.INAPP, model.SLACK} {
		n := model.Notification{
			ID:      "test1",
			To:      "to@example.com",
			From:    "from@example.com",
			Message: model.Content{Subject: "Testing", Body: "Hello! I'm performing Unit test"},
			Channel: ch,
		}

		code, err := svc.Send(n)
		if err != nil || code != 200 {
			t.Fatalf("Expected success, got code=%d, err=%v", code, err)
		}
	}

}

// SUCCESS: Test Immediate with template
func TestSendImmediateWithTemplate(t *testing.T) {
	db := datastore.InitialiseDB()
	db.PutTemplate(model.Template{
		Name:    "welcome",
		Subject: "Hello {0}",
		Message: "Welcome {0} to {1}",
	})
	svc := NewNotificationService(db)

	for _, ch := range []model.Channel{model.EMAIL, model.INAPP, model.SLACK} {
		n := model.Notification{
			ID:      "test2",
			To:      "to@example.com",
			From:    "from@example.com",
			Message: model.Content{Subject: "Testing", Body: "Hello! I'm performing Unit test"},
			Channel: ch,
		}

		code, err := svc.Send(n)
		if err != nil || code != 200 {
			t.Fatalf("Expected success, got code=%d, err=%v", code, err)
		}
	}

}

// SUCCESS: Test Immediate
func TestSendWithTemplate(t *testing.T) {
	db := datastore.InitialiseDB()
	db.PutTemplate(model.Template{
		Name:    "welcome",
		Subject: "Hello {0}",
		Message: "Welcome {0} to {1}",
	})

	for _, ch := range []model.Channel{model.EMAIL, model.INAPP, model.SLACK} {
		template := "welcome"
		svc := NewNotificationService(db)
		sec := 5
		n := model.Notification{
			ID:            "test3",
			To:            "to@example.com",
			From:          "from@example.com",
			Template:      &template,
			SendTimeInSec: &sec,
			Message: model.Content{
				SubPlaceHolder:  []string{"Alice"},
				BodyPlaceHolder: []string{"Alice", "World"},
			},
			Channel: ch,
		}

		code, err := svc.Send(n)
		if err != nil || code != 200 {
			t.Fatalf("Expected success, got code=%d, err=%v", code, err)
		}
	}
}

func TestSendWithTemplateFail(t *testing.T) {
	db := datastore.InitialiseDB()
	db.PutTemplate(model.Template{
		Name:    "welcome",
		Subject: "Hello {0}",
		Message: "Welcome {0} to {1}",
	})

	for _, ch := range []model.Channel{model.EMAIL, model.INAPP, model.SLACK} {
		template := "welcome"
		svc := NewNotificationService(db)
		sec := 5
		n := model.Notification{
			ID:            "test4",
			To:            "to@example.com",
			From:          "from@example.com",
			Template:      &template,
			SendTimeInSec: &sec,
			Message: model.Content{
				SubPlaceHolder:  []string{"Alice"},
				BodyPlaceHolder: []string{"Alice"},
			},
			Channel: ch,
		}

		code, err := svc.Send(n)
		if err == nil || code != 400 {
			t.Fatalf("Expected failure, got code=%d, err=%v", code, err)
		}
	}
}
