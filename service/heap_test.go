package service

import (
	"container/heap"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/model"
)

func TestPopulateHeap(t *testing.T) {
	// Reset heap
	minheap = &NotificationHeap{}
	heap.Init(minheap)

	sendTime := 10 * HEAP_SCAN_INTERVAL
	notification := model.Notification{
		ID:            "test-heap-1",
		To:            "to@example.com",
		From:          "from@example.com",
		SendTimeInSec: &sendTime,
		Message:       model.Content{Subject: "Test", Body: "Test body"},
		Channel:       model.EMAIL,
	}

	initialLen := minheap.Len()
	populateHeap(notification)

	if minheap.Len() != initialLen+1 {
		t.Errorf("Expected heap length %d, got %d", initialLen+1, minheap.Len())
	}

	// Verify the notification was added correctly
	top := (*minheap)[0]
	if top.id != "test-heap-1" {
		t.Errorf("Expected top id 'test-heap-1', got '%s'", top.id)
	}
}

func TestPopulateHeap_Concurrent(t *testing.T) {
	// Reset heap for test
	minheap = &NotificationHeap{}
	heap.Init(minheap)

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent heap population
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sendTime := 10 * HEAP_SCAN_INTERVAL
			notification := model.Notification{
				ID:            fmt.Sprintf("concurrent-test-%d", id),
				To:            "to@example.com",
				From:          "from@example.com",
				SendTimeInSec: &sendTime,
				Message:       model.Content{Subject: "Test", Body: "Test body"},
				Channel:       model.EMAIL,
			}
			populateHeap(notification)
		}(i)
	}

	wg.Wait()

	if minheap.Len() != numGoroutines {
		t.Errorf("Expected heap length %d, got %d", numGoroutines, minheap.Len())
	}
}

func Test_monitorAndPop(t *testing.T) {
	// Initialise db and heap
	db := datastore.InitialiseDB()
	n := NewNotificationService(db)

	// populate heap
	now := time.Now()
	dummy_rt1 := req_time{"monitorAndPop-test1", now.Add(time.Duration(POP_FROM_HEAP_IF_LESS_THAN) * time.Second)}
	heap.Push(minheap, dummy_rt1)
	dummy_rt2 := req_time{"monitorAndPop-test1", now.Add(time.Duration(HEAP_SCAN_INTERVAL+POP_FROM_HEAP_IF_LESS_THAN) * time.Second)}
	heap.Push(minheap, dummy_rt2)

	// check heap size
	go n.(*NotificationService).monitorAndPop()
	time.Sleep(2 * time.Second)
	if minheap.Len() != 1 {
		t.Errorf("Expected heap size to be 1, got %d", minheap.Len())
	}

	time.Sleep(time.Duration(HEAP_SCAN_INTERVAL+2) * time.Second)
	if minheap.Len() != 0 {
		t.Errorf("Expected heap size to be 0, got %d", minheap.Len())
	}
}
