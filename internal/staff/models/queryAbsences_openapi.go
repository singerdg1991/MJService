package models

/*
 * @apiDefine: StaffsQueryAbsencesResponseData
 */
type StaffsQueryAbsencesResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []StaffsCreateAbsencesResponseData `json:"items" openapi:"$ref:StaffsCreateAbsencesResponseData"`
}

/*
 * @apiDefine: StaffsQueryAbsencesResponse
 */
type StaffsQueryAbsencesResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryAbsencesResponseData `json:"data" openapi:"$ref:StaffsQueryAbsencesResponseData"`
}
