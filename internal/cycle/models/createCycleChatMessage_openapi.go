package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateChatMessageResponseData
 */
type CyclesCreateChatMessageResponseData struct {
	ID               uint                         `json:"id" openapi:"example:1"`
	CycleChatID      uint                         `json:"cycleChatId" openapi:"example:1"`
	SenderUserID     uint                         `json:"senderUserId" openapi:"example:1"`
	RecipientUserID  uint                         `json:"recipientUserId" openapi:"example:1"`
	CycleChat        *domain.CycleChat            `json:"cycleChat" openapi:"$ref:CycleChat;type:object;required"`
	CyclePickupShift *domain.CyclePickupShift     `json:"cyclePickupShift" openapi:"$ref:CyclePickupShift;type:object;required"`
	SenderUser       *domain.CycleChatMessageUser `json:"senderUser" openapi:"$ref:CycleChatMessageUser;type:object;required"`
	RecipientUser    *domain.CycleChatMessageUser `json:"recipientUser" openapi:"$ref:CycleChatMessageUser;type:object;required"`
	IsSystem         bool                         `json:"isSystem" openapi:"example:false"`
	Message          *string                      `json:"message" openapi:"example:message"`
	MessageType      string                       `json:"messageType" openapi:"example:text"`
	Attachments      []*types.UploadMetadata      `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt        string                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt        string                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt        *string                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreateChatMessageResponse
 */
type CyclesCreateChatMessageResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateChatMessageResponseData `json:"data" openapi:"$ref:CyclesCreateChatMessageResponseData"`
}
