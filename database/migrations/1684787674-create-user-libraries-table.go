package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type UserLibrariesMigration1684787674 struct {
}

func NewUserLibrariesMigration1684787674() *UserLibrariesMigration1684787674 {
	return &UserLibrariesMigration1684787674{}
}

func (m *UserLibrariesMigration1684787674) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS userLibraries (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			attachments JSONB DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_userLibraries_userId FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
    `)
	if err != nil {
		log.Println("************ Error in creating userLibraries table ************", err.Error())
	}
	return nil
}

func (m *UserLibrariesMigration1684787674) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'UserLibrariesMigration1684787674'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS userLibraries CASCADE;
    `)
	return nil
}
