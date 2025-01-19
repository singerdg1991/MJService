package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: TicketMessagesFilterType
 */
type TicketMessagesFilterType struct {
}

/*
 * @apiDefine: TicketMessagesSortValue
 */
type TicketMessagesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: TicketMessagesSortType
 */
type TicketMessagesSortType struct {
	ID TicketMessagesSortValue `json:"id,omitempty" openapi:"$ref:TicketMessagesSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: TicketsQueryMessagesRequestParams
 */
type TicketsQueryMessagesRequestParams struct {
	ID       int                      `json:"id,string,omitempty" openapi:"example:1"`
	TicketID int                      `json:"ticketId,string,omitempty" openapi:"example:1"`
	Page     int                      `json:"page,string,omitempty" openapi:"example:1"`
	Limit    int                      `json:"limit,string,omitempty" openapi:"example:10"`
	Filters  TicketMessagesFilterType `json:"filters,omitempty" openapi:"$ref:TicketMessagesFilterType;in:query"`
	Sorts    TicketMessagesSortType   `json:"sorts,omitempty" openapi:"$ref:TicketMessagesSortType;in:query"`
}

func (data *TicketsQueryMessagesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Int().Optional(),
		"ticketId": govalidity.New("ticketId").Int().Optional(),
		"page":     govalidity.New("page").Int().Default("1"),
		"limit":    govalidity.New("limit").Int().Default("10"),
		"filters":  govalidity.Schema{},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
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
