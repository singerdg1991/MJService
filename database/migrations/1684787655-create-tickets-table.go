package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type TicketsMigration1684787655 struct {
}

func NewTicketsMigration1684787655() *TicketsMigration1684787655 {
	return &TicketsMigration1684787655{}
}

func (m *TicketsMigration1684787655) MigrateUp() error {
	// Dispatcher: userId or departmentId, title, priority, description is required
	// Staff or Customer: title, priority, description is required
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS tickets (
			id SERIAL PRIMARY KEY,
			code VARCHAR(255) NOT NULL,
			userId INT DEFAULT NULL,
			departmentId INT DEFAULT NULL,
			senderType VARCHAR(255) NOT NULL,
            recipientType VARCHAR(255) NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT NOT NULL,
            status VARCHAR(255) NOT NULL,
            priority VARCHAR(255) NOT NULL,
            attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
            createdBy INT DEFAULT NULL,
			updated_at TIMESTAMP DEFAULT now(),
            updatedBy INT DEFAULT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
            deletedBy INT DEFAULT NULL,
            CONSTRAINT fk_tickets_user_id FOREIGN KEY (userId) REFERENCES users(id),
            CONSTRAINT fk_tickets_created_by FOREIGN KEY (createdBy) REFERENCES users(id),
            CONSTRAINT fk_tickets_updated_by FOREIGN KEY (updatedBy) REFERENCES users(id),
            CONSTRAINT fk_tickets_deleted_by FOREIGN KEY (deletedBy) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS ticketMessages (
			id SERIAL PRIMARY KEY,
			ticketId INT DEFAULT NULL,
			senderId INT DEFAULT NULL,
			recipientId INT DEFAULT NULL,
			senderType VARCHAR(255) NOT NULL,
			recipientType VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			attachments JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT now(),
		    updated_at TIMESTAMP DEFAULT now(),
		    deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_ticket_messages_ticket_id FOREIGN KEY (ticketId) REFERENCES tickets(id),
		    CONSTRAINT fk_ticket_messages_sender_id FOREIGN KEY (senderId) REFERENCES users(id),
		    CONSTRAINT fk_ticket_messages_recipient_id FOREIGN KEY (recipientId) REFERENCES users(id)
		);
    `)
	return nil
}

func (m *TicketsMigration1684787655) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'TicketsMigration1684787655'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS tickets CASCADE;
		DROP TABLE IF EXISTS ticketMessages;
    `)
	return nil
}
