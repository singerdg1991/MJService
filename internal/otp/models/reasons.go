package models

/*
 * @apiDefine: OTPReasonsResponse
 */
type OTPReasonsResponse struct {
	Reasons map[string]string `json:"reasons"`
}
