package server

import (
	"github.com/astoniq/janus/pkg/api"
	"github.com/astoniq/janus/pkg/proxy"
	"net/http"
)

type Server struct {
	server     *http.Server
	repository api.Repository
	register   *proxy.Register
}
