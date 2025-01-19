package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CycleChatsMigration1686787686 struct {
}

func NewCycleChatsMigration1686787686() *CycleChatsMigration1686787686 {
	return &CycleChatsMigration1686787686{}
}

func (m *CycleChatsMigration1686787686) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleChats (
			id SERIAL PRIMARY KEY,
			cycleId INT NOT NULL,
			cyclePickupShiftId INT NOT NULL,
			senderUserId INT NOT NULL,
			recipientUserId INT NOT NULL,
			isSystem BOOLEAN DEFAULT FALSE,
			message TEXT DEFAULT NULL,
			attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycleChats_cycleId FOREIGN KEY (cycleId) REFERENCES cycles(id),
			CONSTRAINT fk_cycleChats_cyclePickupShiftId FOREIGN KEY (cyclePickupShiftId) REFERENCES cyclePickupShifts(id),
			CONSTRAINT fk_cycleChats_senderUserId FOREIGN KEY (senderUserId) REFERENCES users(id),
			CONSTRAINT fk_cycleChats_recipientUserId FOREIGN KEY (recipientUserId) REFERENCES users(id)
		);
		ALTER TABLE cycleChats ALTER COLUMN id SET DEFAULT nextval('cycleChats_id_seq'::regclass);
		ALTER TABLE cycleChats ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleChats ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleChats_id_seq', (SELECT MAX(id) FROM cycleChats));
    `)
	log.Println(err)
	return nil
}

func (m *CycleChatsMigration1686787686) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleChatsMigration1686787686';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration CycleChatsMigration1686787686: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleChats;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete cycleChats: %v\n", err)
	}
	return nil
}
