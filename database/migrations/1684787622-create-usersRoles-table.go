package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type UsersRolesMigration1684787622 struct {
}

func NewUsersRolesMigration1684787622() *UsersRolesMigration1684787622 {
	return &UsersRolesMigration1684787622{}
}

func (m *UsersRolesMigration1684787622) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS usersRoles (
			id SERIAL PRIMARY KEY,
			userId INTEGER NOT NULL,
			roleId INTEGER NOT NULL,
			CONSTRAINT usersRoles_userId_roleId_unique UNIQUE (userId, roleId),
			CONSTRAINT usersRoles_userId_fkey FOREIGN KEY (userId) REFERENCES users (id),
		    CONSTRAINT usersRoles_roleId_fkey FOREIGN KEY (roleId) REFERENCES _roles (id)
		);
		ALTER TABLE usersRoles ALTER COLUMN id SET DEFAULT nextval('usersRoles_id_seq'::regclass);
    `)
	return nil
}

func (m *UsersRolesMigration1684787622) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'UsersRolesMigration1684787622';`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS usersRoles;
    `)
	return nil
}
