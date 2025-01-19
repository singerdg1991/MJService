package migrations

import "github.com/hoitek/Maja-Service/database"

type ServicesMigration1684787646 struct {
}

func NewServicesMigration1684787646() *ServicesMigration1684787646 {
	return &ServicesMigration1684787646{}
}

func (m *ServicesMigration1684787646) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS services (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE services ALTER COLUMN id SET DEFAULT nextval('services_id_seq'::regclass);
		ALTER TABLE services ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE services ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO services
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'service 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'service 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'service 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'service 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'service 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('services_id_seq', (SELECT MAX(id) FROM services));
    `)
	return nil
}

func (m *ServicesMigration1684787646) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ServicesMigration1684787646'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS services CASCADE;
    `)
	return nil
}
