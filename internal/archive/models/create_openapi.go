package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/archive/domain"
)

/*
 * @apiDefine: ArchivesResponseData
 */
type ArchivesResponseData struct {
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
 * @apiDefine: ArchivesCreateResponse
 */
type ArchivesCreateResponse struct {
	StatusCode int                  `json:"statusCode" openapi:"example:200"`
	Data       ArchivesResponseData `json:"data" openapi:"$ref:ArchivesResponseData"`
}
