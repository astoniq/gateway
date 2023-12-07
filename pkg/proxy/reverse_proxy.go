package proxy

import (
	"github.com/astoniq/janus/pkg/observability"
	"github.com/astoniq/janus/pkg/router"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewBalancedReverseProxy(def *Definition, balancer *Balancer) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: createDirector(def, balancer),
	}
}

func createDirector(definition *Definition, balancer *Balancer) func(req *http.Request) {

	matcher := router.NewListenPathMatcher()

	return func(req *http.Request) {
		upstream, err := balancer.Elect(definition.Upstreams.Targets.ToBalancerTargets())
		if err != nil {
			log.Err(err).Msg("Could not elect one upstream")
			return
		}

		target, err := url.Parse(upstream.Target)
		if err != nil {
			log.Err(err).Str("upstream_url", upstream.Target).Msg("Could not parse the target URL")
			return
		}

		originalURI := req.RequestURI
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		path := target.Path

		if definition.AppendPath {
			log.Debug().Msg("Appending listen path to the target url")
			path = singleJoiningSlash(target.Path, req.URL.Path)
		}

		if definition.StripPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)
			listenPath := matcher.Extract(definition.ListenPath)

			log.Debug().Str("listen_path", listenPath).Msg("Stripping listen path")

			path = strings.Replace(path, listenPath, "", 1)

			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		log.Debug().Str("path", path).Msg("Upstream Path")

		req.URL.Path = path

		if definition.PreserveHost {
			log.Debug().Msg("Preserving the host header")
		} else {
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		log.Info().
			Str("request", originalURI).
			Str("request-id", observability.RequestIdFromContext(req.Context())).
			Str("upstream-host", req.URL.Host).
			Str("upstream-request", req.URL.RequestURI()).
			Msg("Proxying request to the following upstream")

	}
}

func singleJoiningSlash(a, b string) string {
	a = cleanSlashes(a)
	b = cleanSlashes(b)

	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")

	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		if len(b) > 0 {
			return a + "/" + b
		}
		return a
	}
	return a + b
}

func cleanSlashes(a string) string {
	endSlash := strings.HasSuffix(a, "//")
	startSlash := strings.HasPrefix(a, "//")

	if startSlash {
		a = "/" + strings.TrimPrefix(a, "//")
	}

	if endSlash {
		a = strings.TrimSuffix(a, "//") + "/"
	}

	return a
}
