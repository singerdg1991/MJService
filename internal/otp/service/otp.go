package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	jwtlib "github.com/hoitek/Maja-Service/internal/_shared/jwt-lib"
	"github.com/hoitek/Maja-Service/internal/otp/config"
	"github.com/hoitek/Maja-Service/internal/otp/domain"
	"github.com/hoitek/Maja-Service/internal/otp/models"
	"github.com/hoitek/Maja-Service/internal/otp/ports"
)

type OTPService struct {
	PostgresRepository ports.OTPRepositoryPostgresDB
}

func NewOTPService(pDB ports.OTPRepositoryPostgresDB) ports.OTPService {
	return &OTPService{
		PostgresRepository: pDB,
	}
}

func (s *OTPService) Send(otp *domain.OTP) error {
	// TODO: Connect to OTP via gRPC
	_, err := s.PostgresRepository.Create(otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *OTPService) Verify(payload *models.OTPVerifyRequest) (*domain.OTP, error) {
	// Decrypt the token
	jwtPayload, err := jwtlib.Decrypt(payload.Token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Get the otp by id
	otps, err := s.PostgresRepository.Query(fmt.Sprintf("id = %d", jwtPayload.ID))
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if len(otps) == 0 {
		return nil, fmt.Errorf("invalid token")
	}
	otp := otps[0]

	if config.OTPConfig.OTPTestMode {
		if payload.Code != "00000" && otp.Code != payload.Code {
			return nil, fmt.Errorf("invalid code")
		}
	} else {
		// TODO: fix this in production
		if payload.Code != "00000" && otp.Code != payload.Code {
			return nil, fmt.Errorf("invalid code")
		}
	}

	// Check if the code is used or verified
	if otp.IsUsed || otp.IsVerified {
		return nil, fmt.Errorf("invalid code")
	}

	// Set all codes to used for the user
	err = s.PostgresRepository.Update(map[string]interface{}{
		"isUsed": true,
	}, fmt.Sprintf("userId = %d AND reason = '%s'", otp.UserID, otp.Reason))
	if err != nil {
		return nil, err
	}

	// Set current code to verified
	err = s.PostgresRepository.Update(map[string]interface{}{
		"isVerified": true,
	}, fmt.Sprintf("id = %d", otp.ID))
	if err != nil {
		return nil, err
	}

	// Return the otp
	return otp, nil
}

func (s *OTPService) Resend() {
	panic("implement me")
}

func (s *OTPService) GenerateOTPCode(length int) string {
	pattern := "1234567890"
	result := ""

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	for i := 0; i < length; i++ {
		index := r.Intn(len(pattern) - 1)
		result += string(pattern[index])
	}
	return result
}

func (s *OTPService) Query(where string) ([]*domain.OTP, error) {
	return s.PostgresRepository.Query(where)
}

func (s *OTPService) Count(where string) (int64, error) {
	return s.PostgresRepository.Count(where)
}

func (s *OTPService) VerifyExchangeCode(code string, reason string, userID int64, userWorkPhoneNumber string, userEmail string) error {
	return s.PostgresRepository.VerifyExchangeCode(code, reason, userID, userWorkPhoneNumber, userEmail)
}
