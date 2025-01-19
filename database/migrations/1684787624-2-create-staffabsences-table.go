package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type StaffAbsencesMigration1684787624 struct {
}

func NewStaffAbsencesMigration1684787624() *StaffAbsencesMigration1684787624 {
	return &StaffAbsencesMigration1684787624{}
}

func (m *StaffAbsencesMigration1684787624) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffAbsences (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP DEFAULT NULL,
			reason TEXT DEFAULT NULL,
			attachments JSONB DEFAULT NULL,
			status VARCHAR(255) DEFAULT 'pending',
		    statusBy INT DEFAULT NULL,
		    status_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_status_by FOREIGN KEY (statusBy) REFERENCES users(id) ON DELETE SET NULL
		);
		ALTER TABLE staffAbsences ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffAbsences ALTER COLUMN updated_at SET DEFAULT now();
    `)
	if err != nil {
		log.Println("****************************************", err)
	}

	return nil
}

func (m *StaffAbsencesMigration1684787624) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffAbsencesMigration1684787624'`)
	_, err := database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffAbsences CASCADE;
    `)
	if err != nil {
		log.Println("****************************************", err)
	}
	return nil
}
