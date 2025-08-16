package db

import (
	"party-games/pkg/nosql"
)

type TRoom struct {
	nosql.Meta
	Started bool     `json:"started"`
	Owner   string   `json:"owner,omitempty"`
	Users   []string `json:"-"`
}
