package internal

import (
	"fmt"
	"os"
	"party-games/config"
	"party-games/internal/api"
	"party-games/internal/db"
	"party-games/internal/route"
	"party-games/web"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
}

func NewApp() App {
	return App{}
}

func (App) init() {
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "06-01-02 15:04:05"})
}

func (app App) Run() error {
	app.init()
	cfg := config.MustNew()
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	db.MustInit(cfg.DATA)

	e := echo.New()
	e.Pre(RewriteRoom)
	e.Use(api.Middleware(cfg))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisablePrintStack: true,
	}))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.DefaultHTTPErrorHandler(err, c)
		//goland:noinspection GoTypeAssertionOnErrors
		if _, ok := err.(*echo.HTTPError); !ok {
			log.Err(err).Send()
		}
	}

	route.Register(e)
	e.StaticFS("/", echo.MustSubFS(web.FS, "dist"))
	return e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}

func RewriteRoom(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if ok, p := roomPath(path); ok {
			c.Request().URL.Path = p
			return next(c)
		}
		return next(c)
	}
}

func roomPath(path string) (bool, string) {
	if path == "/" || strings.HasPrefix(path, "/_") ||
		strings.HasPrefix(path, "/api/") ||
		strings.HasSuffix(path, "/room/") {
		return false, ""
	}
	path = strings.Trim(path, "/")
	sp := strings.Split(path, "/")
	if len(sp) != 2 || !isRoomNo(sp[1]) {
		return false, ""
	}
	return true, fmt.Sprintf("/%s/room/", sp[0])
}

func isRoomNo(room string) bool {
	if len(room) != 4 {
		return false
	}

	for _, r := range room {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
