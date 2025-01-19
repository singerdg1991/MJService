package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

type StaffChatsMigration1686787689 struct {
}

func NewStaffChatsMigration1686787689() *StaffChatsMigration1686787689 {
	return &StaffChatsMigration1686787689{}
}

func (m *StaffChatsMigration1686787689) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS staffChats (
			id SERIAL PRIMARY KEY,
			senderUserId INT NOT NULL,
			recipientUserId INT NOT NULL,
			isSystem BOOLEAN DEFAULT FALSE,
			message TEXT DEFAULT NULL,
			attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffChats_senderUserId FOREIGN KEY (senderUserId) REFERENCES users(id),
			CONSTRAINT fk_staffChats_recipientUserId FOREIGN KEY (recipientUserId) REFERENCES users(id)
		);
		ALTER TABLE staffChats ALTER COLUMN id SET DEFAULT nextval('staffChats_id_seq'::regclass);
		ALTER TABLE staffChats ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffChats ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('staffChats_id_seq', (SELECT MAX(id) FROM staffChats));
    `)
	log.Println(err)
	return nil
}

func (m *StaffChatsMigration1686787689) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffChatsMigration1686787689';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration StaffChatsMigration1686787689: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS staffChats;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete staffChats: %v\n", err)
	}
	return nil
}
