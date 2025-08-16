package route

import (
	"party-games/internal/api"
	"party-games/internal/api/game"
	"party-games/internal/api/room"
	"party-games/internal/api/user"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	g := e.Group("/api")

	u := g.Group("/user")
	u.GET("", api.Wrap(user.All))
	u.POST("", api.Wrap(user.Sync))
	u.GET("/:id", api.Wrap(user.Query))
	u.GET("/:id/cr", api.Wrap(user.QueryOrCreate))

	r := g.Group("/room")
	r.GET("", api.Wrap(room.All))
	r.POST("", api.Wrap(room.Create))
	r.POST("/:id", api.Wrap(room.Update))
	r.GET("/:id", api.Wrap(room.Query))
	r.POST("/:id/cr", api.Wrap(room.CreateOrGet))
	r.GET("/:id/long", api.Wrap(room.LongQuery))

	gm := g.Group("/game")
	gm.POST("/:id/cr", api.Wrap(game.CreateOrGet))
	gm.GET("/:id/long", api.Wrap(game.LongQuery))
	gm.POST("/:id/start", api.Wrap(game.Start))
	gm.POST("/:id/stop", api.Wrap(game.Stop))
}
