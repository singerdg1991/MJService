package domain

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: CyclePickupShiftTodoCreatedBy
 */
type CyclePickupShiftTodoCreatedBy struct {
	ID        int64  `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CyclePickupShiftTodo
 */
type CyclePickupShiftTodo struct {
	ID                 uint                           `json:"id" openapi:"example:1"`
	CyclePickupShiftID uint                           `json:"cyclePickupShiftId" openapi:"example:1"`
	Title              string                         `json:"title" openapi:"example:title;required"`
	TimeValue          time.Time                      `json:"timeValue" openapi:"example:00:00:00"`
	DateValue          time.Time                      `json:"dateValue" openapi:"example:2021-01-01"`
	Description        *string                        `json:"description" openapi:"example:description"`
	Attachments        []*types.UploadMetadata        `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	NotDoneReason      *string                        `json:"notDoneReason" openapi:"example:reason"`
	Status             string                         `json:"status" openapi:"example:done"`
	DoneAt             *time.Time                     `json:"done_at" openapi:"example:2021-01-01T00:00:00Z"`
	NotDoneAt          *time.Time                     `json:"not_done_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt          time.Time                      `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          time.Time                      `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *time.Time                     `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy          *CyclePickupShiftTodoCreatedBy `json:"createdBy" openapi:"$ref:CyclePickupShiftTodoCreatedBy;required"`
}
