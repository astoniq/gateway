package proxy

import (
	"github.com/astoniq/janus/pkg/router"
	"time"
)

type RegisterOption func(register *Register)

func RegisterWithRouter(router router.Router) RegisterOption {
	return func(register *Register) {
		register.router = router
	}
}

func RegisterWithRegistry(registry *Registry) RegisterOption {
	return func(register *Register) {
		register.registry = registry
	}
}

func RegisterWithFlushInterval(d time.Duration) RegisterOption {
	return func(register *Register) {
		register.flushInterval = d
	}
}

func RegisterWithIdleConnectionsPerHost(value int) RegisterOption {
	return func(r *Register) {
		r.idleConnectionsPerHost = value
	}
}

func RegisterWithIdleConnTimeout(d time.Duration) RegisterOption {
	return func(r *Register) {
		r.idleConnTimeout = d
	}
}

func RegisterWithIdleConnPurgeTicker(d time.Duration) RegisterOption {
	var ticker *time.Ticker

	if d != 0 {
		ticker = time.NewTicker(d)
	}

	return func(t *Register) {
		t.idleConnPurgeTicker = ticker
	}
}
