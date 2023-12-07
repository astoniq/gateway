package router

import "github.com/go-chi/chi/v5"

type ChiRouter struct {
	mux chi.Router
}

func NewChiRouterWithOptions(options Options) *ChiRouter {
	router := chi.NewRouter()
	router.NotFound(options.NotFoundHandler)
	return &ChiRouter{
		mux: router,
	}
}
