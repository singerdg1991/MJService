package domain

import (
	"time"
)

/*
 * @apiDefine: CustomerSection
 */
type CustomerSection struct {
	ID          uint               `json:"id" openapi:"example:1"`
	Parent      *CustomerSection   `json:"parent" openapi:"$ref:CustomerSection"`
	ParentID    *int64             `json:"-" openapi:"example:1"`
	Children    []*CustomerSection `json:"children" openapi:"$ref:CustomerSection;type:array;"`
	Name        string             `json:"name" openapi:"example:John;required"`
	Color       *string            `json:"color" openapi:"example:#000000"`
	Description *string            `json:"description" openapi:"example:Description"`
	CreatedAt   time.Time          `json:"-" openapi:"example:2022-01-01T00:00:00Z"`
	UpdatedAt   time.Time          `json:"-" openapi:"example:2023-01-01T00:00:00Z"`
	DeletedAt   *time.Time         `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}
