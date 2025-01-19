package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type AbilitiesMigration1684787620 struct {
}

func NewAbilitiesMigration1684787620() *AbilitiesMigration1684787620 {
	return &AbilitiesMigration1684787620{}
}

func (m *AbilitiesMigration1684787620) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS abilities (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE abilities ALTER COLUMN id SET DEFAULT nextval('abilities_id_seq'::regclass);
		ALTER TABLE abilities ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE abilities ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO abilities
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'ability 01', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'ability 02', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'ability 03', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'ability 04', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'ability 05', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('abilities_id_seq', (SELECT MAX(id) FROM abilities));
    `)
	return nil
}

func (m *AbilitiesMigration1684787620) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'AbilitiesMigration1684787620'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS abilities CASCADE;
    `)
	return nil
}
