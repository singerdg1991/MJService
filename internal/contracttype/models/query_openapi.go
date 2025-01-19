package models

import "github.com/hoitek/Maja-Service/internal/contracttype/domain"

/*
 * @apiDefine: ContractTypesQueryResponseDataItem
 */
type ContractTypesQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ContractTypesQueryResponseData
 */
type ContractTypesQueryResponseData struct {
	Limit      int                                  `json:"limit" openapi:"example:10"`
	Offset     int                                  `json:"offset" openapi:"example:0"`
	Page       int                                  `json:"page" openapi:"example:1"`
	TotalRows  int                                  `json:"totalRows" openapi:"example:1"`
	TotalPages int                                  `json:"totalPages" openapi:"example:1"`
	Items      []ContractTypesQueryResponseDataItem `json:"items" openapi:"$ref:ContractTypesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ContractTypesQueryResponse
 */
type ContractTypesQueryResponse struct {
	StatusCode int                            `json:"statusCode" openapi:"example:200;"`
	Data       ContractTypesQueryResponseData `json:"data" openapi:"$ref:ContractTypesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ContractTypesQueryNotFoundResponse
 */
type ContractTypesQueryNotFoundResponse struct {
	ContractTypes []domain.ContractType `json:"contracttypes" openapi:"$ref:ContractType;type:array"`
}
