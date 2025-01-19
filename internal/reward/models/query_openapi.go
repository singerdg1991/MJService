package models

import "github.com/hoitek/Maja-Service/internal/reward/domain"

/*
 * @apiDefine: RewardsQueryResponseDataItem
 */
type RewardsQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: RewardsQueryResponseData
 */
type RewardsQueryResponseData struct {
	Limit      int                            `json:"limit" openapi:"example:10"`
	Offset     int                            `json:"offset" openapi:"example:0"`
	Page       int                            `json:"page" openapi:"example:1"`
	TotalRows  int                            `json:"totalRows" openapi:"example:1"`
	TotalPages int                            `json:"totalPages" openapi:"example:1"`
	Items      []RewardsQueryResponseDataItem `json:"items" openapi:"$ref:RewardsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: RewardsQueryResponse
 */
type RewardsQueryResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200;"`
	Data       RewardsQueryResponseData `json:"data" openapi:"$ref:RewardsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: RewardsQueryNotFoundResponse
 */
type RewardsQueryNotFoundResponse struct {
	Rewards []domain.Reward `json:"rewards" openapi:"$ref:Reward;type:array"`
}
