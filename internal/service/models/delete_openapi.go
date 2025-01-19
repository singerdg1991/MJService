package models

/*
 * @apiDefine: ServicesDeleteResponseData
 */
type ServicesDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: ServicesDeleteResponse
 */
type ServicesDeleteResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"type:integer;example:200"`
	Data       ServicesDeleteResponseData `json:"data" openapi:"$ref:ServicesDeleteResponseData;type:object;"`
}
