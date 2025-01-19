package models

/*
 * @apiDefine: OTPResendRequest
 */
type OTPResendRequest struct {
	ExchangeCode string `json:"exchangeCode"`
}

/*
 * @apiDefine: OTPResendResponse
 */
type OTPResendResponse struct {
	Success bool `json:"success"`
}
