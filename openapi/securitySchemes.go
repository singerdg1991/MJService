package openapi

import (
	"github.com/hoitek/OpenEngine/engine"
	authtype "github.com/hoitek/OpenEngine/engine/types/authType"
)

var SecuritySchemes = engine.SecuritySchemesTypes{
	ApiKey: engine.ApiKeySecuritySchemesDict{
		"apiKeySecurity": engine.ApiKeySecurityScheme{
			Type:         authtype.TypeHttp,
			Description:  "API Key",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	},
}
