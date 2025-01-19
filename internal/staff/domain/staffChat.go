package domain

import (
	"encoding/json"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/types"
)

/*
 * @apiDefine: StaffChatUser
 */
type StaffChatUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

// NewStaffChatUser returns a new instance of StaffChatUser.
//
// It takes four parameters: id, firstName, lastName, and avatarUrl, which are used to initialize the StaffChatUser struct.
// Returns a pointer to the newly created StaffChatUser instance.
func NewStaffChatUser(id uint, firstName string, lastName string, avatarUrl string) *StaffChatUser {
	return &StaffChatUser{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		AvatarUrl: avatarUrl,
	}
}

/*
 * @apiDefine: StaffChat
 */
type StaffChat struct {
	ID              uint                    `json:"id" openapi:"example:1"`
	SenderUserID    uint                    `json:"senderUserId" openapi:"example:1"`
	RecipientUserID uint                    `json:"recipientUserId" openapi:"example:1"`
	SenderUser      *StaffChatUser          `json:"senderUser" openapi:"$ref:StaffChatUser;type:object;required"`
	RecipientUser   *StaffChatUser          `json:"recipientUser" openapi:"$ref:StaffChatUser;type:object;required"`
	IsSystem        bool                    `json:"isSystem" openapi:"example:false"`
	Message         *string                 `json:"message" openapi:"example:message"`
	Attachments     []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt       time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

// NewStaffChat returns a new instance of StaffChat.
//
// Parameters:
//
//	id (uint): The ID of the StaffChat.
//	senderUserId (uint): The ID of the sender user.
//	recipientUserId (uint): The ID of the recipient user.
//	senderUser (*StaffChatUser): The sender user instance.
//	recipientUser (*StaffChatUser): The recipient user instance.
//	isSystem (bool): Whether the message is system generated.
//	message (*string): The message content.
//	attachments ([]*types.UploadMetadata): The attachments metadata.
//	createdAt (time.Time): The creation time of the StaffChat.
//	updatedAt (time.Time): The last update time of the StaffChat.
//	deletedAt (*time.Time): The deletion time of the StaffChat.
//
// Returns:
//
//	*StaffChat: The newly created StaffChat instance.
func NewStaffChat(
	id uint,
	senderUserId uint,
	recipientUserId uint,
	senderUser *StaffChatUser,
	recipientUser *StaffChatUser,
	isSystem bool,
	message *string,
	attachments []*types.UploadMetadata,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *StaffChat {
	return &StaffChat{
		ID:              id,
		SenderUserID:    senderUserId,
		RecipientUserID: recipientUserId,
		SenderUser:      senderUser,
		RecipientUser:   recipientUser,
		IsSystem:        isSystem,
		Message:         message,
		Attachments:     attachments,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		DeletedAt:       deletedAt,
	}
}

// TableName returns the database table name for the StaffChat struct.
//
// No parameters.
// Returns a string representing the table name.
func (c *StaffChat) TableName() string {
	return "staffChats"
}

// ToJson converts a StaffChat object to a JSON string.
//
// No parameters.
// Returns a string representing the JSON object and an error if the conversion fails.
func (c *StaffChat) ToJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
