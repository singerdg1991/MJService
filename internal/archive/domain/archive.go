package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: ArchiveUser
 */
type ArchiveUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.png;required"`
}

/*
 * @apiDefine: Archive
 */
type Archive struct {
	ID          uint                    `json:"id" openapi:"example:1"`
	UserID      uint                    `json:"userId" openapi:"example:1"`
	User        *ArchiveUser            `json:"user" openapi:"ignored"`
	Title       string                  `json:"title" openapi:"example:title;required"`
	Subject     string                  `json:"subject" openapi:"example:subject;required"`
	Description *string                 `json:"description" openapi:"example:description;required"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Date        string                  `json:"date" openapi:"example:2021-01-01;required"`
	Time        string                  `json:"time" openapi:"example:00:00:00;required"`
	DateTime    time.Time               `json:"-" openapi:"ignored"`
	CreatedAt   time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time              `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Archive) TableName() string {
	return "archives"
}

func (u *Archive) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
