package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultDialTimeout when connecting to a backend server.
	DefaultDialTimeout = 30 * time.Second

	// DefaultIdleConnsPerHost the default value set for http.Transport.MaxIdleConnsPerHost.
	DefaultIdleConnsPerHost = 64

	// DefaultIdleConnTimeout is the default value for the the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing itself.
	DefaultIdleConnTimeout = 90 * time.Second
)

type Transport struct {
	idleConnectionsPerHost int
	insecureSkipVerify     bool
	dialTimeout            time.Duration
	responseHeaderTimeout  time.Duration
	idleConnTimeout        time.Duration
	idleConnPurgeTicker    *time.Ticker
}

func (t Transport) hash() string {
	return strings.Join([]string{
		fmt.Sprintf("idleConnectionsPerHost:%v", t.idleConnectionsPerHost),
		fmt.Sprintf("insecureSkipVerify:%v", t.insecureSkipVerify),
		fmt.Sprintf("dialTimeout:%v", t.dialTimeout),
		fmt.Sprintf("responseHeaderTimeout:%v", t.responseHeaderTimeout),
		fmt.Sprintf("idleConnTimeout:%v", t.idleConnTimeout),
	}, ";")
}

func NewTransport(registry *Registry, opts ...TransportOption) *http.Transport {
	t := Transport{}
	for _, opt := range opts {
		opt(&t)
	}

	if t.dialTimeout <= 0 {
		t.dialTimeout = DefaultDialTimeout
	}

	if t.idleConnectionsPerHost <= 0 {
		t.idleConnectionsPerHost = DefaultIdleConnsPerHost
	}

	if t.idleConnTimeout == 0 {
		t.idleConnTimeout = DefaultIdleConnTimeout
	}

	hash := t.hash()
	if tr, ok := registry.get(hash); ok {
		return tr
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   t.dialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       t.idleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: t.responseHeaderTimeout,
		MaxIdleConnsPerHost:   t.idleConnectionsPerHost,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: t.insecureSkipVerify},
	}

	if t.idleConnPurgeTicker != nil {
		go func(transport *http.Transport) {
			for {
				select {
				case <-t.idleConnPurgeTicker.C:
					transport.DisableKeepAlives = true
					transport.CloseIdleConnections()
					transport.DisableKeepAlives = false
				}
			}
		}(tr)
	}

	registry.put(hash, tr)

	return tr
}
