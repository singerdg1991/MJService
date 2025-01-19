package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

type StaffChatMessagesMigration1687787690 struct {
}

func NewStaffChatMessagesMigration1687787690() *StaffChatMessagesMigration1687787690 {
	return &StaffChatMessagesMigration1687787690{}
}

func (m *StaffChatMessagesMigration1687787690) MigrateUp() error {
	// messageType -- text, image, video, audio
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS staffChatMessages (
			id SERIAL PRIMARY KEY,
			staffChatId INT NOT NULL,
			senderUserId INT NOT NULL,
			recipientUserId INT NOT NULL,
			isSystem BOOLEAN DEFAULT FALSE,
			message TEXT DEFAULT NULL,
			messageType TEXT DEFAULT 'text',
			attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffChatMessages_staffChatId FOREIGN KEY (staffChatId) REFERENCES staffChats(id),
			CONSTRAINT fk_staffChatMessages_senderUserId FOREIGN KEY (senderUserId) REFERENCES users(id),
			CONSTRAINT fk_staffChatMessages_recipientUserId FOREIGN KEY (recipientUserId) REFERENCES users(id)
		);
		ALTER TABLE staffChatMessages ALTER COLUMN id SET DEFAULT nextval('staffChatMessages_id_seq'::regclass);
		ALTER TABLE staffChatMessages ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffChatMessages ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('staffChatMessages_id_seq', (SELECT MAX(id) FROM staffChatMessages));
    `)
	log.Println(err)
	return nil
}

func (m *StaffChatMessagesMigration1687787690) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffChatMessagesMigration1687787690';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration StaffChatMessagesMigration1687787690: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS staffChatMessages;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete staffChatMessages: %v\n", err)
	}
	return nil
}
