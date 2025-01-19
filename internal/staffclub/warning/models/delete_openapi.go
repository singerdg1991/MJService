package models

/*
 * @apiDefine: WarningsDeleteResponseData
 */
type WarningsDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: WarningsDeleteResponse
 */
type WarningsDeleteResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"type:integer;example:200"`
	Data       WarningsDeleteResponseData `json:"data" openapi:"$ref:WarningsDeleteResponseData;type:object;"`
}
