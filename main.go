// Package collider provides a lock-free
// ring buffer for arbitrary Go objects.
package collider

import (
	"sync/atomic"
)

// Ring implements a false-sharing
// resistant ring buffer with slots
// that hold arbitrary Go objects.
type Ring struct {
	p0 [64]byte
	pos uint64
	p1 [64]byte
	mask *uint64
	p2 [64]byte
	slots []interface{}
}

// New takes a size n and returns
// a ring buffer with n slots.
func New(s int) *Ring {
	m := uint64(s-1)
	ring := &Ring{
		slots: make([]interface{}, s),
		mask: &m,
	}

	return ring
}

// Get returns the object at the current
// slot and atomically increments the index. 
func (r *Ring) Get() interface{} {
	return r.slots[(atomic.AddUint64(&r.pos, 1)-1)&(*r.mask)]
}

// Add adds an item to the next open slot.
func (r *Ring) Add(u interface {}) {
	r.slots[(atomic.AddUint64(&r.pos, 1)-1)&(*r.mask)] = u
}