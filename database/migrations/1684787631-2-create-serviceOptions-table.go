package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type ServiceOptionsMigration1684787648 struct {
}

func NewServiceOptionsMigration1684787648() *ServiceOptionsMigration1684787648 {
	return &ServiceOptionsMigration1684787648{}
}

func (m *ServiceOptionsMigration1684787648) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS serviceOptions (
			id SERIAL PRIMARY KEY,
			serviceTypeId INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_serviceTypes FOREIGN KEY (serviceTypeId) REFERENCES serviceTypes(id)
		);
		ALTER TABLE serviceOptions ALTER COLUMN id SET DEFAULT nextval('serviceOptions_id_seq'::regclass);
		ALTER TABLE serviceOptions ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE serviceOptions ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO serviceOptions
			(id, serviceTypeId, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 1, 'service option 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 1, 'service option 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 2, 'service option 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 3, 'service option 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 4, 'service option 05', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(6, 5, 'service option 06', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('serviceOptions_id_seq', (SELECT MAX(id) FROM serviceOptions));
    `)
	log.Println(err)
	return nil
}

func (m *ServiceOptionsMigration1684787648) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ServiceOptionsMigration1684787648';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS serviceOptions CASCADE;
    `)
	return nil
}
