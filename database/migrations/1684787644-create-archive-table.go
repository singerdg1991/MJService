package migrations

import "github.com/hoitek/Maja-Service/database"

type ArchivesMigration1684787644 struct {
}

func NewArchivesMigration1684787644() *ArchivesMigration1684787644 {
	return &ArchivesMigration1684787644{}
}

func (m *ArchivesMigration1684787644) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS archives (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
            subject VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			attachments JSONB NOT NULL DEFAULT '[]',
            datetime TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES users(id)
		);
		ALTER TABLE archives ALTER COLUMN id SET DEFAULT nextval('archives_id_seq'::regclass);
		ALTER TABLE archives ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE archives ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('archives_id_seq', (SELECT MAX(id) FROM archives));
    `)
	return nil
}

func (m *ArchivesMigration1684787644) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ArchivesMigration1684787644';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS archives CASCADE;
    `)
	return nil
}
