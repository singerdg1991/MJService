package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type EmailsMigration1684787667 struct {
}

func NewEmailsMigration1684787667() *EmailsMigration1684787667 {
	return &EmailsMigration1684787667{}
}

func (m *EmailsMigration1684787667) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS emails (
			id SERIAL PRIMARY KEY,
		    senderId INT DEFAULT NULL,
			email VARCHAR(255) NOT NULL,
		    cc JSONB NOT NULL DEFAULT '[]',
		    bcc JSONB NOT NULL DEFAULT '[]',
		    title VARCHAR(255) NOT NULL,
		    subject VARCHAR(255) NOT NULL,
		    message TEXT NOT NULL,
		    attachments JSONB NOT NULL DEFAULT '[]',
		    category VARCHAR(255) NOT NULL DEFAULT 'outbox',
		    starred_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_senderId FOREIGN KEY (senderId) REFERENCES users(id) ON DELETE SET NULL
		);
    `)
	return nil
}

func (m *EmailsMigration1684787667) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'EmailsMigration1684787667'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS emails CASCADE;
    `)
	return nil
}
