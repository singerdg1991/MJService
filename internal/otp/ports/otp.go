package ports

import (
	"github.com/hoitek/Maja-Service/internal/otp/domain"
	"github.com/hoitek/Maja-Service/internal/otp/models"
)

type OTPService interface {
	Send(otp *domain.OTP) error
	Verify(payload *models.OTPVerifyRequest) (*domain.OTP, error)
	Resend()
	GenerateOTPCode(length int) string
	Query(where string) ([]*domain.OTP, error)
	Count(where string) (int64, error)
	VerifyExchangeCode(code string, reason string, userID int64, userWorkPhoneNumber string, userEmail string) error
}

type OTPRepositoryPostgresDB interface {
	Query(where string) ([]*domain.OTP, error)
	Count(where string) (int64, error)
	Create(otp *domain.OTP) (int64, error)
	Delete(otp *domain.OTP) error
	Update(fields map[string]interface{}, where string) error
	VerifyExchangeCode(code string, reason string, userID int64, userWorkPhoneNumber string, userEmail string) error
}
