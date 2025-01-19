package models

import (
	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersQueryOtherAttachmentsSortValue
 */
type CustomersQueryOtherAttachmentsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryOtherAttachmentsSortType
 */
type CustomersQueryOtherAttachmentsSortType struct {
	Title     CustomersQueryOtherAttachmentsSortValue `json:"title,omitempty" openapi:"$ref:CustomersQueryOtherAttachmentsSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt CustomersQueryOtherAttachmentsSortValue `json:"created_at,omitempty" openapi:"$ref:CustomersQueryOtherAttachmentsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryOtherAttachmentsFilterType
 */
type CustomersQueryOtherAttachmentsFilterType struct {
	Title filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"title\"}}"`
}

/*
 * @apiDefine: CustomersQueryOtherAttachmentsRequestParams
 */
type CustomersQueryOtherAttachmentsRequestParams struct {
	ID      int                                      `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int                                      `json:"userId,string,omitempty" openapi:"example:1;required"`
	Page    int                                      `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                                      `json:"limit,string,omitempty" openapi:"example:10"`
	Filters CustomersQueryOtherAttachmentsFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryOtherAttachmentsFilterType;in:query"`
	Sorts   CustomersQueryOtherAttachmentsSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryOtherAttachmentsSortType;in:query"`
}

func (data *CustomersQueryOtherAttachmentsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"userId": govalidity.New("userId").Int().Required(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
			},
			"created_at": govalidity.Schema{
				"op": govalidity.New("sort.created_at.op"),
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
