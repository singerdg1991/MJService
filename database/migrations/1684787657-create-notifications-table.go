package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type NotificationsMigration1684787657 struct {
}

func NewNotificationsMigration1684787657() *NotificationsMigration1684787657 {
	return &NotificationsMigration1684787657{}
}

func (m *NotificationsMigration1684787657) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT NOT NULL,
            read_at TIMESTAMP DEFAULT NULL,
            readBy INT DEFAULT NULL,
            extra JSONB DEFAULT NULL,
            isForSystem BOOLEAN DEFAULT FALSE,
            status VARCHAR(255) DEFAULT NULL, -- if null, it means notification otherwise it is request
            status_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
            CONSTRAINT fk_readBy FOREIGN KEY (readBy) REFERENCES users(id) ON DELETE SET NULL,
            CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *NotificationsMigration1684787657) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'NotificationsMigration1684787657'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS notifications CASCADE;
    `)
	return nil
}
