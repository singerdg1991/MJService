package domain

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/*
 * @apiDefine: Section
 */
type Section struct {
	ID          uint               `json:"id" openapi:"example:1"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" json:"-" openapi:"example:5f7b5f5b9b9b9b9b9b9b9b9b"`
	Parent      *Section           `json:"parent" openapi:"$ref:Section"`
	ParentID    *int64             `json:"-" openapi:"example:1"`
	Children    []*Section         `json:"children" openapi:"$ref:Section;type:array;"`
	Name        string             `json:"name" openapi:"example:John;required"`
	Color       *string            `json:"color" openapi:"example:#000000"`
	Description *string            `json:"description" openapi:"example:Description"`
	CreatedAt   time.Time          `json:"-" openapi:"example:2022-01-01T00:00:00Z"`
	UpdatedAt   time.Time          `json:"-" openapi:"example:2023-01-01T00:00:00Z"`
	DeletedAt   *time.Time         `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewSection() *Section {
	return &Section{}
}

func (u *Section) TableName() string {
	return "sections"
}

func (u *Section) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
