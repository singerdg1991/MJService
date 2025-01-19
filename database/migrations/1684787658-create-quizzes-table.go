package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type QuizzesMigration1684787658 struct {
}

func NewQuizzesMigration1684787658() *QuizzesMigration1684787658 {
	return &QuizzesMigration1684787658{}
}

func (m *QuizzesMigration1684787658) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS quizzes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
            startDateTime TIMESTAMP DEFAULT NULL,
            endDateTime TIMESTAMP DEFAULT NULL,
            durationInMinute INT DEFAULT NULL,
            status VARCHAR(255) DEFAULT 'disable',
            availableParticipantType VARCHAR(255) DEFAULT 'all',
            isLock BOOLEAN DEFAULT false,
            password VARCHAR(255) DEFAULT NULL,
            description TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL
		);
		CREATE TABLE IF NOT EXISTS quizAvailableParticipants (
			id SERIAL PRIMARY KEY,
			quizId INT NOT NULL,
			userId INT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_quizId FOREIGN KEY (quizId) REFERENCES quizzes(id) ON DELETE CASCADE,
			CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *QuizzesMigration1684787658) MigrateDown() error {
	tx, err := database.PostgresDB.Begin()
	if err != nil {
		return err
	}
	tx.Exec(`
		DROP TABLE IF EXISTS quizzes CASCADE;
		DROP TABLE IF EXISTS quizAvailableParticipants CASCADE;
	    DROP TABLE IF EXISTS quizParticipants CASCADE;
		DROP TABLE IF EXISTS quizQuestions CASCADE;
		DROP TABLE IF EXISTS quizQuestionOptions CASCADE;
		DROP TABLE IF EXISTS quizQuestionAnswers CASCADE;
    `)
	tx.Exec(`DELETE FROM _migrations WHERE name = 'QuizzesMigration1684787658'`)
	tx.Exec(`DELETE FROM _migrations WHERE name = 'QuizParticipantsMigration1684787661'`)
	tx.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionsMigration1684787659'`)
	tx.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionOptionsMigration1684787660'`)
	tx.Exec(`DELETE FROM _migrations WHERE name = 'QuizQuestionAnswersMigration1684787666'`)
	tx.Commit()
	return nil
}
