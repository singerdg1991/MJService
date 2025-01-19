package models

/*
 * @apiDefine: EquipmentsDeleteResponseData
 */
type EquipmentsDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: EquipmentsDeleteResponse
 */
type EquipmentsDeleteResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"type:integer;example:200"`
	Data       EquipmentsDeleteResponseData `json:"data" openapi:"$ref:EquipmentsDeleteResponseData;type:object;"`
}
