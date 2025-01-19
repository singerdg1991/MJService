package openengine

import "github.com/hoitek/OpenEngine/engine"

func flatSecuritySchemes(securitySchemas engine.SecuritySchemesTypes) engine.SecuritySchemes {
	flat := engine.SecuritySchemes{}
	for name, schema := range securitySchemas.ApiKey {
		flat[name] = schema
	}
	for name, schema := range securitySchemas.Http {
		flat[name] = schema
	}
	for name, schema := range securitySchemas.OAuth2 {
		flat[name] = schema
	}
	for name, schema := range securitySchemas.OpenId {
		flat[name] = schema
	}

	return flat
}

func (p *openEngine) AddSecuritySchemes(securitySchemas engine.SecuritySchemesTypes) OpenEngine {
	p.Components.SecuritySchemes = flatSecuritySchemes(securitySchemas)
	return p
}
