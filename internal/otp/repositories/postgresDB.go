package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/otp"
	"github.com/hoitek/Maja-Service/internal/otp/domain"
	"log"
	"strconv"
)

type OTPRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewOTPRepositoryPostgresDB(d *sql.DB) *OTPRepositoryPostgresDB {
	return &OTPRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *OTPRepositoryPostgresDB) Query(where string) ([]*domain.OTP, error) {
	q := `
		SELECT *
		FROM otps
	`
	if len(where) > 0 {
		q += ` WHERE ` + where
	}
	q += `;`
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var otps []*domain.OTP
	for rows.Next() {
		var otp domain.OTP
		if err := rows.Scan(&otp.ID, &otp.UserID, &otp.Username, &otp.Password, &otp.Code, &otp.ExchangeCode, &otp.Ip, &otp.UserAgent, &otp.Type, &otp.IsUsed, &otp.IsVerified, &otp.Reason, &otp.ExpiredAt, &otp.CreatedAt, &otp.UpdatedAt, &otp.DeletedAt); err != nil {
			return nil, err
		}
		otps = append(otps, &otp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return otps, nil
}

func (r *OTPRepositoryPostgresDB) Create(otp *domain.OTP) (int64, error) {
	q := `
		INSERT INTO otps
		(userId, username, password, code, exchangeCode, ip, userAgent, otpType, reason, expired_at, created_at, updated_at, deleted_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id;
	`
	err := r.PostgresDB.QueryRow(q, otp.UserID, otp.Username, otp.Password, otp.Code, otp.ExchangeCode, otp.Ip, otp.UserAgent, otp.Type, otp.Reason, otp.ExpiredAt, otp.CreatedAt, otp.UpdatedAt, otp.DeletedAt).Scan(&otp.ID)
	if err != nil {
		return 0, err
	}

	return otp.ID, nil
}

func (r *OTPRepositoryPostgresDB) Update(fields map[string]interface{}, where string) error {
	setClause := ""
	index := 0
	args := make([]interface{}, 0)
	for field, value := range fields {
		index++
		setClause += field + "=$" + strconv.Itoa(index)
		args = append(args, value)
		if index != len(fields) {
			setClause += ","
		}
	}
	q := `
		UPDATE otps
		SET
		   ` + setClause + `
		WHERE ` + where + `;
	`
	log.Println("aaaaaaaaaaaaaaaaa", q)
	_, err := r.PostgresDB.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *OTPRepositoryPostgresDB) Delete(otp *domain.OTP) error {
	q := `
		DELETE FROM otps
		WHERE id = $1;
	`
	_, err := r.PostgresDB.Exec(q, otp.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OTPRepositoryPostgresDB) Count(where string) (int64, error) {
	q := `
		SELECT COUNT(*)
		FROM otps
	`
	if len(where) > 0 {
		q += ` WHERE ` + where
	}
	q += `;`
	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *OTPRepositoryPostgresDB) VerifyExchangeCode(code string, reason string, userID int64, userWorkPhoneNumber string, userEmail string) error {
	otps, err := r.Query(fmt.Sprintf("userId = '%d' AND exchangeCode = '%s' AND isVerified = '%t' AND isUsed = '%t' AND reason = '%s'", userID, code, true, true, reason))
	if err != nil {
		return err
	}
	if len(otps) == 0 {
		return errors.New("invalid exchange code")
	}
	otpIdentity := otps[0]
	if otpIdentity.Type == otp.TypePhone {
		if otpIdentity.Username != userWorkPhoneNumber {
			return errors.New("invalid exchange code")
		}
	} else if otpIdentity.Type == otp.TypeEmail {
		if otpIdentity.Username != userEmail {
			return errors.New("invalid exchange code")
		}
	} else {
		return errors.New("invalid exchange code")
	}
	return nil
}
