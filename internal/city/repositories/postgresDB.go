package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/city/domain"
	"github.com/hoitek/Maja-Service/internal/city/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type CityRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewCityRepositoryPostgresDB(d *sql.DB) *CityRepositoryPostgresDB {
	return &CityRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.CitiesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" c.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" c.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *CityRepositoryPostgresDB) Query(queries *models.CitiesQueryRequestParams) ([]*domain.City, error) {
	q := `SELECT * FROM cities c `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var cities []*domain.City
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var city domain.City
		err := rows.Scan(&city.ID, &city.Name, &city.CreatedAt, &city.UpdatedAt, &city.DeletedAt)
		if err != nil {
			return nil, err
		}
		cities = append(cities, &city)
	}
	return cities, nil
}

func (r *CityRepositoryPostgresDB) Count(queries *models.CitiesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM cities c `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CityRepositoryPostgresDB) Create(payload *models.CitiesCreateRequestBody) (*domain.City, error) {
	var city domain.City

	// Current time
	currentTime := time.Now()

	// Find the city by name
	var foundCity domain.City
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM cities
WHERE name = $1
`, payload.Name).Scan(&foundCity.ID, &foundCity.Name, &foundCity.CreatedAt, &foundCity.UpdatedAt, &foundCity.DeletedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the city
	err = r.PostgresDB.QueryRow(`
INSERT INTO cities (name, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4)
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, currentTime, nil).Scan(&city.ID, &city.Name, &city.CreatedAt, &city.UpdatedAt, &city.DeletedAt)
	// If the city does not insert, return an error
	if err != nil {
		return nil, err
	}

	// Return the city
	return &city, nil
}

func (r *CityRepositoryPostgresDB) Delete(payload *models.CitiesDeleteRequestBody) ([]int64, error) {
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
DELETE FROM cities
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

func (r *CityRepositoryPostgresDB) Update(payload *models.CitiesCreateRequestBody, name string) (*domain.City, error) {
	var city domain.City

	// Current time
	currentTime := time.Now()

	// Find the city by name
	var foundCity domain.City
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM cities
WHERE name = $1
`, name).Scan(&foundCity.ID, &foundCity.Name, &foundCity.CreatedAt, &foundCity.UpdatedAt, &foundCity.DeletedAt)

	// If the city is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the city
	err = r.PostgresDB.QueryRow(`
UPDATE cities
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, foundCity.ID).Scan(&city.ID, &city.Name, &city.CreatedAt, &city.UpdatedAt, &city.DeletedAt)

	// If the city does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the city
	return &city, nil
}
