package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryChatsResponseDataItem
 */
type CyclesQueryChatsResponseDataItem struct {
	ID                 uint                     `json:"id" openapi:"example:1"`
	CyclePickupShiftID uint                     `json:"cyclePickupShiftId" openapi:"example:1"`
	SenderUserID       uint                     `json:"senderUserId" openapi:"example:1"`
	RecipientUserID    uint                     `json:"recipientUserId" openapi:"example:1"`
	CyclePickupShift   *domain.CyclePickupShift `json:"cyclePickupShift" openapi:"$ref:CyclePickupShift;type:object;required"`
	SenderUser         *domain.CycleChatUser    `json:"senderUser" openapi:"$ref:CycleChatUser;type:object;required"`
	RecipientUser      *domain.CycleChatUser    `json:"recipientUser" openapi:"$ref:CycleChatUser;type:object;required"`
	IsSystem           bool                     `json:"isSystem" openapi:"example:false"`
	Message            *string                  `json:"message" openapi:"example:message"`
	Attachments        []*types.UploadMetadata  `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt          string                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          string                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *string                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryChatsResponseData
 */
type CyclesQueryChatsResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []CyclesQueryChatsResponseDataItem `json:"items" openapi:"$ref:CyclesQueryChatsResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesQueryChatsResponse
 */
type CyclesQueryChatsResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200"`
	Data       CyclesQueryChatsResponseData `json:"data" openapi:"$ref:CyclesQueryChatsResponseData"`
}

/*
 * @apiDefine: CyclesQueryChatsNotFoundResponse
 */
type CyclesQueryChatsNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
