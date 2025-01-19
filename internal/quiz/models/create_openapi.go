package models

import (
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
)

/*
 * @apiDefine: QuizzesCreateResponseDataPermission
 */
type QuizzesCreateResponseDataPermission struct {
	ID    int64  `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:John;required"`
	Title string `json:"title" openapi:"example:John;required"`
}

/*
 * @apiDefine: QuizzesCreateResponseData
 */
type QuizzesCreateResponseData struct {
	ID                                     uint                               `json:"id" openapi:"example:1"`
	Title                                  string                             `json:"title" openapi:"example:title;required"`
	StartDateTime                          string                             `json:"startDateTime" openapi:"example:2021-01-01T00:00:00Z;required"`
	EndDateTime                            string                             `json:"endDateTime" openapi:"example:2021-01-01T00:00:00Z;required"`
	DurationInMinute                       *int                               `json:"durationInMinute" openapi:"example:60;required"`
	Status                                 string                             `json:"status" openapi:"example:disable;required"`
	AvailableParticipantType               *string                            `json:"availableParticipantType" openapi:"example:all;required"`
	AvailableParticipants                  []*domain.QuizAvailableParticipant `json:"availableParticipants" openapi:"$ref:QuizAvailableParticipant;type:array;"`
	IsLock                                 *bool                              `json:"isLock" openapi:"example:false;required"`
	Password                               *string                            `json:"password" openapi:"example:password;required"`
	IsCurrentAuthorizedUserStartedThisQuiz bool                               `json:"isCurrentAuthorizedUserStartedThisQuiz" openapi:"example:false;required"`
	IsCurrentAuthorizedUserEndedThisQuiz   bool                               `json:"isCurrentAuthorizedUserEndedThisQuiz" openapi:"example:false;required"`
	QuestionsCount                         int                                `json:"questionsCount" openapi:"example:0;required"`
	Description                            *string                            `json:"description" openapi:"example:description;required"`
	CreatedAt                              string                             `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                              string                             `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                              *string                            `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesCreateResponse
 */
type QuizzesCreateResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200"`
	Data       []QuizzesCreateResponseData `json:"data" openapi:"$ref:QuizzesCreateResponseData;type:array;required"`
}
