package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type TodosMigration1684787656 struct {
}

func NewTodosMigration1684787656() *TodosMigration1684787656 {
	return &TodosMigration1684787656{}
}

func (m *TodosMigration1684787656) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS Todos (
			id SERIAL PRIMARY KEY,
			userId INT DEFAULT NULL,
            title VARCHAR(255) NOT NULL,
            timeValue TIME NOT NULL,
            dateValue DATE NOT NULL,
            description TEXT DEFAULT NULL,
            status VARCHAR(255) NOT NULL,
            done_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			createdBy INT DEFAULT NULL,
            CONSTRAINT fk_todos_userId FOREIGN KEY (userId) REFERENCES users(id),
			CONSTRAINT fk_todos_createdBy FOREIGN KEY (createdBy) REFERENCES users(id)
		);
    `)
	return nil
}

func (m *TodosMigration1684787656) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'TodosMigration1684787656'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS Todos CASCADE;
    `)
	return nil
}
