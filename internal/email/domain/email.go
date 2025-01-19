package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: EmailSender
 */
type EmailSender struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:John;required"`
	Email     string `json:"email" openapi:"example:John;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:John;required"`
}

/*
 * @apiDefine: Email
 */
type Email struct {
	ID          uint                    `json:"id" openapi:"example:1"`
	SenderID    *uint                   `json:"senderId" openapi:"example:1"`
	Sender      *EmailSender            `json:"sender" openapi:"$ref:EmailSender;required"`
	ToEmail     string                  `json:"email" openapi:"example:sgh370@yahoo.com;required"`
	Cc          []string                `json:"cc" openapi:"example:[];type:array;required"`
	Bcc         []string                `json:"bcc" openapi:"example:[];type:array;required"`
	Title       string                  `json:"title" openapi:"example:John;required"`
	Subject     string                  `json:"subject" openapi:"example:John;required"`
	Message     string                  `json:"message" openapi:"example:John;required"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Category    string                  `json:"category" openapi:"example:outbox;required"`
	StarredAt   *time.Time              `json:"starred_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt   time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Email) TableName() string {
	return "emails"
}

func (u *Email) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
