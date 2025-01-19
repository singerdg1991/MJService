package domain

import (
	"encoding/json"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/types"
)

/*
 * @apiDefine: StaffChatMessageUser
 */
type StaffChatMessageUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

// NewStaffChatMessageUser creates a new StaffChatMessageUser instance.
//
// id: The ID of the user.
// firstName: The first name of the user.
// lastName: The last name of the user.
// avatarUrl: The URL of the user's avatar.
// Returns a pointer to a StaffChatMessageUser instance.
func NewStaffChatMessageUser(id uint, firstName string, lastName string, avatarUrl string) *StaffChatMessageUser {
	return &StaffChatMessageUser{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		AvatarUrl: avatarUrl,
	}
}

/*
 * @apiDefine: StaffChatMessage
 */
type StaffChatMessage struct {
	ID              uint                    `json:"id" openapi:"example:1"`
	StaffChatID     uint                    `json:"staffChatId" openapi:"example:1"`
	SenderUserID    uint                    `json:"senderUserId" openapi:"example:1"`
	RecipientUserID uint                    `json:"recipientUserId" openapi:"example:1"`
	StaffChat       *StaffChat              `json:"staffChat" openapi:"$ref:StaffChat;type:object;required"`
	SenderUser      *StaffChatMessageUser   `json:"senderUser" openapi:"$ref:StaffChatMessageUser;type:object;required"`
	RecipientUser   *StaffChatMessageUser   `json:"recipientUser" openapi:"$ref:StaffChatMessageUser;type:object;required"`
	IsSystem        bool                    `json:"isSystem" openapi:"example:false"`
	Message         *string                 `json:"message" openapi:"example:message"`
	MessageType     string                  `json:"messageType" openapi:"example:text"`
	Attachments     []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt       time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

// NewStaffChatMessage returns a new instance of StaffChatMessage.
//
// Parameters:
// - id: The unique identifier of the message.
// - staffChatID: The identifier of the staff chat the message belongs to.
// - senderUserID: The identifier of the user sending the message.
// - recipientUserID: The identifier of the user receiving the message.
// - staffChat: The staff chat the message belongs to.
// - staffPickupShift: The staff pickup shift associated with the message.
// - senderUser: The user sending the message.
// - recipientUser: The user receiving the message.
// - isSystem: Whether the message is a system message.
// - message: The content of the message.
// - messageType: The type of the message.
// - attachments: The attachments associated with the message.
// - createdAt: The timestamp when the message was created.
// - updatedAt: The timestamp when the message was last updated.
// - deletedAt: The timestamp when the message was deleted.
//
// Returns:
//
//	*StaffChatMessage: The newly created StaffChatMessage instance.
func NewStaffChatMessage(
	id uint,
	staffChatID uint,
	senderUserID uint,
	recipientUserID uint,
	staffChat *StaffChat,
	senderUser *StaffChatMessageUser,
	recipientUser *StaffChatMessageUser,
	isSystem bool,
	message *string,
	messageType string,
	attachments []*types.UploadMetadata,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *StaffChatMessage {
	return &StaffChatMessage{
		ID:              id,
		StaffChatID:     staffChatID,
		SenderUserID:    senderUserID,
		RecipientUserID: recipientUserID,
		StaffChat:       staffChat,
		SenderUser:      senderUser,
		RecipientUser:   recipientUser,
		IsSystem:        isSystem,
		Message:         message,
		MessageType:     messageType,
		Attachments:     attachments,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		DeletedAt:       deletedAt,
	}
}

// TableName returns the database table name for the StaffChatMessage struct.
//
// No parameters.
// Returns a string representing the table name.
func (c *StaffChatMessage) TableName() string {
	return "staffChatMessages"
}

// ToJson converts a StaffChatMessage object to a JSON string.
//
// No parameters.
// Returns a string representing the JSON object and an error if the conversion fails.
func (c *StaffChatMessage) ToJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
