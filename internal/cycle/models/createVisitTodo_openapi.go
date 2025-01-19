package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateVisitTodoResponseData
 */
type CyclesCreateVisitTodoResponseData struct {
	ID                 uint                                  `json:"id" openapi:"example:1"`
	CyclePickupShiftID uint                                  `json:"cyclePickupShiftId" openapi:"example:1"`
	Title              string                                `json:"title" openapi:"example:title;required"`
	TimeValue          string                                `json:"timeValue" openapi:"example:00:00:00"`
	DateValue          string                                `json:"dateValue" openapi:"example:2021-01-01"`
	Description        *string                               `json:"description" openapi:"example:description"`
	Attachments        []*types.UploadMetadata               `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	NotDoneReason      *string                               `json:"notDoneReason" openapi:"example:reason"`
	Status             string                                `json:"status" openapi:"example:done"`
	DoneAt             *string                               `json:"done_at" openapi:"example:2021-01-01T00:00:00Z"`
	NotDoneAt          *string                               `json:"not_done_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt          string                                `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          string                                `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *string                               `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy          *domain.CyclePickupShiftTodoCreatedBy `json:"createdBy" openapi:"$ref:CyclePickupShiftTodoCreatedBy;required"`
}

/*
 * @apiDefine: CyclesCreateVisitTodoResponse
 */
type CyclesCreateVisitTodoResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateVisitTodoResponseData `json:"data" openapi:"$ref:CyclesCreateVisitTodoResponseData"`
}
