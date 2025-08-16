package types

import "party-games/internal/db"

type Room struct {
	*db.TRoom
	Users []RoomUser `json:"users"`
}

type RoomUser struct {
	*db.TUser
	Online bool `json:"online,omitempty"`
}
