package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: HolidayUserRole
 */
type HolidayUserRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:name;required"`
}

/*
 * @apiDefine: HolidayUser
 */
type HolidayUser struct {
	ID        uint             `json:"id" openapi:"example:1"`
	FirstName string           `json:"firstName" openapi:"example:firstName;required"`
	LastName  string           `json:"lastName" openapi:"example:lastName;required"`
	Email     string           `json:"email" openapi:"example:email;required"`
	AvatarUrl string           `json:"avatarUrl" openapi:"example:avatarUrl;required"`
	Role      *HolidayUserRole `json:"role" openapi:"$ref:HolidayUserRole"`
}

/*
 * @apiDefine: Holiday
 */
type Holiday struct {
	ID             uint        `json:"id" openapi:"example:1"`
	StartDate      time.Time   `json:"start_date" openapi:"example:2021-01-01"`
	EndDate        time.Time   `json:"end_date" openapi:"example:2021-01-01"`
	Title          string      `json:"title" openapi:"example:title;required"`
	PaymentType    string      `json:"paymentType" openapi:"example:withSalary;required"`
	Description    *string     `json:"description" openapi:"example:description;required"`
	Status         string      `json:"status" openapi:"example:pending;required"`
	RejectedReason *string     `json:"rejectedReason" openapi:"example:rejectedReason"`
	AcceptedAt     *time.Time  `json:"accepted_at" openapi:"example:2021-01-01T00:00:00Z"`
	RejectedAt     *time.Time  `json:"rejected_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt      time.Time   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedByID    uint        `json:"-" openapi:"ignored"`
	UpdatedByID    uint        `json:"-" openapi:"ignored"`
	CreatedBy      HolidayUser `json:"createdBy" openapi:"$ref:HolidayUser"`
	UpdatedBy      HolidayUser `json:"updatedBy" openapi:"$ref:HolidayUser"`
}

func (u *Holiday) TableName() string {
	return "staffClubHolidays"
}

func (u *Holiday) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
