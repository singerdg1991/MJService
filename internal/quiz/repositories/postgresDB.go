package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
	"github.com/hoitek/Maja-Service/internal/quiz/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type QuizRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewQuizRepositoryPostgresDB(d *sql.DB) *QuizRepositoryPostgresDB {
	return &QuizRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.QuizzesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" q.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" q.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" q.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *QuizRepositoryPostgresDB) Query(queries *models.QuizzesQueryRequestParams) ([]*domain.Quiz, error) {
	q := `SELECT * FROM quizzes q `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" q.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" q.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var quizzes []*domain.Quiz
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizzes, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			quiz             domain.Quiz
			endDateTime      sql.NullTime
			durationInMinute sql.NullInt64
			password         sql.NullString
			description      sql.NullString
		)
		err := rows.Scan(
			&quiz.ID,
			&quiz.Title,
			&quiz.StartDateTime,
			&endDateTime,
			&durationInMinute,
			&quiz.Status,
			&quiz.AvailableParticipantType,
			&quiz.IsLock,
			&password,
			&description,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
			&quiz.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// EndDateTime
		if endDateTime.Valid {
			quiz.EndDateTime = &endDateTime.Time
		}

		// DurationInMinute
		if durationInMinute.Valid {
			dim := int(durationInMinute.Int64)
			quiz.DurationInMinute = &dim
		}

		// Password
		if password.Valid {
			quiz.Password = &password.String
		}

		// Description
		if description.Valid {
			quiz.Description = &description.String
		}

		// Find availableParticipants
		pRows, err := r.PostgresDB.Query(`
			SELECT
			    qap.id,
			    qap.quizId,
			    qap.userId,
			    u.id AS userId,
			    u.firstName AS userFirstName,
			    u.lastName AS userLastName,
			    u.email AS userEmail,
			    u.avatarUrl AS userAvatarUrl
			FROM quizAvailableParticipants qap
			INNER JOIN users u ON u.id = qap.userId
			WHERE qap.quizId = $1
		`, quiz.ID)
		if err != nil {
			pRows.Close()
			return nil, err
		}
		var qaps []*domain.QuizAvailableParticipant
		for pRows.Next() {
			var (
				qap           domain.QuizAvailableParticipant
				userId        sql.NullInt64
				userFirstName sql.NullString
				userLastName  sql.NullString
				userEmail     sql.NullString
				userAvatarUrl sql.NullString
			)
			err := pRows.Scan(
				&qap.ID,
				&qap.QuizID,
				&qap.UserID,
				&userId,
				&userFirstName,
				&userLastName,
				&userEmail,
				&userAvatarUrl,
			)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			if userId.Valid {
				qap.UserID = uint(userId.Int64)
				qap.User = &domain.QuizAvailableParticipantUser{
					ID: qap.UserID,
				}
				if userFirstName.Valid {
					qap.User.FirstName = userFirstName.String
				}
				if userLastName.Valid {
					qap.User.LastName = userLastName.String
				}
				if userEmail.Valid {
					qap.User.Email = userEmail.String
				}
				if userAvatarUrl.Valid {
					qap.User.AvatarUrl = userAvatarUrl.String
				}
			}
			qaps = append(qaps, &qap)
		}
		pRows.Close()
		quiz.AvailableParticipants = qaps

		// Check if user already started the quiz
		quiz.IsCurrentAuthorizedUserStartedThisQuiz = false
		if queries.AuthenticatedUser != nil {
			var quizParticipantID sql.NullInt64
			err := r.PostgresDB.QueryRow(`
				SELECT qp.id FROM quizParticipants qp
				WHERE qp.quizId = $1 AND qp.userId = $2
			`, quiz.ID, queries.AuthenticatedUser.ID).Scan(&quizParticipantID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					quiz.IsCurrentAuthorizedUserStartedThisQuiz = false
				} else {
					return nil, err
				}
			}
			if quizParticipantID.Valid {
				quiz.IsCurrentAuthorizedUserStartedThisQuiz = true
			}
		}

		// Check if user already ended the quiz
		quiz.IsCurrentAuthorizedUserEndedThisQuiz = false
		if queries.AuthenticatedUser != nil {
			var quizParticipantID sql.NullInt64
			err := r.PostgresDB.QueryRow(`
				SELECT qp.id FROM quizParticipants qp
				WHERE qp.quizId = $1 AND qp.userId = $2 AND qp.ended_at IS NOT NULL
			`, quiz.ID, queries.AuthenticatedUser.ID).Scan(&quizParticipantID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					quiz.IsCurrentAuthorizedUserEndedThisQuiz = false
				} else {
					return nil, err
				}
			}
			if quizParticipantID.Valid {
				quiz.IsCurrentAuthorizedUserEndedThisQuiz = true
			}
		}

		// Find questions count
		var questionsCount int
		err = r.PostgresDB.QueryRow(`
			SELECT COUNT(qq.id) FROM quizQuestions qq
			WHERE qq.quizId = $1
		`, quiz.ID).Scan(&questionsCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				questionsCount = 0
			} else {
				return nil, err
			}
		}
		quiz.QuestionsCount = questionsCount

		// Append to quizzes
		quizzes = append(quizzes, &quiz)
	}
	return quizzes, nil
}

