package repositories

import (
	"github.com/hoitek/Maja-Service/internal/otp/domain"
)

type OTPRepositoryStub struct {
	OTPs []*domain.OTP
}

type otpTestCondition struct {
	HasError bool
}

var UserTestCondition *otpTestCondition = &otpTestCondition{}

func NewOTPRepositoryStub() *OTPRepositoryStub {
	return &OTPRepositoryStub{
		OTPs: []*domain.OTP{
			{
				ID: 1,
			},
		},
	}
}

func (r *OTPRepositoryStub) Query(where string) ([]*domain.OTP, error) {
	panic("implement me")
}

func (r *OTPRepositoryStub) Count(where string) (int64, error) {
	panic("implement me")
}

func (r *OTPRepositoryStub) Create(otp *domain.OTP) (int64, error) {
	panic("implement me")
}

func (r *OTPRepositoryStub) Delete(otp *domain.OTP) error {
	panic("implement me")
}

func (r *OTPRepositoryStub) Update(fields map[string]interface{}, where string) error {
	panic("implement me")
}

func (r *OTPRepositoryStub) VerifyExchangeCode(code string, reason string, userID int64, userWorkPhoneNumber string, userEmail string) error {
	panic("implement me")
}
