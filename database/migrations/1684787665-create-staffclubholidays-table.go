package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffClubHolidaysMigration1684787665 struct {
}

func NewStaffClubHolidaysMigration1684787665() *StaffClubHolidaysMigration1684787665 {
	return &StaffClubHolidaysMigration1684787665{}
}

func (m *StaffClubHolidaysMigration1684787665) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffClubHolidays (
			id SERIAL PRIMARY KEY,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			title VARCHAR(255) NOT NULL,
		    paymentType VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT NULL,
		    status VARCHAR(255) DEFAULT 'pending',
		    rejectedReason TEXT DEFAULT NULL,
		    accepted_at TIMESTAMP DEFAULT NULL,
		    rejected_at TIMESTAMP DEFAULT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    deleted_at TIMESTAMP DEFAULT NULL,
		    createdBy INTEGER NOT NULL,
		    updatedBy INTEGER NOT NULL,
		    CONSTRAINT staffClubHolidays_createdBy_foreign FOREIGN KEY (createdBy) REFERENCES users (id) ON DELETE CASCADE,
		    CONSTRAINT staffClubHolidays_updatedBy_foreign FOREIGN KEY (updatedBy) REFERENCES users (id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffClubHolidaysMigration1684787665) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffClubHolidaysMigration1684787665'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffClubHolidays CASCADE;
    `)
	return nil
}
