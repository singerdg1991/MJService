package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreateOtherAttachmentsResponseData
 */
type CustomersCreateOtherAttachmentsResponseData struct {
	ID          uint                                `json:"id" openapi:"example:1"`
	UserID      uint                                `json:"userId" openapi:"example:1"`
	CustomerID  uint                                `json:"customerId" openapi:"example:1"`
	Title       string                              `json:"title" openapi:"example:Title"`
	User        *domain.CustomerOtherAttachmentUser `json:"user" openapi:"$ref:CustomerOtherAttachmentUser"`
	Attachments []*types.UploadMetadata             `json:"attachments" openapi:"$ref:UploadMetadata;type:array;required"`
	CreatedAt   string                              `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt   string                              `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt   *string                             `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersCreateOtherAttachmentsResponse
 */
type CustomersCreateOtherAttachmentsResponse struct {
	StatusCode int                                         `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateOtherAttachmentsResponseData `json:"data" openapi:"$ref:CustomersCreateOtherAttachmentsResponseData"`
}
