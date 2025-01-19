package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomerContractualMobilityRestrictionLogsMigration1684787676 struct {
}

func NewCustomerContractualMobilityRestrictionLogsMigration1684787676() *CustomerContractualMobilityRestrictionLogsMigration1684787676 {
	return &CustomerContractualMobilityRestrictionLogsMigration1684787676{}
}

func (m *CustomerContractualMobilityRestrictionLogsMigration1684787676) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerContractualMobilityRestrictionLogs (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			beforeValue TEXT DEFAULT NULL,
			afterValue TEXT DEFAULT NULL,
			createdBy INT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customerContractualMobilityRestrictionLogs_customerId FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_customerContractualMobilityRestrictionLogs_createdBy FOREIGN KEY (createdBy) REFERENCES users(id) ON DELETE SET NULL
		);
    `)
	if err != nil {
		log.Println("************ Error in creating customerContractualMobilityRestrictionLogs table ************", err.Error())
	}
	return nil
}

func (m *CustomerContractualMobilityRestrictionLogsMigration1684787676) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomerContractualMobilityRestrictionLogsMigration1684787676'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerContractualMobilityRestrictionLogs CASCADE;
    `)
	return nil
}
