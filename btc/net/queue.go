package net

import (
	"fmt"
)

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		hashes: make([][32]byte, size),
		invs: 	make(map[[32]byte]*inv),
		size:  	size,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	hashes 	[][32]byte
	invs		map[[32]byte]*inv
	size  	int
	head  	int
	tail  	int
	count 	int
}

func (q *Queue) Exists(hash [32]byte) bool {
	_, exists := q.invs[hash]
	return exists
}

// Push adds a node to the queue.
func (q *Queue) Push(hash [32]byte, inventory *inv) error {
	// TODO: Check if already exists
	if _, exists := q.invs[hash]; exists {
		return fmt.Errorf("Hash already exists")
	}
	if q.head == q.tail && q.count > 0 {
		if err := q.Pop(); err != nil {
			return err
		}
	}
	q.hashes[q.tail] = hash
	q.invs[hash] = inventory
	q.tail = (q.tail + 1) % len(q.hashes)
	q.count++
	return nil
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() error {
	if q.count == 0 {
		return fmt.Errorf("Queue is empty")
	}
	hash := q.hashes[q.head]
	delete(q.invs, hash)
	q.head = (q.head + 1) % len(q.hashes)
	q.count--
	return nil
}
