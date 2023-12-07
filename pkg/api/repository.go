package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/astoniq/janus/pkg/config"
	"net/url"
)

type Repository interface {
	All() ([]*Definition, error)
}

type Watcher interface {
	Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged)
}

type Listener interface {
	Listen(ctx context.Context, cfgChan chan<- ConfigurationMessage)
}

func BuildRepository(cfg config.Config) (Repository, error) {
	dsnUrl, err := url.Parse(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("error parsing the DSN: %w", err)
	}

	switch dsnUrl.Scheme {
	default:
		return nil, errors.New("selected scheme is not supported to load API definition")
	}
}
