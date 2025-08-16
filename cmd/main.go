package main

import (
	"party-games/internal"

	"github.com/rs/zerolog/log"
)

func main() {
	err := internal.NewApp().Run()
	if err != nil {
		log.Panic().Err(err).Send()
	}
}
