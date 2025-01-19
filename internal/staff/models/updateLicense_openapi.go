package models

import "github.com/hoitek/Maja-Service/internal/_shared/types"

/*
 * @apiDefine: StaffsUpdateLicenseResponseDataLicense
 */
type StaffsUpdateLicenseResponseDataLicense struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsUpdateLicenseResponseData
 */
type StaffsUpdateLicenseResponseData struct {
	ID          int                                    `json:"id" openapi:"example:1"`
	StaffID     int                                    `json:"staffId" openapi:"example:1"`
	LicenseID   int                                    `json:"licenseId" openapi:"example:1"`
	License     StaffsUpdateLicenseResponseDataLicense `json:"license" openapi:"$ref:StaffsUpdateLicenseResponseDataLicense"`
	ExpireDate  string                                 `json:"expire_date" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments []*types.UploadMetadata                `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
}

/*
 * @apiDefine: StaffsUpdateLicenseResponse
 */
type StaffsUpdateLicenseResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200"`
	Data       StaffsUpdateLicenseResponseData `json:"data" openapi:"$ref:StaffsUpdateLicenseResponseData"`
}
