package migrations

import "github.com/hoitek/Maja-Service/database"

type PunishmentsMigration1684787643 struct {
}

func NewPunishmentsMigration1684787643() *PunishmentsMigration1684787643 {
	return &PunishmentsMigration1684787643{}
}

func (m *PunishmentsMigration1684787643) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS punishments (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE punishments ALTER COLUMN id SET DEFAULT nextval('punishments_id_seq'::regclass);
		ALTER TABLE punishments ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE punishments ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO punishments
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'punishment 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'punishment 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'punishment 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'punishment 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'punishment 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('punishments_id_seq', (SELECT MAX(id) FROM punishments));
    `)
	return nil
}

func (m *PunishmentsMigration1684787643) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'PunishmentsMigration1684787643';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS punishments CASCADE;
    `)
	return nil
}
