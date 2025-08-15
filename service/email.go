package service

import (
	"fmt"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

type Email struct {
	db datastore.DataStore
}

//  1. Check template validity if template is provide in request
//  2. Check if placeholder count expected by template is same as placeholder count provided in request
//  3. If no template, then simply post subject and message provided in request
//  4. Check is channel needs to be populated now
//  5. If not, then check after long how channel needs to be populated
//     a. If population time < POP_FROM_HEAP_IF_LESS_THAN, spawn a ticker for that notification
//     b. Else, push notification to min-heap
func (e *Email) Send(nDetail model.Notification) (int, error) {
	template := &model.Template{}
	var err error
	// If specific template needs to be used then create the template
	if nDetail.Template != nil {
		// check if mentioned template is correct
		template = e.db.GetTemplate(*nDetail.Template)
		if template == nil {
			fmt.Println("Provide a valid template name. No template named", *nDetail.Template, "found")
			return 400, fmt.Errorf("invalid template")
		}

		// Populate placeholders in template
		template, err = populatePlaceholders(*template, nDetail.Message)
		if err != nil {
			fmt.Println("Failed populating placeholders in template. Reason:", err)
			return 400, err
		}
	} else {
		template.Subject = nDetail.Message.Subject
		template.Message = nDetail.Message.Body
	}

	if shouldBeSentNow(nDetail.SendTimeInSec) {
		// Check if end user needs to be notified now
		populateChannel(nDetail.To, nDetail.From, template.Subject, template.Message, model.EMAIL)
	} else {
		nDetail.Message.Subject = template.Subject
		nDetail.Message.Body = template.Message
		if *nDetail.SendTimeInSec < POP_FROM_HEAP_IF_LESS_THAN {
			fmt.Println("Notification needs to sent in ", *nDetail.SendTimeInSec, "(<", POP_FROM_HEAP_IF_LESS_THAN, ") seconds, hence starting a ticker")
			go startTicker(&nDetail)
		} else {
			fmt.Println("Notification needs to sent in ", *nDetail.SendTimeInSec, "(>", POP_FROM_HEAP_IF_LESS_THAN, ") seconds, hence pushing in heap")
			e.db.PutNotification(nDetail)
			go populateHeap(nDetail)
		}
	}

	return 200, nil
}
