package migrations

import "github.com/hoitek/Maja-Service/database"

type RewardsMigration1684787642 struct {
}

func NewRewardsMigration1684787642() *RewardsMigration1684787642 {
	return &RewardsMigration1684787642{}
}

func (m *RewardsMigration1684787642) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS rewards (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE rewards ALTER COLUMN id SET DEFAULT nextval('rewards_id_seq'::regclass);
		ALTER TABLE rewards ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE rewards ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO rewards
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'reward 01', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'reward 02', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'reward 03', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'reward 04', '', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'reward 05', '', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('rewards_id_seq', (SELECT MAX(id) FROM rewards));
    `)
	return nil
}

func (m *RewardsMigration1684787642) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'RewardsMigration1684787642';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS rewards CASCADE;
    `)
	return nil
}
