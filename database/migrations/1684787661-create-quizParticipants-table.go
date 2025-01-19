package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type QuizParticipantsMigration1684787661 struct {
}

func NewQuizParticipantsMigration1684787661() *QuizParticipantsMigration1684787661 {
	return &QuizParticipantsMigration1684787661{}
}

func (m *QuizParticipantsMigration1684787661) MigrateUp() error {
	database.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS quizParticipants (
			id SERIAL PRIMARY KEY,
			quizId INT NOT NULL,
			userId INT NOT NULL,
			ended_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_quizId FOREIGN KEY (quizId) REFERENCES quizzes(id) ON DELETE CASCADE,
			CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
    `)
	return nil
}

func (m *QuizParticipantsMigration1684787661) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'QuizParticipantsMigration1684787661'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS quizParticipants CASCADE;
    `)
	return nil
}
