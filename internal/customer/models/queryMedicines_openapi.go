package models

import "github.com/hoitek/Maja-Service/internal/customer/domain"

/*
 * @apiDefine: CustomersQueryMedicinesResponseData
 */
type CustomersQueryMedicinesResponseData struct {
	Limit      int                                    `json:"limit" openapi:"example:10"`
	Offset     int                                    `json:"offset" openapi:"example:0"`
	Page       int                                    `json:"page" openapi:"example:1"`
	TotalRows  int                                    `json:"totalRows" openapi:"example:1"`
	TotalPages int                                    `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreateMedicinesResponseData `json:"items" openapi:"$ref:CustomersCreateMedicinesResponseData;type:array;"`
}

/*
 * @apiDefine: CustomersQueryMedicinesResponse
 */
type CustomersQueryMedicinesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryMedicinesResponseData `json:"data" openapi:"$ref:CustomersQueryMedicinesResponseData"`
}

/*
 * @apiDefine: CustomersQueryMedicinesNotFoundResponse
 */
type CustomersQueryMedicinesNotFoundResponse struct {
	Customers []domain.Customer `json:"customers" openapi:"$ref:Customer;type:array"`
}
