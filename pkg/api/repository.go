package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
)

const (
	file = "file"
)

type Repository interface {
	io.Closer
	All() ([]*Definition, error)
}

type Watcher interface {
	Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged)
}

type Listener interface {
	Listen(ctx context.Context, cfgChan chan<- ConfigurationMessage)
}

func BuildRepository(dsn string, refreshTime time.Duration) (Repository, error) {
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing the DSN: %w", err)
	}

	switch dsnUrl.Scheme {
	default:
		return nil, errors.New("selected scheme is not supported to load API definition")
	}
}
