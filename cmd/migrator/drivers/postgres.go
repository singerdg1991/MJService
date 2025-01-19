package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

type Postgres struct {
	DB         *sql.DB
	tableName  string
	migrations []Migration
}

type Migration interface {
	MigrateUp() error
	MigrateDown() error
}

func NewPostgresDriver(db *sql.DB) *Postgres {
	p := &Postgres{
		DB:        db,
		tableName: "_migrations",
	}
	p.setupPrerequisite()
	return p
}

func (m *Postgres) AddMigrations(migrations ...Migration) {
	m.migrations = append(migrations, migrations...)
}

func (m *Postgres) DownAll() {
	for _, migration := range m.migrations {
		if err := migration.MigrateDown(); err != nil {
			log.Printf("Error migrating down: %v\n", err)
		}
	}
}

func (m *Postgres) setupPrerequisite() error {
	// Create migration table if not exists
	log.Printf("Creating migration table %s if not exists", m.tableName)
	_, err := m.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ` + m.tableName + ` (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating migration table %s: %v", m.tableName, err)
	}

	return nil
}

func (m *Postgres) MigrateUp(forceMigrate bool) error {
	for _, migration := range m.migrations {
		// Get struct name from migration
		migrationName := reflect.TypeOf(migration).Elem().Name()

		// Check force migrate
		if forceMigrate {
			// Delete migration name from database
			if _, err := m.DB.Exec("DELETE FROM "+m.tableName+" WHERE name = $1", migrationName); err != nil {
				return fmt.Errorf("error deleting migration name %s from database: %v", migrationName, err)
			}

			// Run migration down
			if err := migration.MigrateDown(); err != nil {
				return fmt.Errorf("error running migration %s: %v", migrationName, err)
			}
		}

		// Check if migration name already exists in database
		var name string
		if err := m.DB.QueryRow("SELECT name FROM "+m.tableName+" WHERE name = $1", migrationName).Scan(&name); err == nil {
			log.Printf("Migration %s already migrated\n", migrationName)
			continue
		}

		// Insert migration name to database
		if _, err := m.DB.Exec("INSERT INTO "+m.tableName+" (name) VALUES ($1)", migrationName); err != nil {
			return fmt.Errorf("error inserting migration name %s to database: %v", migrationName, err)
		}

		// Run migration
		log.Printf("Running migration %s\n", migrationName)
		if err := migration.MigrateUp(); err != nil {
			return fmt.Errorf("error running migration %s: %v", migrationName, err)
		}

		log.Printf("Migration %s ran successfully\n", migrationName)
	}

	return nil
}
