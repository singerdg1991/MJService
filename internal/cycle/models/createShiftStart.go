package models

import (
	"fmt"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: CyclesCreateShiftStartRequestBody
 */
type CyclesCreateShiftStartRequestBody struct {
	CycleID             int               `json:"cycleId,string" openapi:"example:1;required;"`
	ShiftID             int               `json:"shiftId,string" openapi:"example:1;required;"`
	VehicleType         string            `json:"vehicleType" openapi:"example:own;required;"`
	StartLocation       *string           `json:"startLocation" openapi:"example:office;required;"`
	Reason              *string           `json:"reason" openapi:"example:This is a reason;"`
	VerificationPicture *govalidityt.File `json:"verificationPicture" openapi:"format:binary"`
}

// ValidateBody validates the body of the HTTP request for creating a shift start.
//
// It checks the request body against a predefined schema and returns any validation errors.
// The function takes an HTTP request as a parameter and returns a ValidityResponseErrors object.
func (data *CyclesCreateShiftStartRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId": govalidity.New("cycleId").Int().Required(),
		"shiftId": govalidity.New("shiftId").Int().Required(),
		"vehicleType": govalidity.New("vehicleType").In([]string{
			constants.SHIFT_VEHICLE_TYPE_OWN,
			constants.SHIFT_VEHICLE_TYPE_COMPANY,
			constants.SHIFT_VEHICLE_TYPE_PUBLIC_TRANSPORTATION,
		}),
		"startLocation":       govalidity.New("startLocation").Optional(),
		"reason":              govalidity.New("reason").Optional(),
		"verificationPicture": govalidity.New("verificationPicture").File(),
	}
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	switch data.VehicleType {
	case constants.SHIFT_VEHICLE_TYPE_OWN:
		if data.StartLocation == nil {
			return govalidity.ValidityResponseErrors{
				"startLocation": []string{
					"startLocation is required",
					fmt.Sprintf("startLocation should be one of %s, %s", constants.SHIFT_START_LOCATION_OFFICE, constants.SHIFT_START_LOCATION_OTHER_LOCATION),
				},
			}
		}
		if *data.StartLocation != constants.SHIFT_START_LOCATION_OFFICE && *data.StartLocation != constants.SHIFT_START_LOCATION_OTHER_LOCATION {
			return govalidity.ValidityResponseErrors{
				"startLocation": []string{
					fmt.Sprintf("startLocation should be one of %s, %s", constants.SHIFT_START_LOCATION_OFFICE, constants.SHIFT_START_LOCATION_OTHER_LOCATION),
				},
			}
		}
		if data.Reason == nil || *data.Reason == "" {
			return govalidity.ValidityResponseErrors{
				"reason": []string{
					"reason is required",
				},
			}
		}
		if *data.StartLocation == constants.SHIFT_START_LOCATION_OFFICE {
			if data.VerificationPicture == nil {
				return govalidity.ValidityResponseErrors{
					"verificationPicture": []string{
						"verificationPicture is required",
					},
				}
			}
		}
	case constants.SHIFT_VEHICLE_TYPE_COMPANY:
		if data.StartLocation != nil {
			return govalidity.ValidityResponseErrors{
				"startLocation": []string{
					"You don't need to provide startLocation",
				},
			}
		}
		if data.VerificationPicture == nil {
			return govalidity.ValidityResponseErrors{
				"verificationPicture": []string{
					"verificationPicture is required",
				},
			}
		}
		if data.Reason != nil && *data.Reason != "" {
			return govalidity.ValidityResponseErrors{
				"reason": []string{
					"You don't need to provide reason",
				},
			}
		}
	case constants.SHIFT_VEHICLE_TYPE_PUBLIC_TRANSPORTATION:
		if data.StartLocation != nil {
			return govalidity.ValidityResponseErrors{
				"startLocation": []string{
					"You don't need to provide startLocation",
				},
			}
		}
		if data.Reason == nil || *data.Reason == "" {
			return govalidity.ValidityResponseErrors{
				"reason": []string{
					"reason is required",
				},
			}
		}
		if data.VerificationPicture == nil {
			return govalidity.ValidityResponseErrors{
				"verificationPicture": []string{
					"verificationPicture is required",
				},
			}
		}
	}

	return nil
}
