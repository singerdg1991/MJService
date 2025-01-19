package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type EvaluationsMigration1684787669 struct {
}

func NewEvaluationsMigration1684787669() *EvaluationsMigration1684787669 {
	return &EvaluationsMigration1684787669{}
}

func (m *EvaluationsMigration1684787669) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS evaluations (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			evaluationType VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staffId FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *EvaluationsMigration1684787669) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'EvaluationsMigration1684787669'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS evaluations CASCADE;
    `)
	return nil
}
