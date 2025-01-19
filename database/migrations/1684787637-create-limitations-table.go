package migrations

import "github.com/hoitek/Maja-Service/database"

type LimitationsMigration1684787637 struct {
}

func NewLimitationsMigration1684787637() *LimitationsMigration1684787637 {
	return &LimitationsMigration1684787637{}
}

func (m *LimitationsMigration1684787637) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS limitations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE limitations ALTER COLUMN id SET DEFAULT nextval('limitations_id_seq'::regclass);
		ALTER TABLE limitations ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE limitations ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO limitations
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'limitation 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'limitation 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'limitation 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'limitation 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'limitation 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('limitations_id_seq', (SELECT MAX(id) FROM limitations));
    `)
	return nil
}

func (m *LimitationsMigration1684787637) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'LimitationsMigration1684787637'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS limitations CASCADE;
    `)
	return nil
}
