package user

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"party-games/internal/api"
	"party-games/internal/constant"
	"party-games/internal/db"
	"party-games/internal/utils"
	"strings"
)

func All(ctx *api.Context) error {
	list, err := db.User.All()
	if err != nil {
		return ctx.Send(api.NewErrMsg(err, "查询用户列表异常"))
	}
	return ctx.Data(list)
}

func Query(ctx *api.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.BadParam()
	}
	exist, err := db.User.Exist(id)
	if !exist {
		return ctx.ErrMsgf(err, "查询用户%s不存在", id)
	}

	user, err := db.User.Get(id)
	if err != nil {
		return ctx.ErrMsg(err, "查询用户异常")
	}
	user.R()
	return ctx.Data(user)
}

func QueryOrCreate(ctx *api.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.BadParam()
	}
	exist, err := db.User.Exist(id)
	if !exist {
		user := *randomNew()
		if err = db.User.Create(&user); err != nil {
			return ctx.ErrMsg(err, "创建用户记录异常")
		}
		user.R()
		return ctx.Data(user)
	}

	user, err := db.User.Get(id)
	if err != nil {
		return ctx.ErrMsg(err, "查询用户异常")
	}
	user.R()
	return ctx.Data(user)
}

func Sync(ctx *api.Context) error {
	var user db.TUser
	if err := ctx.Bind(&user); err != nil {
		return ctx.BadParam()
	}
	if user.Id == "" {
		user = *randomNew()
		if err := db.User.Create(&user); err != nil {
			return ctx.ErrMsg(err, "创建用户记录异常")
		}
		user.R()
	} else {
		if err := db.User.Update(user.Id, func(t *db.TUser) error {
			t.Icon = user.Icon
			return nil
		}); err != nil {
			return ctx.ErrMsg(err, "更新用户记录异常")
		}
	}

	return ctx.Data(user)
}

func randomNew() *db.TUser {
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	return &db.TUser{
		Base: db.Base{Id: id},
		Icon: rand.IntN(constant.UserIconTotal) + 1,
		Name: utils.RandomNickname(),
	}
}
