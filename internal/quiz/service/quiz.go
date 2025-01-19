package service

import (
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/quiz/constants"
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
	"github.com/hoitek/Maja-Service/internal/quiz/models"
	"github.com/hoitek/Maja-Service/internal/quiz/ports"
	"github.com/hoitek/Maja-Service/storage"
	"math"

	"github.com/hoitek/Kit/exp"
)

type QuizService struct {
	PostgresRepository ports.QuizRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewQuizService(pDB ports.QuizRepositoryPostgresDB, m *storage.MinIO) QuizService {
	go minio.SetupMinIOStorage(constants.QUIZ_BUCKET_NAME, m)
	return QuizService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *QuizService) Query(q *models.QuizzesQueryRequestParams, authenticatedUser *sharedmodels.AuthenticatedUser) (*restypes.QueryResponse, error) {
	q.AuthenticatedUser = authenticatedUser
	quizzes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.QuizzesQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      quizzes,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizService) Create(payload *models.QuizzesCreateRequestBody) (*domain.Quiz, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *QuizService) Delete(payload *models.QuizzesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// TODO this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *QuizService) Update(payload *models.QuizzesCreateRequestBody, id int64) (*domain.Quiz, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *QuizService) GetQuizByTitle(quizName string, authenticatedUser *sharedmodels.AuthenticatedUser) *domain.Quiz {
	r, err := s.Query(&models.QuizzesQueryRequestParams{
		Filters: models.QuizFilterType{
			Title: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: quizName,
			},
		},
	}, authenticatedUser)
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	quizzes, ok := r.Items.([]*domain.Quiz)
	if !ok {
		return nil
	}
	return quizzes[0]
}

func (s *QuizService) GetQuizByID(id int, authenticatedUser *sharedmodels.AuthenticatedUser) *domain.Quiz {
	r, err := s.Query(&models.QuizzesQueryRequestParams{
		ID: id,
	}, authenticatedUser)
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	quizzes, ok := r.Items.([]*domain.Quiz)
	if !ok {
		return nil
	}
	return quizzes[0]
}

func (s *QuizService) UpdateStatus(payload *models.QuizzesUpdateStatusRequestBody, id int64) (*domain.Quiz, error) {
	return s.PostgresRepository.UpdateStatus(payload, id)
}

func (s *QuizService) CreateQuestion(payload *models.QuizzesCreateQuestionRequestBody) (*domain.QuizQuestion, error) {
	return s.PostgresRepository.CreateQuestion(payload)
}

func (s *QuizService) QueryQuestions(q *models.QuizzesQueryQuestionsRequestParams) (*restypes.QueryResponse, error) {
	quizzes, err := s.PostgresRepository.QueryQuestions(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountQuestions(&models.QuizzesQueryQuestionsRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      quizzes,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizService) GetQuizQuestionByID(id int64) *domain.QuizQuestion {
	r, err := s.QueryQuestions(&models.QuizzesQueryQuestionsRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	if r.Items == nil {
		return nil
	}
	quizzes, ok := r.Items.([]*domain.QuizQuestion)
	if !ok {
		return nil
	}
	if len(quizzes) == 0 {
		return nil
	}
	return quizzes[0]
}

func (s *QuizService) UpdateQuizQuestion(payload *models.QuizzesCreateQuestionRequestBody, id int64) (*domain.QuizQuestion, error) {
	return s.PostgresRepository.UpdateQuizQuestion(payload, id)
}

func (s *QuizService) StartQuiz(payload *models.QuizzesCreateStartRequestBody) (*domain.QuizParticipant, error) {
	return s.PostgresRepository.StartQuiz(payload)
}

func (s *QuizService) QueryParticipants(queries *models.QuizzesQueryParticipantsRequestParams) (*restypes.QueryResponse, error) {
	quizzes, err := s.PostgresRepository.QueryParticipants(queries)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountParticipants(&models.QuizzesQueryParticipantsRequestParams{
		ID:      queries.ID,
		Page:    queries.Page,
		Limit:   0,
		Filters: queries.Filters,
	})
	if err != nil {
		return nil, err
	}

	queries.Page = exp.TerIf(queries.Page < 1, 1, queries.Page)
	queries.Limit = exp.TerIf(queries.Limit < 10, 1, queries.Limit)

	page := queries.Page
	limit := queries.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      quizzes,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizService) GetQuizParticipantByQuizIDAndUserID(quizId int64, userId int64) *domain.QuizParticipant {
	r, err := s.QueryParticipants(&models.QuizzesQueryParticipantsRequestParams{
		QuizID: quizId,
		UserID: userId,
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	if r.Items == nil {
		return nil
	}
	quizParticipants, ok := r.Items.([]*domain.QuizParticipant)
	if !ok {
		return nil
	}
	if len(quizParticipants) == 0 {
		return nil
	}
	return quizParticipants[0]
}

func (s *QuizService) QueryQuestionOptions(q *models.QuizzesQueryQuestionOptionsRequestParams) (*restypes.QueryResponse, error) {
	quizzes, err := s.PostgresRepository.QueryQuestionOptions(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountQuestionOptions(&models.QuizzesQueryQuestionOptionsRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      quizzes,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizService) GetQuizQuestionOptionByID(id int64) *domain.QuizQuestionOption {
	r, err := s.QueryQuestionOptions(&models.QuizzesQueryQuestionOptionsRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	if r.Items == nil {
		return nil
	}
	options, ok := r.Items.([]*domain.QuizQuestionOption)
	if !ok {
		return nil
	}
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

func (s *QuizService) QueryAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) (*restypes.QueryResponse, error) {
	quizzes, err := s.PostgresRepository.QueryAnswers(queries)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountAnswers(&models.QuizzesQueryQuestionAnswersRequestParams{
		ID:      queries.ID,
		Page:    queries.Page,
		Limit:   0,
		Filters: queries.Filters,
	})
	if err != nil {
		return nil, err
	}

	queries.Page = exp.TerIf(queries.Page < 1, 1, queries.Page)
	queries.Limit = exp.TerIf(queries.Limit < 10, 1, queries.Limit)

	page := queries.Page
	limit := queries.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      quizzes,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizService) CreateAnswer(payload *models.QuizzesCreateAnswerRequestBody) (*domain.QuizQuestionAnswer, error) {
	return s.PostgresRepository.CreateAnswer(payload)
}

func (s *QuizService) GetQuizQuestionsByQuizID(quizId int64) []*domain.QuizQuestion {
	r, err := s.QueryQuestions(&models.QuizzesQueryQuestionsRequestParams{
		QuizID: quizId,
		Limit:  -1,
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	if r.Items == nil {
		return nil
	}
	quizQuestions, ok := r.Items.([]*domain.QuizQuestion)
	if !ok {
		return nil
	}
	return quizQuestions
}

func (s *QuizService) GetQuizAnswersByQuizIDAndUserID(quizId int64, userId int64) []*domain.QuizQuestionAnswer {
	r, err := s.QueryAnswers(&models.QuizzesQueryQuestionAnswersRequestParams{
		QuizID: quizId,
		UserID: userId,
		Limit:  -1,
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	if r.Items == nil {
		return nil
	}
	quizAnswers, ok := r.Items.([]*domain.QuizQuestionAnswer)
	if !ok {
		return nil
	}
	return quizAnswers
}

func (s *QuizService) UpdateQuizEnd(payload *models.QuizzesUpdateEndRequestBody, quizId int64) (*domain.QuizParticipant, error) {
	return s.PostgresRepository.UpdateQuizEnd(payload, quizId)
}
