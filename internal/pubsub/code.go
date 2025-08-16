package pubsub

import "fmt"

type Code int

const (
	RoomUpdate Code = iota
	GameUpdate
)

func (c Code) Key(id string) string {
	return fmt.Sprintf("%s_%d", id, int(c))
}
