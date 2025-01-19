package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: OTP
 */
type OTP struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Code         string    `json:"code"`
	ExchangeCode string    `json:"exchangeCode"`
	Ip           string    `json:"ip"`
	UserAgent    string    `json:"userAgent"`
	Type         string    `json:"type"`
	IsUsed       bool      `json:"isUsed"`
	IsVerified   bool      `json:"isVerified"`
	Reason       string    `json:"reason"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (u *OTP) TableName() string {
	return "otps"
}

func (u *OTP) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
