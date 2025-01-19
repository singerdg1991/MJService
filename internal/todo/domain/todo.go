package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: TodoUser
 */
type TodoUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName;required"`
	LastName  string `json:"lastName" openapi:"example:lastName;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:avatarUrl;required"`
}

/*
 * @apiDefine: Todo
 */
type Todo struct {
	ID            uint       `json:"id" openapi:"example:1"`
	UserID        uint       `json:"userId" openapi:"example:1"`
	Title         string     `json:"title" openapi:"example:title;required"`
	Date          time.Time  `json:"-" openapi:"ignored"`
	Time          time.Time  `json:"-" openapi:"ignored"`
	DateStr       string     `json:"date" openapi:"example:2021-01-01;required"`
	TimeStr       string     `json:"time" openapi:"example:00:00;required"`
	User          *TodoUser  `json:"user" openapi:"$ref:TodoUser"`
	Description   *string    `json:"description" openapi:"example:description;required"`
	Status        string     `json:"status" openapi:"example:active;required"`
	DoneAt        *time.Time `json:"done_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt     time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy     uint       `json:"-" openapi:"example:1"`
	CreatedByUser *TodoUser  `json:"createdBy" openapi:"$ref:TodoUser"`
}

func (u *Todo) TableName() string {
	return "todos"
}

func (u *Todo) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
