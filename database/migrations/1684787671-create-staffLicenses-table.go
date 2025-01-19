package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type StaffLicensesMigration1684787671 struct {
}

func NewStaffLicensesMigration1684787671() *StaffLicensesMigration1684787671 {
	return &StaffLicensesMigration1684787671{}
}

func (m *StaffLicensesMigration1684787671) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS staffLicenses (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			licenseId INT NOT NULL,
			expire_date TIMESTAMP DEFAULT NULL,
			attachments JSONB DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_license_id FOREIGN KEY (licenseId) REFERENCES licenses(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *StaffLicensesMigration1684787671) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffLicensesMigration1684787671'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS staffLicenses CASCADE;
    `)
	return nil
}
