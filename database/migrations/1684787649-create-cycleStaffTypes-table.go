package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CycleStaffTypesMigration1684787649 struct {
}

func NewCycleStaffTypesMigration1684787649() *CycleStaffTypesMigration1684787649 {
	return &CycleStaffTypesMigration1684787649{}
}

func (m *CycleStaffTypesMigration1684787649) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleStaffTypes (
			id SERIAL PRIMARY KEY,
			cycleId INT NOT NULL,
			roleId INT NOT NULL,
			datetime TIMESTAMP NOT NULL,
			shiftName VARCHAR(255) NOT NULL,
            neededStaffCount INT DEFAULT 0,
            startHour TIME NOT NULL,
			endHour TIME NOT NULL,
			isUnplanned BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycle FOREIGN KEY (cycleId) REFERENCES cycles(id),
			CONSTRAINT fk_role FOREIGN KEY (roleId) REFERENCES _roles(id)
		);
		ALTER TABLE cycleStaffTypes ALTER COLUMN id SET DEFAULT nextval('cycleStaffTypes_id_seq'::regclass);
		ALTER TABLE cycleStaffTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleStaffTypes ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleStaffTypes_id_seq', (SELECT MAX(id) FROM cycleStaffTypes));
    `)
	log.Println(err)
	return nil
}

func (m *CycleStaffTypesMigration1684787649) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleStaffTypesMigration1684787649';`)
	if err != nil {
		log.Printf("-------- Error deleting migration: %v", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleStaffTypes;
    `)
	if err != nil {
		log.Printf("-------- Error deleting table: %v", err)
	}
	return nil
}
