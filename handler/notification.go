package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/priykumar/notification-service/model"
	"github.com/priykumar/notification-service/service"
)

type NotificationHandler struct {
	svc service.NotificationSender
}

func NewNotificationHandler(s service.NotificationSender) *NotificationHandler {
	return &NotificationHandler{svc: s}
}

// Validates provided channel
func validateChannel(ch model.Channel) bool {
	if ch == "" {
		return false
	} else if ch != model.EMAIL && ch != model.INAPP && ch != model.SLACK {
		return false
	}

	return true
}

// Validates request
func validateNotification(ncation model.Notification) (bool, string) {
	if strings.TrimSpace(ncation.To) == "" {
		return false, "Provide a valid receiver address"
	} else if strings.TrimSpace(ncation.From) == "" {
		return false, "Provide a valid sender address"
	} else if !validateChannel(ncation.Channel) {
		return false, "Provide a valid channel [email, inapp, slack]"
	} else if ncation.Template == nil && (ncation.Message.Body == "" || ncation.Message.Subject == "") {
		return false, "Provide a valid template or provide valid subject and body to be sent"
	}

	return true, ""
}

func (n *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var ncation model.Notification

	// Parse request
	if err := json.NewDecoder(r.Body).Decode(&ncation); err != nil {
		fmt.Printf("Failed to decode the notification body")
		GenerateResponse(w, 500, "Failed to decode the notification body")
		return
	}

	// Validate request
	if isValid, msg := validateNotification(ncation); !isValid {
		fmt.Println("Failed in notification validation. Reason: ", msg)
		GenerateResponse(w, 400, "Failed in notification validation. "+msg)
		return
	}

	// Create unique for each notification and send notification to channel
	ncation.ID = uuid.New().String()
	if code, err := n.svc.Send(ncation); err != nil {
		fmt.Println("Failed to generate notification. Reason: ", err.Error())
		GenerateResponse(w, code, "Failed in notification validation. "+err.Error())
		return
	}

	GenerateResponse(w, 200, "")
}
