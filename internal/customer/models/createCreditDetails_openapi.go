package models

import "github.com/hoitek/Maja-Service/internal/customer/domain"

/*
 * @apiDefine: CustomersCreateCreditDetailsResponseData
 */
type CustomersCreateCreditDetailsResponseData struct {
	ID                int                                       `json:"id" openapi:"example:1"`
	CustomerID        int                                       `json:"customerId" openapi:"example:1"`
	BillingAddressId  int                                       `json:"billingAddressId" openapi:"example:1"`
	BillingAddress    domain.CustomerCreditDetailBillingAddress `json:"billingAddress" openapi:"$ref:CustomerCreditDetailBillingAddress"`
	BankAccountNumber string                                    `json:"bankAccountNumber" openapi:"example:234567890"`
	CreatedAt         string                                    `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt         string                                    `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt         string                                    `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersCreateCreditDetailsResponse
 */
type CustomersCreateCreditDetailsResponse struct {
	StatusCode int                                      `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateCreditDetailsResponseData `json:"data" openapi:"$ref:CustomersCreateCreditDetailsResponseData"`
}
