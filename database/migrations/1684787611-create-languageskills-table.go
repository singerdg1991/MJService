package migrations

import "github.com/hoitek/Maja-Service/database"

type LanguageSkillsMigration1684787611 struct {
}

func NewLanguageSkillsMigration1684787611() *LanguageSkillsMigration1684787611 {
	return &LanguageSkillsMigration1684787611{}
}

func (m *LanguageSkillsMigration1684787611) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS languageskills (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE languageskills ALTER COLUMN id SET DEFAULT nextval('languageskills_id_seq'::regclass);
		ALTER TABLE languageskills ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE languageskills ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO languageskills
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'English', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'Spanish', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'French', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'German', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'Italian', '2020-01-01', '2020-01-01', '2020-01-01'),
			(6, 'Japanese', '2020-01-01', '2020-01-01', '2020-01-01'),
			(7, 'Chinese', '2020-01-01', '2020-01-01', '2020-01-01'),
			(8, 'Russian', '2020-01-01', '2020-01-01', '2020-01-01'),
			(9, 'Arabic', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('languageskills_id_seq', (SELECT MAX(id) FROM languageskills));
    `)
	return nil
}

func (m *LanguageSkillsMigration1684787611) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'LanguageSkillsMigration1684787611'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS languageskills CASCADE;
    `)
	return nil
}
