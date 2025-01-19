package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/todo/constants"
	"log"
	"net/http"
)

/*
 * @apiDefine: TodosUpdateStatusRequestParams
 */
type TodosUpdateStatusRequestParams struct {
	ID     int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
	UserID int `json:"userid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *TodosUpdateStatusRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Required(),
		"userid": govalidity.New("userid").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}
	log.Printf("data: %+v", data)
	return nil
}

/*
 * @apiDefine: TodosUpdateStatusRequestBody
 */
type TodosUpdateStatusRequestBody struct {
	Status string `json:"status" openapi:"example:done;required;maxLen:100;minLen:2;"`
}

func (data *TodosUpdateStatusRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"status": govalidity.New("status").In([]string{constants.TODO_STATUS_ACTIVE, constants.TODO_STATUS_DONE}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
