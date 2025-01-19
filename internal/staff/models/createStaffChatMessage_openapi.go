package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
)

/*
 * @apiDefine: StaffsCreateChatMessageResponseData
 */
type StaffsCreateChatMessageResponseData struct {
	ID              uint                         `json:"id" openapi:"example:1"`
	StaffChatID     uint                         `json:"staffChatId" openapi:"example:1"`
	SenderUserID    uint                         `json:"senderUserId" openapi:"example:1"`
	RecipientUserID uint                         `json:"recipientUserId" openapi:"example:1"`
	StaffChat       *domain.StaffChat            `json:"staffChat" openapi:"$ref:StaffChat;type:object;required"`
	SenderUser      *domain.StaffChatMessageUser `json:"senderUser" openapi:"$ref:StaffChatMessageUser;type:object;required"`
	RecipientUser   *domain.StaffChatMessageUser `json:"recipientUser" openapi:"$ref:StaffChatMessageUser;type:object;required"`
	IsSystem        bool                         `json:"isSystem" openapi:"example:false"`
	Message         *string                      `json:"message" openapi:"example:message"`
	MessageType     string                       `json:"messageType" openapi:"example:text"`
	Attachments     []*types.UploadMetadata      `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt       string                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       string                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *string                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsCreateChatMessageResponse
 */
type StaffsCreateChatMessageResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateChatMessageResponseData `json:"data" openapi:"$ref:StaffsCreateChatMessageResponseData"`
}
