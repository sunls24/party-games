package pubsub

import (
	"sync"
)

var (
	idx         uint
	mu          sync.RWMutex
	subscribers = map[string]map[uint]chan any{}
)

func Subscribe(key string) (<-chan any, func()) {
	ch := make(chan any, 1)

	mu.Lock()
	defer mu.Unlock()
	idx++
	id := idx
	_, ok := subscribers[key]
	if !ok {
		subscribers[key] = map[uint]chan any{}
	}
	subscribers[key][id] = ch
	return ch, func() {
		mu.Lock()
		defer mu.Unlock()
		close(subscribers[key][id])
		delete(subscribers[key], id)
		if len(subscribers[key]) == 0 {
			delete(subscribers, key)
		}
	}
}

func Publish(key string, data any) {
	mu.RLock()
	defer mu.RUnlock()
	if len(subscribers) == 0 {
		return
	}
	if m, ok := subscribers[key]; ok {
		for _, ch := range m {
			ch <- data
		}
	}
}
