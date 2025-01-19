package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: ServiceOptionServiceTypeService
 */
type ServiceOptionServiceTypeService struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServiceOptionServiceType
 */
type ServiceOptionServiceType struct {
	ID          uint                            `json:"id" openapi:"example:1"`
	ServiceID   uint                            `json:"serviceId" openapi:"example:1"`
	Service     ServiceOptionServiceTypeService `json:"service" openapi:"$ref:ServiceOptionServiceTypeService"`
	Name        string                          `json:"name" openapi:"example:John;required"`
	Description string                          `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServiceOption
 */
type ServiceOption struct {
	ID            uint                      `json:"id" openapi:"example:1"`
	ServiceTypeID uint                      `json:"serviceTypeId" openapi:"example:1"`
	ServiceType   *ServiceOptionServiceType `json:"serviceType" openapi:"$ref:ServiceOptionServiceType"`
	Name          string                    `json:"name" openapi:"example:John;required"`
	Description   string                    `json:"description" openapi:"example:John;required"`
	CreatedAt     time.Time                 `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time                 `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time                `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *ServiceOption) TableName() string {
	return "serviceOptions"
}

func (u *ServiceOption) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
