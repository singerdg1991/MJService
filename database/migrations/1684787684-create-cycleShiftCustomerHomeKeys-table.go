package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CycleShiftCustomerHomeKeysMigration1684787684 struct {
}

func NewCycleShiftCustomerHomeKeysMigration1684787684() *CycleShiftCustomerHomeKeysMigration1684787684 {
	return &CycleShiftCustomerHomeKeysMigration1684787684{}
}

func (m *CycleShiftCustomerHomeKeysMigration1684787684) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleShiftCustomerHomeKeys (
			id SERIAL PRIMARY KEY,
			shiftId INT NOT NULL,
			keyNo VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			reason VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			created_by INT DEFAULT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_shift FOREIGN KEY (shiftId) REFERENCES cycleShifts(id),
			CONSTRAINT fk_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
		);
		ALTER TABLE cycleShiftCustomerHomeKeys ALTER COLUMN id SET DEFAULT nextval('cycleShiftCustomerHomeKeys_id_seq'::regclass);
		ALTER TABLE cycleShiftCustomerHomeKeys ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleShiftCustomerHomeKeys ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleShiftCustomerHomeKeys_id_seq', (SELECT MAX(id) FROM cycleShiftCustomerHomeKeys));
    `)
	log.Println(err)
	return nil
}

func (m *CycleShiftCustomerHomeKeysMigration1684787684) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleShiftCustomerHomeKeysMigration1684787684';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration: %v", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleShiftCustomerHomeKeys;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete cycleShiftCustomerHomeKeys: %v", err)
	}
	return nil
}
