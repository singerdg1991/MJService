package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type CustomersAbsencesMigration1684787653 struct {
}

func NewCustomersAbsencesMigration1684787653() *CustomersAbsencesMigration1684787653 {
	return &CustomersAbsencesMigration1684787653{}
}

func (m *CustomersAbsencesMigration1684787653) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerAbsences (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP DEFAULT NULL,
			reason TEXT DEFAULT NULL,
			attachments JSONB NOT NULL DEFAULT '[]',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *CustomersAbsencesMigration1684787653) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersAbsencesMigration1684787653'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerAbsences;
    `)
	return nil
}
