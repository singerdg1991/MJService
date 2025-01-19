package openapi

import "github.com/hoitek/OpenEngine/engine"

/*
 * @apiDefine: ErrDataType
 */

type ErrDataType struct {
	Key   string `json:"key" openapi:"example:key;"`
	Value string `json:"value" openapi:"example:value;"`
}

/*
 * @apiDefine: ErrorResponse
 */
type ErrorResponse struct {
	Errors  []ErrDataType `json:"errors" yaml:"errors" openapi:"$ref:ErrDataType"`
	Message string        `json:"message" yaml:"message"`
}

var DataAndMessageSwaggerErrorResponses = engine.ErrorResponses{
	"400": engine.Response{
		Description: "Bad Request2",
		Content: engine.Content{
			ApplicationJson: engine.MediaType{
				Schema: engine.DataSchema{
					Ref: "ErrorResponse",
				},
			},
			ApplicationXWwwFormUrlencoded: engine.MediaType{
				Schema: engine.DataSchema{
					Ref: "ErrorResponse",
				},
			},
		},
	},
	"401": engine.Response{
		Description: "Bad Request",
		Content: engine.Content{
			ApplicationJson: engine.MediaType{
				Schema: engine.DataSchema{
					Ref: "ErrorResponse",
				},
			},
			ApplicationXWwwFormUrlencoded: engine.MediaType{
				Schema: engine.DataSchema{
					Ref: "ErrorResponse",
				},
			},
		},
	},
}
