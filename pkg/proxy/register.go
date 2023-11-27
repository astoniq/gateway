package proxy

import (
	"github.com/astoniq/janus/pkg/router"
	"net/http"
	"strings"
)

const (
	methodAll = "ALL"
)

type Register struct {
	router router.Router
}

func (p *Register) Add(definition *RouterDefinition) error {
	return nil
}

func (p *Register) doRegister(listenPath string, def *RouterDefinition, handler http.Handler) {
	if strings.Index(listenPath, "/") != 0 {

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
