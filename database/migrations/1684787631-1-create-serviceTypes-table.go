package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type ServiceTypesMigration1684787647 struct {
}

func NewServiceTypesMigration1684787647() *ServiceTypesMigration1684787647 {
	return &ServiceTypesMigration1684787647{}
}

func (m *ServiceTypesMigration1684787647) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS serviceTypes (
			id SERIAL PRIMARY KEY,
			serviceId INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_services FOREIGN KEY (serviceId) REFERENCES services(id)
		);
		ALTER TABLE serviceTypes ALTER COLUMN id SET DEFAULT nextval('serviceTypes_id_seq'::regclass);
		ALTER TABLE serviceTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE serviceTypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO serviceTypes
			(id, serviceId, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 1, 'service type 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 1, 'service type 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 2, 'service type 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 3, 'service type 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 4, 'service type 05', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(6, 5, 'service type 06', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('serviceTypes_id_seq', (SELECT MAX(id) FROM serviceTypes));
    `)
	log.Println(err)
	return nil
}

func (m *ServiceTypesMigration1684787647) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ServiceTypesMigration1684787647';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS serviceTypes CASCADE;
    `)
	return nil
}
