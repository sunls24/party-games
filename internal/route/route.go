package route

import (
	"party-games/internal/api"
	"party-games/internal/api/base"
	"party-games/internal/api/game"
	"party-games/internal/api/game/sswd"
	"party-games/internal/api/game/wzq"
	"party-games/internal/api/room"
	"party-games/internal/api/user"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	g := e.Group("/api")
	g.GET("/:route/long", api.Wrap(base.Long))
	g.POST("/:route/init", api.Wrap(base.Init))

	userR := g.Group("/user")
	userR.GET("/name", api.Wrap(user.Name))
	userR.POST("/save", api.Wrap(user.Save))

	roomR := g.Group("/room")
	roomR.GET("/no", api.Wrap(room.No))
	roomR.POST("/seat", api.Wrap(room.Seat))

	gameR := g.Group("/game")
	gameR.POST("/start", api.Wrap(game.Start))
	gameR.POST("/stop", api.Wrap(game.Stop))
	gameR.POST("/restart", api.Wrap(game.Restart))

	sswdR := g.Group("/sswd")
	sswdR.POST("/vote", api.Wrap(sswd.Vote))
	sswdR.POST("/stage", api.Wrap(sswd.NextStage))
	sswdR.GET("/help", api.Wrap(sswd.Help))

	wzqR := g.Group("/wzq")
	wzqR.POST("/chess", api.Wrap(wzq.Chess))
	wzqR.POST("/regret", api.Wrap(wzq.Regret))
	wzqR.POST("/regret/action", api.Wrap(wzq.RegretAction))
}
