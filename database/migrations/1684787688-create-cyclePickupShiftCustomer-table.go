package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

type CyclePickupShiftCustomersMigration1688787688 struct {
}

func NewCyclePickupShiftCustomersMigration1688787688() *CyclePickupShiftCustomersMigration1688787688 {
	return &CyclePickupShiftCustomersMigration1688787688{}
}

func (m *CyclePickupShiftCustomersMigration1688787688) MigrateUp() error {
	// messageType -- text, image, video, audio
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cyclePickupShiftCustomers (
			id SERIAL PRIMARY KEY,
			cyclePickupShiftId INT NOT NULL,
			customerId INT NOT NULL,
			customerServiceId INT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			CONSTRAINT fk_CyclePickupShiftCustomers_cyclePickupShiftId FOREIGN KEY (cyclePickupShiftId) REFERENCES cyclePickupShifts(id),
			CONSTRAINT fk_CyclePickupShiftCustomers_customerId FOREIGN KEY (customerId) REFERENCES customers(id),
			CONSTRAINT fk_CyclePickupShiftCustomers_customerServiceId FOREIGN KEY (customerServiceId) REFERENCES customerServices(id)
		);
		ALTER TABLE cyclePickupShiftCustomers ALTER COLUMN id SET DEFAULT nextval('CyclePickupShiftCustomers_id_seq'::regclass);
		SELECT setval('CyclePickupShiftCustomers_id_seq', (SELECT MAX(id) FROM cyclePickupShiftCustomers));
    `)
	log.Println(err)
	return nil
}

func (m *CyclePickupShiftCustomersMigration1688787688) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CyclePickupShiftCustomersMigration1688787688';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration CyclePickupShiftCustomersMigration1688787688: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cyclePickupShiftCustomers;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete CyclePickupShiftCustomers: %v\n", err)
	}
	return nil
}
