package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CyclePickupShiftTodosMigration1685787685 struct {
}

func NewCyclePickupShiftTodosMigration1685787685() *CyclePickupShiftTodosMigration1685787685 {
	return &CyclePickupShiftTodosMigration1685787685{}
}

func (m *CyclePickupShiftTodosMigration1685787685) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS cyclePickupShiftTodos (
			id SERIAL PRIMARY KEY,
			cyclePickupShiftId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
            timeValue TIME NOT NULL,
            dateValue DATE NOT NULL,
            description TEXT DEFAULT NULL,
			attachments JSONB DEFAULT '[]',
			notDoneReason TEXT DEFAULT NULL,
            status VARCHAR(255) NOT NULL,
            done_at TIMESTAMP DEFAULT NULL,
			not_done_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			createdBy INT DEFAULT NULL,
			CONSTRAINT fk_todos_cyclePickupShiftId FOREIGN KEY (cyclePickupShiftId) REFERENCES cyclePickupShifts(id),
			CONSTRAINT fk_todos_createdBy FOREIGN KEY (createdBy) REFERENCES users(id)
		);
		ALTER TABLE cyclePickupShiftTodos ALTER COLUMN id SET DEFAULT nextval('cyclePickupShiftTodos_id_seq'::regclass);
		ALTER TABLE cyclePickupShiftTodos ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE cyclePickupShiftTodos ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('cyclePickupShiftTodos_id_seq', (SELECT MAX(id) FROM cyclePickupShiftTodos));
    `)
	log.Println(err)
	return nil
}

func (m *CyclePickupShiftTodosMigration1685787685) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CyclePickupShiftTodosMigration1685787685';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration CyclePickupShiftTodosMigration1685787685: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS cyclePickupShiftTodos;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete cyclePickupShiftTodos: %v\n", err)
	}
	return nil
}
