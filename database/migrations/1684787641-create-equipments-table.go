package migrations

import "github.com/hoitek/Maja-Service/database"

type EquipmentsMigration1684787641 struct {
}

func NewEquipmentsMigration1684787641() *EquipmentsMigration1684787641 {
	return &EquipmentsMigration1684787641{}
}

func (m *EquipmentsMigration1684787641) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS equipments (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE equipments ALTER COLUMN id SET DEFAULT nextval('equipments_id_seq'::regclass);
		ALTER TABLE equipments ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE equipments ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO equipments
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'equipment 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'equipment 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'equipment 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'equipment 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'equipment 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('equipments_id_seq', (SELECT MAX(id) FROM equipments));
    `)
	return nil
}

func (m *EquipmentsMigration1684787641) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'EquipmentsMigration1684787641';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS equipments CASCADE;
    `)
	return nil
}
