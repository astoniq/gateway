package proxy

import (
	"net/http"
	"sync"
)

type Registry struct {
	sync.RWMutex
	store map[string]*http.Transport
}

func NewRegistry() *Registry {
	r := new(Registry)
	r.store = make(map[string]*http.Transport)
	return r
}

func (r *Registry) get(key string) (*http.Transport, bool) {
	r.RLock()
	defer r.RUnlock()
	tr, ok := r.store[key]
	return tr, ok
}

func (r *Registry) put(key string, tr *http.Transport) {
	r.Lock()
	defer r.Unlock()
	r.store[key] = tr
}
