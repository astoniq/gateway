package server

import (
	"github.com/astoniq/janus/pkg/api"
	"github.com/astoniq/janus/pkg/config"
)

type ServerOption func(server *Server)

func WithConfig(config *config.Config) ServerOption {
	return func(s *Server) {
		s.config = config
	}
}

func WithRepository(repository *api.Repository) ServerOption {
	return func(s *Server) {
		s.repository = repository
	}
}
