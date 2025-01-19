package models

import (
	"github.com/hoitek/Maja-Service/internal/keikkala/domain"
)

/*
 * @apiDefine: KeikkalasResponseData
 */
type KeikkalasResponseData struct {
	ID              uint                      `json:"id" openapi:"example:1"`
	RoleID          *uint                     `json:"roleId" openapi:"example:1"`
	Role            *domain.KeikkalaRole      `json:"role" openapi:"$ref:KeikkalaRole;required"`
	StartDate       string                    `json:"start_date" openapi:"example:2021-01-01"`
	EndDate         string                    `json:"end_date" openapi:"example:2021-01-01"`
	StartTime       string                    `json:"start_time" openapi:"example:00:00:00"`
	EndTime         string                    `json:"end_time" openapi:"example:00:00:00"`
	KaupunkiAddress *string                   `json:"kaupunkiAddress" openapi:"example:address;required"`
	Sections        []*domain.KeikkalaSection `json:"sections" openapi:"$ref:KeikkalaSection;type:array;required"`
	PaymentType     string                    `json:"paymentType" openapi:"example:paySoon;required"`
	ShiftName       string                    `json:"shiftName" openapi:"example:morning;required"`
	Description     *string                   `json:"description" openapi:"example:John;required"`
	Status          string                    `json:"status" openapi:"example:open;required"`
	PickedAt        *string                   `json:"picked_at" openapi:"example:2021-01-01T00:00:00Z"`
	PickedBy        *domain.KeikkalaUser      `json:"pickedBy" openapi:"$ref:KeikkalaUser;required"`
	CreatedBy       *domain.KeikkalaUser      `json:"createdBy" openapi:"$ref:KeikkalaUser;required"`
	UpdatedBy       *domain.KeikkalaUser      `json:"updatedBy" openapi:"$ref:KeikkalaUser;required"`
	CreatedAt       string                    `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       string                    `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *string                   `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: KeikkalasCreateResponse
 */
type KeikkalasCreateResponse struct {
	StatusCode int                   `json:"statusCode" openapi:"example:200"`
	Data       KeikkalasResponseData `json:"data" openapi:"$ref:KeikkalasResponseData"`
}
