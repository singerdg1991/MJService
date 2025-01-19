package models

import "github.com/hoitek/Maja-Service/internal/report/domain"

/*
 * @apiDefine: ReportsQueryVisitsTableResponseDataItem
 */
type ReportsQueryVisitsTableResponseDataItem struct {
	ID                      uint                          `json:"id" openapi:"example:1"`
	CycleID                 uint                          `json:"cycleId" openapi:"example:1"`
	Staff                   *domain.ReportVisitTableStaff `json:"staff" openapi:"$ref:ReportVisitTableStaff"`
	Status                  string                        `json:"status" openapi:"example:not-started"`
	PrevStatus              string                        `json:"prevStatus" openapi:"example:not-started"`
	StartKilometer          *string                       `json:"startKilometer" openapi:"example:0"`
	ReasonOfTheCancellation *string                       `json:"reasonOfTheCancellation" openapi:"example:reason"`
	ReasonOfTheReactivation *string                       `json:"reasonOfTheReactivation" openapi:"example:reason"`
	ReasonOfTheResume       *string                       `json:"reasonOfTheResume" openapi:"example:reason"`
	ReasonOfThePause        *string                       `json:"reasonOfThePause" openapi:"example:reason"`
	IsUnplanned             bool                          `json:"isUnplanned" openapi:"example:false"`
	DateTime                string                        `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt               string                        `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt               string                        `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt               *string                       `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	StartedAt               *string                       `json:"started_at" openapi:"example:2021-01-01T00:00:00Z"`
	EndedAt                 *string                       `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CancelledAt             *string                       `json:"cancelled_at" openapi:"example:2021-01-01T00:00:00Z"`
	DelayedAt               *string                       `json:"delayed_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: ReportsQueryVisitsTableResponseData
 */
type ReportsQueryVisitsTableResponseData struct {
	Limit      int                                       `json:"limit" openapi:"example:10"`
	Offset     int                                       `json:"offset" openapi:"example:0"`
	Page       int                                       `json:"page" openapi:"example:1"`
	TotalRows  int                                       `json:"totalRows" openapi:"example:1"`
	TotalPages int                                       `json:"totalPages" openapi:"example:1"`
	Items      []ReportsQueryVisitsTableResponseDataItem `json:"items" openapi:"$ref:ReportsQueryVisitsTableResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryVisitsTableResponse
 */
type ReportsQueryVisitsTableResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryVisitsTableResponseData `json:"data" openapi:"$ref:ReportsQueryVisitsTableResponseData"`
}

/*
 * @apiDefine: ReportsQueryVisitsTableNotFoundResponse
 */
type ReportsQueryVisitsTableNotFoundResponse struct {
	ReportsQueryVisitsTable []domain.ReportVisitTable `json:"reportsQueryVisitsTable" openapi:"$ref:ReportVisitTable;type:array"`
}
