package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type CyclesMigration1684787632 struct {
}

func NewCyclesMigration1684787632() *CyclesMigration1684787632 {
	return &CyclesMigration1684787632{}
}

func (m *CyclesMigration1684787632) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE cycles ALTER COLUMN id SET DEFAULT nextval('cycles_id_seq'::regclass);
		ALTER TABLE cycles ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycles ALTER COLUMN updated_at SET DEFAULT now();
    `)
	return nil
}

func (m *CyclesMigration1684787632) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CyclesMigration1684787632'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycles CASCADE;
    `)
	return nil
}
