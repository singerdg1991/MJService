package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type UserOtherAttachmentsMigration1684787673 struct {
}

func NewUserOtherAttachmentsMigration1684787673() *UserOtherAttachmentsMigration1684787673 {
	return &UserOtherAttachmentsMigration1684787673{}
}

func (m *UserOtherAttachmentsMigration1684787673) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS userOtherAttachments (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			attachments JSONB DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_userOtherAttachments_userId FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
    `)
	if err != nil {
		log.Println("************ Error in creating userOtherAttachments table ************", err.Error())
	}
	return nil
}

func (m *UserOtherAttachmentsMigration1684787673) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'UserOtherAttachmentsMigration1684787673'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS userOtherAttachments CASCADE;
    `)
	return nil
}
