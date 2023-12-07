package proxy

import (
	"errors"
	"sync"
)

var (
	ErrEmptyBackendList = errors.New("can not elect backend, Backends empty")
)

type Balancer struct {
	current int
	mu      sync.RWMutex
}

type BalancerTarget struct {
	Target string
}

func NewBalancer() *Balancer {
	return &Balancer{}
}

func (b *Balancer) Elect(hosts []*BalancerTarget) (*BalancerTarget, error) {
	if len(hosts) == 0 {
		return nil, ErrEmptyBackendList
	}

	if len(hosts) == 1 {
		return hosts[0], nil
	}

	if b.current >= len(hosts) {
		b.current = 0
	}

	host := hosts[b.current]

	b.mu.Lock()
	defer b.mu.Unlock()
	b.current++

	return host, nil
}
