package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffClubWarningsMigration1684787664 struct {
}

func NewStaffClubWarningsMigration1684787664() *StaffClubWarningsMigration1684787664 {
	return &StaffClubWarningsMigration1684787664{}
}

func (m *StaffClubWarningsMigration1684787664) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffClubWarnings (
			id SERIAL PRIMARY KEY,
			punishmentId INT NOT NULL,
			isAutoRewardSetEnable BOOLEAN NOT NULL DEFAULT FALSE,
			warningNumber INT NOT NULL,
			title VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffClubWarnings_punishmentId FOREIGN KEY (punishmentId) REFERENCES punishments(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffClubWarningsMigration1684787664) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffClubWarningsMigration1684787664'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffClubWarnings CASCADE;
    `)
	return nil
}
