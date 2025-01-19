package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: CycleChatMessageUser
 */
type CycleChatMessageUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

// NewCycleChatMessageUser creates a new CycleChatMessageUser instance.
//
// id: The ID of the user.
// firstName: The first name of the user.
// lastName: The last name of the user.
// avatarUrl: The URL of the user's avatar.
// Returns a pointer to a CycleChatMessageUser instance.
func NewCycleChatMessageUser(id uint, firstName string, lastName string, avatarUrl string) *CycleChatMessageUser {
	return &CycleChatMessageUser{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		AvatarUrl: avatarUrl,
	}
}

/*
 * @apiDefine: CycleChatMessage
 */
type CycleChatMessage struct {
	ID               uint                    `json:"id" openapi:"example:1"`
	CycleChatID      uint                    `json:"cycleChatId" openapi:"example:1"`
	SenderUserID     uint                    `json:"senderUserId" openapi:"example:1"`
	RecipientUserID  uint                    `json:"recipientUserId" openapi:"example:1"`
	CycleChat        *CycleChat              `json:"cycleChat" openapi:"$ref:CycleChat;type:object;required"`
	CyclePickupShift *CyclePickupShift       `json:"cyclePickupShift" openapi:"$ref:CyclePickupShift;type:object;required"`
	SenderUser       *CycleChatMessageUser   `json:"senderUser" openapi:"$ref:CycleChatMessageUser;type:object;required"`
	RecipientUser    *CycleChatMessageUser   `json:"recipientUser" openapi:"$ref:CycleChatMessageUser;type:object;required"`
	IsSystem         bool                    `json:"isSystem" openapi:"example:false"`
	Message          *string                 `json:"message" openapi:"example:message"`
	MessageType      string                  `json:"messageType" openapi:"example:text"`
	Attachments      []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt        time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt        time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt        *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

// NewCycleChatMessage returns a new instance of CycleChatMessage.
//
// Parameters:
// - id: The unique identifier of the message.
// - cycleChatID: The identifier of the cycle chat the message belongs to.
// - senderUserID: The identifier of the user sending the message.
// - recipientUserID: The identifier of the user receiving the message.
// - cycleChat: The cycle chat the message belongs to.
// - cyclePickupShift: The cycle pickup shift associated with the message.
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
//	*CycleChatMessage: The newly created CycleChatMessage instance.
func NewCycleChatMessage(
	id uint,
	cycleChatID uint,
	senderUserID uint,
	recipientUserID uint,
	cycleChat *CycleChat,
	cyclePickupShift *CyclePickupShift,
	senderUser *CycleChatMessageUser,
	recipientUser *CycleChatMessageUser,
	isSystem bool,
	message *string,
	messageType string,
	attachments []*types.UploadMetadata,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *CycleChatMessage {
	return &CycleChatMessage{
		ID:               id,
		CycleChatID:      cycleChatID,
		SenderUserID:     senderUserID,
		RecipientUserID:  recipientUserID,
		CycleChat:        cycleChat,
		CyclePickupShift: cyclePickupShift,
		SenderUser:       senderUser,
		RecipientUser:    recipientUser,
		IsSystem:         isSystem,
		Message:          message,
		MessageType:      messageType,
		Attachments:      attachments,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		DeletedAt:        deletedAt,
	}
}

// TableName returns the database table name for the CycleChatMessage struct.
//
// No parameters.
// Returns a string representing the table name.
func (c *CycleChatMessage) TableName() string {
	return "cycleChatMessages"
}

// ToJson converts a CycleChatMessage object to a JSON string.
//
// No parameters.
// Returns a string representing the JSON object and an error if the conversion fails.
func (c *CycleChatMessage) ToJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
