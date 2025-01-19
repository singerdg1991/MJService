package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomersDiagnosesMigration1684787678 struct {
}

func NewCustomersDiagnosesMigration1684787678() *CustomersDiagnosesMigration1684787678 {
	return &CustomersDiagnosesMigration1684787678{}
}

func (m *CustomersDiagnosesMigration1684787678) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customersDiagnoses (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			diagnoseId INT NOT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_diagnose_id FOREIGN KEY (diagnoseId) REFERENCES diagnoses(id) ON DELETE CASCADE
		);
    `)
	if err != nil {
		log.Println("************ Error in creating customersDiagnoses table ************", err.Error())
	}
	return nil
}

func (m *CustomersDiagnosesMigration1684787678) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersDiagnosesMigration1684787678'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customersDiagnoses CASCADE;
    `)
	return nil
}
