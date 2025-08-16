package user

import (
	"math/rand/v2"
	"party-games/internal/api"
	"party-games/internal/constant"
	"party-games/internal/db"
)

func Init(id string) (*db.TUser, error) {
	user, err := db.User.Get(id)
	if user != nil {
		return user, err
	}

	user = &db.TUser{
		Name: constant.RandomNickname(),
		Icon: rand.IntN(constant.UserIconTotal) + 1,
	}
	err = db.User.Create(id, user)
	if err != nil {
		return nil, err
	}
	user.Version = 0
	return user, nil
}
func Name(ctx *api.Context) error {
	return ctx.Data(constant.RandomNickname())
}

func Save(ctx *api.Context) error {
	var req struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Icon int    `json:"icon"`
	}
	if err := ctx.Bind(&req); err != nil || req.Id == "" || req.Name == "" || req.Icon <= 0 || req.Icon > constant.UserIconTotal {
		return err
	}
	var ret *db.TUser
	err := db.User.Update(req.Id, func(user *db.TUser) error {
		user.Name = req.Name
		user.Icon = req.Icon
		ret = user
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(ret)
}
