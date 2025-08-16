package room

import (
	"party-games/internal/pubsub"
	"sync"
	"time"
)

var (
	online = map[string]map[string]chan struct{}{}
	mu     sync.RWMutex
)

func addOnline(roomId, userId string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := online[roomId]; !ok {
		online[roomId] = make(map[string]chan struct{})
	}
	ch, ok := online[roomId][userId]
	if !ok {
		online[roomId][userId] = nil
		pubsub.Publish(pubsub.RoomUpdate, roomId)
		return
	}
	if ch == nil {
		return
	}
	ch <- struct{}{}
	online[roomId][userId] = nil
}

func delOnline(roomId, userId string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := online[roomId]; !ok {
		return
	}
	ch, ok := online[roomId][userId]
	if !ok {
		return
	}
	if ch != nil {
		ch <- struct{}{}
	}
	ch = make(chan struct{}, 1)
	online[roomId][userId] = ch
	go func() {
		defer close(ch)
		select {
		case <-ch:
			return
		case <-time.After(time.Second * 5):
			mu.Lock()
			defer mu.Unlock()
			delete(online[roomId], userId)
			pubsub.Publish(pubsub.RoomUpdate, roomId)
		}
	}()
}

func isOnline(roomId, userId string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := online[roomId]; !ok {
		return false
	}
	_, ok := online[roomId][userId]
	return ok
}
