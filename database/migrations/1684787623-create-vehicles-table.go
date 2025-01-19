package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type VehiclesMigration1684787623 struct {
}

func NewVehiclesMigration1684787623() *VehiclesMigration1684787623 {
	return &VehiclesMigration1684787623{}
}

func (m *VehiclesMigration1684787623) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS vehicles (
			id SERIAL PRIMARY KEY,
			vehicleType VARCHAR(255) NOT NULL,
			userId INT NOT NULL,
			brand VARCHAR(255) NOT NULL,
			model VARCHAR(255) NOT NULL,
			year VARCHAR(255) NOT NULL,
			variant VARCHAR(255) NOT NULL,
			fuelType VARCHAR(255) NOT NULL,
			vehicleNo VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_user_id FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
		ALTER TABLE vehicles ALTER COLUMN id SET DEFAULT nextval('vehicles_id_seq'::regclass);
		ALTER TABLE vehicles ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE vehicles ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE vehicles ALTER COLUMN deleted_at SET DEFAULT NULL;
		SELECT setval('vehicles_id_seq', (SELECT MAX(id) FROM vehicles));
    `)

	return nil
}

func (m *VehiclesMigration1684787623) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'VehiclesMigration1684787623'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS vehicles CASCADE;
    `)
	return nil
}
