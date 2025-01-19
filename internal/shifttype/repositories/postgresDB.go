package repositories

import (
	"database/sql"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/shifttype/domain"
	"github.com/hoitek/Maja-Service/internal/shifttype/models"
	"github.com/lib/pq"
)

type ShiftTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewShiftTypeRepositoryPostgresDB(d *sql.DB) *ShiftTypeRepositoryPostgresDB {
	return &ShiftTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *ShiftTypeRepositoryPostgresDB) Query(queries *models.ShiftTypesQueryRequestParams) ([]*domain.ShiftType, error) {
	q := `SELECT * FROM shiftTypes `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit

		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return nil, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var shiftTypes []*domain.ShiftType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ShiftType domain.ShiftType
		err := rows.Scan(&ShiftType.ID, &ShiftType.Name, &ShiftType.CreatedAt, &ShiftType.UpdatedAt, &ShiftType.DeletedAt)
		if err != nil {
			return nil, err
		}
		shiftTypes = append(shiftTypes, &ShiftType)
	}
	return shiftTypes, nil
}

func (r *ShiftTypeRepositoryPostgresDB) Count(queries *models.ShiftTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM shiftTypes `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit
		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return 0, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ShiftTypeRepositoryPostgresDB) GetShiftTypesByIds(ids []int64) ([]*domain.ShiftType, error) {
	q := `SELECT * FROM shiftTypes WHERE id = ANY($1)`
	var shiftTypes []*domain.ShiftType
	rows, err := r.PostgresDB.Query(q, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ShiftType domain.ShiftType
		err := rows.Scan(&ShiftType.ID, &ShiftType.Name, &ShiftType.CreatedAt, &ShiftType.UpdatedAt, &ShiftType.DeletedAt)
		if err != nil {
			return nil, err
		}
		shiftTypes = append(shiftTypes, &ShiftType)
	}
	return shiftTypes, nil
}
