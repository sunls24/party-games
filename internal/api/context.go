package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"party-games/config"
)

type Context struct {
	echo.Context
	Cfg *config.Config
}

func Middleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Context{c, cfg})
		}
	}
}

func Wrap(fn func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return fn(c.(*Context))
	}
}

func (c *Context) Bad(msg string) error {
	return echo.NewHTTPError(http.StatusBadRequest, msg)
}

//goland:noinspection SpellCheckingInspection
func (c *Context) Badf(f string, args ...any) error {
	return c.Bad(fmt.Sprintf(f, args...))
}

func (c *Context) BadParam() error {
	return c.Bad("请求参数异常")
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
	return Message{
		Code:    -1,
		Message: msg,
	}
}

func NewErrMsg(err error, msg string) Message {
	return Message{
		Code:    -1,
		Message: msg,
		Err:     err,
	}
}

func (c *Context) Send(err error) error {
	var m Message
	if !errors.As(err, &m) {
		const unknown = "未知错误"
		m = NewMsg(unknown)
	}

	if m.Err != nil {
		log.Err(m.Err).Msg(m.Message)
	}
	return c.JSON(http.StatusOK, m)
}

func (c *Context) Msg(msg string) error {
	return c.Send(NewMsg(msg))
}

//goland:noinspection SpellCheckingInspection
func (c *Context) ErrMsgf(err error, f string, args ...any) error {
	return c.Send(NewErrMsg(err, fmt.Sprintf(f, args...)))
}

func (c *Context) ErrMsg(err error, msg string) error {
	return c.Send(NewErrMsg(err, msg))
}

func (c *Context) Data(data any) error {
	return c.Send(NewData(data))
}
