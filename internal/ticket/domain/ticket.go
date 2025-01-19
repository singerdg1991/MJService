package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: TicketUser
 */
type TicketUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	Email     string `json:"email" openapi:"example:sgh370@yahoo.com"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50"`
}

/*
 * @apiDefine: Ticket
 */
type Ticket struct {
	ID            uint                    `json:"id" openapi:"example:1"`
	Code          string                  `json:"code" openapi:"example:13256987"`
	UserID        uint                    `json:"userId" openapi:"example:1"`
	User          *TicketUser             `json:"user" openapi:"$ref:TicketUser;example:[];required"`
	DepartmentID  uint                    `json:"departmentId" openapi:"example:1"`
	SenderType    string                  `json:"senderType" openapi:"example:customer"`
	RecipientType string                  `json:"recipientType" openapi:"example:customer"`
	Title         string                  `json:"title" openapi:"example:title;required"`
	Description   string                  `json:"description" openapi:"example:description;required"`
	Status        string                  `json:"status" openapi:"example:open;required"`
	Priority      string                  `json:"priority" openapi:"example:low;required"`
	Attachments   []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedByID   uint                    `json:"-" openapi:"ignored"`
	CreatedBy     TicketUser              `json:"createdBy" openapi:"$ref:TicketUser;example:[];required"`
	UpdatedByID   *uint                   `json:"-" openapi:"ignored"`
	UpdatedBy     *TicketUser             `json:"updatedBy" openapi:"$ref:TicketUser;example:[];required"`
	DeletedByID   *uint                   `json:"-" openapi:"ignored"`
	DeletedBy     *TicketUser             `json:"deletedBy" openapi:"$ref:TicketUser;example:[];required"`
	CreatedAt     time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time              `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Ticket) TableName() string {
	return "tickets"
}

func (u *Ticket) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
