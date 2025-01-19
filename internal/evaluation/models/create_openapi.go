package models

/*
 * @apiDefine: EvaluationsResponseData
 */
type EvaluationsResponseData struct {
	ID             uint    `json:"id" openapi:"example:1"`
	StaffID        uint    `json:"staffId" openapi:"example:1"`
	EvaluationType string  `json:"evaluationType" openapi:"example:grace"`
	Title          string  `json:"title" openapi:"example:saeed"`
	Description    *string `json:"description" openapi:"example:test"`
	CreatedAt      string  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      string  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *string `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: EvaluationsCreateResponse
 */
type EvaluationsCreateResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       EvaluationsResponseData `json:"data" openapi:"$ref:EvaluationsResponseData"`
}
