package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type QuizQuestionOptionsMigration1684787660 struct {
}

func NewQuizQuestionOptionsMigration1684787660() *QuizQuestionOptionsMigration1684787660 {
	return &QuizQuestionOptionsMigration1684787660{}
}

func (m *QuizQuestionOptionsMigration1684787660) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS quizQuestionOptions (
			id SERIAL PRIMARY KEY,
			quizQuestionId INT NOT NULL,
			title VARCHAR(255) NOT NULL,
		    score INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_quizQuestionId FOREIGN KEY (quizQuestionId) REFERENCES quizQuestions(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *QuizQuestionOptionsMigration1684787660) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionOptionsMigration1684787660'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS quizQuestionOptions CASCADE;
    `)
	return nil
}
