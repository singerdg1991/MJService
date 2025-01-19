package models

import "github.com/hoitek/Maja-Service/internal/_shared/types"

/*
 * @apiDefine: StaffsCreateLicensesResponseDataLicense
 */
type StaffsCreateLicensesResponseDataLicense struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateLicensesResponseData
 */
type StaffsCreateLicensesResponseData struct {
	ID          int                                     `json:"id" openapi:"example:1"`
	StaffID     int                                     `json:"staffId" openapi:"example:1"`
	LicenseID   int                                     `json:"licenseId" openapi:"example:1"`
	License     StaffsCreateLicensesResponseDataLicense `json:"license" openapi:"$ref:StaffsCreateLicensesResponseDataLicense"`
	ExpireDate  string                                  `json:"expire_date" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments []*types.UploadMetadata                 `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
}

/*
 * @apiDefine: StaffsCreateLicensesResponse
 */
type StaffsCreateLicensesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateLicensesResponseData `json:"data" openapi:"$ref:StaffsCreateLicensesResponseData"`
}
