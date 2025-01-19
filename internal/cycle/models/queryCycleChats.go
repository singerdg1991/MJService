package models

import (
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CycleChatsFilterType
 */
type CycleChatsFilterType struct {
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: CycleChatsSortValue
 */
type CycleChatsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CycleChatsSortType
 */
type CycleChatsSortType struct {
	ID        CycleChatsSortValue `json:"id,omitempty" openapi:"$ref:CycleChatsSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	CreatedAt CycleChatsSortValue `json:"created_at,omitempty" openapi:"$ref:CycleChatsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CyclesQueryChatsRequestParams
 */
type CyclesQueryChatsRequestParams struct {
	ID                 int                  `json:"id,string,omitempty" openapi:"example:1"`
	CycleID            int                  `json:"cycleid,string,omitempty" openapi:"example:1"`
	CyclePickupShiftID int                  `json:"cyclePickupShiftId,string,omitempty" openapi:"example:1"`
	SenderUserID       int                  `json:"senderUserId,string,omitempty" openapi:"example:1"`
	RecipientUserID    int                  `json:"recipientUserId,string,omitempty" openapi:"example:1"`
	Page               int                  `json:"page,string,omitempty" openapi:"example:1"`
	Limit              int                  `json:"limit,string,omitempty" openapi:"example:10"`
	Filters            CycleChatsFilterType `json:"filters,omitempty" openapi:"$ref:CycleChatsFilterType;in:query"`
	Sorts              CycleChatsSortType   `json:"sorts,omitempty" openapi:"$ref:CycleChatsSortType;in:query"`
}

// ValidateQueries validates the CyclesQueryChatsRequestParams based on the provided schema and request.
//
// It takes an http.Request as a parameter to validate the query parameters.
// It returns a govalidity.ValidityResponseErrors if the validation fails, otherwise nil.
func (data *CyclesQueryChatsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":                 govalidity.New("id").Int().Optional(),
		"cycleid":            govalidity.New("cycleid").Int().Optional(),
		"cyclePickupShiftId": govalidity.New("cyclePickupShiftId").Int().Optional(),
		"senderUserId":       govalidity.New("senderUserId").Int().Optional(),
		"recipientUserId":    govalidity.New("recipientUserId").Int().Optional(),
		"page":               govalidity.New("page").Int().Default("1"),
		"limit":              govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"created_at": govalidity.Schema{
				"op": govalidity.New("sort.created_at.op"),
			},
		},
	}

	// Validate queries
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
