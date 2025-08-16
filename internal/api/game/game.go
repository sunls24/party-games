package game

import (
	"net/http"
	"party-games/internal/api"
	"party-games/internal/api/room"
	"party-games/internal/db"
	"party-games/internal/pubsub"
	"strconv"
	"time"
)

func CreateOrGet(ctx *api.Context) error {
	id, rt := room.GetRoomId(ctx)
	if rt == "" || id == "" {
		return ctx.BadParam()
	}
	exist, err := db.Game.Exist(id)
	if err != nil {
		return ctx.Send(err)
	}
	if !exist {
		err = db.Game.Create(&db.TGame{
			Base:    db.Base{Id: id},
			Started: false,
		})
		if err != nil {
			return ctx.Send(err)
		}
	}

	game, err := db.Game.Get(id)
	if err != nil {
		return ctx.Send(err)
	}
	game.R()
	return ctx.Data(game)
}

func LongQuery(ctx *api.Context) error {
	id, rt := room.GetRoomId(ctx)
	if rt == "" || id == "" {
		return ctx.BadParam()
	}

	versionStr := ctx.QueryParam("version")
	if id == "" || versionStr == "" {
		return ctx.BadParam()
	}
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return ctx.BadParam()
	}

	game, err := db.Game.Get(id)
	if err != nil {
		return ctx.ErrMsg(err, "查询房间异常")
	}
	if game.Version > version {
		game.R()
		return ctx.Data(game)
	}

	ch, cancel := pubsub.Subscribe(pubsub.GameUpdate, id)
	defer cancel()
	select {
	case <-ctx.Request().Context().Done():
		return nil
	case <-time.After(time.Minute):
		return ctx.NoContent(http.StatusNoContent)
	case <-ch:
		game, err = db.Game.Get(id)
		if err != nil {
			return ctx.ErrMsg(err, "查询游戏信息异常")
		}

		game.R()
		return ctx.Data(game)
	}
}

func Start(ctx *api.Context) error {
	id, rt := room.GetRoomId(ctx)
	if rt == "" || id == "" {
		return ctx.BadParam()
	}

	err := db.Game.Update(id, func(game *db.TGame) error {
		game.Started = true
		return nil
	})
	if err != nil {
		return ctx.Send(err)
	}
	pubsub.Publish(pubsub.GameUpdate, id)
	return ctx.Data(nil)
}

func Stop(ctx *api.Context) error {
	id, rt := room.GetRoomId(ctx)
	if rt == "" || id == "" {
		return ctx.BadParam()
	}

	err := db.Game.Update(id, func(game *db.TGame) error {
		game.Started = false
		return nil
	})
	if err != nil {
		return ctx.Send(err)
	}
	pubsub.Publish(pubsub.GameUpdate, id)
	return ctx.Data(nil)
}
