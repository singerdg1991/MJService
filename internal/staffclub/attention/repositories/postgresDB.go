package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type AttentionRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewAttentionRepositoryPostgresDB(d *sql.DB) *AttentionRepositoryPostgresDB {
	return &AttentionRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.AttentionsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" g.id = %d", queries.ID))
		}
		if queries.Filters.AttentionNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AttentionNumber.Op, fmt.Sprintf("%v", queries.Filters.AttentionNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.attentionNumber %s %s", opValue.Operator, val))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *AttentionRepositoryPostgresDB) Query(queries *models.AttentionsQueryRequestParams) ([]*domain.Attention, error) {
	q := `
		SELECT
		    g.id,
		    g.punishmentId,
			g.isAutoRewardSetEnable,
			g.attentionNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS punishmentId,
			r.name AS punishmentName,
			r.description AS punishmentDescription
		FROM staffClubAttentions g
		LEFT JOIN punishments r ON g.punishmentId = r.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var attentions []*domain.Attention
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			attention             domain.Attention
			description           sql.NullString
			deletedAt             pq.NullTime
			punishmentId          sql.NullInt64
			punishmentName        sql.NullString
			punishmentDescription sql.NullString
		)
		err := rows.Scan(
			&attention.ID,
			&attention.PunishmentID,
			&attention.IsAutoRewardSetEnable,
			&attention.AttentionNumber,
			&attention.Title,
			&description,
			&attention.CreatedAt,
			&attention.UpdatedAt,
			&deletedAt,
			&punishmentId,
			&punishmentName,
			&punishmentDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			attention.Description = &description.String
		}
		if deletedAt.Valid {
			attention.DeletedAt = &deletedAt.Time
		}
		if punishmentId.Valid {
			attention.Punishment = &domain.AttentionPunishment{
				ID: uint(punishmentId.Int64),
			}
			if punishmentName.Valid {
				attention.Punishment.Name = punishmentName.String
			}
			if punishmentDescription.Valid {
				attention.Punishment.Description = punishmentDescription.String
			}
		}
		attentions = append(attentions, &attention)
	}
	return attentions, nil
}

func (r *AttentionRepositoryPostgresDB) Count(queries *models.AttentionsQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(g.id)
		FROM staffClubAttentions g
		LEFT JOIN punishments r ON g.punishmentId = r.id
	`
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

func (r *AttentionRepositoryPostgresDB) Create(payload *models.AttentionsCreateRequestBody) (*domain.Attention, error) {
	// Current time
	currentTime := time.Now()

	// Insert the attention
	var insertedId int
	err := r.PostgresDB.QueryRow(`
		INSERT INTO staffClubAttentions (punishmentId, isAutoRewardSetEnable, attentionNumber, title, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.PunishmentID,
		payload.IsAutoRewardSetEnableAsBool,
		payload.AttentionNumber,
		payload.Title,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the attention
	attentions, err := r.Query(&models.AttentionsQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(attentions) == 0 {
		return nil, errors.New("no rows affected")
	}
	return attentions[0], nil
}

func (r *AttentionRepositoryPostgresDB) Delete(payload *models.AttentionsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM staffClubAttentions
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

func (r *AttentionRepositoryPostgresDB) Update(payload *models.AttentionsCreateRequestBody, id int64) (*domain.Attention, error) {
	// Current time
	currentTime := time.Now()

	// Find the attention by id
	foundAttentions, _ := r.Query(&models.AttentionsQueryRequestParams{ID: int(id)})
	if len(foundAttentions) == 0 {
		return nil, errors.New("attention not found")
	}
	foundAttention := foundAttentions[0]
	if foundAttention.DeletedAt != nil {
		return nil, errors.New("attention is deleted")
	}

	// Update the attention
	var updatedId int
	err := r.PostgresDB.QueryRow(`
		UPDATE staffClubAttentions
		SET punishmentId = $1, isAutoRewardSetEnable = $2, attentionNumber = $3, title = $4, description = $5, updated_at = $6
		WHERE id = $7
		RETURNING id
	`, payload.PunishmentID, payload.IsAutoRewardSetEnableAsBool, payload.AttentionNumber, payload.Title, payload.Description, currentTime, id).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the attention
	attentions, err := r.Query(&models.AttentionsQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(attentions) == 0 {
		return nil, errors.New("no rows affected")
	}
	return attentions[0], nil
}

func (r *AttentionRepositoryPostgresDB) GetAttentionsByIds(ids []int64) ([]*domain.Attention, error) {
	var attentions []*domain.Attention
	rows, err := r.PostgresDB.Query(`
		SELECT
		    g.id,
		    g.punishmentId,
			g.isAutoRewardSetEnable,
			g.attentionNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS punishmentId,
			r.name AS punishmentName,
			r.description AS punishmentDescription
		FROM staffClubAttentions g
		LEFT JOIN punishments r ON g.punishmentId = r.id
		WHERE g.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			attention             domain.Attention
			description           sql.NullString
			deletedAt             pq.NullTime
			punishmentId          sql.NullInt64
			punishmentName        sql.NullString
			punishmentDescription sql.NullString
		)
		err := rows.Scan(
			&attention.ID,
			&attention.PunishmentID,
			&attention.IsAutoRewardSetEnable,
			&attention.AttentionNumber,
			&attention.Title,
			&description,
			&attention.CreatedAt,
			&attention.UpdatedAt,
			&deletedAt,
			&punishmentId,
			&punishmentName,
			&punishmentDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			attention.Description = &description.String
		}
		if deletedAt.Valid {
			attention.DeletedAt = &deletedAt.Time
		}
		if punishmentId.Valid {
			attention.Punishment = &domain.AttentionPunishment{
				ID: uint(punishmentId.Int64),
			}
			if punishmentName.Valid {
				attention.Punishment.Name = punishmentName.String
			}
			if punishmentDescription.Valid {
				attention.Punishment.Description = punishmentDescription.String
			}
		}
		attentions = append(attentions, &attention)
	}
	return attentions, nil
}
