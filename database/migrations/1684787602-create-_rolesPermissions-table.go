package migrations

import (
	"fmt"
	"log"
	"strings"

	"github.com/hoitek/Maja-Service/database"
)

type RolesPermissionsMigration1684787602 struct {
}

func NewRolesPermissionsMigration1684787602() *RolesPermissionsMigration1684787602 {
	return &RolesPermissionsMigration1684787602{}
}

func generateRolePermValues() string {
	var values = []string{}
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 60; j++ {
			values = append(values, fmt.Sprintf("(%d, %d, %d, '2020-01-01', '2020-01-01', '2020-01-01')", j+((i-1)*60), i, j))
		}
	}
	return strings.Join(values, ",\n")
}

func (m *RolesPermissionsMigration1684787602) MigrateUp() error {
	database.PostgresDB.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS _rolesPermissions (
			id SERIAL PRIMARY KEY,
			roleId INT NOT NULL,
			permissionId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE _rolesPermissions ALTER COLUMN id SET DEFAULT nextval('_rolesPermissions_id_seq'::regclass);
		ALTER TABLE _rolesPermissions ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE _rolesPermissions ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO _rolesPermissions
			(id, roleId, permissionId, created_at, updated_at, deleted_at)
		VALUES
			%s
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('_rolesPermissions_id_seq', (SELECT MAX(id) FROM _rolesPermissions));
    `, generateRolePermValues()))
	return nil
}

func (m *RolesPermissionsMigration1684787602) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'RolesPermissionsMigration1684787602';`)
	if err != nil {
		log.Printf("-------- Error deleting migration RolesPermissionsMigration1684787602: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS _rolesPermissions;
    `)
	if err != nil {
		log.Printf("-------- Error deleting table _rolesPermissions: %v\n", err)
	}
	return nil
}
