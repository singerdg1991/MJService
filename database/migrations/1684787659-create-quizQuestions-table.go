package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type QuizQuestionsMigration1684787659 struct {
}

func NewQuizQuestionsMigration1684787659() *QuizQuestionsMigration1684787659 {
	return &QuizQuestionsMigration1684787659{}
}

func (m *QuizQuestionsMigration1684787659) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS quizQuestions (
			id SERIAL PRIMARY KEY,
			quizId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_quizId FOREIGN KEY (quizId) REFERENCES quizzes(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *QuizQuestionsMigration1684787659) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionsMigration1684787659'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS quizQuestions CASCADE;
    `)
	return nil
}
