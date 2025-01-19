package repositories

import (
	"database/sql"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/stafftype/domain"
	"github.com/hoitek/Maja-Service/internal/stafftype/models"
	"github.com/lib/pq"
	"log"
)

type StaffTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewStaffTypeRepositoryPostgresDB(d *sql.DB) *StaffTypeRepositoryPostgresDB {
	return &StaffTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *StaffTypeRepositoryPostgresDB) Query(queries *models.StaffTypesQueryRequestParams) ([]*domain.StaffType, error) {
	q := `SELECT * FROM staffTypes `
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

	var staffTypes []*domain.StaffType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var staffType domain.StaffType
		err := rows.Scan(&staffType.ID, &staffType.Name, &staffType.CreatedAt, &staffType.UpdatedAt, &staffType.DeletedAt)
		if err != nil {
			return nil, err
		}
		staffTypes = append(staffTypes, &staffType)
	}
	return staffTypes, nil
}

func (r *StaffTypeRepositoryPostgresDB) Count(queries *models.StaffTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM staffTypes `
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

func (r *StaffTypeRepositoryPostgresDB) GetStaffTypesByIds(ids []int64) ([]*domain.StaffType, error) {
	q := `SELECT * FROM staffTypes WHERE id = ANY($1)`
	var staffTypes []*domain.StaffType
	rows, err := r.PostgresDB.Query(q, pq.Int64Array(ids))
	log.Printf("staffTypes: %#v", err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var staffType domain.StaffType
		err := rows.Scan(&staffType.ID, &staffType.Name, &staffType.CreatedAt, &staffType.UpdatedAt, &staffType.DeletedAt)
		if err != nil {
			return nil, err
		}
		staffTypes = append(staffTypes, &staffType)
	}
	return staffTypes, nil
}
