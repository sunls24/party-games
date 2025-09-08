package wzq

import (
	"encoding/gob"
	"encoding/json"
	"math/rand/v2"
	"party-games/internal/api"
	"party-games/internal/api/game"
	"party-games/internal/db"
	"party-games/internal/utils"

	"github.com/sunls24/gox"
)

func init() {
	game.RegisterStartGame("wzq", startGame)
	gob.Register(Conf{})
	gob.Register(&State{})
}

func startGame(ctx *api.Context, _ int, game *db.TGame, confRaw json.RawMessage) (err error) {
	if confRaw != nil {
		game.Conf, err = gox.Unmarshal[Conf](confRaw)
		if err != nil {
			return err
		}
	}
	state := &State{
		Black:  int8(rand.IntN(2)),
		Board:  utils.RepeatSlice(196, []int8{0}),
		LastCt: -1,
		Regret: -1,
	}
	state.Current = state.Black
	game.State = state
	return nil
}

func Chess(ctx *api.Context) error {
	var req struct {
		Idx uint8 `json:"idx"`
		Val int8  `json:"val"`
		Ct  int   `json:"ct"`
	}
	id := ctx.QueryParam("id")
	if err := ctx.Bind(&req); err != nil || req.Idx < 0 || req.Ct < 0 || id == "" {
		return ctx.BadParam()
	}
	err := game.UpdateState[*State](id, func(state *State) error {
		if req.Val < 0 && state.Current != state.Black {
			return api.NewMsg("有点问题，请刷新一下")
		}
		if req.Ct >= len(state.Board[req.Idx]) {
			ns := make([]int8, req.Ct+1)
			copy(ns, state.Board[req.Idx])
			state.Board[req.Idx] = ns
		}
		state.Board[req.Idx][req.Ct] = req.Val
		state.Current = int8(gox.If(state.Current == 0, 1, 0))
		state.LastIdx = req.Idx
		state.LastCt = int8(req.Ct)
		state.Regret = -1
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func Regret(ctx *api.Context) error {
	var req struct {
		Idx int8 `json:"idx"`
	}
	id := ctx.QueryParam("id")
	if err := ctx.Bind(&req); err != nil || req.Idx < 0 || id == "" {
		return ctx.BadParam()
	}
	err := game.UpdateState[*State](id, func(state *State) error {
		if state.Regret >= 0 {
			return api.NewMsg("点击太快了")
		}
		state.Regret = req.Idx
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func RegretAction(ctx *api.Context) error {
	var req struct {
		Action int8 `json:"action"`
	}
	id := ctx.QueryParam("id")
	if err := ctx.Bind(&req); err != nil || req.Action < 0 || id == "" {
		return ctx.BadParam()
	}
	err := game.UpdateState[*State](id, func(state *State) error {
		if state.Regret < 0 {
			return api.NewMsg("点击太快了")
		}
		if req.Action == 1 {
			state.Current = int8(gox.If(state.Current == 0, 1, 0))
			state.Board[state.LastIdx][state.LastCt] = 0
			state.LastCt = -1
		}
		state.Regret = -1
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}
