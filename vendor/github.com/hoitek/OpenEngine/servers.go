package openengine

import "github.com/hoitek/OpenEngine/engine"

func (p *openEngine) AddServers(servers engine.ApiServers) OpenEngine {
	p.Servers = servers
	return p
}
