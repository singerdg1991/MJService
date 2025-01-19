package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
)

/*
 * @apiDefine: VehicleFilterType
 */
type VehicleFilterType struct {
	UserID      filters.FilterValue[int]    `json:"userId,omitempty" openapi:"$ref:FilterValueInt;example:{\"userId\":{\"op\":\"equals\",\"value\":1}"`
	VehicleType filters.FilterValue[string] `json:"VehicleType,omitempty" openapi:"$ref:FilterValueString;example:{\"VehicleType\":{\"op\":\"equals\",\"value\":\"car\"}"`
	Brand       filters.FilterValue[string] `json:"brand,omitempty" openapi:"$ref:FilterValueString;example:{\"brand\":{\"op\":\"equals\",\"value\":\"Toyota\"}"`
	Model       filters.FilterValue[string] `json:"model,omitempty" openapi:"$ref:FilterValueString;example:{\"model\":{\"op\":\"equals\",\"value\":\"Vios\"}"`
	Year        filters.FilterValue[string] `json:"year,omitempty" openapi:"$ref:FilterValueString;example:{\"year\":{\"op\":\"equals\",\"value\":\"2019\"}"`
	Variant     filters.FilterValue[string] `json:"variant,omitempty" openapi:"$ref:FilterValueString;example:{\"variant\":{\"op\":\"equals\",\"value\":\"1.5G\"}"`
	FuelType    filters.FilterValue[string] `json:"fuelType,omitempty" openapi:"$ref:FilterValueString;example:{\"fuelType\":{\"op\":\"equals\",\"value\":\"Gasoline\"}"`
	VehicleNo   filters.FilterValue[string] `json:"vehicleNo,omitempty" openapi:"$ref:FilterValueString;example:{\"vehicleNo\":{\"op\":\"equals\",\"value\":\"ABC123\"}"`
	CreatedAt   filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: VehiclesQueryRequestParams
 */
type VehiclesQueryRequestParams struct {
	ID      int               `json:"id,string,omitempty" openapi:"example:1"`
	Page    int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters VehicleFilterType `json:"filters,omitempty" openapi:"$ref:VehicleFilterType;in:query"`
}

func (data *VehiclesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: VehiclesQueryResponseDataItemVehicleType
 */
type VehiclesQueryResponseDataItemVehicleType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:car"`
}

/*
 * @apiDefine: VehiclesQueryResponseDataItemLanguageSkill
 */
type VehiclesQueryResponseDataItemLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: VehiclesQueryResponseDataItemRole
 */
type VehiclesQueryResponseDataItemRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: VehiclesQueryResponseDataItemUser
 */
type VehiclesQueryResponseDataItemUser struct {
	ID                      uint                                         `json:"id" openapi:"example:1"`
	RoleID                  uint                                         `json:"roleId" openapi:"example:1"`
	FirstName               string                                       `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName                string                                       `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username                string                                       `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Email                   string                                       `json:"email" openapi:"example:sgh370@yahoo.com"`
	Password                string                                       `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Phone                   string                                       `json:"phone" openapi:"example:09123456789;required;"`
	Telephone               string                                       `json:"telephone" openapi:"example:02112345678"`
	LanguageSkills          []VehiclesQueryResponseDataItemLanguageSkill `json:"languageSkills" openapi:"$ref:VehiclesQueryResponseDataItemLanguageSkill;type:array;"`
	RegistrationNumber      string                                       `json:"registrationNumber" openapi:"example:1234567890"`
	WorkPhoneNumber         string                                       `json:"workPhoneNumber" openapi:"example:02112345678"`
	Gender                  string                                       `json:"gender" openapi:"example:male"`
	AccountNumber           string                                       `json:"accountNumber" openapi:"example:02112345678"`
	NationalCode            string                                       `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate               string                                       `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z"`
	AvatarUrl               string                                       `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	ForcedChangePassword    bool                                         `json:"forcedChangePassword" openapi:"example:true"`
	Role                    VehiclesQueryResponseDataItemRole            `json:"role" openapi:"$ref:VehiclesQueryResponseDataItemRole"`
	PrivacyPolicyAcceptedAt string                                       `json:"privacy_policy_accepted_at" openapi:"example:2023-05-06T12:34:56Z"`
	SuspendedAt             string                                       `json:"suspended_at" openapi:"example:2023-05-06T12:34:56Z"`
	CreatedAt               string                                       `json:"created_at" openapi:"example:2023-05-06T12:34:56Z"`
	UpdatedAt               string                                       `json:"updated_at" openapi:"example:2023-05-06T12:34:56Z"`
	DeletedAt               string                                       `json:"deleted_at" openapi:"example:2023-05-06T12:34:56Z"`
}

/*
 * @apiDefine: VehiclesQueryResponseDataItemCompany
 */
type VehiclesQueryResponseDataItemCompany struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: VehiclesQueryResponseDataItem
 */
type VehiclesQueryResponseDataItem struct {
	ID          int                                  `json:"id" openapi:"example:1"`
	VehicleType string                               `json:"vehicleType" openapi:"example:car"`
	Owner       int                                  `json:"owner" openapi:"example:1"`
	OwnerType   string                               `json:"ownerType" openapi:"example:user"`
	User        VehiclesQueryResponseDataItemUser    `json:"user" openapi:"$ref:VehiclesQueryResponseDataItemUser"`
	Company     VehiclesQueryResponseDataItemCompany `json:"company" openapi:"$ref:VehiclesQueryResponseDataItemCompany"`
	Brand       string                               `json:"brand" openapi:"example:peugeot"`
	Model       string                               `json:"model" openapi:"example:206"`
	Year        int                                  `json:"year" openapi:"example:2020"`
	Variant     string                               `json:"variant" openapi:"example:206"`
	FuelType    string                               `json:"fuelType" openapi:"example:gasoline"`
	VehicleNo   string                               `json:"vehicleNo" openapi:"example:1234567890"`
	CreatedAt   string                               `json:"createdAt" openapi:"example:2023-05-06T12:34:56Z"`
	UpdatedAt   string                               `json:"updatedAt" openapi:"example:2023-05-06T12:34:56Z"`
	DeletedAt   string                               `json:"deletedAt" openapi:"example:2023-05-06T12:34:56Z"`
}

/*
 * @apiDefine: VehiclesQueryResponseData
 */
type VehiclesQueryResponseData struct {
	Limit      int                             `json:"limit" openapi:"example:10"`
	Offset     int                             `json:"offset" openapi:"example:0"`
	Page       int                             `json:"page" openapi:"example:1"`
	TotalRows  int                             `json:"totalRows" openapi:"example:1"`
	TotalPages int                             `json:"totalPages" openapi:"example:1"`
	Items      []VehiclesQueryResponseDataItem `json:"items" openapi:"$ref:VehiclesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: VehiclesQueryResponse
 */
type VehiclesQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200"`
	Data       VehiclesQueryResponseData `json:"data" openapi:"$ref:VehiclesQueryResponseData"`
}

/*
 * @apiDefine: VehiclesQueryNotFoundResponse
 */
type VehiclesQueryNotFoundResponse struct {
	Vehicles []domain.Vehicle `json:"vehicles" openapi:"$ref:Vehicle;type:array"`
}
