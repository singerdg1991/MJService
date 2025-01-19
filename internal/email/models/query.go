package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: EmailFilterType
 */
type EmailFilterType struct {
	Title     filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Email     filters.FilterValue[string] `json:"email,omitempty" openapi:"$ref:FilterValueString;example:{\"email\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	SenderID  filters.FilterValue[int]    `json:"senderId,omitempty" openapi:"$ref:FilterValueInt;example:{\"senderId\":{\"op\":\"equals\",\"value\":1}"`
	Category  filters.FilterValue[string] `json:"category,omitempty" openapi:"$ref:FilterValueString;example:{\"category\":{\"op\":\"equals\",\"value\":\"outbox\"}"`
	StarredAt filters.FilterValue[string] `json:"starred_at,omitempty" openapi:"$ref:FilterValueString;example:{\"starred_at\":{\"op\":\"equals\",\"value\":\"2021-01-01T00:00:00Z\"}"`
}

/*
 * @apiDefine: EmailSortValue
 */
type EmailSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: EmailSortType
 */
type EmailSortType struct {
	ID        EmailSortValue `json:"id,omitempty" openapi:"$ref:EmailSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title     EmailSortValue `json:"title,omitempty" openapi:"$ref:EmailSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt EmailSortValue `json:"created_at,omitempty" openapi:"$ref:EmailSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: EmailsQueryRequestParams
 */
type EmailsQueryRequestParams struct {
	ID      int             `json:"id,string,omitempty" openapi:"example:1"`
	Page    int             `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters EmailFilterType `json:"filters,omitempty" openapi:"$ref:EmailFilterType;in:query"`
	Sorts   EmailSortType   `json:"sorts,omitempty" openapi:"$ref:EmailSortType;in:query"`
}

func (data *EmailsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"email": govalidity.Schema{
				"op":    govalidity.New("filter.email.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.email.value").Optional(),
			},
			"senderId": govalidity.Schema{
				"op":    govalidity.New("filter.senderId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.senderId.value").Optional(),
			},
			"category": govalidity.Schema{
				"op":    govalidity.New("filter.category.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category.value").Optional(),
			},
			"starred_at": govalidity.Schema{
				"op":    govalidity.New("filter.starred_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.starred_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
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
