package repositories

import (
	"database/sql"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/paymenttype/domain"
	"github.com/hoitek/Maja-Service/internal/paymenttype/models"
)

type PaymentTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewPaymentTypeRepositoryPostgresDB(d *sql.DB) *PaymentTypeRepositoryPostgresDB {
	return &PaymentTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *PaymentTypeRepositoryPostgresDB) Query(queries *models.PaymentTypesQueryRequestParams) ([]*domain.PaymentType, error) {
	q := `SELECT * FROM paymentTypes `
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

	var paymentTypes []*domain.PaymentType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var paymentType domain.PaymentType
		err := rows.Scan(&paymentType.ID, &paymentType.Name, &paymentType.CreatedAt, &paymentType.UpdatedAt, &paymentType.DeletedAt)
		if err != nil {
			return nil, err
		}
		paymentTypes = append(paymentTypes, &paymentType)
	}
	return paymentTypes, nil
}

func (r *PaymentTypeRepositoryPostgresDB) Count(queries *models.PaymentTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM paymentTypes `
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
