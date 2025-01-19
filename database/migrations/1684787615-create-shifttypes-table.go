package migrations

import "github.com/hoitek/Maja-Service/database"

type ShiftTypesMigration1684787615 struct {
}

func NewShiftTypesMigration1684787615() *ShiftTypesMigration1684787615 {
	return &ShiftTypesMigration1684787615{}
}

func (m *ShiftTypesMigration1684787615) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS shiftTypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE shiftTypes ALTER COLUMN id SET DEFAULT nextval('shiftTypes_id_seq'::regclass);
		ALTER TABLE shiftTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE shiftTypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO shiftTypes
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'morning', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'evening', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'night', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('shiftTypes_id_seq', (SELECT MAX(id) FROM shiftTypes));
    `)
	return nil
}

func (m *ShiftTypesMigration1684787615) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ShiftTypesMigration1684787615'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS shiftTypes CASCADE;
    `)
	return nil
}
