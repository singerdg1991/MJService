package models

import (
	"fmt"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/notification/constants"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: NotificationFilterType
 */
type NotificationFilterType struct {
	Title       filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	ReadAt      filters.FilterValue[string] `json:"read_at,omitempty" openapi:"$ref:FilterValueString;example:{\"read_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	IsForSystem filters.FilterValue[string] `json:"isForSystem,omitempty" openapi:"$ref:FilterValueString;example:{\"isForSystem\":{\"op\":\"equals\",\"value\":\"true\"}"`
}

/*
 * @apiDefine: NotificationSortValue
 */
type NotificationSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: NotificationSortType
 */
type NotificationSortType struct {
	Title       NotificationSortValue `json:"title,omitempty" openapi:"$ref:NotificationSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	Description NotificationSortValue `json:"description,omitempty" openapi:"$ref:NotificationSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	ReadAt      NotificationSortValue `json:"read_at,omitempty" openapi:"$ref:NotificationSortValue;example:{\"read_at\":{\"op\":\"asc\"}}"`
	CreatedAt   NotificationSortValue `json:"created_at,omitempty" openapi:"$ref:NotificationSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: NotificationsQueryRequestParams
 */
type NotificationsQueryRequestParams struct {
	ID      int                    `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int                    `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int                    `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                    `json:"limit,string,omitempty" openapi:"example:10"`
	Type    string                 `json:"type,omitempty" openapi:"example:all"`
	Filters NotificationFilterType `json:"filters,omitempty" openapi:"$ref:NotificationFilterType;in:query"`
	Sorts   NotificationSortType   `json:"sorts,omitempty" openapi:"$ref:NotificationSortType;in:query"`
}

func (data *NotificationsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"userId": govalidity.New("userId").Int().Optional(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
		"type":   govalidity.New("type"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"read_at": govalidity.Schema{
				"op":    govalidity.New("filter.read_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.read_at.value").Optional(),
			},
			"isForSystem": govalidity.Schema{
				"op":    govalidity.New("filter.isForSystem.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.isForSystem.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
			},
			"description": govalidity.Schema{
				"op": govalidity.New("sort.description.op"),
			},
			"read_at": govalidity.Schema{
				"op": govalidity.New("sort.read_at.op"),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Check type
	if data.Type == "" {
		data.Type = constants.NOTIFICATION_TYPE_ALL
	}
	if data.Type != constants.NOTIFICATION_TYPE_ALL && data.Type != constants.NOTIFICATION_TYPE_NOTIFICATION && data.Type != constants.NOTIFICATION_TYPE_REQUEST {
		return govalidity.ValidityResponseErrors{
			"type": []string{fmt.Sprintf("The type must be one of %s, %s, %s", constants.NOTIFICATION_TYPE_ALL, constants.NOTIFICATION_TYPE_NOTIFICATION, constants.NOTIFICATION_TYPE_REQUEST)},
		}
	}

	return nil
}
