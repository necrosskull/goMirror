package main

import (
	"goMirror/cfg"
	"goMirror/routers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg.Init()
	routers.Init(cfg.Settings.Url)

	engine := gin.Default()
	routers.Router(engine)

	log.Info().Str("url", cfg.Settings.Url).Str("port", cfg.Settings.Port).Msg("Starting server")
	log.Info().Msg("Server started on http://127.0.0.1:" + cfg.Settings.Port)

	err := engine.Run(":" + cfg.Settings.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
		return
	}
}
