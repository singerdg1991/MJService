package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/archive/domain"
)

/*
 * @apiDefine: ArchivesQueryResponseDataItem
 */
type ArchivesQueryResponseDataItem struct {
	ID          uint                   `json:"id" openapi:"example:1"`
	UserID      uint                   `json:"userId" openapi:"example:1"`
	User        domain.ArchiveUser     `json:"user" openapi:"$ref:ArchiveUser"`
	Title       string                 `json:"title" openapi:"example:title"`
	Subject     string                 `json:"subject" openapi:"example:subject"`
	Description *string                `json:"description" openapi:"example:description"`
	Attachments []types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;type:array"`
	Date        string                 `json:"date" openapi:"example:2021-01-01"`
	Time        string                 `json:"time" openapi:"example:00:00:00"`
}

/*
 * @apiDefine: ArchivesQueryResponseData
 */
type ArchivesQueryResponseData struct {
	Limit      int                             `json:"limit" openapi:"example:10"`
	Offset     int                             `json:"offset" openapi:"example:0"`
	Page       int                             `json:"page" openapi:"example:1"`
	TotalRows  int                             `json:"totalRows" openapi:"example:1"`
	TotalPages int                             `json:"totalPages" openapi:"example:1"`
	Items      []ArchivesQueryResponseDataItem `json:"items" openapi:"$ref:ArchivesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ArchivesQueryResponse
 */
type ArchivesQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       ArchivesQueryResponseData `json:"data" openapi:"$ref:ArchivesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ArchivesQueryNotFoundResponse
 */
type ArchivesQueryNotFoundResponse struct {
	Archives []domain.Archive `json:"archives" openapi:"$ref:Archive;type:array"`
}
