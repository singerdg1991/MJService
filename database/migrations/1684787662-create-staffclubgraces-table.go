package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffClubGracesMigration1684787662 struct {
}

func NewStaffClubGracesMigration1684787662() *StaffClubGracesMigration1684787662 {
	return &StaffClubGracesMigration1684787662{}
}

func (m *StaffClubGracesMigration1684787662) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffClubGraces (
			id SERIAL PRIMARY KEY,
			rewardId INT NOT NULL,
			isAutoRewardSetEnable BOOLEAN NOT NULL DEFAULT FALSE,
			graceNumber INT NOT NULL,
			title VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffClubGraces_rewardId FOREIGN KEY (rewardId) REFERENCES rewards(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffClubGracesMigration1684787662) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffClubGracesMigration1684787662'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffClubGraces CASCADE;
    `)
	return nil
}
