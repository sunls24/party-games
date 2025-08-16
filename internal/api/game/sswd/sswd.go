package sswd

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"party-games/internal/api"
	"party-games/internal/api/game"
	"party-games/internal/db"
	"party-games/internal/utils"
	"strings"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/openai"
)

const (
	temperature = 0.88
)

func init() {
	game.RegisterStartGame("sswd", startGame)
	gob.Register(Conf{})
	gob.Register(&State{})
}

func startGame(ctx *api.Context, count int, game *db.TGame, confRaw json.RawMessage) (err error) {
	if confRaw != nil {
		game.Conf, err = gox.Unmarshal[Conf](confRaw)
		if err != nil {
			return err
		}
	}
	state, ok := game.State.(*State)
	if !ok {
		state = &State{}
	}

	resp, err := ctx.OAI.Chat(ctx.Request().Context(), openai.ReqChat{
		Model:       ctx.C.OAI.Model,
		Messages:    buildPrompt(state.Played),
		Temperature: temperature,
	})
	if err != nil {
		return api.NewErrMsg(err, "调用 AI 生成游戏词异常")
	}
	word, err := gox.Unmarshal[GameWord](gox.Str2Bytes(strings.TrimSpace(resp)))
	if err != nil {
		return api.NewErrMsg(err, "解析 AI 输出异常，请重试")
	}

	state.Undercover = rand.IntN(count)
	state.Word = word.Word
	state.SpyWord = word.SpyWord
	state.Players = utils.RepeatSlice(count, Player{Vote: -1})

	played := fmt.Sprintf("%s/%s", word.Word, word.SpyWord)
	if state.Played != "" {
		state.Played += "," + played
	} else {
		state.Played = played
	}

	state.RandomStart()
	game.State = state
	return nil
}

func Vote(ctx *api.Context) error {
	var req struct {
		Idx     int `json:"idx"`
		IdxVote int `json:"idxVote"`
	}
	id := ctx.QueryParam("id")
	if err := ctx.Bind(&req); err != nil || id == "" {
		return ctx.BadParam()
	}
	err := game.UpdateState[*State](id, func(state *State) error {
		if state.Stage != vote {
			return api.NewMsg("当前不是投票阶段")
		}
		if len(state.Players) <= req.Idx || len(state.Players) <= req.IdxVote {
			return nil
		}
		if state.Players[req.Idx].Vote != -1 {
			return ctx.BadMsg()
		}
		state.Players[req.Idx].Vote = req.IdxVote

		m := make(map[int]int) // <index,vote>
		for _, p := range state.Players {
			if p.Out || p.Tie {
				continue
			}
			if p.Vote < 0 {
				return nil
			}
			m[p.Vote]++
		}
		state.VoteDone = true
		for i := range state.Players { // TODO: for
			state.Players[i].Tie = false
		}
		maxVote := 0
		mr := make(map[int][]int) // <vote,[]index>
		for k, v := range m {
			if v > maxVote {
				maxVote = v
			}
			mr[v] = append(mr[v], k)
		}
		if len(mr[maxVote]) == 1 {
			state.Players[mr[maxVote][0]].Out = true
			state.Out = mr[maxVote][0]
			return nil
		}
		for _, i := range mr[maxVote] {
			state.Players[i].Tie = true
		}
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func NextStage(ctx *api.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.BadParam()
	}
	err := game.UpdateState[*State](id, func(state *State) error {
		switch state.Stage {
		case speech:
			state.Stage = vote
		case vote:
			if !state.VoteDone {
				return api.NewMsg("等待所有玩家投票完成")
			}
			var allTie = true
			for _, p := range state.Players {
				if p.Out || p.Tie {
					continue
				}
				allTie = false
				break
			}
			if allTie {
				for i, p := range state.Players {
					if p.Tie {
						state.Players[i].Tie = false
					}
				}
			}

			state.RandomStart()
		default:
			state.Stage = speech
		}
		return nil
	})
	if err != nil {
		return ctx.Err(err)
	}
	return ctx.Data(nil)
}

func Help(ctx *api.Context) error {
	word := ctx.QueryParam("word")
	_, err := ctx.OAI.Chat(ctx.Request().Context(), openai.ReqChat{
		Stream:      true,
		Model:       ctx.C.OAI.Model,
		Messages:    openai.SimplePrompt(fmt.Sprintf("我在玩谁是卧底，给我几条关于它的描述语句：%s\n（不要多余的说明，返回有序列表）", word), ""),
		Temperature: temperature,
		OnStart: func() {
			ctx.Response().Header().Add("Content-Type", "text/stream")
			ctx.Response().WriteHeader(http.StatusOK)
		},
		OnStream: func(data []byte) error {
			if _, err := ctx.Response().Write(data); err != nil {
				return err
			}
			ctx.Response().Flush()
			return nil
		},
	})
	if err != nil {
		return ctx.ErrMsg(err, "调用 AI 生成描述异常")
	}
	return nil
}
