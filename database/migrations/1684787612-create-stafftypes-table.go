package migrations

import "github.com/hoitek/Maja-Service/database"

type StaffTypesMigration1684787612 struct {
}

func NewStaffTypesMigration1684787612() *StaffTypesMigration1684787612 {
	return &StaffTypesMigration1684787612{}
}

func (m *StaffTypesMigration1684787612) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS staffTypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE staffTypes ALTER COLUMN id SET DEFAULT nextval('staffTypes_id_seq'::regclass);
		ALTER TABLE staffTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffTypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO staffTypes
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'contract', '2020-01-01', '2020-01-01', NULL),
			(2, 'kaikala', '2020-01-01', '2020-01-01', NULL)
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('staffTypes_id_seq', (SELECT MAX(id) FROM staffTypes));
    `)
	return nil
}

func (m *StaffTypesMigration1684787612) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffTypesMigration1684787612'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS staffTypes CASCADE;
    `)
	return nil
}
