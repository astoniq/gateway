package proxy

import "github.com/astoniq/janus/pkg/router"

type Definition struct {
	Upstreams *Upstreams
	Methods   []string
}

type Upstreams struct {
	Balancing string
	Targets   Targets
}

type Target struct {
	Target string
	Weight int
}

type Targets []*Target

type RouterDefinition struct {
	*Definition
	middleware []router.Constructor
}
