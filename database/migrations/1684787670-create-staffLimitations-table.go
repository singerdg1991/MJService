package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffLimitationsMigration1684787669 struct {
}

func NewStaffLimitationsMigration1684787669() *StaffLimitationsMigration1684787669 {
	return &StaffLimitationsMigration1684787669{}
}

func (m *StaffLimitationsMigration1684787669) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffLimitations (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			limitationId INT NOT NULL,
			description TEXT DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_limitation_id FOREIGN KEY (limitationId) REFERENCES limitations(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffLimitationsMigration1684787669) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffLimitationsMigration1684787669'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffLimitations CASCADE;
    `)
	return nil
}
