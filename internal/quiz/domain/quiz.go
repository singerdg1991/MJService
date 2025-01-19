package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: QuizAvailableParticipantUser
 */
type QuizAvailableParticipantUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	Email     string `json:"email" openapi:"example:sgh370@yahoo.com;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:John;required"`
}

/*
 * @apiDefine: QuizAvailableParticipant
 */
type QuizAvailableParticipant struct {
	ID     uint                          `json:"id" openapi:"example:1"`
	QuizID uint                          `json:"quizId" openapi:"example:1;required"`
	UserID uint                          `json:"userId" openapi:"example:1;required"`
	User   *QuizAvailableParticipantUser `json:"user" openapi:"$ref:QuizAvailableParticipantUser;required"`
}

/*
 * @apiDefine: Quiz
 */
type Quiz struct {
	ID                                     uint                        `json:"id" openapi:"example:1"`
	Title                                  string                      `json:"title" openapi:"example:title;required"`
	StartDateTime                          time.Time                   `json:"startDateTime" openapi:"example:2021-01-01T00:00:00Z;required"`
	EndDateTime                            *time.Time                  `json:"endDateTime" openapi:"example:2021-01-01T00:00:00Z;required"`
	DurationInMinute                       *int                        `json:"durationInMinute" openapi:"example:60;required"`
	Status                                 string                      `json:"status" openapi:"example:disable;required"`
	AvailableParticipantType               *string                     `json:"availableParticipantType" openapi:"example:all;required"`
	AvailableParticipants                  []*QuizAvailableParticipant `json:"availableParticipants" openapi:"$ref:QuizAvailableParticipant;type:array;"`
	IsLock                                 *bool                       `json:"isLock" openapi:"example:false;required"`
	Password                               *string                     `json:"password" openapi:"example:password;required"`
	Description                            *string                     `json:"description" openapi:"example:description;required"`
	IsCurrentAuthorizedUserStartedThisQuiz bool                        `json:"isCurrentAuthorizedUserStartedThisQuiz" openapi:"example:false;required"`
	IsCurrentAuthorizedUserEndedThisQuiz   bool                        `json:"isCurrentAuthorizedUserEndedThisQuiz" openapi:"example:false;required"`
	QuestionsCount                         int                         `json:"questionsCount" openapi:"example:0;required"`
	CreatedAt                              time.Time                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                              time.Time                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                              *time.Time                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Quiz) TableName() string {
	return "quizzes"
}

func (u *Quiz) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *Quiz) ValidatePassword(password string) bool {
	if !*u.IsLock {
		return true
	}
	return *u.Password == password
}
