package main

import (
	"github.com/astoniq/janus/pkg/api"
	"github.com/astoniq/janus/pkg/config"
	"github.com/astoniq/janus/pkg/server"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.NewConfig()

	repository, err := api.BuildRepository(cfg)

	if err != nil {
		log.Fatal().Err(err).Msg("init: could not build a repository for the database. Shutting down.")
	}

	svr := server.New(
		server.WithConfig(&cfg),
		server.WithRepository(&repository))

	if err := svr.Start(); err != nil {
		log.Fatal().Err(err).Msg("init: could not start server. Shutting down.")
	}

	defer svr.Close()

	svr.Wait()
}
