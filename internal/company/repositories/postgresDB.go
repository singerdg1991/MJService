package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/company/domain"
	"github.com/hoitek/Maja-Service/internal/company/models"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type CompanyRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewCompanyRepositoryPostgresDB(d *sql.DB) *CompanyRepositoryPostgresDB {
	return &CompanyRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *CompanyRepositoryPostgresDB) Query(queries *models.CompaniesQueryRequestParams) ([]*domain.Company, error) {
	q := `SELECT * FROM companies `
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

	var companies []*domain.Company
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var company domain.Company
		err := rows.Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt, &company.DeletedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}
	return companies, nil
}

func (r *CompanyRepositoryPostgresDB) Count(queries *models.CompaniesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM companies `
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

func (r *CompanyRepositoryPostgresDB) Create(payload *models.CompaniesCreateRequestBody) (*domain.Company, error) {
	var company domain.Company

	// Current time
	currentTime := time.Now()

	// Find the company by name
	var foundCompany domain.Company
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM companies
		WHERE name = $1
	`, payload.Name).Scan(&foundCompany.ID, &foundCompany.Name, &foundCompany.CreatedAt, &foundCompany.UpdatedAt, &foundCompany.DeletedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the company
	err = r.PostgresDB.QueryRow(`
		INSERT INTO companies (name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, currentTime, nil).Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt, &company.DeletedAt)
	// If the company does not insert, return an error
	if err != nil {
		return nil, err
	}

	// Return the company
	return &company, nil
}

func (r *CompanyRepositoryPostgresDB) Delete(payload *models.CompaniesDeleteRequestBody) ([]int64, error) {
	var ids []int64
	idsStr := strings.Split(payload.IDs, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int64(id))
	}
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM companies
		WHERE id = ANY ($1)
		RETURNING id
	`, pq.Int64Array(ids)).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return ids, nil
}

func (r *CompanyRepositoryPostgresDB) Update(payload *models.CompaniesCreateRequestBody, name string) (*domain.Company, error) {
	var company domain.Company

	// Current time
	currentTime := time.Now()

	// Find the company by name
	var foundCompany domain.Company
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM companies
		WHERE name = $1
	`, name).Scan(&foundCompany.ID, &foundCompany.Name, &foundCompany.CreatedAt, &foundCompany.UpdatedAt, &foundCompany.DeletedAt)

	// If the company is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the company
	err = r.PostgresDB.QueryRow(`
		UPDATE companies
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, foundCompany.ID).Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt, &company.DeletedAt)

	// If the company does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the company
	return &company, nil
}
