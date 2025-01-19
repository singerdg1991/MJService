package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CycleNextStaffTypesMigration1684787680 struct {
}

func NewCycleNextStaffTypesMigration1684787680() *CycleNextStaffTypesMigration1684787680 {
	return &CycleNextStaffTypesMigration1684787680{}
}

func (m *CycleNextStaffTypesMigration1684787680) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cycleNextStaffTypes (
			id SERIAL PRIMARY KEY,
			currentCycleId INT NOT NULL,
			roleId INT NOT NULL,
			datetime TIMESTAMP NOT NULL,
			shiftName VARCHAR(255) NOT NULL,
            neededStaffCount INT DEFAULT 0,
            startHour TIME NOT NULL,
			endHour TIME NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_cycle FOREIGN KEY (currentCycleId) REFERENCES cycles(id),
			CONSTRAINT fk_role FOREIGN KEY (roleId) REFERENCES _roles(id)
		);
		ALTER TABLE cycleNextStaffTypes ALTER COLUMN id SET DEFAULT nextval('cycleNextStaffTypes_id_seq'::regclass);
		ALTER TABLE cycleNextStaffTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cycleNextStaffTypes ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cycleNextStaffTypes_id_seq', (SELECT MAX(id) FROM cycleNextStaffTypes));
    `)
	log.Println(err)
	return nil
}

func (m *CycleNextStaffTypesMigration1684787680) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CycleNextStaffTypesMigration1684787680';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cycleNextStaffTypes;
    `)
	return nil
}
