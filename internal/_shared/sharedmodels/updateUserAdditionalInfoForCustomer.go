package sharedmodels

type UpdateUserAdditionalInfoForCustomer struct {
	UserID     int64 `json:"userId" openapi:"example:1"`
	CustomerID int64 `json:"customerId" openapi:"example:1"`
}
