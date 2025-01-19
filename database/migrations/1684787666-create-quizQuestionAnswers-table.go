package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type QuizQuestionAnswersMigration1684787666 struct {
}

func NewQuizQuestionAnswersMigration1684787666() *QuizQuestionAnswersMigration1684787666 {
	return &QuizQuestionAnswersMigration1684787666{}
}

func (m *QuizQuestionAnswersMigration1684787666) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS quizQuestionAnswers (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
			questionId INT NOT NULL,
			quizQuestionOptionId INT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_quizQuestionAnswers_userId FOREIGN KEY (userId) REFERENCES users(id),
			CONSTRAINT fk_quizQuestionAnswers_questionId FOREIGN KEY (questionId) REFERENCES quizQuestions(id),
			CONSTRAINT fk_quizQuestionAnswers_quizQuestionOptionId FOREIGN KEY (quizQuestionOptionId) REFERENCES quizQuestionOptions(id),
			CONSTRAINT fk_quizQuestionAnswers_unique UNIQUE (userId, questionId, quizQuestionOptionId)
		);
    `)
	return nil
}

func (m *QuizQuestionAnswersMigration1684787666) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionAnswersMigration1684787666'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS quizQuestionAnswers CASCADE;
    `)
	return nil
}
