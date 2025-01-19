package migrations

import "github.com/hoitek/Maja-Service/database"

type PermissionsAddDescriptionColumnMigration1684787636 struct {
}

func NewPermissionsAddDescriptionColumnMigration1684787636() *PermissionsAddDescriptionColumnMigration1684787636 {
	return &PermissionsAddDescriptionColumnMigration1684787636{}
}

func (m *PermissionsAddDescriptionColumnMigration1684787636) MigrateUp() error {
	database.PostgresDB.Exec(`
		ALTER TABLE IF EXISTS permissions ADD COLUMN IF NOT EXISTS description VARCHAR(255) NOT NULL DEFAULT '';
    `)
	return nil
}

func (m *PermissionsAddDescriptionColumnMigration1684787636) MigrateDown() error {
	database.PostgresDB.Exec(`
		DELETE FROM _migrations WHERE name = 'PermissionsAddDescriptionColumnMigration1684787636';
	`)
	database.PostgresDB.Exec(`
		ALTER TABLE permissions DROP COLUMN description;
    `)
	return nil
}
