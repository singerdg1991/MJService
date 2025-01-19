package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type IncomingCycleShiftsMigration1684787682 struct {
}

func NewIncomingCycleShiftsMigration1684787682() *IncomingCycleShiftsMigration1684787682 {
	return &IncomingCycleShiftsMigration1684787682{}
}

func (m *IncomingCycleShiftsMigration1684787682) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS incomingCyclePickupShifts (
			id SERIAL PRIMARY KEY,
			cycleId INT NOT NULL,
			staffId INT NOT NULL,
			cycleNextStaffTypeId INT NOT NULL,
			datetime TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycle FOREIGN KEY (cycleId) REFERENCES cycles(id),
			CONSTRAINT fk_staff FOREIGN KEY (staffId) REFERENCES staffs(id),
			CONSTRAINT fk_cycleStaffType FOREIGN KEY (cycleNextStaffTypeId) REFERENCES cycleNextStaffTypes(id)
		);
		ALTER TABLE incomingCyclePickupShifts ALTER COLUMN id SET DEFAULT nextval('incomingCyclePickupShifts_id_seq'::regclass);
		ALTER TABLE incomingCyclePickupShifts ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE incomingCyclePickupShifts ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('incomingCyclePickupShifts_id_seq', (SELECT MAX(id) FROM incomingCyclePickupShifts));
    `)
	log.Println(err)
	return nil
}

func (m *IncomingCycleShiftsMigration1684787682) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'IncomingCycleShiftsMigration1684787682';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS incomingCyclePickupShifts;
    `)
	return nil
}
