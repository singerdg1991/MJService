package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: CycleChatUser
 */
type CycleChatUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

// NewCycleChatUser returns a new instance of CycleChatUser.
//
// It takes four parameters: id, firstName, lastName, and avatarUrl, which are used to initialize the CycleChatUser struct.
// Returns a pointer to the newly created CycleChatUser instance.
func NewCycleChatUser(id uint, firstName string, lastName string, avatarUrl string) *CycleChatUser {
	return &CycleChatUser{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		AvatarUrl: avatarUrl,
	}
}

/*
 * @apiDefine: CycleChat
 */
type CycleChat struct {
	ID                 uint                    `json:"id" openapi:"example:1"`
	CyclePickupShiftID uint                    `json:"cyclePickupShiftId" openapi:"example:1"`
	SenderUserID       uint                    `json:"senderUserId" openapi:"example:1"`
	RecipientUserID    uint                    `json:"recipientUserId" openapi:"example:1"`
	CyclePickupShift   *CyclePickupShift       `json:"cyclePickupShift" openapi:"$ref:CyclePickupShift;type:object;required"`
	SenderUser         *CycleChatUser          `json:"senderUser" openapi:"$ref:CycleChatUser;type:object;required"`
	RecipientUser      *CycleChatUser          `json:"recipientUser" openapi:"$ref:CycleChatUser;type:object;required"`
	IsSystem           bool                    `json:"isSystem" openapi:"example:false"`
	Message            *string                 `json:"message" openapi:"example:message"`
	Attachments        []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt          time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

// NewCycleChat returns a new instance of CycleChat.
//
// Parameters:
//
//	id (uint): The ID of the CycleChat.
//	cyclePickupShiftId (uint): The ID of the CyclePickupShift.
//	senderUserId (uint): The ID of the sender user.
//	recipientUserId (uint): The ID of the recipient user.
//	cyclePickupShift (*CyclePickupShift): The CyclePickupShift instance.
//	senderUser (*CycleChatUser): The sender user instance.
//	recipientUser (*CycleChatUser): The recipient user instance.
//	isSystem (bool): Whether the message is system generated.
//	message (*string): The message content.
//	attachments ([]*types.UploadMetadata): The attachments metadata.
//	createdAt (time.Time): The creation time of the CycleChat.
//	updatedAt (time.Time): The last update time of the CycleChat.
//	deletedAt (*time.Time): The deletion time of the CycleChat.
//
// Returns:
//
//	*CycleChat: The newly created CycleChat instance.
func NewCycleChat(
	id uint,
	cyclePickupShiftId uint,
	senderUserId uint,
	recipientUserId uint,
	cyclePickupShift *CyclePickupShift,
	senderUser *CycleChatUser,
	recipientUser *CycleChatUser,
	isSystem bool,
	message *string,
	attachments []*types.UploadMetadata,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *CycleChat {
	return &CycleChat{
		ID:                 id,
		CyclePickupShiftID: cyclePickupShiftId,
		SenderUserID:       senderUserId,
		RecipientUserID:    recipientUserId,
		CyclePickupShift:   cyclePickupShift,
		SenderUser:         senderUser,
		RecipientUser:      recipientUser,
		IsSystem:           isSystem,
		Message:            message,
		Attachments:        attachments,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		DeletedAt:          deletedAt,
	}
}

// TableName returns the database table name for the CycleChat struct.
//
// No parameters.
// Returns a string representing the table name.
func (c *CycleChat) TableName() string {
	return "cycleChats"
}

// ToJson converts a CycleChat object to a JSON string.
//
// No parameters.
// Returns a string representing the JSON object and an error if the conversion fails.
func (c *CycleChat) ToJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
