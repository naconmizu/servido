package main

import (
	"fmt"
	"strings"
	"sync"
)

// Queueable interface defines types that can be queued
type Queueable interface {
	String() string
}

type node struct {
	value string
	next  *node
}

type Queue struct {
	head *node
	tail *node
	size uint
}

// // Enqueue adds an element to the end of the queue
// func (q *Queue) Enqueue(val string) {
// 	newNode := &node{value: val}

// 	if q.tail != nil {
// 		q.tail.next = newNode
// 	}
// 	q.tail = newNode

// 	if q.head == nil {
// 		q.head = newNode
// 	}
// 	q.size++
// }

// // EnqueueMultiple adds multiple elements to the queue
// func (q *Queue) EnqueueMultiple(values ...string) {
// 	for _, val := range values {
// 		q.Enqueue(val)
// 	}
// }

// // Dequeue removes and returns the first element's value
// func (q *Queue) Dequeue() (string, bool) {
// 	if q.head == nil {
// 		return "", false
// 	}

// 	removedNode := q.head
// 	q.head = q.head.next

// 	if q.head == nil {
// 		q.tail = nil
// 	}
// 	q.size--

// 	return removedNode.value, true
// }

// // Peek returns the first element without removing it
// func (q *Queue) Peek() (string, bool) {
// 	if q.head == nil {
// 		return "", false
// 	}
// 	return q.head.value, true
// }

// // PeekLast returns the last element without removing it
// func (q *Queue) PeekLast() (string, bool) {
// 	if q.tail == nil {
// 		return "", false
// 	}
// 	return q.tail.value, true
// }

// // IsEmpty checks if the queue is empty
// func (q *Queue) IsEmpty() bool {
// 	return q.size == 0
// }

// // Size returns the number of elements in the queue
// func (q *Queue) Size() uint {
// 	return q.size
// }

// // Clear removes all elements from the queue
// func (q *Queue) Clear() {
// 	q.head = nil
// 	q.tail = nil
// 	q.size = 0
// }

// // Contains checks if a value exists in the queue
// func (q *Queue) Contains(val string) bool {
// 	current := q.head
// 	for current != nil {
// 		if current.value == val {
// 			return true
// 		}
// 		current = current.next
// 	}
// 	return false
// }

// // ToSlice converts the queue to a slice
// func (q *Queue) ToSlice() []string {
// 	result := make([]string, 0, q.size)
// 	current := q.head
// 	for current != nil {
// 		result = append(result, current.value)
// 		current = current.next
// 	}
// 	return result
// }

// // PrintQueue prints the queue elements with size
// func (q *Queue) PrintQueue() {
// 	if q.IsEmpty() {
// 		fmt.Println("Queue is empty")
// 		return
// 	}

// 	fmt.Printf("Queue (size: %d): ", q.size)
// 	current := q.head
// 	for current != nil {
// 		fmt.Print(current.value)
// 		if current.next != nil {
// 			fmt.Print(" -> ")
// 		}
// 		current = current.next
// 	}
// 	fmt.Println()

// }

// // FromSlice creates a new queue from a slice of strings
// func FromSlice(slice []string) *Queue {
// 	q := NewQueue()
// 	for _, val := range slice {
// 		q.Enqueue(val)
// 	}
// 	return q
// }

//----------------------------------------------------------Claude3.7Queue-----------------------------------------------------------------------

// FileQueue represents a thread-safe queue of filenames
type FileQueue struct {
	items []string
	mutex sync.RWMutex
}

// NewQueue creates a new empty queue
func NewQueue() *FileQueue {
	return &FileQueue{
		items: make([]string, 0),
	}
}

// Enqueue adds a file to the queue
func (q *FileQueue) Enqueue(filename string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = append(q.items, filename)
}

// Dequeue removes and returns the first file in the queue
func (q *FileQueue) Dequeue() (string, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// PeekLast returns the last item without removing it
func (q *FileQueue) PeekLast() (string, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.items) == 0 {
		return "", false
	}

	return q.items[len(q.items)-1], true
}

// IsEmpty checks if the queue is empty
func (q *FileQueue) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.items) == 0
}

// ToSlice returns a copy of the queue as a slice
func (q *FileQueue) ToSlice() []string {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	result := make([]string, len(q.items))
	copy(result, q.items)
	return result
}

// PrintQueue logs the current queue contents
func (q *FileQueue) PrintQueue() {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.items) == 0 {
		fmt.Println("Queue is empty")
		return
	}

	fmt.Printf("Queue contents: %s\n", strings.Join(q.items, ", "))
}



// FromSliceToQueue creates a new FileQueue from a slice of strings.
// It initializes the queue with the provided items.
func FromSliceToQueue(files []string) *FileQueue {
    queue := &FileQueue{
        items: make([]string, len(files)),
        mutex: sync.RWMutex{},
    }
    
    // Acquire lock before modifying the queue
    queue.mutex.Lock()
    defer queue.mutex.Unlock()
    
    // Copy all elements from the input slice
    copy(queue.items, files)
    
    return queue
}


/*
func main() {
    q := NewQueue()

    // Test basic operations
    q.EnqueueMultiple("A", "B", "C")
    q.PrintQueue()

    if val, ok := q.Peek(); ok {
        fmt.Printf("Front element: %s\n", val)
    }

    if val, ok := q.Dequeue(); ok {
        fmt.Printf("Dequeued: %s\n", val)
    }

    q.PrintQueue()
    fmt.Printf("Contains 'B': %t\n", q.Contains("B"))
    fmt.Printf("Contains 'Z': %t\n", q.Contains("Z"))

    slice := q.ToSlice()
    fmt.Printf("As slice: %v\n", slice)

    q.Clear()
    q.PrintQueue()
}
*/