func (r *QuizRepositoryPostgresDB) Count(queries *models.QuizzesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(q.id) FROM quizzes q `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *QuizRepositoryPostgresDB) Create(payload *models.QuizzesCreateRequestBody) (*domain.Quiz, error) {
	var quiz domain.Quiz

	// Current time
	currentTime := time.Now()

	// Find the quiz by title
	var foundQuiz domain.Quiz
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM quizzes
		WHERE title = $1
	`, payload.Title).Scan(&foundQuiz.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the quiz
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(`
		INSERT INTO quizzes (title, startDateTime, endDateTime, durationInMinute, status, availableParticipantType, isLock, password, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 , $12)
		RETURNING id
	`,
		payload.Title,
		payload.StartDateTime,
		payload.EndDateTime,
		payload.DurationInMinute,
		payload.Status,
		payload.AvailableParticipantType,
		payload.IsLockAsBool,
		payload.Password,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(&quiz.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, participantUserID := range payload.ParticipantUserIDsAsInt64 {
		_, err = tx.Exec(`
			INSERT INTO quizAvailableParticipants (quizId, userId, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5)
		`, quiz.ID, participantUserID, currentTime, currentTime, nil)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	quizzes, err := r.Query(&models.QuizzesQueryRequestParams{ID: int(quiz.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		return nil, errors.New("quiz not found")
	}

	// Return the quiz
	return quizzes[0], nil
}

func (r *QuizRepositoryPostgresDB) Delete(payload *models.QuizzesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM quizzes
		WHERE id = ANY ($1)
		RETURNING id
	`, pq.Int64Array(payload.IDsInt64)).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}

func (r *QuizRepositoryPostgresDB) Update(payload *models.QuizzesCreateRequestBody, id int64) (*domain.Quiz, error) {
	var quiz domain.Quiz

	// Current time
	currentTime := time.Now()

	// Find the quiz by id
	var foundQuiz domain.Quiz
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM quizzes
		WHERE id = $1
	`, id).Scan(&foundQuiz.ID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Update the quiz
	err = tx.QueryRow(`
		UPDATE quizzes
		SET title = $1, startDateTime = $2, endDateTime = $3, durationInMinute = $4, status = $5, availableParticipantType = $6, isLock = $7, password = $8, description = $9, updated_at = $10
		WHERE id = $11
		RETURNING id
	`,
		payload.Title,
		payload.StartDateTime,
		payload.EndDateTime,
		payload.DurationInMinute,
		payload.Status,
		payload.AvailableParticipantType,
		payload.IsLockAsBool,
		payload.Password,
		payload.Description,
		currentTime,
		foundQuiz.ID,
	).Scan(&quiz.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete all available participants
	_, err = tx.Exec(`
		DELETE FROM quizAvailableParticipants
		WHERE quizId = $1
	`, foundQuiz.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert new available participants
	for _, participantUserID := range payload.ParticipantUserIDsAsInt64 {
		_, err = tx.Exec(`
			INSERT INTO quizAvailableParticipants (quizId, userId, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5)
		`, quiz.ID, participantUserID, currentTime, currentTime, nil)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz
	quizzes, err := r.Query(&models.QuizzesQueryRequestParams{ID: int(quiz.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		return nil, errors.New("quiz not found")
	}

	// Return the quiz
	return quizzes[0], nil
}

func (r *QuizRepositoryPostgresDB) UpdateStatus(payload *models.QuizzesUpdateStatusRequestBody, id int64) (*domain.Quiz, error) {
	var quiz domain.Quiz

	// Current time
	currentTime := time.Now()

	// Find the quiz by id
	var foundQuiz domain.Quiz
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM quizzes
		WHERE id = $1
	`, id).Scan(&foundQuiz.ID)
	if err != nil {
		return nil, err
	}

	// Update the quiz
	err = r.PostgresDB.QueryRow(`
		UPDATE quizzes
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		payload.Status,
		currentTime,
		foundQuiz.ID,
	).Scan(&quiz.ID)
	if err != nil {
		return nil, err
	}

	// Query the quiz
	quizzes, err := r.Query(&models.QuizzesQueryRequestParams{ID: int(quiz.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		return nil, errors.New("quiz not found")
	}

	// Return the quiz
	return quizzes[0], nil
}

func makeQuizQuestionsWhereFilters(queries *models.QuizzesQueryQuestionsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" qq.id = %d", queries.ID))
		}
		if queries.QuizID != 0 {
			where = append(where, fmt.Sprintf(" qq.quizId = %d", queries.QuizID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" qq.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" qq.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *QuizRepositoryPostgresDB) QueryQuestions(queries *models.QuizzesQueryQuestionsRequestParams) ([]*domain.QuizQuestion, error) {
	q := `
		SELECT
			qq.id,
			qq.quizId,
			qq.title,
			qq.description,
			qq.created_at,
			qq.updated_at,
			qq.deleted_at,
			q.id AS quizId,
			q.title AS quizTitle
		FROM quizQuestions qq
		LEFT JOIN quizzes q ON qq.quizId = q.id
	`
	if queries != nil {
		where := makeQuizQuestionsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" qq.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" qq.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit >= 0 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"
	log.Println(q)

	var quizQuestions []*domain.QuizQuestion
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizQuestions, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			quizQuestion domain.QuizQuestion
			description  sql.NullString
			deletedAt    sql.NullTime
			quizId       sql.NullInt64
			quizTitle    sql.NullString
		)
		err := rows.Scan(
			&quizQuestion.ID,
			&quizQuestion.QuizID,
			&quizQuestion.Title,
			&description,
			&quizQuestion.CreatedAt,
			&quizQuestion.UpdatedAt,
			&deletedAt,
			&quizId,
			&quizTitle,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			qd := quizQuestion.Description
			quizQuestion.Description = qd
		}
		if deletedAt.Valid {
			quizQuestion.DeletedAt = &deletedAt.Time
		}
		if quizId.Valid {
			qid := quizId.Int64
			quizQuestion.Quiz = &domain.QuizQuestionQuiz{
				ID: uint(qid),
			}
			if quizTitle.Valid {
				quizQuestion.Quiz.Title = quizTitle.String
			}
		}

		// Fetch options
		options, err := r.PostgresDB.Query(`
			SELECT id, title, score
			FROM quizQuestionOptions
			WHERE quizQuestionId = $1
		`, quizQuestion.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				options.Close()
				continue
			}
			return nil, err
		}
		if quizQuestion.Options == nil {
			quizQuestion.Options = []*domain.QuizQuestionOptionItem{}
		}
		for options.Next() {
			var option domain.QuizQuestionOptionItem
			err := options.Scan(
				&option.ID,
				&option.Title,
				&option.Score,
			)
			if err != nil {
				options.Close()
				return nil, err
			}
			quizQuestion.Options = append(quizQuestion.Options, &option)
		}

		quizQuestions = append(quizQuestions, &quizQuestion)
	}
	return quizQuestions, nil
}

func (r *QuizRepositoryPostgresDB) CountQuestions(queries *models.QuizzesQueryQuestionsRequestParams) (int64, error) {
	q := `SELECT COUNT(qq.id) FROM quizQuestions qq `
	if queries != nil {
		where := makeQuizQuestionsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *QuizRepositoryPostgresDB) CreateQuestion(payload *models.QuizzesCreateQuestionRequestBody) (*domain.QuizQuestion, error) {
	var quizQuestion domain.QuizQuestion

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Create the quiz question
	err = tx.QueryRow(`
		INSERT INTO quizQuestions (quizId, title, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		payload.QuizID,
		payload.Title,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(&quizQuestion.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert all the quiz question Options
	for _, choice := range payload.Options {
		_, err = tx.Exec(`
			INSERT INTO quizQuestionOptions (quizQuestionId, title, score, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			quizQuestion.ID,
			choice.Title,
			choice.Score,
			currentTime,
			currentTime,
			nil,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz question
	quizQuestions, err := r.QueryQuestions(&models.QuizzesQueryQuestionsRequestParams{ID: int(quizQuestion.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizQuestions) == 0 {
		return nil, errors.New("quiz question not found")
	}

	// Return the quiz question
	return quizQuestions[0], nil
}

func (r *QuizRepositoryPostgresDB) UpdateQuizQuestion(payload *models.QuizzesCreateQuestionRequestBody, questionID int64) (*domain.QuizQuestion, error) {
	var quizQuestion domain.QuizQuestion

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Update the quiz question
	err = tx.QueryRow(`
		UPDATE quizQuestions
		SET quizId = $1, title = $2, description = $3, updated_at = $4
		WHERE id = $5
		RETURNING id
	`,
		payload.QuizID,
		payload.Title,
		payload.Description,
		currentTime,
		questionID,
	).Scan(&quizQuestion.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete all the quiz question Options
	_, err = tx.Exec(`
		DELETE FROM quizQuestionOptions
		WHERE quizQuestionId = $1
	`,
		questionID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert all the quiz question Options
	for _, option := range payload.Options {
		_, err = tx.Exec(`
			INSERT INTO quizQuestionOptions (quizQuestionId, title, score, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			questionID,
			option.Title,
			option.Score,
			currentTime,
			currentTime,
			nil,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz question
	quizQuestions, err := r.QueryQuestions(&models.QuizzesQueryQuestionsRequestParams{ID: int(quizQuestion.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizQuestions) == 0 {
		return nil, errors.New("quiz question not found")
	}

	// Return the quiz question
	return quizQuestions[0], nil
}

func makeQuizParticipantsWhereFilters(queries *models.QuizzesQueryParticipantsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" qp.id = %d", queries.ID))
		}
		if queries.QuizID != 0 {
			where = append(where, fmt.Sprintf(" qp.quizId = %d", queries.QuizID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" qp.userId = %d", queries.UserID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" q.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" qp.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *QuizRepositoryPostgresDB) QueryParticipants(queries *models.QuizzesQueryParticipantsRequestParams) ([]*domain.QuizParticipant, error) {
	q := `
		SELECT
			qp.id,
			qp.quizId,
			qp.userId,
			qp.ended_at,
			qp.created_at,
			qp.updated_at,
			qp.deleted_at,
			q.id AS qId,
			q.title AS qTitle,
			u.id AS uId,
			u.firstName AS uFirstName,
			u.lastName AS uLastName,
			u.email AS uEmail,
			u.avatarUrl AS uAvatarUrl
		FROM quizParticipants qp
		INNER JOIN quizzes q ON q.id = qp.quizId
		INNER JOIN users u ON u.id = qp.userId
	`
	if queries != nil {
		where := makeQuizParticipantsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" q.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" qp.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"
	log.Println(q)

	var quizParticipants []*domain.QuizParticipant
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizParticipants, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			quizParticipant domain.QuizParticipant
			endedAt         sql.NullTime
			deletedAt       sql.NullTime
			quizId          sql.NullInt64
			quizTitle       sql.NullString
			userId          sql.NullInt64
			userFirstName   sql.NullString
			userLastName    sql.NullString
			userEmail       sql.NullString
			userAvatarUrl   sql.NullString
		)
		err := rows.Scan(
			&quizParticipant.ID,
			&quizParticipant.QuizID,
			&quizParticipant.UserID,
			&endedAt,
			&quizParticipant.CreatedAt,
			&quizParticipant.UpdatedAt,
			&deletedAt,
			&quizId,
			&quizTitle,
			&userId,
			&userFirstName,
			&userLastName,
			&userEmail,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if endedAt.Valid {
			quizParticipant.EndedAt = &endedAt.Time
		}
		if deletedAt.Valid {
			quizParticipant.DeletedAt = &deletedAt.Time
		}
		if quizId.Valid {
			quizParticipant.Quiz = &domain.QuizParticipantQuiz{
				ID: uint(quizId.Int64),
			}
			if quizTitle.Valid {
				quizParticipant.Quiz.Title = quizTitle.String
			}
		}
		if userId.Valid {
			quizParticipant.User = &domain.QuizParticipantUser{
				ID: uint(userId.Int64),
			}
			if userFirstName.Valid {
				quizParticipant.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				quizParticipant.User.LastName = userLastName.String
			}
			if userEmail.Valid {
				quizParticipant.User.Email = userEmail.String
			}
			if userAvatarUrl.Valid {
				quizParticipant.User.AvatarUrl = userAvatarUrl.String
			}
		}
		quizParticipants = append(quizParticipants, &quizParticipant)
	}
	return quizParticipants, nil
}

func (r *QuizRepositoryPostgresDB) CountParticipants(queries *models.QuizzesQueryParticipantsRequestParams) (int64, error) {
	q := `SELECT COUNT(qp.id) FROM quizParticipants qp `
	if queries != nil {
		where := makeQuizParticipantsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *QuizRepositoryPostgresDB) StartQuiz(payload *models.QuizzesCreateStartRequestBody) (*domain.QuizParticipant, error) {
	var quizParticipant domain.QuizParticipant

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Create the quiz participant
	err = tx.QueryRow(`
		INSERT INTO quizParticipants (quizId, userId, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		payload.QuizID,
		payload.UserID,
		currentTime,
		currentTime,
		nil,
	).Scan(&quizParticipant.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz participant
	quizParticipants, err := r.QueryParticipants(&models.QuizzesQueryParticipantsRequestParams{ID: int(quizParticipant.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizParticipants) == 0 {
		return nil, errors.New("quiz participant not found")
	}

	// Return the quiz participant
	return quizParticipants[0], nil
}

func makeQuestionOptionsWhereFilters(queries *models.QuizzesQueryQuestionOptionsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" qqo.id = %d", queries.ID))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" qqo.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *QuizRepositoryPostgresDB) QueryQuestionOptions(queries *models.QuizzesQueryQuestionOptionsRequestParams) ([]*domain.QuizQuestionOption, error) {
	q := `
		SELECT
			qqo.id,
			qqo.quizQuestionId,
			qqo.title,
			qqo.score,
			qqo.created_at,
			qqo.updated_at,
			qqo.deleted_at
		FROM quizQuestionOptions qqo
		INNER JOIN quizQuestions qq ON qq.id = qqo.quizQuestionId
		INNER JOIN quizzes q ON q.id = qq.quizId
	`
	if queries != nil {
		where := makeQuestionOptionsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" q.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var quizQuestionOptions []*domain.QuizQuestionOption
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizQuestionOptions, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			quizQuestionOption domain.QuizQuestionOption
			deletedAt          sql.NullTime
		)
		err := rows.Scan(
			&quizQuestionOption.ID,
			&quizQuestionOption.QuizQuestionId,
			&quizQuestionOption.Title,
			&quizQuestionOption.Score,
			&quizQuestionOption.CreatedAt,
			&quizQuestionOption.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			quizQuestionOption.DeletedAt = &deletedAt.Time
		}
		quizQuestionOptions = append(quizQuestionOptions, &quizQuestionOption)
	}
	return quizQuestionOptions, nil
}

func (r *QuizRepositoryPostgresDB) CountQuestionOptions(queries *models.QuizzesQueryQuestionOptionsRequestParams) (int64, error) {
	q := `
	SELECT
		COUNT(qqo.id)
	FROM quizQuestionOptions qqo
	INNER JOIN quizQuestions qq ON qq.id = qqo.quizQuestionId
	INNER JOIN quizzes q ON q.id = qq.quizId
`
	if queries != nil {
		where := makeQuestionOptionsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func makeQuestionAnswersWhereFilters(queries *models.QuizzesQueryQuestionAnswersRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" qqa.id = %d", queries.ID))
		}
		if queries.QuizID != 0 {
			where = append(where, fmt.Sprintf(" qq.quizId = %d", queries.QuizID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" qqa.userId = %d", queries.UserID))
		}
		if queries.QuestionID != 0 {
			where = append(where, fmt.Sprintf(" qqa.questionId = %d", queries.QuestionID))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" qqa.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *QuizRepositoryPostgresDB) QueryAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) ([]*domain.QuizQuestionAnswer, error) {
	q := `
		SELECT
			qqa.id,
			qqa.userId,
			qqa.questionId,
			qqa.quizQuestionOptionId,
			qqa.created_at,
			qqa.updated_at,
			qqa.deleted_at,
			u.id AS userId,
			u.firstName AS userFirstName,
			u.lastName AS userLastName,
			u.email AS userEmail,
			u.avatarUrl AS userAvatarUrl,
			qq.id AS qqQuestionId,
			qq.quizId AS qqQuizId,
			qq.title AS qqTitle,
			qq.description AS qqDescription,
			q.id AS qId,
			q.title AS qTitle,
			q.description AS qDescription,
			qqo.id AS qqoId,
			qqo.quizQuestionId AS qqoQuizQuestionId,
			qqo.title AS qqoTitle,
			qqo.score AS qqoScore
		FROM quizQuestionAnswers qqa
		INNER JOIN users u ON u.id = qqa.userId
		INNER JOIN quizQuestions qq ON qq.id = qqa.questionId
		INNER JOIN quizzes q ON q.id = qq.quizId
		INNER JOIN quizQuestionOptions qqo ON qqo.id = qqa.quizQuestionOptionId
	`
	if queries != nil {
		where := makeQuestionAnswersWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" q.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit >= 0 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"
	log.Printf("QueryAnswers: %s", q)

	var quizQuestionAnswers []*domain.QuizQuestionAnswer
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizQuestionAnswers, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			quizQuestionAnswer domain.QuizQuestionAnswer
			deletedAt          sql.NullTime
			userId             sql.NullInt64
			userFirstName      sql.NullString
			userLastName       sql.NullString
			userEmail          sql.NullString
			userAvatarUrl      sql.NullString
			qqQuestionId       sql.NullInt64
			qqQuizId           sql.NullInt64
			qqTitle            sql.NullString
			qqDescription      sql.NullString
			qId                sql.NullInt64
			qTitle             sql.NullString
			qDescription       sql.NullString
			qqoId              sql.NullInt64
			qqoQuizQuestionId  sql.NullInt64
			qqoTitle           sql.NullString
			qqoScore           sql.NullFloat64
		)
		err := rows.Scan(
			&quizQuestionAnswer.ID,
			&quizQuestionAnswer.UserID,
			&quizQuestionAnswer.QuestionID,
			&quizQuestionAnswer.QuizQuestionOptionID,
			&quizQuestionAnswer.CreatedAt,
			&quizQuestionAnswer.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userEmail,
			&userAvatarUrl,
			&qqQuestionId,
			&qqQuizId,
			&qqTitle,
			&qqDescription,
			&qId,
			&qTitle,
			&qDescription,
			&qqoId,
			&qqoQuizQuestionId,
			&qqoTitle,
			&qqoScore,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			quizQuestionAnswer.DeletedAt = &deletedAt.Time
		}
		if userId.Valid {
			quizQuestionAnswer.User = &domain.QuizQuestionAnswerUser{
				ID: uint(userId.Int64),
			}
			if userFirstName.Valid {
				quizQuestionAnswer.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				quizQuestionAnswer.User.LastName = userLastName.String
			}
			if userEmail.Valid {
				quizQuestionAnswer.User.Email = userEmail.String
			}
			if userAvatarUrl.Valid {
				quizQuestionAnswer.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if qqQuestionId.Valid {
			quizQuestionAnswer.Question = &domain.QuizQuestionAnswerQuestion{
				ID: uint(qqQuestionId.Int64),
			}
			if qqQuizId.Valid {
				quizQuestionAnswer.Question.QuizID = uint(qqQuizId.Int64)
			}
			if qqTitle.Valid {
				quizQuestionAnswer.Question.Title = qqTitle.String
			}
			if qqDescription.Valid {
				quizQuestionAnswer.Question.Description = &qqDescription.String
			}
			if qId.Valid {
				quizQuestionAnswer.Question.Quiz = &domain.QuizQuestionAnswerQuestionQuiz{
					ID: uint(qId.Int64),
				}
				if qTitle.Valid {
					quizQuestionAnswer.Question.Quiz.Title = qTitle.String
				}
				if qDescription.Valid {
					quizQuestionAnswer.Question.Quiz.Description = &qDescription.String
				}
			}
		}
		if qqoId.Valid {
			quizQuestionAnswer.QuizQuestionOption = &domain.QuizQuestionAnswerQuizQuestionOption{
				ID: uint(qqoId.Int64),
			}
			if qqoQuizQuestionId.Valid {
				quizQuestionAnswer.QuizQuestionOption.QuizQuestionID = uint(qqoQuizQuestionId.Int64)
			}
			if qqoTitle.Valid {
				quizQuestionAnswer.QuizQuestionOption.Title = qqoTitle.String
			}
			if qqoScore.Valid {
				quizQuestionAnswer.QuizQuestionOption.Score = uint(qqoScore.Float64)
			}
		}
		quizQuestionAnswers = append(quizQuestionAnswers, &quizQuestionAnswer)
	}
	return quizQuestionAnswers, nil
}

func (r *QuizRepositoryPostgresDB) CountAnswers(queries *models.QuizzesQueryQuestionAnswersRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(qqa.id)
		FROM quizQuestionAnswers qqa
		INNER JOIN users u ON u.id = qqa.userId
		INNER JOIN quizQuestions qq ON qq.id = qqa.questionId
		INNER JOIN quizzes q ON q.id = qq.quizId
		INNER JOIN quizQuestionOptions qqo ON qqo.id = qqa.quizQuestionOptionId
	`
	if queries != nil {
		where := makeQuestionAnswersWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *QuizRepositoryPostgresDB) CreateAnswer(payload *models.QuizzesCreateAnswerRequestBody) (*domain.QuizQuestionAnswer, error) {
	var quizQuestionAnswer domain.QuizQuestionAnswer

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Find quizQuestionAnswers and update if exists or create if not exists
	quizQuestionAnswers, err := r.QueryAnswers(&models.QuizzesQueryQuestionAnswersRequestParams{QuestionID: payload.QuestionID, UserID: int64(payload.User.ID)})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(quizQuestionAnswers) > 0 {
		quizQuestionAnswer = *quizQuestionAnswers[0]
		quizQuestionAnswer.QuizQuestionOptionID = uint(payload.QuizQuestionOptionID)
		quizQuestionAnswer.UpdatedAt = currentTime
		quizQuestionAnswer.DeletedAt = nil
		_, err = tx.Exec(`
			UPDATE quizQuestionAnswers
			SET quizQuestionOptionId = $1, updated_at = $2, deleted_at = $3
			WHERE id = $4
		`,
			quizQuestionAnswer.QuizQuestionOptionID,
			quizQuestionAnswer.UpdatedAt,
			quizQuestionAnswer.DeletedAt,
			quizQuestionAnswer.ID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		quizQuestionAnswer = domain.QuizQuestionAnswer{
			UserID:               payload.User.ID,
			QuestionID:           uint(payload.QuestionID),
			QuizQuestionOptionID: uint(payload.QuizQuestionOptionID),
			CreatedAt:            currentTime,
			UpdatedAt:            currentTime,
			DeletedAt:            nil,
		}
		_, err = tx.Exec(`
			INSERT INTO quizQuestionAnswers (userId, questionId, quizQuestionOptionId, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			quizQuestionAnswer.UserID,
			quizQuestionAnswer.QuestionID,
			quizQuestionAnswer.QuizQuestionOptionID,
			quizQuestionAnswer.CreatedAt,
			quizQuestionAnswer.UpdatedAt,
			quizQuestionAnswer.DeletedAt,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz question answer
	quizQuestionAnswers, err = r.QueryAnswers(&models.QuizzesQueryQuestionAnswersRequestParams{ID: int(quizQuestionAnswer.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizQuestionAnswers) == 0 {
		return nil, errors.New("quiz question answer not found")
	}

	// Return the quiz question answer
	return quizQuestionAnswers[0], nil
}

func (r *QuizRepositoryPostgresDB) UpdateQuizEnd(payload *models.QuizzesUpdateEndRequestBody, quizId int64) (*domain.QuizParticipant, error) {
	var quizParticipant domain.QuizParticipant

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Find quizParticipants and update if exists or create if not exists
	quizParticipants, err := r.QueryParticipants(&models.QuizzesQueryParticipantsRequestParams{QuizID: quizId, UserID: int64(payload.AuthenticatedUser.ID)})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(quizParticipants) <= 0 {
		tx.Rollback()
		return nil, errors.New("you are not a participant of this quiz")
	}

	quizParticipant = *quizParticipants[0]
	quizParticipant.EndedAt = &currentTime
	quizParticipant.UpdatedAt = currentTime
	quizParticipant.DeletedAt = nil
	_, err = tx.Exec(`
		UPDATE quizParticipants
		SET ended_at = $1, updated_at = $2, deleted_at = $3
		WHERE id = $4
	`,
		quizParticipant.EndedAt,
		quizParticipant.UpdatedAt,
		quizParticipant.DeletedAt,
		quizParticipant.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the quiz participant
	quizParticipants, err = r.QueryParticipants(&models.QuizzesQueryParticipantsRequestParams{ID: int(quizParticipant.ID)})
	if err != nil {
		return nil, err
	}
	if len(quizParticipants) == 0 {
		return nil, errors.New("quiz participant not found")
	}

	// Return the quiz participant
	return quizParticipants[0], nil
}
