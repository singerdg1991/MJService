package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type OTPsMigration1684787633 struct {
}

func NewOTPsMigration1684787633() *OTPsMigration1684787633 {
	return &OTPsMigration1684787633{}
}

func (m *OTPsMigration1684787633) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			userId INTEGER NOT NULL,
			username VARCHAR(256) NOT NULL,
		    password VARCHAR(256) NOT NULL,
			code VARCHAR(256) NOT NULL,
			exchangeCode VARCHAR(256) NOT NULL,
			ip VARCHAR(256) NOT NULL,
			userAgent VARCHAR(256) NOT NULL,
			otpType VARCHAR(256) NOT NULL,
			isUsed BOOLEAN NOT NULL DEFAULT FALSE,
		    isVerified BOOLEAN NOT NULL DEFAULT FALSE,
		    reason VARCHAR(256) NOT NULL,
			expired_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_user_id FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE SEQUENCE IF NOT EXISTS otps_id_seq START 1;
		ALTER TABLE otps ALTER COLUMN id SET DEFAULT nextval('otps_id_seq');
		ALTER TABLE otps ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE otps ALTER COLUMN updated_at SET DEFAULT now();
    `)
	return nil
}

func (m *OTPsMigration1684787633) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'OTPsMigration1684787633'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS otps CASCADE;
    `)
	return nil
}
