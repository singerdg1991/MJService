package models

import "github.com/hoitek/Maja-Service/internal/_shared/types"

/*
 * @apiDefine: StaffsQueryLicensesResponseDataItemLicense
 */
type StaffsQueryLicensesResponseDataItemLicense struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:license 01"`
}

/*
 * @apiDefine: StaffsQueryLicensesResponseDataItem
 */
type StaffsQueryLicensesResponseDataItem struct {
	ID          int                                        `json:"id" openapi:"example:1"`
	StaffID     int                                        `json:"staffId" openapi:"example:1"`
	LicenseID   int                                        `json:"licenseId" openapi:"example:1"`
	License     StaffsQueryLicensesResponseDataItemLicense `json:"license" openapi:"$ref:StaffsQueryLicensesResponseDataItemLicense"`
	ExpireDate  string                                     `json:"expire_date" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments []*types.UploadMetadata                    `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
}

/*
 * @apiDefine: StaffsQueryLicensesResponseData
 */
type StaffsQueryLicensesResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []StaffsQueryLicensesResponseDataItem `json:"items" openapi:"$ref:StaffsQueryLicensesResponseDataItem"`
}

/*
 * @apiDefine: StaffsQueryLicensesResponse
 */
type StaffsQueryLicensesResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryLicensesResponseData `json:"data" openapi:"$ref:StaffsQueryLicensesResponseData"`
}
