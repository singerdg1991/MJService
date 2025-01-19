package migrations

import "github.com/hoitek/Maja-Service/database"

type GearTypesMigration1684787610 struct {
}

func NewGearTypesMigration1684787610() *GearTypesMigration1684787610 {
	return &GearTypesMigration1684787610{}
}

func (m *GearTypesMigration1684787610) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS geartypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE geartypes ALTER COLUMN id SET DEFAULT nextval('geartypes_id_seq'::regclass);
		ALTER TABLE geartypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE geartypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO geartypes
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'gearType 01', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'gearType 02', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'gearType 03', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'gearType 04', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'gearType 05', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('geartypes_id_seq', (SELECT MAX(id) FROM geartypes));
    `)
	return nil
}

func (m *GearTypesMigration1684787610) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'GearTypesMigration1684787610'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS geartypes CASCADE;
    `)
	return nil
}
