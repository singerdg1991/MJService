package models

/*
 * @apiDefine: OTPSettingsResponse
 */
type OTPSettingsResponse struct {
	Length int  `json:"length" openapi:"example:5;"`
	Enable bool `json:"enable" openapi:"example:true;"`
}
