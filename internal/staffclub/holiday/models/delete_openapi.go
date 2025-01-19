package models

/*
 * @apiDefine: HolidaysDeleteResponseData
 */
type HolidaysDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: HolidaysDeleteResponse
 */
type HolidaysDeleteResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"type:integer;example:200"`
	Data       HolidaysDeleteResponseData `json:"data" openapi:"$ref:HolidaysDeleteResponseData;type:object;"`
}
