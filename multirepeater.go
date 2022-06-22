package repeater

import (
	"sync"
	"time"
)

// Holds any amount of repeaters. Uses StartRepeater. Use this to manage several repeaters by ids. Use NewMultiRepeater to instance MultiRepeater.
type MultiRepeater[ID comparable] struct {
	repeaters map[ID]func()
	mux       sync.Mutex
}

// Instances MultiRepeater. Pass ID type you'd like to use to manage repeaters.
func NewMultiRepeater[ID comparable]() *MultiRepeater[ID] {
	return &MultiRepeater[ID]{
		repeaters: make(map[ID]func()),
		mux:       sync.Mutex{},
	}
}

// Works almost same as original StartRepeater, but delegates stop function managment to MultiRepeater and requires id for repeater.
// Returns `true` on success, and `false` if passed id is occupied.
func (mr *MultiRepeater[ID]) StartRepeater(id ID, frequency time.Duration, fnToCall func()) bool {
	mr.mux.Lock()
	if mr.repeaters[id] != nil {
		mr.mux.Unlock()
		return false
	}
	mr.repeaters[id] = StartRepeater(frequency, fnToCall)
	mr.mux.Unlock()
	return true
}

// Stops repeater by ID. Return `true` on success, and `false` if passed id is not found.
func (mr *MultiRepeater[ID]) StopRepeater(id ID) bool {
	mr.mux.Lock()
	if mr.repeaters[id] == nil {
		mr.mux.Unlock()
		return false
	}
	mr.repeaters[id]()
	delete(mr.repeaters, id)
	mr.mux.Unlock()
	return true
}

// Stops all repeaters managed by this MultiRepeater. Waits until all repeaters are stopped.
func (mr *MultiRepeater[ID]) StopAllRepeaters() {
	wg := &sync.WaitGroup{}

	mr.mux.Lock()
	wg.Add(len(mr.repeaters))
	for _, stop := range mr.repeaters {
		go func(stop func()) {
			stop()
			wg.Done()
		}(stop)
	}
	mr.repeaters = make(map[ID]func())
	mr.mux.Unlock()

	wg.Wait()
}
