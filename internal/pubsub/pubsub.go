package pubsub

import (
	"sync"
)

var (
	index       uint
	mu          sync.RWMutex
	subscribers = map[string]map[uint]chan struct{}{}
)

func Subscribe(code Code, id string) (<-chan struct{}, func()) {
	key := code.Key(id)
	ch := make(chan struct{}, 1)

	mu.Lock()
	defer mu.Unlock()
	index++
	uid := index
	_, ok := subscribers[key]
	if !ok {
		subscribers[key] = map[uint]chan struct{}{}
	}
	subscribers[key][index] = ch
	return ch, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(subscribers[key], uid)
		if len(subscribers[key]) == 0 {
			delete(subscribers, key)
		}
	}
}

func Publish(code Code, id string) {
	mu.RLock()
	defer mu.RUnlock()
	if len(subscribers) == 0 {
		return
	}

	if m, ok := subscribers[code.Key(id)]; ok {
		for _, ch := range m {
			ch <- struct{}{}
		}
	}
}
