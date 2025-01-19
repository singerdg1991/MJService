package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CycleShiftsMigration1684787683 struct {
}

func NewCycleShiftsMigration1684787683() *CycleShiftsMigration1684787683 {
	return &CycleShiftsMigration1684787683{}
}

func (m *CycleShiftsMigration1684787683) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleShifts (
			id SERIAL PRIMARY KEY,
			exchangeKey VARCHAR(255) NOT NULL,
			cycleId INT NOT NULL,
			staffTypeIds JSONB NOT NULL,
			shiftName VARCHAR(255) NOT NULL,
			vehicleType VARCHAR(255) DEFAULT NULL,
			startLocation VARCHAR(255) DEFAULT NULL,
			datetime DATE NOT NULL,
			status VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycle FOREIGN KEY (cycleId) REFERENCES cycles(id)
		);
		ALTER TABLE cycleShifts ALTER COLUMN id SET DEFAULT nextval('cycleShifts_id_seq'::regclass);
		ALTER TABLE cycleShifts ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleShifts ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleShifts_id_seq', (SELECT MAX(id) FROM cycleShifts));
    `)
	log.Println(err)
	return nil
}

func (m *CycleShiftsMigration1684787683) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleShiftsMigration1684787683';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration: %v", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleShifts;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete cycleShifts: %v", err)
	}
	return nil
}
