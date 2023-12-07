package proxy

import (
	"github.com/astoniq/janus/pkg/router"
	"strconv"
	"time"
)

type Definition struct {
	Upstreams          *Upstreams
	ListenPath         string
	Methods            []string
	InsecureSkipVerify bool
	ForwardingTimeouts ForwardingTimeouts
	Hosts              []string
	AppendPath         bool
	StripPath          bool
	PreserveHost       bool
}

type Duration time.Duration

// MarshalJSON implements marshalling from JSON
func (d *Duration) MarshalJSON() ([]byte, error) {
	s := (*time.Duration)(d).String()
	s = strconv.Quote(s)

	return []byte(s), nil
}

// UnmarshalJSON implements unmarshalling from JSON
func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	if s == "null" {
		return
	}

	// if Unquote returns error - assume that string is not quoted at all
	if sUnquoted, err := strconv.Unquote(s); err == nil {
		s = sUnquoted
	}

	t, err := time.ParseDuration(s)
	if err != nil {
		return
	}

	*d = Duration(t)
	return
}

type ForwardingTimeouts struct {
	DialTimeout           Duration `bson:"dial_timeout" json:"dial_timeout"`
	ResponseHeaderTimeout Duration `bson:"response_header_timeout" json:"response_header_timeout"`
}

type Upstreams struct {
	Balancing string
	Targets   Targets
}

type Target struct {
	Target string
}

type Targets []*Target

type RouterDefinition struct {
	*Definition
	middleware []router.Constructor
}

func (t Targets) ToBalancerTargets() []*BalancerTarget {
	var balancerTargets []*BalancerTarget
	for _, t := range t {
		balancerTargets = append(balancerTargets, &BalancerTarget{
			Target: t.Target,
		})
	}

	return balancerTargets
}
