package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: ServiceTypeService
 */
type ServiceTypeService struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServiceType
 */
type ServiceType struct {
	ID          uint                `json:"id" openapi:"example:1"`
	ServiceID   uint                `json:"serviceId" openapi:"example:1"`
	Service     *ServiceTypeService `json:"service" openapi:"$ref:ServiceTypeService"`
	Name        string              `json:"name" openapi:"example:John;required"`
	Description string              `json:"description" openapi:"example:John;required"`
	CreatedAt   time.Time           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *ServiceType) TableName() string {
	return "serviceTypes"
}

func (u *ServiceType) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
