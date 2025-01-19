package models

/*
 * @apiDefine: CustomersDeleteMedicinesResponseData
 */
type CustomersDeleteMedicinesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteMedicinesResponse
 */
type CustomersDeleteMedicinesResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteMedicinesResponseData `json:"data" openapi:"$ref:CustomersDeleteMedicinesResponseData"`
}
