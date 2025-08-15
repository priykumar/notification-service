package service

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/priykumar/notification-service/model"
)

type req_time struct {
	id string
	t  time.Time
}

var heaplock sync.Mutex

// NotificationHeap is a min-heap based on SendTime
type NotificationHeap []req_time

var minheap *NotificationHeap

func (h NotificationHeap) Len() int { return len(h) }

func (h NotificationHeap) Less(i, j int) bool { return h[i].t.Before(h[j].t) }

func (h NotificationHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *NotificationHeap) Push(x interface{}) { *h = append(*h, x.(req_time)) }

func (h *NotificationHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Populate heap with notification that needs to send in future
func populateHeap(ncation model.Notification) {
	rt := req_time{
		id: ncation.ID,
		t:  time.Now().Add(time.Duration(*ncation.SendTimeInSec) * time.Second),
	}

	fmt.Println("Pushing into heap, Acquiring lock...")
	heaplock.Lock()
	defer fmt.Println("Pushed into heap, Releasing lock...")
	defer heaplock.Unlock()

	heap.Push(minheap, rt)
}

// 1. Scan heap every HEAP_SCAN_INTERVAL
// 2. Keeping popping element from heap and create ticker those element if time for element
// to pushed to channel is less than POP_FROM_HEAP_IF_LESS_THAN
func (not *NotificationService) MonitorAndPop() {
	ticker := time.NewTicker(time.Duration(HEAP_SCAN_INTERVAL) * time.Second)
	defer ticker.Stop()

	for {
		fmt.Println("Scanning heap, Acquiring lock...")
		heaplock.Lock()
		for minheap.Len() > 0 {
			earliest := (*minheap)[0]
			if earliest.t.After(time.Now().Add(time.Duration(POP_FROM_HEAP_IF_LESS_THAN) * time.Second)) {
				// stop if earliest is beyond POP_FROM_HEAP_IF_LESS_THAN
				break
			}

			// Pop the earliest
			n := heap.Pop(minheap).(req_time)
			diff := int(time.Until(earliest.t).Seconds())
			ncation := not.db.GetNotification(n.id)
			ncation.SendTimeInSec = &diff
			go startTicker(&ncation)
		}
		fmt.Println("Scanned heap, Releasing lock...")
		heaplock.Unlock()
		<-ticker.C
	}
}

// StartTicker waits until SendTime and then populates the channel
func startTicker(n *model.Notification) {
	fmt.Println("Starting ticker for", *n.SendTimeInSec, "seconds")
	duration := *n.SendTimeInSec
	if duration <= 0 {
		populateChannel(n.To, n.From, n.Message.Subject, n.Message.Body, n.Channel)
		return
	}

	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	defer ticker.Stop()

	<-ticker.C
	populateChannel(n.To, n.From, n.Message.Subject, n.Message.Body, n.Channel)
}
