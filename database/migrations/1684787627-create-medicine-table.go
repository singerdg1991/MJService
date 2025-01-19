package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type MedicinesMigration1684787627 struct {
}

func NewMedicinesMigration1684787627() *MedicinesMigration1684787627 {
	return &MedicinesMigration1684787627{}
}

func (m *MedicinesMigration1684787627) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS medicines (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			code VARCHAR(255) DEFAULT NULL,
            availability VARCHAR(255) DEFAULT NULL,
            manufacturer VARCHAR(255) DEFAULT NULL,
            purposeOfUse TEXT DEFAULT NULL,
            instruction TEXT DEFAULT NULL,
            sideEffects TEXT DEFAULT NULL,
            conditions TEXT DEFAULT NULL,
            description TEXT DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE medicines ALTER COLUMN id SET DEFAULT nextval('medicines_id_seq'::regclass);
		ALTER TABLE medicines ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE medicines ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('medicines_id_seq', (SELECT MAX(id) FROM medicines));
    `)

	return nil
}

func (m *MedicinesMigration1684787627) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'MedicinesMigration1684787627'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS medicines;
    `)
	return nil
}
