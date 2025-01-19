package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type LicensesMigration1684787638 struct {
}

func NewLicensesMigration1684787638() *LicensesMigration1684787638 {
	return &LicensesMigration1684787638{}
}

func (m *LicensesMigration1684787638) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS licenses (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE licenses ALTER COLUMN id SET DEFAULT nextval('licenses_id_seq'::regclass);
		ALTER TABLE licenses ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE licenses ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO licenses
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'license 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'license 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'license 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'license 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'license 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('licenses_id_seq', (SELECT MAX(id) FROM licenses));
    `)
	return nil
}

func (m *LicensesMigration1684787638) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'LicensesMigration1684787638'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS licenses CASCADE;
    `)
	return nil
}
