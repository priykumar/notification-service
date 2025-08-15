package service

import (
	"container/heap"
	"fmt"
	"regexp"
	"strings"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

const HEAP_SCAN_INTERVAL int = 15
const POP_FROM_HEAP_IF_LESS_THAN int = 20

type Notification interface {
	Send(model.Notification) (int, error)
}

type NotificationService struct {
	db datastore.DataStore
}

func NewNotificationService(d datastore.DataStore) *NotificationService {
	// Initialse heap and start ticket that montiors the heap
	minheap = &NotificationHeap{}
	heap.Init(minheap)

	return &NotificationService{db: d}
}

func (n *NotificationService) SendNotification(ch model.Channel, detail model.Notification) (int, error) {
	var notify Notification
	switch ch {
	case model.EMAIL:
		notify = &Email{db: n.db}
	}

	return notify.Send(detail)
}

func shouldBeSentNow(t *int) bool {
	return t == nil
}

func getPlaceholderCount(str string) int {
	re := regexp.MustCompile(`\{\d+\}`)
	matches := re.FindAllString(str, -1)
	return len(matches)
}

// Populate placeholders in saved template
func populatePlaceholders(template model.Template, content model.Content) (*model.Template, error) {
	phCountSubject := getPlaceholderCount(template.Subject)
	phCountMsg := getPlaceholderCount(template.Message)
	if len(content.SubPlaceHolder) != phCountSubject || (len(content.BodyPlaceHolder) != phCountMsg) {
		return nil, fmt.Errorf("placeholder needed in template deosn't match with provided count")
	}

	// Populate placeholders in subject
	for i, val := range content.SubPlaceHolder {
		ph := fmt.Sprintf("{%d}", i)
		template.Subject = strings.ReplaceAll(template.Subject, ph, val)
	}

	// Populate placeholders in message
	for i, val := range content.BodyPlaceHolder {
		ph := fmt.Sprintf("{%d}", i)
		template.Message = strings.ReplaceAll(template.Message, ph, val)
	}

	return &template, nil
}

func populateChannel(to, from, sub, msg string, ch model.Channel) {
	fmt.Printf(`
MODE: %s
TO: %s
FROM: %s
SUBJECT: %s

Message: %s`,
		ch, to, from, sub, msg)
	fmt.Println()
	fmt.Println()
}
