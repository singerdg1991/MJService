package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type KeikkalaShiftsMigration1684787668 struct {
}

func NewKeikkalaShiftsMigration1684787668() *KeikkalaShiftsMigration1684787668 {
	return &KeikkalaShiftsMigration1684787668{}
}

func (m *KeikkalaShiftsMigration1684787668) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS keikkalaShifts (
			id SERIAL PRIMARY KEY,
			roleId INT DEFAULT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			kaupunkiAddress TEXT DEFAULT NULL,
			sections JSONB NOT NULL DEFAULT '[]',
			paymentType VARCHAR(255) NOT NULL,
		    shiftName VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT NULL,
		    status VARCHAR(255) NOT NULL DEFAULT 'open',
		    picked_at TIMESTAMP DEFAULT NULL,
		    pickedBy INT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
		    createdBy INT DEFAULT NULL,
		    updatedBy INT DEFAULT NULL,
		    CONSTRAINT fk_roleId FOREIGN KEY (roleId) REFERENCES _roles(id) ON DELETE SET NULL,
		    CONSTRAINT fk_pickedBy FOREIGN KEY (pickedBy) REFERENCES users(id) ON DELETE SET NULL,
		    CONSTRAINT fk_createdBy FOREIGN KEY (createdBy) REFERENCES users(id) ON DELETE SET NULL,
		    CONSTRAINT fk_updatedBy FOREIGN KEY (updatedBy) REFERENCES users(id) ON DELETE SET NULL
		);
    `)
	return nil
}

func (m *KeikkalaShiftsMigration1684787668) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'KeikkalaShiftsMigration1684787668'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS keikkalaShifts CASCADE;
    `)
	return nil
}
