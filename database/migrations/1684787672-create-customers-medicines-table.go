package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type CustomersMedicinesMigration1684787672 struct {
}

func NewCustomersMedicinesMigration1684787672() *CustomersMedicinesMigration1684787672 {
	return &CustomersMedicinesMigration1684787672{}
}

func (m *CustomersMedicinesMigration1684787672) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerMedicines (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			prescriptionId INT NOT NULL,
			medicineId INT NOT NULL,
			dosageAmount INT NOT NULL,
			dosageUnit VARCHAR(255) NOT NULL,
			days JSONB NOT NULL,
            isJustOneTime BOOLEAN NOT NULL,
			hours JSONB NOT NULL,
            start_date TIMESTAMP DEFAULT NULL,
            end_date TIMESTAMP DEFAULT NULL,
			warning TEXT DEFAULT NULL,
			isUseAsNeeded BOOLEAN NOT NULL,
			attachments JSONB DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
            CONSTRAINT fk_customerMedicines_customerId FOREIGN KEY (customerId) REFERENCES customers(id),
            CONSTRAINT fk_customerMedicines_prescriptionId FOREIGN KEY (prescriptionId) REFERENCES prescriptions(id),
            CONSTRAINT fk_customerMedicines_medicineId FOREIGN KEY (medicineId) REFERENCES medicines(id)
		);
    `)
	return nil
}

func (m *CustomersMedicinesMigration1684787672) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersMedicinesMigration1684787672'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerMedicines CASCADE;
    `)
	return nil
}
