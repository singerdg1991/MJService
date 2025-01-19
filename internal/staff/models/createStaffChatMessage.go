package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
)

/*
 * @apiDefine: StaffsCreateChatMessageRequestBody
 */
type StaffsCreateChatMessageRequestBody struct {
	StaffChatID       int                             `json:"staffChatId,string" openapi:"example:1;required;"`
	SenderUserID      int                             `json:"senderUserId,string" openapi:"example:1;required;"`
	RecipientUserID   int                             `json:"recipientUserId,string" openapi:"example:1;required;"`
	Message           *string                         `json:"message" openapi:"example:message;required;"`
	Attachments       []*govalidityt.File             `json:"attachments" openapi:"format:binary;type:array"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the StaffsCreateChatMessageRequestBody based on the provided schema and request.
//
// It takes an http.Request as a parameter and returns govalidity.ValidityResponseErrors.
func (data *StaffsCreateChatMessageRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffChatId":     govalidity.New("staffChatId").Int().Required(),
		"senderUserId":    govalidity.New("senderUserId").Int().Required(),
		"recipientUserId": govalidity.New("recipientUserId").Int().Required(),
		"message":         govalidity.New("message").Optional(),
		"attachments":     govalidity.New("attachments").Files(),
	}

	// Validate request body
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Check one of the attachments or message is required
	if data.Message == nil && len(data.Attachments) == 0 {
		return govalidity.ValidityResponseErrors{
			"message": []string{"message or attachments is required"},
		}
	}

	// Check one of the attachments or message is required
	if data.Message != nil && len(data.Attachments) > 0 {
		return govalidity.ValidityResponseErrors{
			"message": []string{"message or attachments is required"},
		}
	}

	// Check if sender and recipient are not the same
	if data.SenderUserID == data.RecipientUserID {
		return govalidity.ValidityResponseErrors{
			"senderUserId": []string{"sender and recipient can not be the same"},
		}
	}

	return nil
}
