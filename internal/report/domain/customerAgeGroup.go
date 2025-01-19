package domain

/*
 * @apiDefine: CustomerAgeGroup
 */
type CustomerAgeGroup struct {
	AgeGroup      string `json:"ageGroup" openapi:"example:0-20"`     // Age group label (e.g., "0-20", "21-40", etc.)
	CustomerCount int    `json:"customerCount" openapi:"example: 10"` // Number of customers in this age group
	GroupOrder    int    `json:"groupOrder" openapi:"example: 1"`     // Order for display (1 for 0-20, 2 for 21-40, etc.)
}
