package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: KeikkalaUser
 */
type KeikkalaUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:John;required"`
	Email     string `json:"email" openapi:"example:sgh370@yahoo.com;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:John;required"`
}

/*
 * @apiDefine: KeikkalaRole
 */
type KeikkalaRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: KeikkalaSection
 */
type KeikkalaSection struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: Keikkala
 */
type Keikkala struct {
	ID              uint               `json:"id" openapi:"example:1"`
	RoleID          *uint              `json:"roleId" openapi:"example:1"`
	Role            *KeikkalaRole      `json:"role" openapi:"$ref:KeikkalaRole;required"`
	StartDate       time.Time          `json:"start_date" openapi:"example:2021-01-01"`
	EndDate         time.Time          `json:"end_date" openapi:"example:2021-01-01"`
	StartTime       time.Time          `json:"start_time" openapi:"example:00:00:00"`
	EndTime         time.Time          `json:"end_time" openapi:"example:00:00:00"`
	KaupunkiAddress *string            `json:"kaupunkiAddress" openapi:"example:address;required"`
	Sections        []*KeikkalaSection `json:"sections" openapi:"$ref:KeikkalaSection;type:array;required"`
	PaymentType     string             `json:"paymentType" openapi:"example:paySoon;required"`
	ShiftName       string             `json:"shiftName" openapi:"example:morning;required"`
	Description     *string            `json:"description" openapi:"example:John;required"`
	Status          string             `json:"status" openapi:"example:open;required"`
	PickedAt        *time.Time         `json:"picked_at" openapi:"example:2021-01-01T00:00:00Z"`
	PickedBy        *KeikkalaUser      `json:"pickedBy" openapi:"$ref:KeikkalaUser;required"`
	CreatedAt       time.Time          `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time          `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *time.Time         `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy       *KeikkalaUser      `json:"createdBy" openapi:"$ref:KeikkalaUser;required"`
	UpdatedBy       *KeikkalaUser      `json:"updatedBy" openapi:"$ref:KeikkalaUser;required"`
}

func (u *Keikkala) TableName() string {
	return "keikkalaShifts"
}

func (u *Keikkala) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
