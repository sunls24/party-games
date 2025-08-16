package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"party-games/internal/api"
	"party-games/internal/db"
	"time"
)

func Long(ctx *api.Context, id string, version int) error {
	game, err := db.Game.Get(id)
	if err != nil {
		return ctx.Err(err)
	}
	if game == nil {
		return ctx.BadMsg()
	}
	if game.Version > version {
		return ctx.Data(game)
	}
	ch, cancel := db.Game.Subscribe(id)
	defer cancel()
	select {
	case <-ctx.Request().Context().Done():
		return nil
	case <-time.After(time.Minute):
		return ctx.NoContent(http.StatusNoContent)
	case game = <-ch:
		return ctx.Data(game)
	}
}

func checkRoom(id string) (*db.TRoom, error) {
	room, err := db.Room.Get(id)
	if room == nil {
		return nil, err
	}
	for i, uid := range room.Users {
		if uid == "" {
			return nil, api.NewMsg(fmt.Sprintf("位置%d没有玩家", i+1))
		}
	}
	return room, nil
}

func Start(ctx *api.Context) error {
	id := ctx.QueryParam("id")
	var conf json.RawMessage
	if err := ctx.Bind(&conf); err != nil || id == "" {
		return ctx.BadParam()
	}
	startGame, ok := startById(id)
	if !ok {
		return ctx.BadMsg()
	}

	room, err := checkRoom(id)
	if err != nil {
		return err
	}
	if room == nil {
		return ctx.BadMsg()
	}

	game, err := db.Game.Get(id)
	if err != nil {
		return ctx.Err(err)
	}
	if game == nil {
		game = &db.TGame{}
		err = startGame(ctx, len(room.Users), game, conf)
		if err != nil {
			return ctx.Err(err)
		}
		err = db.Game.Create(id, game)
	} else {
		err = db.Game.Update(id, func(tg *db.TGame) error {
			game = tg
			return startGame(ctx, len(room.Users), game, conf)
		})
	}
	if err != nil {
		return ctx.Err(err)
	}

	err = db.Room.UpdatePublish(id, func(room *db.TRoom) error {
		room.Started = true
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(game)
}

func Stop(ctx *api.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.BadParam()
	}
	err := db.Room.UpdatePublish(id, func(room *db.TRoom) error {
		room.Started = false
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	err = db.Game.Delete(id)
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func Restart(ctx *api.Context) error {
	id := ctx.QueryParam("id")
	startGame, ok := startById(id)
	if !ok {
		return ctx.BadMsg()
	}
	room, err := checkRoom(id)
	if err != nil {
		return err
	}
	if room == nil {
		return ctx.BadMsg()
	}
	err = db.Game.UpdatePublish(id, func(game *db.TGame) error {
		return startGame(ctx, len(room.Users), game, nil)
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func UpdateState[T any](id string, update func(state T) error) error {
	return db.Game.UpdatePublish(id, func(game *db.TGame) error {
		state, ok := game.State.(T)
		if !ok {
			return errors.New("update game state type mismatch")
		}
		err := update(state)
		if err != nil {
			return err
		}
		return nil
	})
}
