package proxy

import "time"

type TransportOption func(*Transport)

func TransportWithInsecureSkipVerify(value bool) TransportOption {
	return func(t *Transport) {
		t.insecureSkipVerify = value
	}
}

func TransportWithDialTimeout(d time.Duration) TransportOption {
	return func(t *Transport) {
		t.dialTimeout = d
	}
}

func TransportWithResponseHeaderTimeout(d time.Duration) TransportOption {
	return func(t *Transport) {
		t.responseHeaderTimeout = d
	}
}

func TransportWithIdleConnTimeout(d time.Duration) TransportOption {
	return func(t *Transport) {
		t.idleConnTimeout = d
	}
}

func TransportWithIdleConnPurgeTicker(ticker *time.Ticker) TransportOption {
	return func(t *Transport) {
		t.idleConnPurgeTicker = ticker
	}
}
