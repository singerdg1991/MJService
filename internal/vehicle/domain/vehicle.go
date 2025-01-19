package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: VehicleUser
 */
type VehicleUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: Vehicle
 */
type Vehicle struct {
	ID          uint         `json:"id" openapi:"example:1"`
	VehicleType string       `json:"vehicleType" openapi:"example:car"`
	UserID      uint         `json:"userId" openapi:"example:1"`
	User        *VehicleUser `json:"user" openapi:"$ref:VehicleUser"`
	Brand       string       `json:"brand" openapi:"example:1"`
	Model       string       `json:"model" openapi:"example:1"`
	Year        string       `json:"year" openapi:"example:1"`
	Variant     string       `json:"variant" openapi:"example:1"`
	FuelType    string       `json:"fuelType" openapi:"example:1"`
	VehicleNo   string       `json:"vehicleNo" openapi:"example:1"`
	CreatedAt   time.Time    `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time    `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time   `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Vehicle) TableName() string {
	return "vehicles"
}

func (u *Vehicle) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
