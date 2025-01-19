package migrations

import "github.com/hoitek/Maja-Service/database"

type PermissionsMigration1684787613 struct {
}

func NewPermissionsMigration1684787613() *PermissionsMigration1684787613 {
	return &PermissionsMigration1684787613{}
}

func (m *PermissionsMigration1684787613) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE permissions ALTER COLUMN id SET DEFAULT nextval('permissions_id_seq'::regclass);
		ALTER TABLE permissions ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE permissions ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO permissions
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'permission 01', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'permission 02', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'permission 03', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'permission 04', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'permission 05', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));
    `)
	return nil
}

func (m *PermissionsMigration1684787613) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'PermissionsMigration1684787613'`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS permissions CASCADE;
    `)
	return nil
}
