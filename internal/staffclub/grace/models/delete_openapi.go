package models

/*
 * @apiDefine: GracesDeleteResponseData
 */
type GracesDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: GracesDeleteResponse
 */
type GracesDeleteResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"type:integer;example:200"`
	Data       GracesDeleteResponseData `json:"data" openapi:"$ref:GracesDeleteResponseData;type:object;"`
}
