package domain

/*
 * @apiDefine: CustomerLocation
 */
type CustomerLocation struct {
	CityID        int    `json:"cityId" openapi:"example: 1"`          // City ID
	CityName      string `json:"cityName" openapi:"example: Helsinki"` // City name
	CustomerCount int    `json:"customerCount" openapi:"example: 1"`   // Number of customers in this city
}
