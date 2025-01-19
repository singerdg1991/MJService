package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

type CycleChatMessagesMigration1687787687 struct {
}

func NewCycleChatMessagesMigration1687787687() *CycleChatMessagesMigration1687787687 {
	return &CycleChatMessagesMigration1687787687{}
}

func (m *CycleChatMessagesMigration1687787687) MigrateUp() error {
	// messageType -- text, image, video, audio
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleChatMessages (
			id SERIAL PRIMARY KEY,
			cycleChatId INT NOT NULL,
			senderUserId INT NOT NULL,
			recipientUserId INT NOT NULL,
			isSystem BOOLEAN DEFAULT FALSE,
			message TEXT DEFAULT NULL,
			messageType TEXT DEFAULT 'text',
			attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycleChatMessages_cycleChatId FOREIGN KEY (cycleChatId) REFERENCES cycleChats(id),
			CONSTRAINT fk_cycleChatMessages_senderUserId FOREIGN KEY (senderUserId) REFERENCES users(id),
			CONSTRAINT fk_cycleChatMessages_recipientUserId FOREIGN KEY (recipientUserId) REFERENCES users(id)
		);
		ALTER TABLE cycleChatMessages ALTER COLUMN id SET DEFAULT nextval('cycleChatMessages_id_seq'::regclass);
		ALTER TABLE cycleChatMessages ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleChatMessages ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleChatMessages_id_seq', (SELECT MAX(id) FROM cycleChatMessages));
    `)
	log.Println(err)
	return nil
}

func (m *CycleChatMessagesMigration1687787687) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleChatMessagesMigration1687787687';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration CycleChatMessagesMigration1687787687: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleChatMessages;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete cycleChatMessages: %v\n", err)
	}
	return nil
}
