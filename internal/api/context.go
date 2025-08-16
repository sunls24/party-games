package api

import (
	"fmt"
	"net/http"
	"party-games/config"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sunls24/gox/openai"
)

type Context struct {
	echo.Context
	C   *config.Config
	OAI *openai.OpenAI
}

func Middleware(cfg *config.Config) echo.MiddlewareFunc {
	oai := openai.New(cfg.OAI.BaseURL, cfg.OAI.APIKey)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Context{c, cfg, oai})
		}
	}
}

func Wrap(fn func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return fn(c.(*Context))
	}
}

func (c *Context) BadParam() error {
	return echo.NewHTTPError(http.StatusBadRequest, "请求参数异常")
}

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Err     error  `json:"-"`
}

func (m Message) Error() string {
	return fmt.Sprintf("%s: %v", m.Message, m.Err)
}

func NewData(data any) Message {
	return Message{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func NewMsg(msg string) Message {
	return NewErrMsg(nil, msg)
}

func NewErrMsg(err error, msg string) Message {
	return Message{
		Code:    -1,
		Message: msg,
		Err:     err,
	}
}

func (c *Context) json(v any) error {
	return c.JSON(http.StatusOK, v)
}

func (c *Context) Err(err error) error {
	const unknown = "未知错误"

	//goland:noinspection GoTypeAssertionOnErrors
	m, ok := err.(Message)
	if !ok {
		m = NewErrMsg(err, unknown)
	}

	if m.Err != nil {
		log.Err(m.Err).Msg(m.Message)
	}
	return c.json(m)
}

func (c *Context) BadMsg() error {
	return c.Msg("不要瞎搞了好吧")
}

func (c *Context) Msg(msg string) error {
	return c.json(NewMsg(msg))
}

//goland:noinspection SpellCheckingInspection
func (c *Context) Msgf(f string, args ...any) error {
	return c.json(NewMsg(fmt.Sprintf(f, args...)))
}

func (c *Context) ErrMsg(err error, msg string) error {
	return c.Err(NewErrMsg(err, msg))
}

//goland:noinspection SpellCheckingInspection
func (c *Context) ErrMsgf(err error, f string, args ...any) error {
	return c.Err(NewErrMsg(err, fmt.Sprintf(f, args...)))
}

func (c *Context) Data(data any) error {
	return c.json(NewData(data))
}
