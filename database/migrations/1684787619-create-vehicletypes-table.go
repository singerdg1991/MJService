package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type VehicleTypesMigration1684787619 struct {
}

func NewVehicleTypesMigration1684787619() *VehicleTypesMigration1684787619 {
	return &VehicleTypesMigration1684787619{}
}

func (m *VehicleTypesMigration1684787619) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS vehicletypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE vehicletypes ALTER COLUMN id SET DEFAULT nextval('vehicletypes_id_seq'::regclass);
		ALTER TABLE vehicletypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE vehicletypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO vehicletypes
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'car', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'public transportation', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'bicycle', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('vehicletypes_id_seq', (SELECT MAX(id) FROM vehicletypes));
    `)
	return nil
}

func (m *VehicleTypesMigration1684787619) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'VehicleTypesMigration1684787619'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS vehicletypes CASCADE;
    `)
	return nil
}
