package models

/*
 * @apiDefine: ServiceOptionsDeleteResponseData
 */
type ServiceOptionsDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: ServiceOptionsDeleteResponse
 */
type ServiceOptionsDeleteResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"type:integer;example:200"`
	Data       ServiceOptionsDeleteResponseData `json:"data" openapi:"$ref:ServiceOptionsDeleteResponseData;type:object;"`
}
