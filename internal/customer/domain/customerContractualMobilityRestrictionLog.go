package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: CustomerContractualMobilityRestrictionLogCustomer
 */
type CustomerContractualMobilityRestrictionLogCustomer struct {
	ID        int64  `json:"id" openapi:"example:1"`
	UserID    int64  `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CustomerContractualMobilityRestrictionLogCreatedBy
 */
type CustomerContractualMobilityRestrictionLogCreatedBy struct {
	ID        int64  `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CustomerContractualMobilityRestrictionLog
 */
type CustomerContractualMobilityRestrictionLog struct {
	ID          int64                                               `json:"id" openapi:"example:1"`
	CustomerID  int64                                               `json:"customerId" openapi:"example:1"`
	BeforeValue *string                                             `json:"beforeValue" openapi:"example:BeforeValue"`
	AfterValue  *string                                             `json:"afterValue" openapi:"example:AfterValue"`
	Customer    *CustomerContractualMobilityRestrictionLogCustomer  `json:"customer" openapi:"$ref:CustomerContractualMobilityRestrictionLogCustomer"`
	CreatedBy   *CustomerContractualMobilityRestrictionLogCreatedBy `json:"createdBy" openapi:"$ref:CustomerContractualMobilityRestrictionLogCreatedBy"`
	CreatedAt   time.Time                                           `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                                           `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                                          `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerContractualMobilityRestrictionLog() *CustomerContractualMobilityRestrictionLog {
	return &CustomerContractualMobilityRestrictionLog{}
}

func (ns *CustomerContractualMobilityRestrictionLog) TableName() string {
	return "customerContractualMobilityRestrictionLogs"
}

func (ns *CustomerContractualMobilityRestrictionLog) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
