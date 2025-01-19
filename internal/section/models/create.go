package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: SectionsCreateRequestBody
 */
type SectionsCreateRequestBody struct {
	ParentID    *int    `json:"parentId,omitempty" openapi:"example:1"`
	Color       *string `json:"color,omitempty" openapi:"example:#000000"`
	Description *string `json:"description,omitempty" openapi:"example:description"`
	Name        string  `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *SectionsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"parentId":    govalidity.New("parentId").Int().Optional(),
		"color":       govalidity.New("color").HexColor().Optional(),
		"description": govalidity.New("description").MaxLength(1000).Optional(),
		"name":        govalidity.New("name").MinMaxLength(3, 55).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
