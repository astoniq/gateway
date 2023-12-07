package proxy

import (
	"github.com/astoniq/janus/pkg/router"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

const (
	methodAll = "ALL"
)

type Register struct {
	router                 router.Router
	matcher                *router.ListenPathMatcher
	registry               *Registry
	idleConnectionsPerHost int
	idleConnTimeout        time.Duration
	idleConnPurgeTicker    *time.Ticker
	flushInterval          time.Duration
}

func NewRegister(opts ...RegisterOption) *Register {
	r := Register{matcher: router.NewListenPathMatcher()}

	for _, opt := range opts {
		opt(&r)
	}

	return &r
}

func (p *Register) UpdateRouter(router router.Router) {
	p.router = router
}

func (p *Register) Add(definition *RouterDefinition) error {
	balancer := NewBalancer()
	handler := NewBalancedReverseProxy(definition.Definition, balancer)
	handler.FlushInterval = p.flushInterval
	handler.Transport = NewTransport(p.registry,
		TransportWithIdleConnTimeout(p.idleConnTimeout),
		TransportWithIdleConnPurgeTicker(p.idleConnPurgeTicker),
		TransportWithInsecureSkipVerify(definition.InsecureSkipVerify),
		TransportWithDialTimeout(time.Duration(definition.ForwardingTimeouts.DialTimeout)),
		TransportWithResponseHeaderTimeout(time.Duration(definition.ForwardingTimeouts.ResponseHeaderTimeout)),
	)

	if p.matcher.Match(definition.ListenPath) {
		p.doRegister(p.matcher.Extract(definition.ListenPath), definition, handler)
	}

	p.doRegister(definition.ListenPath, definition, handler)

	return nil
}

func (p *Register) doRegister(listenPath string, def *RouterDefinition, handler http.Handler) {
	log.Debug().Str("listen_path", listenPath).Msg("Registering a route")

	if strings.Index(listenPath, "/") != 0 {
		log.Error().Str("listen_path", listenPath).Msg("Route listen path must begin with '/'. Skipping invalid route")
	} else {
		for _, method := range def.Methods {
			if strings.ToUpper(method) == methodAll {
				p.router.Any(listenPath, handler.ServeHTTP, def.middleware...)
			} else {
				p.router.Handle(strings.ToUpper(method), listenPath, handler.ServeHTTP, def.middleware...)
			}
		}
	}
}
