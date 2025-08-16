package db

import (
	"party-games/pkg/nosql"
)

type TGame struct {
	nosql.Meta
	Conf  any `json:"conf"`
	State any `json:"state"`
}
