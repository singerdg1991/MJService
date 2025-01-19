package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: QuizParticipantQuiz
 */
type QuizParticipantQuiz struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
}

/*
 * @apiDefine: QuizParticipantUser
 */
type QuizParticipantUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName;required"`
	LastName  string `json:"lastName" openapi:"example:lastName;required"`
	Email     string `json:"email" openapi:"example:email;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:avatarUrl;required"`
}

/*
 * @apiDefine: QuizParticipant
 */
type QuizParticipant struct {
	ID        uint                 `json:"id" openapi:"example:1"`
	QuizID    uint                 `json:"quizId" openapi:"example:1"`
	UserID    uint                 `json:"userId" openapi:"example:1"`
	Quiz      *QuizParticipantQuiz `json:"quiz" openapi:"$ref:QuizParticipantQuiz"`
	User      *QuizParticipantUser `json:"user" openapi:"$ref:QuizParticipantUser"`
	EndedAt   *time.Time           `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt time.Time            `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt time.Time            `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *time.Time           `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *QuizParticipant) TableName() string {
	return "quizParticipants"
}

func (u *QuizParticipant) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
