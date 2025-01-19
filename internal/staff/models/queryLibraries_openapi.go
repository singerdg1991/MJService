package models

/*
 * @apiDefine: StaffsQueryLibrariesResponseData
 */
type StaffsQueryLibrariesResponseData struct {
	Limit      int                                 `json:"limit" openapi:"example:10"`
	Offset     int                                 `json:"offset" openapi:"example:0"`
	Page       int                                 `json:"page" openapi:"example:1"`
	TotalRows  int                                 `json:"totalRows" openapi:"example:1"`
	TotalPages int                                 `json:"totalPages" openapi:"example:1"`
	Items      []StaffsCreateLibrariesResponseData `json:"items" openapi:"$ref:StaffsCreateLibrariesResponseData;type:array;required"`
}

/*
 * @apiDefine: StaffsQueryLibrariesResponse
 */
type StaffsQueryLibrariesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryLibrariesResponseData `json:"data" openapi:"$ref:StaffsQueryLibrariesResponseData"`
}
