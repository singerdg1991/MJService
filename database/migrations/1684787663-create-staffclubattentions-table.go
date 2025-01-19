package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffClubAttentionsMigration1684787663 struct {
}

func NewStaffClubAttentionsMigration1684787663() *StaffClubAttentionsMigration1684787663 {
	return &StaffClubAttentionsMigration1684787663{}
}

func (m *StaffClubAttentionsMigration1684787663) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffClubAttentions (
			id SERIAL PRIMARY KEY,
			punishmentId INT NOT NULL,
			isAutoRewardSetEnable BOOLEAN NOT NULL DEFAULT FALSE,
			attentionNumber INT NOT NULL,
			title VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffClubAttentions_punishmentId FOREIGN KEY (punishmentId) REFERENCES punishments(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffClubAttentionsMigration1684787663) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffClubAttentionsMigration1684787663'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffClubAttentions CASCADE;
    `)
	return nil
}
