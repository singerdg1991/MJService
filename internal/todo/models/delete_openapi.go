package models

/*
 * @apiDefine: TodosDeleteResponseData
 */
type TodosDeleteResponseData struct {
	IDs interface{} `json:"ids" openapi:"type:array;example:[1,2,3];"`
}

// FIXME: This is a workaround for a bug in OpenEngine we need to support arrays of basic types like int, string etc.

/*
 * @apiDefine: TodosDeleteResponse
 */
type TodosDeleteResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"type:integer;example:200"`
	Data       TodosDeleteResponseData `json:"data" openapi:"$ref:TodosDeleteResponseData;type:object;"`
}
