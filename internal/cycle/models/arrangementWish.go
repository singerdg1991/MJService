package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesArrangementWishRequestBody
 */
type CyclesArrangementWishRequestBody struct {
	CycleID int64         `json:"cycleId" openapi:"example:1"`
	Cycle   *domain.Cycle `json:"cycle" openapi:"ignored"`
}

// ValidateBody validates the request body of CyclesArrangementWishRequestBody.
//
// It takes an http.Request as a parameter and returns a govalidity.ValidityResponseErrors.
func (data *CyclesArrangementWishRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId": govalidity.New("cycleId").Int().Min(1).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
