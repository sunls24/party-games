package db

import "party-games/pkg/nosql"

type TUser struct {
	nosql.Meta
	Name string `json:"name,omitempty"`
	Icon int    `json:"icon,omitempty"`
}
