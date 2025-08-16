package base

import (
	"party-games/internal/api"
	"party-games/internal/api/game"
	"party-games/internal/api/room"
	"party-games/internal/api/types"
	"party-games/internal/api/user"
	"strconv"
)

type param struct {
	r  types.Route
	id string
}

func getParam(ctx *api.Context) (p param, bad bool) {
	p.id = ctx.QueryParam("id")
	p.r = types.Route(ctx.Param("route"))
	if p.id == "" || len(p.id) > 64 || p.r == "" {
		return p, true
	}
	return p, false
}

func Init(ctx *api.Context) error {
	p, bad := getParam(ctx)
	if bad {
		return ctx.BadParam()
	}
	var data any
	var err error
	switch p.r {
	case types.RUser:
		data, err = user.Init(p.id)
	case types.RRoom:
		uid := ctx.QueryParam("uid")
		var req struct {
			Seat int `json:"seat"`
		}
		if err = ctx.Bind(&req); err != nil || req.Seat <= 1 || uid == "" {
			return ctx.BadParam()
		}

		data, err = room.Init(p.id, uid, req.Seat)
	default:
		return ctx.BadMsg()
	}
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(data)
}

func Long(ctx *api.Context) error {
	p, bad := getParam(ctx)
	if bad {
		return ctx.BadParam()
	}

	version, err := strconv.Atoi(ctx.QueryParam("version"))
	if err != nil {
		return ctx.BadParam()
	}

	switch p.r {
	case types.RRoom:
		uid := ctx.QueryParam("uid")
		if uid == "" {
			return ctx.BadParam()
		}
		return room.Long(ctx, p.id, uid, version)
	case types.RGame:
		return game.Long(ctx, p.id, version)
	default:
		return ctx.BadMsg()
	}
}
