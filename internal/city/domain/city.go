package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: City
 */
type City struct {
	ID        uint       `json:"id" openapi:"example:1"`
	Name      string     `json:"name" openapi:"example:John;required"`
	CreatedAt time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *City) TableName() string {
	return "cities"
}

func (u *City) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
