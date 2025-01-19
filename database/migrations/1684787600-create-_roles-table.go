package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

var (
	RoleOwner          uint = 1
	RoleStateManager   uint = 2
	RoleCityManager    uint = 3
	RoleSectionManager uint = 4
	RoleTeamManager    uint = 5
)

type RoleMigration1684787600 struct {
}

func NewRoleMigration1684787600() *RoleMigration1684787600 {
	return &RoleMigration1684787600{}
}

func (m *RoleMigration1684787600) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS _roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(255) DEFAULT 'internal', -- "core", "internal"
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE _roles ALTER COLUMN id SET DEFAULT nextval('_roles_id_seq'::regclass);
		ALTER TABLE _roles ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE _roles ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO _roles
			(id, name, type, created_at, updated_at, deleted_at)
		VALUES
			(1, 'Owner', 'core', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'State-Manager', 'core', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'City-Manager', 'core', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'Section-Manager', 'core', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'Team-Manager', 'core', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('_roles_id_seq', (SELECT MAX(id) FROM _roles));
    `)
	return nil
}

func (m *RoleMigration1684787600) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'RoleMigration1684787600';`)
	if err != nil {
		log.Printf("-------- Error deleting migration RoleMigration1684787600: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS _roles;
    `)
	if err != nil {
		log.Printf("-------- Error deleting table _roles: %v\n", err)
	}
	return nil
}
