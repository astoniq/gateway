package server

import (
	"context"
	"github.com/astoniq/janus/pkg/api"
	"github.com/astoniq/janus/pkg/config"
	"github.com/astoniq/janus/pkg/proxy"
	"github.com/astoniq/janus/pkg/router"
	"github.com/astoniq/janus/pkg/web"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server            *http.Server
	repository        *api.Repository
	register          *proxy.Register
	webServer         *web.Server
	config            *config.Config
	configurationChan chan api.ConfigurationChanged
	stopChan          chan struct{}
	apiConfigurations *api.Configuration
}

func New(opts ...ServerOption) *Server {
	s := Server{
		configurationChan: make(chan api.ConfigurationChanged, 100),
		stopChan:          make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s *Server) Start() error {
	newCtx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signals:
			cancel()
			close(signals)
		}
	}()
	return s.StartWithContext(newCtx)
}

func (s *Server) StartWithContext(ctx context.Context) error {
	go func() {
		defer s.Close()
		<-ctx.Done()
		reqAcceptGraceTimeOut := time.Duration(s.config.GraceTimeOut)
		if reqAcceptGraceTimeOut > 0 {
			log.Info().Msgf("Waiting %s for incoming requests to cease", reqAcceptGraceTimeOut)
			time.Sleep(reqAcceptGraceTimeOut)
		}
		log.Info().Msg("Stopping server gracefully")
	}()

}

func (s *Server) Wait() {
	<-s.stopChan
}

func (s *Server) Stop() {

}

func (s *Server) Close() error {

}

func (s *Server) createRouter() router.Router {

}
