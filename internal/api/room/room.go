package room

import (
	"net/http"
	"party-games/internal/api"
	"party-games/internal/db"
	"party-games/internal/pubsub"
	"party-games/internal/utils"
	"strconv"
	"strings"
	"time"
)

func All(ctx *api.Context) error {
	list, err := db.Room.All()
	if err != nil {
		return ctx.ErrMsg(err, "查询房间列表异常")
	}
	return ctx.Data(list)
}

func GetRoomId(ctx *api.Context) (string, string) {
	id := ctx.Param("id")
	rt := ctx.QueryParam("type")
	if rt == "" || id == "" {
		return id, rt
	}
	return rt + id, rt
}

func Query(ctx *api.Context) error {
	id, _ := GetRoomId(ctx)
	if id == "" {
		return ctx.BadParam()
	}

	exist, err := db.Room.Exist(id)
	if !exist {
		return ctx.ErrMsgf(err, "查询房间%s不存在", id)
	}

	room, err := db.Room.Get(id)
	if err != nil {
		return ctx.ErrMsg(err, "查询房间异常")
	}

	return ctx.Data(room.R())
}

func LongQuery(ctx *api.Context) error {
	id, _ := GetRoomId(ctx)
	versionStr := ctx.QueryParam("version")
	if id == "" || versionStr == "" {
		return ctx.BadParam()
	}
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return ctx.BadParam()
	}

	exist, err := db.Room.Exist(id)
	if !exist {
		return ctx.ErrMsgf(err, "查询房间%s不存在", id)
	}

	userId := ctx.QueryParam("userId")
	if userId != "" {
		addOnline(id, userId)
		defer func() {
			delOnline(id, userId)
		}()
	}

	room, err := db.Room.Get(id)
	if err != nil {
		return ctx.ErrMsg(err, "查询房间异常")
	}
	if room.Version > version {
		return ctx.Data(room.SetOnline(isOnline).R())
	}

	ch, cancel := pubsub.Subscribe(pubsub.RoomUpdate, id)
	defer cancel()
	select {
	case <-ctx.Request().Context().Done():
		return nil
	case <-time.After(time.Minute):
		return ctx.NoContent(http.StatusNoContent)
	case <-ch:
		room, err = db.Room.Get(id)
		if err != nil {
			return ctx.ErrMsg(err, "查询房间异常")
		}

		return ctx.Data(room.SetOnline(isOnline).R())
	}
}

type ReqSeatUpdate struct {
	Index  int    `json:"index"`
	UserId string `json:"userId"`
	Count  int    `json:"count"`
}

func Update(ctx *api.Context) (err error) {
	id, _ := GetRoomId(ctx)
	var req ReqSeatUpdate
	if err = ctx.Bind(&req); err != nil || id == "" {
		return ctx.BadParam()
	}
	var user *db.TUser
	if req.UserId != "" && req.Index >= 0 {
		user, err = db.User.Get(req.UserId)
		if err != nil {
			return err
		}
	}
	err = db.Room.Update(id, func(room *db.TRoom) error {
		if req.Count != 0 {
			if req.Count <= len(room.Users) {
				room.Users = room.Users[:req.Count]
			} else {
				ns := make([]db.TRUser, req.Count)
				copy(ns, room.Users)
				room.Users = ns
			}
			return nil
		}
		if req.Index == -1 && req.UserId != "" {
			var otherUsers []db.TRUser
			for i, v := range room.Users {
				if v.Id == req.UserId {
					room.Users[i] = db.TRUser{}
					continue
				}
				if v.Id != "" {
					otherUsers = append(otherUsers, v)
				}
			}
			if len(otherUsers) == 0 {
				room.Owner = ""
			} else if len(otherUsers) == 1 || req.UserId == room.Owner {
				room.Owner = otherUsers[0].Id
			}
			return nil
		}
		if user == nil || len(room.Users) <= req.Index {
			return nil
		}
		if room.Users[req.Index].Id != "" {
			return api.NewMsg("此位置已经有人了")
		}
		empty := true
		for i, u := range room.Users {
			if u.Id != "" {
				empty = false
				if u.Id == user.Id {
					room.Users[i] = db.TRUser{}
				}
			}
		}
		if empty {
			room.Owner = user.Id
		}
		room.Users[req.Index] = db.TRUser{TUser: *user}
		return nil
	})

	if err != nil {
		return ctx.Send(err)
	}
	pubsub.Publish(pubsub.RoomUpdate, id)
	return ctx.Data(nil)
}

type ReqRoom struct {
	Seat int `json:"seat"`
}

func Create(ctx *api.Context) error {
	rt := ctx.QueryParam("type")
	var req ReqRoom
	if err := ctx.Bind(&req); err != nil || req.Seat == 0 || rt == "" {
		return ctx.BadParam()
	}

	room, err := randomNew(rt, req)
	if err != nil {
		return ctx.Send(err)
	}
	err = db.Room.Create(room)
	if err != nil {
		return ctx.ErrMsg(err, "创建房间异常")
	}
	return ctx.Data(room.R())
}

func CreateOrGet(ctx *api.Context) error {
	id, rt := GetRoomId(ctx)
	var req ReqRoom
	if err := ctx.Bind(&req); err != nil || req.Seat == 0 || rt == "" || id == "" {
		return ctx.BadParam()
	}
	exist, err := db.Room.Exist(id)
	if err != nil {
		return ctx.Send(err)
	}

	userId := ctx.QueryParam("userId")
	if userId != "" {
		addOnline(id, userId)
	}

	if exist {
		room, err := db.Room.Get(id)
		if err != nil {
			return ctx.Send(err)
		}
		return ctx.Data(room.SetOnline(isOnline).R())
	}

	room := createNew(id, req.Seat)
	err = db.Room.Create(room)
	if err != nil {
		return ctx.ErrMsg(err, "创建房间异常")
	}
	return ctx.Data(room.SetOnline(isOnline).R())
}

func randomNew(rt string, req ReqRoom) (*db.TRoom, error) {
	var id string
	for {
		id = rt + strings.ToLower(utils.RandomStr(4))
		exist, err := db.Room.Exist(id)
		if err != nil {
			return nil, err
		}
		if !exist {
			break
		}
	}
	return createNew(id, req.Seat), nil
}

func createNew(id string, seat int) *db.TRoom {
	return &db.TRoom{
		Base:  db.Base{Id: id},
		Users: make([]db.TRUser, seat)}
}
