package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: TrashCreatedBy
 */
type TrashCreatedBy struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
}

/*
 * @apiDefine: Trash
 */
type Trash struct {
	ID        uint           `json:"id" openapi:"example:1"`
	ModelName     string         `json:"modelName" openapi:"example:User;required"`
	ModelID   uint           `json:"modelId" openapi:"example:1;required"`
	Reason    string         `json:"reason" openapi:"example:Deleted by mistake"`
	CreatedAt time.Time      `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy TrashCreatedBy `json:"createdBy" openapi:"example:{\"id\":1,\"firstName\":\"John\",\"lastName\":\"Doe\"}"`
}

func NewTrash() *Trash {
	return &Trash{}
}

func (u *Trash) TableName() string {
	return "trashes"
}

func (u *Trash) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
