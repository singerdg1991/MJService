package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type CompaniesMigration1684787618 struct {
}

func NewCompaniesMigration1684787618() *CompaniesMigration1684787618 {
	return &CompaniesMigration1684787618{}
}

func (m *CompaniesMigration1684787618) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS companies (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE companies ALTER COLUMN id SET DEFAULT nextval('companies_id_seq'::regclass);
		ALTER TABLE companies ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE companies ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO companies
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'company 01', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'company 02', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'company 03', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'company 04', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'company 05', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('companies_id_seq', (SELECT MAX(id) FROM companies));
    `)
	return nil
}

func (m *CompaniesMigration1684787618) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CompaniesMigration1684787618'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS companies CASCADE;
    `)
	return nil
}
