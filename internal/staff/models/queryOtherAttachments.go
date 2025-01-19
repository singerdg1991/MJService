package models

import (
	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsQueryOtherAttachmentsSortValue
 */
type StaffsQueryOtherAttachmentsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: StaffsQueryOtherAttachmentsSortType
 */
type StaffsQueryOtherAttachmentsSortType struct {
	Title     StaffsQueryOtherAttachmentsSortValue `json:"title,omitempty" openapi:"$ref:StaffsQueryOtherAttachmentsSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt StaffsQueryOtherAttachmentsSortValue `json:"created_at,omitempty" openapi:"$ref:StaffsQueryOtherAttachmentsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: StaffsQueryOtherAttachmentsFilterType
 */
type StaffsQueryOtherAttachmentsFilterType struct {
	Title filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"title\"}}"`
}

/*
 * @apiDefine: StaffsQueryOtherAttachmentsRequestParams
 */
type StaffsQueryOtherAttachmentsRequestParams struct {
	ID      int                                   `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int                                   `json:"userId,string,omitempty" openapi:"example:1;required"`
	Page    int                                   `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                                   `json:"limit,string,omitempty" openapi:"example:10"`
	Filters StaffsQueryOtherAttachmentsFilterType `json:"filters,omitempty" openapi:"$ref:StaffsQueryOtherAttachmentsFilterType;in:query"`
	Sorts   StaffsQueryOtherAttachmentsSortType   `json:"sorts,omitempty" openapi:"$ref:StaffsQueryOtherAttachmentsSortType;in:query"`
}

func (data *StaffsQueryOtherAttachmentsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
