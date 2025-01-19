package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/trash"
)

type TrashesCreateRequestBodyCreatedBy struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
}

/*
 * @apiDefine: TrashesCreateRequestBody
 */
type TrashesCreateRequestBody struct {
	ModelName string                            `json:"modelName" openapi:"type:string;required;"`
	ModelID   uint                              `json:"modelId" openapi:"type:integer;required;"`
	Reason    string                            `json:"reason" openapi:"type:string;required;"`
	CreatedBy TrashesCreateRequestBodyCreatedBy `json:"createdBy" openapi:"ignored"`
}

func (data *TrashesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	var models []string
	for model := range trash.TrashModels {
		models = append(models, model)
	}
	schema := govalidity.Schema{
		"modelName": govalidity.New("modelName").MinMaxLength(3, 25).In(models).Required(),
		"modelId":   govalidity.New("modelId").Int().Min(1).Required(),
		"reason":    govalidity.New("reason").MinMaxLength(3, 255).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
