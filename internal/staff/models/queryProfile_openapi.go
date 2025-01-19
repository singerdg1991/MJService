package models

import "github.com/hoitek/Maja-Service/internal/staff/domain"

/*
 * @apiDefine: StaffsQueryProfileResponse
 */
type StaffsQueryProfileResponse struct {
	StatusCode int                 `json:"statusCode" openapi:"example:200"`
	Data       domain.StaffProfile `json:"data" openapi:"$ref:StaffProfile"`
}
