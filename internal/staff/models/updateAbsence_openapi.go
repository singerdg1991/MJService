package models

/*
 * @apiDefine: StaffsUpdateAbsenceResponseData
 */
type StaffsUpdateAbsenceResponseData struct {
	ID        int    `json:"id" openapi:"example:1"`
	StartDate string `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate   string `json:"end_date" openapi:"example:2020-01-01T00:00:00Z"`
	Reason    string `json:"reason" openapi:"example:reason"`
}

/*
 * @apiDefine: StaffsUpdateAbsenceResponse
 */
type StaffsUpdateAbsenceResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200"`
	Data       StaffsUpdateAbsenceResponseData `json:"data" openapi:"$ref:StaffsUpdateAbsenceResponseData"`
}
