package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CyclePickupShiftsMigration1684787681 struct {
}

func NewCyclePickupShiftsMigration1684787681() *CyclePickupShiftsMigration1684787681 {
	return &CyclePickupShiftsMigration1684787681{}
}

func (m *CyclePickupShiftsMigration1684787681) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cyclePickupShifts (
			id SERIAL PRIMARY KEY,
			cycleId INT NOT NULL,
			staffId INT NOT NULL,
			shiftId INT NOT NULL,
			cycleStaffTypeId INT NOT NULL,
			datetime TIMESTAMP NOT NULL,
			status VARCHAR(255) DEFAULT 'not-started',
			prevStatus VARCHAR(255) DEFAULT 'not-started',
			startKilometer VARCHAR(255) DEFAULT NULL,
			reasonOfTheCancellation TEXT DEFAULT NULL,
			reasonOfTheReactivation TEXT DEFAULT NULL,
			reasonOfTheResume TEXT DEFAULT NULL,
			reasonOfThePause TEXT DEFAULT NULL,
			isUnplanned BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			started_at TIMESTAMP DEFAULT NULL,
			ended_at TIMESTAMP DEFAULT NULL,
			cancelled_at TIMESTAMP DEFAULT NULL,
			delayed_at TIMESTAMP DEFAULT NULL,
			paused_at TIMESTAMP DEFAULT NULL,
			resumed_at TIMESTAMP DEFAULT NULL,
			reactivated_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycle FOREIGN KEY (cycleId) REFERENCES cycles(id),
			CONSTRAINT fk_staff FOREIGN KEY (staffId) REFERENCES staffs(id),
			CONSTRAINT fk_shift FOREIGN KEY (shiftId) REFERENCES cycleShifts(id),
			CONSTRAINT fk_cycleStaffType FOREIGN KEY (cycleStaffTypeId) REFERENCES cycleStaffTypes(id)
		);
		ALTER TABLE cyclePickupShifts ALTER COLUMN id SET DEFAULT nextval('cyclePickupShifts_id_seq'::regclass);
		ALTER TABLE cyclePickupShifts ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cyclePickupShifts ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cyclePickupShifts_id_seq', (SELECT MAX(id) FROM cyclePickupShifts));
    `)
	log.Println(err)
	return nil
}

func (m *CyclePickupShiftsMigration1684787681) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CyclePickupShiftsMigration1684787681';`)
	if err != nil {
		log.Printf("-------- Error deleting migration: %v", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cyclePickupShifts;
    `)
	if err != nil {
		log.Printf("-------- Error dropping table: %v", err)
	}
	return nil
}
