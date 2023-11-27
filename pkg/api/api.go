package api

import "github.com/astoniq/janus/pkg/proxy"

type Plugin struct {
	Name    string
	Enabled bool
	Config  map[string]interface{}
}

type Definition struct {
	Name        string
	Active      bool
	Proxy       *proxy.Definition
	Plugins     []Plugin
	HealthCheck HealthCheck
}

type HealthCheck struct {
	Url     string
	Timeout int
}

type Configuration struct {
	Definitions []*Definition
}

type ConfigurationChanged struct {
	Configurations *Configuration
}

type ConfigurationOperation int

type ConfigurationMessage struct {
	Operation     ConfigurationOperation
	Configuration *Definition
}

const (
	RemoveOperation ConfigurationOperation = iota
	UpdateOperation
	AddedOperation
)
