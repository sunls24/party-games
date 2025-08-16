package game

import (
	"encoding/json"
	"party-games/internal/api"
	"party-games/internal/db"
	"strings"
)

type StartGameFunc func(ctx *api.Context, count int, game *db.TGame, conf json.RawMessage) error

var startGameHandler = map[string]StartGameFunc{}

func RegisterStartGame(key string, fn StartGameFunc) {
	startGameHandler[key] = fn
}

func startById(id string) (StartGameFunc, bool) {
	fn, ok := startGameHandler[strings.Split(id, "/")[0]]
	return fn, ok
}
