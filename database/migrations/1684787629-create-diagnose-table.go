package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type DiagnosesMigration1684787629 struct {
}

func NewDiagnosesMigration1684787629() *DiagnosesMigration1684787629 {
	return &DiagnosesMigration1684787629{}
}

func (m *DiagnosesMigration1684787629) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS diagnoses (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
            code VARCHAR(255) NOT NULL,
            description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE diagnoses ALTER COLUMN id SET DEFAULT nextval('diagnoses_id_seq'::regclass);
		ALTER TABLE diagnoses ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE diagnoses ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('diagnoses_id_seq', (SELECT MAX(id) FROM diagnoses));
    `)

	return nil
}

func (m *DiagnosesMigration1684787629) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'DiagnosesMigration1684787629';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS diagnoses CASCADE;
    `)
	return nil
}
