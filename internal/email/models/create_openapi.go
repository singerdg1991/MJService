package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/email/domain"
)

/*
 * @apiDefine: EmailsResponseData
 */
type EmailsResponseData struct {
	ID          uint                    `json:"id" openapi:"example:1"`
	SenderID    *uint                   `json:"senderId" openapi:"example:1"`
	Sender      *domain.EmailSender     `json:"sender" openapi:"$ref:EmailSender;required"`
	ToEmail     string                  `json:"email" openapi:"example:sgh370@yahoo.com;required"`
	Cc          []string                `json:"cc" openapi:"example:[];type:array;required"`
	Bcc         []string                `json:"bcc" openapi:"example:[];type:array;required"`
	Title       string                  `json:"title" openapi:"example:John;required"`
	Subject     string                  `json:"subject" openapi:"example:John;required"`
	Message     string                  `json:"message" openapi:"example:John;required"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Category    string                  `json:"category" openapi:"example:outbox;required"`
	StarredAt   *string                 `json:"starred_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt   string                  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   string                  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *string                 `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: EmailsCreateResponse
 */
type EmailsCreateResponse struct {
	StatusCode int                `json:"statusCode" openapi:"example:200"`
	Data       EmailsResponseData `json:"data" openapi:"$ref:EmailsResponseData"`
}
