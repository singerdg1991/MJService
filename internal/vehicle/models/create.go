package models

import (
	"github.com/hoitek/Maja-Service/internal/staff/constants"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	ud "github.com/hoitek/Maja-Service/internal/user/domain"

	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
)

/*
 * @apiDefine: VehiclesCreateRequestBody
 */
type VehiclesCreateRequestBody struct {
	VehicleType string   `json:"vehicleType" openapi:"example:car"`
	Brand       string   `json:"brand" openapi:"example:Toyota"`
	Model       string   `json:"model" openapi:"example:Corolla"`
	Year        string   `json:"year" openapi:"example:2019"`
	Variant     string   `json:"variant" openapi:"example:15G"`
	FuelType    string   `json:"fuelType" openapi:"example:Gasoline"`
	VehicleNo   string   `json:"vehicleNo" openapi:"example:ABC123"`
	UserID      int      `json:"userId" openapi:"example:1"`
	User        *ud.User `json:"user" openapi:"ignored"`
}

func (data *VehiclesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"vehicleType": govalidity.New("vehicleType").In([]string{
			constants.STAFF_VEHICLE_TYPE_CAR,
			constants.STAFF_VEHICLE_TYPE_BICYCLE,
			constants.STAFF_VEHICLE_TYPE_PUBLIC_TRANSPORTATION,
		}).Required(),
		"brand":     govalidity.New("brand").Required(),
		"model":     govalidity.New("model").Required(),
		"year":      govalidity.New("year").Int().Required(),
		"variant":   govalidity.New("variant").Required(),
		"fuelType":  govalidity.New("fuelType").Required(),
		"vehicleNo": govalidity.New("vehicleNo").Required(),
		"userId":    govalidity.New("userId").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: VehiclesCreateResponse
 */
type VehiclesCreateResponse struct {
	Vehicle domain.Vehicle `json:"vehicle" openapi:"$ref:Vehicle;type:object;"`
}
