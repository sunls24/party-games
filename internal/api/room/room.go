package room

import (
	"net/http"
	"party-games/internal/api"
	"party-games/internal/api/types"
	"party-games/internal/db"
	"path"
	"strings"
	"time"

	"github.com/sunls24/gox"
)

func Init(id, uid string, seat int) (types.Room, error) {
	room, err := db.Room.Get(id)
	if err != nil {
		return types.Room{}, err
	}

	addOnline(id, uid, nil)
	defer delOnline(id, uid)

	if room != nil {
		return FillUser(room)
	}
	room = &db.TRoom{
		Users: make([]string, seat),
	}
	err = db.Room.Create(id, room)
	if err != nil {
		return types.Room{}, err
	}
	return FillUser(room)
}

func Long(ctx *api.Context, id, uid string, version int) error {
	room, err := db.Room.Get(id)
	if err != nil {
		return ctx.Err(err)
	}
	if room == nil {
		return ctx.BadMsg()
	}
	addOnline(id, uid, room)
	defer delOnline(id, uid)

	var resp types.Room
	if room.Version > version {
		resp, err = FillUser(room)
		if err != nil {
			return ctx.Err(err)
		}
		return ctx.Data(resp)
	}
	ch, cancel := db.Room.Subscribe(id)
	defer cancel()
	select {
	case <-ctx.Request().Context().Done():
		return nil
	case <-time.After(time.Minute):
		return ctx.NoContent(http.StatusNoContent)
	case room = <-ch:
		resp, err = FillUser(room)
		if err != nil {
			return ctx.Err(err)
		}

		return ctx.Data(resp)
	}
}

func No(ctx *api.Context) error {
	game := ctx.QueryParam("game")
	if game == "" {
		return ctx.BadParam()
	}
	var id string
	for {
		id = strings.ToLower(gox.RandomStr(4))
		r, err := db.Room.Get(path.Join(game, id))
		if err != nil {
			return ctx.Err(err)
		}
		if r == nil {
			break
		}
	}
	return ctx.Data(id)
}

func Seat(ctx *api.Context) (err error) {
	id := ctx.QueryParam("id")
	var req struct {
		Idx    int    `json:"idx"`
		UserId string `json:"uid"`
		Count  int    `json:"count"`
	}
	if err = ctx.Bind(&req); err != nil || id == "" || (req.UserId == "" && req.Count < 2) || req.Count > 12 {
		return ctx.BadParam()
	}

	err = db.Room.UpdatePublish(id, func(room *db.TRoom) error {
		if room.Started {
			return api.NewMsg("游戏进行中，无法操作")
		}
		if req.Count != 0 {
			if req.Count <= len(room.Users) {
				room.Users = room.Users[:req.Count]
			} else {
				ns := make([]string, req.Count)
				copy(ns, room.Users)
				room.Users = ns
			}
			return nil
		}
		if req.Idx == -1 && req.UserId != "" {
			var otherUsers []string
			for i, uid := range room.Users {
				if uid == req.UserId {
					room.Users[i] = ""
					continue
				}
				if uid != "" {
					otherUsers = append(otherUsers, uid)
				}
			}
			if len(otherUsers) == 0 {
				room.Owner = ""
			} else if len(otherUsers) == 1 || req.UserId == room.Owner {
				room.Owner = otherUsers[0]
			}
			return nil
		}
		if len(room.Users) <= req.Idx {
			return nil
		}
		if room.Users[req.Idx] != "" {
			return api.NewMsg("此位置已经有人了")
		}
		empty := true
		for i, u := range room.Users {
			if u != "" {
				empty = false
				if u == req.UserId {
					room.Users[i] = ""
					break
				}
			}
		}
		if empty {
			room.Owner = req.UserId
		}
		room.Users[req.Idx] = req.UserId
		return nil
	})

	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func FillUser(r *db.TRoom) (types.Room, error) {
	var ids []string
	for _, uid := range r.Users {
		if uid == "" {
			continue
		}
		ids = append(ids, uid)
	}
	rr := types.Room{TRoom: r, Users: make([]types.RoomUser, len(r.Users))}
	if len(ids) == 0 {
		return rr, nil
	}
	m, err := db.User.InMap(ids)
	if err != nil {
		return rr, err
	}
	for i, uid := range r.Users {
		if uid == "" {
			continue
		}
		rr.Users[i].TUser = m[r.Users[i]]
		rr.Users[i].TUser.Version = 0
		rr.Users[i].Online = isOnline(r.Id, uid)
	}
	return rr, nil
}
