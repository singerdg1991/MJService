package domain

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: StaffOtherAttachmentUser
 */
type StaffOtherAttachmentUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: StaffOtherAttachment
 */
type StaffOtherAttachment struct {
	ID          uint                      `json:"id" openapi:"example:1"`
	UserID      uint                      `json:"userId" openapi:"example:1"`
	StaffID     uint                      `json:"staffId" openapi:"example:1"`
	User        *StaffOtherAttachmentUser `json:"user" openapi:"$ref:StaffOtherAttachmentUser"`
	Title       string                    `json:"title" openapi:"example:attachment title"`
	Attachments []*types.UploadMetadata   `json:"attachments" openapi:"$ref:UploadMetadata;type:array;required"`
	CreatedAt   time.Time                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}
