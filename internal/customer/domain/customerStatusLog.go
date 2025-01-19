package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: CustomerStatusLogCustomer
 */
type CustomerStatusLogCustomer struct {
	ID        int64  `json:"id" openapi:"example:1"`
	UserID    int64  `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CustomerStatusLogCreatedBy
 */
type CustomerStatusLogCreatedBy struct {
	ID        int64  `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CustomerStatusLog
 */
type CustomerStatusLog struct {
	ID          int64                       `json:"id" openapi:"example:1"`
	CustomerID  int64                       `json:"customerId" openapi:"example:1"`
	StatusValue string                      `json:"statusValue" openapi:"example:StatusValue"`
	Customer    *CustomerStatusLogCustomer  `json:"customer" openapi:"$ref:CustomerStatusLogCustomer"`
	CreatedBy   *CustomerStatusLogCreatedBy `json:"createdBy" openapi:"$ref:CustomerStatusLogCreatedBy"`
	CreatedAt   time.Time                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerStatusLog() *CustomerStatusLog {
	return &CustomerStatusLog{}
}

func (ns *CustomerStatusLog) TableName() string {
	return "customerStatusLogs"
}

func (ns *CustomerStatusLog) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
