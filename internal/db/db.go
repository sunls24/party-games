package db

import (
	"party-games/internal/utils"
	"party-games/pkg/nosql"
)

var (
	User nosql.Table[*TUser]
	Room nosql.Table[*TRoom]
	Game nosql.Table[*TGame]
)

func MustInit(path string) {
	db, err := nosql.NewDB(path)
	utils.Must(err)
	User, err = nosql.NewTable[*TUser](db)
	utils.Must(err)
	Room, err = nosql.NewTable[*TRoom](db)
	utils.Must(err)
	Game, err = nosql.NewTable[*TGame](db)
	utils.Must(err)
}
