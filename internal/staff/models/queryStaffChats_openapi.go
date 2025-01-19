package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
)

/*
 * @apiDefine: StaffsQueryChatsResponseDataItem
 */
type StaffsQueryChatsResponseDataItem struct {
	ID              uint                    `json:"id" openapi:"example:1"`
	SenderUserID    uint                    `json:"senderUserId" openapi:"example:1"`
	RecipientUserID uint                    `json:"recipientUserId" openapi:"example:1"`
	SenderUser      *domain.StaffChatUser   `json:"senderUser" openapi:"$ref:StaffChatUser;type:object;required"`
	RecipientUser   *domain.StaffChatUser   `json:"recipientUser" openapi:"$ref:StaffChatUser;type:object;required"`
	IsSystem        bool                    `json:"isSystem" openapi:"example:false"`
	Message         *string                 `json:"message" openapi:"example:message"`
	Attachments     []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt       string                  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       string                  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *string                 `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsQueryChatsResponseData
 */
type StaffsQueryChatsResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []StaffsQueryChatsResponseDataItem `json:"items" openapi:"$ref:StaffsQueryChatsResponseDataItem;type:array"`
}

/*
 * @apiDefine: StaffsQueryChatsResponse
 */
type StaffsQueryChatsResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryChatsResponseData `json:"data" openapi:"$ref:StaffsQueryChatsResponseData"`
}

/*
 * @apiDefine: StaffsQueryChatsNotFoundResponse
 */
type StaffsQueryChatsNotFoundResponse struct {
	Staffs []domain.Staff `json:"staffs" openapi:"$ref:Staff;type:array"`
}
