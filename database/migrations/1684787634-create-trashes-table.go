package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type TrashesMigration1684787634 struct {
}

func NewTrashesMigration1684787634() *TrashesMigration1684787634 {
	return &TrashesMigration1684787634{}
}

func (m *TrashesMigration1684787634) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS trashes (
			id SERIAL PRIMARY KEY,
			modelName VARCHAR(255) NOT NULL,
            modelId BIGINT NOT NULL,
            reason TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
            created_by BIGINT NOT NULL,
            CONSTRAINT trashes_model_modelId_unique UNIQUE (modelName, modelId),
            CONSTRAINT trashes_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
		);
		ALTER TABLE trashes ALTER COLUMN id SET DEFAULT nextval('trashes_id_seq'::regclass);
		ALTER TABLE trashes ALTER COLUMN created_at SET DEFAULT now();
    `)

	return nil
}

func (m *TrashesMigration1684787634) MigrateDown() error {
	database.PostgresDB.Exec(`
		DELETE FROM _migrations WHERE name = 'TrashesMigration1684787634'
	`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS trashes CASCADE;
    `)
	return nil
}
