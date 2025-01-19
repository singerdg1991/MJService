package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
	"github.com/hoitek/Maja-Service/internal/quiz/models"
)

type QuizService interface {
	Query(dataModel *models.QuizzesQueryRequestParams, authenticatedUser *sharedmodels.AuthenticatedUser) (*restypes.QueryResponse, error)
	Create(payload *models.QuizzesCreateRequestBody) (*domain.Quiz, error)
	Delete(payload *models.QuizzesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.QuizzesCreateRequestBody, id int64) (*domain.Quiz, error)
	GetQuizByTitle(quizName string, authenticatedUser *sharedmodels.AuthenticatedUser) *domain.Quiz
	GetQuizByID(id int, authenticatedUser *sharedmodels.AuthenticatedUser) *domain.Quiz
	UpdateStatus(payload *models.QuizzesUpdateStatusRequestBody, id int64) (*domain.Quiz, error)
	CreateQuestion(payload *models.QuizzesCreateQuestionRequestBody) (*domain.QuizQuestion, error)
	QueryQuestions(q *models.QuizzesQueryQuestionsRequestParams) (*restypes.QueryResponse, error)
	GetQuizQuestionByID(id int64) *domain.QuizQuestion
	UpdateQuizQuestion(payload *models.QuizzesCreateQuestionRequestBody, questionID int64) (*domain.QuizQuestion, error)
	StartQuiz(payload *models.QuizzesCreateStartRequestBody) (*domain.QuizParticipant, error)
	QueryParticipants(queries *models.QuizzesQueryParticipantsRequestParams) (*restypes.QueryResponse, error)
	GetQuizParticipantByQuizIDAndUserID(quizId int64, userId int64) *domain.QuizParticipant
	GetQuizQuestionOptionByID(id int64) *domain.QuizQuestionOption
	QueryQuestionOptions(q *models.QuizzesQueryQuestionOptionsRequestParams) (*restypes.QueryResponse, error)
	CreateAnswer(payload *models.QuizzesCreateAnswerRequestBody) (*domain.QuizQuestionAnswer, error)
	QueryAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) (*restypes.QueryResponse, error)
	GetQuizQuestionsByQuizID(quizId int64) []*domain.QuizQuestion
	GetQuizAnswersByQuizIDAndUserID(quizId int64, userId int64) []*domain.QuizQuestionAnswer
	UpdateQuizEnd(payload *models.QuizzesUpdateEndRequestBody, quizId int64) (*domain.QuizParticipant, error)
}

type QuizRepositoryPostgresDB interface {
	Query(dataModel *models.QuizzesQueryRequestParams) ([]*domain.Quiz, error)
	Count(dataModel *models.QuizzesQueryRequestParams) (int64, error)
	Create(payload *models.QuizzesCreateRequestBody) (*domain.Quiz, error)
	Delete(payload *models.QuizzesDeleteRequestBody) ([]int64, error)
	Update(payload *models.QuizzesCreateRequestBody, id int64) (*domain.Quiz, error)
	UpdateStatus(payload *models.QuizzesUpdateStatusRequestBody, id int64) (*domain.Quiz, error)
	CreateQuestion(payload *models.QuizzesCreateQuestionRequestBody) (*domain.QuizQuestion, error)
	QueryQuestions(queries *models.QuizzesQueryQuestionsRequestParams) ([]*domain.QuizQuestion, error)
	CountQuestions(queries *models.QuizzesQueryQuestionsRequestParams) (int64, error)
	UpdateQuizQuestion(payload *models.QuizzesCreateQuestionRequestBody, questionID int64) (*domain.QuizQuestion, error)
	StartQuiz(payload *models.QuizzesCreateStartRequestBody) (*domain.QuizParticipant, error)
	QueryParticipants(queries *models.QuizzesQueryParticipantsRequestParams) ([]*domain.QuizParticipant, error)
	CountParticipants(queries *models.QuizzesQueryParticipantsRequestParams) (int64, error)
	QueryQuestionOptions(q *models.QuizzesQueryQuestionOptionsRequestParams) ([]*domain.QuizQuestionOption, error)
	CountQuestionOptions(dataModel *models.QuizzesQueryQuestionOptionsRequestParams) (int64, error)
	CreateAnswer(payload *models.QuizzesCreateAnswerRequestBody) (*domain.QuizQuestionAnswer, error)
	QueryAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) ([]*domain.QuizQuestionAnswer, error)
	CountAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) (int64, error)
	UpdateQuizEnd(payload *models.QuizzesUpdateEndRequestBody, quizId int64) (*domain.QuizParticipant, error)
}
