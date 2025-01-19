package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
)

/*
 * @apiDefine: StaffsCreateAbsencesResponseData
 */
type StaffsCreateAbsencesResponseData struct {
	ID          uint                                  `json:"id" openapi:"example:1"`
	StaffID     uint                                  `json:"staffId" openapi:"ignored"`
	StartDate   string                                `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate     *string                               `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Reason      *string                               `json:"reason" openapi:"example:reason"`
	Attachments []*types.UploadMetadata               `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Status      *string                               `json:"status" openapi:"example:status"`
	StatusBy    *domain.StaffAbsencesQueryResStatusBy `json:"statusBy" openapi:"$ref:StaffAbsencesQueryResStatusBy"`
	StatusAt    *string                               `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsCreateAbsencesResponse
 */
type StaffsCreateAbsencesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateAbsencesResponseData `json:"data" openapi:"$ref:StaffsCreateAbsencesResponseData"`
}
