package room

import (
	"party-games/internal/db"
	"sync"
	"time"
)

var (
	online = map[string]map[string]int{}
	mu     sync.RWMutex
)

func addOnline(id, uid string, room *db.TRoom) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := online[id]; !ok {
		online[id] = make(map[string]int)
	}
	if online[id][uid] == 0 && room != nil {
		for _, u := range room.Users {
			if u == uid {
				db.Room.Publish(id, room)
				break
			}
		}
	}

	online[id][uid]++
}

func delOnline(id, uid string) {
	go func() {
		<-time.After(time.Second * 3)
		mu.Lock()
		defer mu.Unlock()
		if _, ok := online[id]; !ok {
			return
		}
		online[id][uid]--
		if online[id][uid] <= 0 {
			delete(online[id], uid)
		}
	}()
}

func isOnline(id, uid string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := online[id]; !ok {
		return false
	}
	c, ok := online[id][uid]
	return ok && c > 0
}
