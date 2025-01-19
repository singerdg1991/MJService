package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomersRelativesMigration1684787677 struct {
}

func NewCustomersRelativesMigration1684787677() *CustomersRelativesMigration1684787677 {
	return &CustomersRelativesMigration1684787677{}
}

func (m *CustomersRelativesMigration1684787677) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customersRelatives (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			relativeId INT NOT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_relative_id FOREIGN KEY (relativeId) REFERENCES customerRelatives(id) ON DELETE CASCADE
		);
    `)
	if err != nil {
		log.Println("************ Error in creating customersRelatives table ************", err.Error())
	}
	return nil
}

func (m *CustomersRelativesMigration1684787677) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersRelativesMigration1684787677'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customersRelatives CASCADE;
    `)
	return nil
}
