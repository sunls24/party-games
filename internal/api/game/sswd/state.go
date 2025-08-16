package sswd

import (
	"fmt"
	"math/rand/v2"

	"github.com/sunls24/gox/openai"
)

type GameWord struct {
	Word    string `json:"word"`
	SpyWord string `json:"spyWord"`
}

type Conf struct {
}

type Stage int

const (
	speech Stage = iota + 1
	vote
)

type Player struct {
	Out  bool `json:"out"`
	Tie  bool `json:"tie"`
	Vote int  `json:"vote"` // index
}

type State struct {
	Undercover int    `json:"undercover"`
	Word       string `json:"word"`
	SpyWord    string `json:"spyWord"`

	Stage    Stage    `json:"stage"`
	Players  []Player `json:"players"`
	VoteDone bool     `json:"voteDone"`
	Out      int      `json:"out"` // index

	Start     int  `json:"start"`
	Clockwise bool `json:"clockwise"`

	Played string `json:"-"`
}

func buildPrompt(played string) []openai.Message {
	const systemPrompt = `你是一个专业的“谁是卧底”游戏词对生成AI，你的任务是生成一对高质量的游戏词，基于以下高质量标准：
- 相似度高：词对在类别、功能、外观或文化联想上高度重叠，但不是同义词（相似度70-80%）
- 差异性明显但隐蔽：至少有2-3个关键区别，这些区别应通过多轮描述逐步暴露，而非第一轮即暴露
- 创新与多样性：从不同主题随机生成（如食物、动物、科技、名人、自然、历史、交通等）优先选择有趣，能引发笑点或讨论的词对
- 描述空间丰富：词有多个维度可描述（如颜色、用途、情感、文化、历史等）鼓励玩家使用创意、幽默或误导性描述，以增加游戏乐趣

示例输出：{"word":"咖啡","spyWord":"茶"}（两者都是饮料，但来源和口味不同）
每次只输出一个JSON对象，不要添加任何额外文本`
	const example = `{"word":"咖啡","spyWord":"茶"}`
	ret := []openai.Message{
		openai.SystemMessage(systemPrompt),
		openai.AssistantMessage(example),
	}
	if played != "" {
		ret = append(ret, openai.UserMessage(fmt.Sprintf("下面是已经玩过的词，避免重复生成:\n%s", played)))
	} else {
		ret = append(ret, openai.UserMessage("生成一对游戏词"))
	}
	return ret
}

func (s *State) RandomStart() {
	s.Stage = speech
	s.VoteDone = false
	s.Out = -1

	var ties []int
	for i, p := range s.Players {
		if p.Tie {
			ties = append(ties, i)
		}
		s.Players[i].Vote = -1
	}
	if len(ties) == 0 {
		s.Start = rand.IntN(len(s.Players))
	} else {
		s.Start = ties[rand.IntN(len(ties))]
	}
	s.Clockwise = rand.IntN(100) >= 50
}
