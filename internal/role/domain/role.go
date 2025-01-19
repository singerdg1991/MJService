package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: RolePermission
 */
type RolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:John;required"`
	Title string `json:"title" openapi:"example:John;required"`
}

/*
 * @apiDefine: Role
 */
type Role struct {
	ID          uint              `json:"id" openapi:"example:1"`
	Name        string            `json:"name" openapi:"example:John;required"`
	Type        string            `json:"type" openapi:"example:core;required"`
	Permissions []*RolePermission `json:"permissions" openapi:"$ref:RolePermission;type:array;required"`
	CreatedAt   time.Time         `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time         `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time        `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Role) TableName() string {
	return "_roles"
}

func (u *Role) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
