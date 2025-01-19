package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: QuizQuestionAnswerUser
 */
type QuizQuestionAnswerUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName;required"`
	LastName  string `json:"lastName" openapi:"example:lastName;required"`
	Email     string `json:"email" openapi:"example:email;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:avatarUrl;required"`
}

/*
 * @apiDefine: QuizQuestionAnswerQuestionQuiz
 */
type QuizQuestionAnswerQuestionQuiz struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Title       string  `json:"title" openapi:"example:title;required"`
	Description *string `json:"description" openapi:"example:description;required"`
}

/*
 * @apiDefine: QuizQuestionAnswerQuestion
 */
type QuizQuestionAnswerQuestion struct {
	ID          uint                            `json:"id" openapi:"example:1"`
	QuizID      uint                            `json:"quizId" openapi:"example:1;required"`
	Quiz        *QuizQuestionAnswerQuestionQuiz `json:"quiz" openapi:"$ref:QuizQuestionAnswerQuestionQuiz"`
	Title       string                          `json:"title" openapi:"example:title;required"`
	Description *string                         `json:"description" openapi:"example:description;required"`
}

/*
 * @apiDefine: QuizQuestionAnswerQuizQuestionOption
 */
type QuizQuestionAnswerQuizQuestionOption struct {
	ID             uint   `json:"id" openapi:"example:1"`
	QuizQuestionID uint   `json:"quizQuestionId" openapi:"example:1"`
	Title          string `json:"title" openapi:"example:title;required"`
	Score          uint   `json:"score" openapi:"example:1;required"`
}

/*
 * @apiDefine: QuizQuestionAnswer
 */
type QuizQuestionAnswer struct {
	ID                   uint                                  `json:"id" openapi:"example:1"`
	UserID               uint                                  `json:"userId" openapi:"example:1"`
	User                 *QuizQuestionAnswerUser               `json:"user" openapi:"$ref:QuizQuestionAnswerUser"`
	QuestionID           uint                                  `json:"questionId" openapi:"example:1"`
	Question             *QuizQuestionAnswerQuestion           `json:"question" openapi:"$ref:QuizQuestionAnswerQuestion"`
	QuizQuestionOptionID uint                                  `json:"quizQuestionOptionId" openapi:"example:1"`
	QuizQuestionOption   *QuizQuestionAnswerQuizQuestionOption `json:"quizQuestionOption" openapi:"$ref:QuizQuestionAnswerQuizQuestionOption"`
	CreatedAt            time.Time                             `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt            time.Time                             `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt            *time.Time                            `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *QuizQuestionAnswer) TableName() string {
	return "quizQuestionAnswers"
}

func (u *QuizQuestionAnswer) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
