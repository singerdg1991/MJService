package domain

/*
 * @apiDefine: CustomerLimitationLimitation
 */
type CustomerLimitationLimitation struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:english"`
	Description *string `json:"description" openapi:"example:english"`
}

/*
 * @apiDefine: CustomerLimitation
 */
type CustomerLimitation struct {
	ID           uint                         `json:"id" openapi:"example:1"`
	CustomerID   uint                         `json:"customerId" openapi:"example:1"`
	LimitationID uint                         `json:"limitationId" openapi:"example:1"`
	Limitation   CustomerLimitationLimitation `json:"limitation" openapi:"$ref:CustomerLimitationLimitation"`
	Description  *string                      `json:"description" openapi:"example:english"`
}
