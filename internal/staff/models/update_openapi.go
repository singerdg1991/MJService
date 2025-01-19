package models

import "github.com/hoitek/Maja-Service/internal/staff/domain"

/*
 * @apiDefine: StaffsUpdateResponse
 */
type StaffsUpdateResponse struct {
	Staff domain.Staff `json:"staff" openapi:"$ref:Staff;type:object;"`
}
