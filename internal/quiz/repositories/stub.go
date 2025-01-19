package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
	"github.com/hoitek/Maja-Service/internal/quiz/models"
)

type QuizRepositoryStub struct {
	Quizzes []*domain.Quiz
}

type quizTestCondition struct {
	HasError bool
}

var UserTestCondition *quizTestCondition = &quizTestCondition{}

func NewQuizRepositoryStub() *QuizRepositoryStub {
	return &QuizRepositoryStub{
		Quizzes: []*domain.Quiz{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *QuizRepositoryStub) Query(dataModel *models.QuizzesQueryRequestParams) ([]*domain.Quiz, error) {
	var quizzes []*domain.Quiz
	for _, v := range r.Quizzes {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			quizzes = append(quizzes, v)
			break
		}
	}
	return quizzes, nil
}

func (r *QuizRepositoryStub) Count(dataModel *models.QuizzesQueryRequestParams) (int64, error) {
	var quizzes []*domain.Quiz
	for _, v := range r.Quizzes {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			quizzes = append(quizzes, v)
			break
		}
	}
	return int64(len(quizzes)), nil
}

func (r *QuizRepositoryStub) Migrate() {
	// do stuff
}

func (r *QuizRepositoryStub) Seed() {
	// do stuff
}

func (r *QuizRepositoryStub) Create(payload *models.QuizzesCreateRequestBody) (*domain.Quiz, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) Delete(payload *models.QuizzesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) Update(payload *models.QuizzesCreateRequestBody, id int64) (*domain.Quiz, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) UpdateStatus(payload *models.QuizzesUpdateStatusRequestBody, id int64) (*domain.Quiz, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CreateQuestion(payload *models.QuizzesCreateQuestionRequestBody) (*domain.QuizQuestion, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) QueryQuestions(queries *models.QuizzesQueryQuestionsRequestParams) ([]*domain.QuizQuestion, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CountQuestions(queries *models.QuizzesQueryQuestionsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) UpdateQuizQuestion(payload *models.QuizzesCreateQuestionRequestBody, questionID int64) (*domain.QuizQuestion, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) StartQuiz(payload *models.QuizzesCreateStartRequestBody) (*domain.QuizParticipant, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) QueryParticipants(queries *models.QuizzesQueryParticipantsRequestParams) ([]*domain.QuizParticipant, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CountParticipants(queries *models.QuizzesQueryParticipantsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) QueryQuestionOptions(q *models.QuizzesQueryQuestionOptionsRequestParams) ([]*domain.QuizQuestionOption, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CountQuestionOptions(dataModel *models.QuizzesQueryQuestionOptionsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CreateAnswer(payload *models.QuizzesCreateAnswerRequestBody) (*domain.QuizQuestionAnswer, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) QueryAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) ([]*domain.QuizQuestionAnswer, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) CountAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) (int64, error) {
	panic("implement me")
}

func (r *QuizRepositoryStub) UpdateQuizEnd(payload *models.QuizzesUpdateEndRequestBody, quizId int64) (*domain.QuizParticipant, error) {
	panic("implement me")
}
