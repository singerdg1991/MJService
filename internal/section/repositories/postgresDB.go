package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/section/constants"
	"github.com/hoitek/Maja-Service/internal/section/domain"
	"github.com/hoitek/Maja-Service/internal/section/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type SectionRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewSectionRepositoryPostgresDB(d *sql.DB) *SectionRepositoryPostgresDB {
	return &SectionRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.SectionsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" s.id = %d", queries.ID))
		}
		if queries.Type != constants.SECTION_TYPE_PARENT {
			if queries.ParentID != 0 {
				where = append(where, fmt.Sprintf(" s.parentId = %d", queries.ParentID))
			}
		}
		if queries.Type == constants.SECTION_TYPE_PARENT {
			where = append(where, " s.parentId IS NULL")
		}
		if queries.Type == constants.SECTION_TYPE_CHILDREN {
			where = append(where, " s.parentId IS NOT NULL")
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" s.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" s.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *SectionRepositoryPostgresDB) Query(queries *models.SectionsQueryRequestParams) ([]*domain.Section, error) {
	q := `
		SELECT * FROM sections s
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" s.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var sections []*domain.Section
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			section     domain.Section
			parentID    sql.NullInt64
			deletedAt   sql.NullTime
			color       sql.NullString
			description sql.NullString
		)
		err := rows.Scan(
			&section.ID,
			&parentID,
			&section.Name,
			&color,
			&description,
			&section.CreatedAt,
			&section.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			section.DeletedAt = &deletedAt.Time
		}
		if color.Valid {
			section.Color = &color.String
		}
		if description.Valid {
			section.Description = &description.String
		}
		if parentID.Valid && parentID.Int64 > 0 {
			var (
				parent           domain.Section
				childParentID    sql.NullInt64
				childColor       sql.NullString
				childDescription sql.NullString
				childDeletedAt   sql.NullTime
			)
			err := r.PostgresDB.QueryRow(`
				SELECT * FROM sections WHERE id = $1
			`, parentID.Int64).Scan(
				&parent.ID,
				&childParentID,
				&parent.Name,
				&childColor,
				&childDescription,
				&parent.CreatedAt,
				&parent.UpdatedAt,
				&childDeletedAt,
			)
			if err != nil {
				return nil, err
			}
			parentId := int64(parent.ID)
			if parentId > 0 {
				if childParentID.Valid && childParentID.Int64 > 0 {
					parent.ParentID = &childParentID.Int64
				}
				if childDeletedAt.Valid {
					parent.DeletedAt = &childDeletedAt.Time
				}
				if childColor.Valid {
					parent.Color = &childColor.String
				}
				if childDescription.Valid {
					parent.Description = &childDescription.String
				}
				section.ParentID = &parentId
				section.Parent = &parent
			}
		}

		// Get children
		rows, err := r.PostgresDB.Query(`
			SELECT * FROM sections WHERE parentId = $1
		`, section.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var (
				child       domain.Section
				parentID    sql.NullInt64
				color       sql.NullString
				description sql.NullString
				deletedAt   sql.NullTime
			)
			err := rows.Scan(
				&child.ID,
				&parentID,
				&child.Name,
				&color,
				&description,
				&child.CreatedAt,
				&child.UpdatedAt,
				&deletedAt,
			)
			if err != nil {
				return nil, err
			}
			if parentID.Valid && parentID.Int64 > 0 {
				child.ParentID = &parentID.Int64
				// Get parent
				var (
					parent           domain.Section
					childParentID    sql.NullInt64
					childColor       sql.NullString
					childDescription sql.NullString
					childDeletedAt   sql.NullTime
				)
				err := r.PostgresDB.QueryRow(`
					SELECT * FROM sections WHERE id = $1
				`, parentID.Int64).Scan(&parent.ID, &childParentID, &parent.Name, &childColor, &childDescription, &parent.CreatedAt, &parent.UpdatedAt, &childDeletedAt)
				if err != nil {
					return nil, err
				}
				if childParentID.Valid && childParentID.Int64 > 0 {
					parent.ParentID = &childParentID.Int64
				}
				if childDeletedAt.Valid {
					parent.DeletedAt = &childDeletedAt.Time
				}
				if childColor.Valid {
					parent.Color = &childColor.String
				}
				if childDescription.Valid {
					parent.Description = &description.String
				}
				child.Parent = &parent
			}
			if color.Valid {
				child.Color = &color.String
			}
			if description.Valid {
				child.Description = &description.String
			}
			if deletedAt.Valid {
				child.DeletedAt = &deletedAt.Time
			}
			section.Children = append(section.Children, &child)
		}
		rows.Close()
		sections = append(sections, &section)
	}
	return sections, nil
}

func (r *SectionRepositoryPostgresDB) Count(queries *models.SectionsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(s.id) FROM sections s `
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

func (r *SectionRepositoryPostgresDB) Create(payload *models.SectionsCreateRequestBody) (*domain.Section, error) {
	var section domain.Section

	// Current time
	currentTime := time.Now()

	// Insert the section
	err := r.PostgresDB.QueryRow(`
		INSERT INTO sections (parentId, name, color, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, parentId, name, color, description, created_at, updated_at, deleted_at
	`, payload.ParentID, payload.Name, payload.Color, payload.Description, currentTime, currentTime, nil).Scan(&section.ID, &section.ParentID, &section.Name, &section.Color, &section.Description, &section.CreatedAt, &section.UpdatedAt, &section.DeletedAt)

	// If the section does not insert, return an error
	if err != nil {
		return nil, err
	}

	parentId := 0
	if section.ParentID != nil {
		parentId = int(*section.ParentID)
	}
	sections, err := r.Query(&models.SectionsQueryRequestParams{
		ID:       int(section.ID),
		ParentID: parentId,
	})
	if err != nil {
		return nil, err
	}
	if len(sections) == 0 {
		return nil, errors.New("section not found")
	}
	section = *sections[0]

	// Return the section
	return &section, nil
}

func (r *SectionRepositoryPostgresDB) Delete(payload *models.SectionsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM sections
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

func (r *SectionRepositoryPostgresDB) Update(payload *models.SectionsCreateRequestBody, id int) (*domain.Section, error) {
	var section domain.Section

	// Current time
	currentTime := time.Now()

	// Find the section by name
	var foundSection domain.Section
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM sections
		WHERE id = $1
	`, id).Scan(&foundSection.ID, &foundSection.ParentID, &foundSection.Name, &foundSection.Color, &foundSection.Description, &foundSection.CreatedAt, &foundSection.UpdatedAt, &foundSection.DeletedAt)

	// If the section is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	var (
		parentId    sql.NullInt64
		color       sql.NullString
		description sql.NullString
	)

	// Update the section
	err = r.PostgresDB.QueryRow(`
		UPDATE sections
		SET parentId = $1, name = $2, color = $3, description = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, parentId, name, color, description, created_at, updated_at, deleted_at
	`, payload.ParentID, payload.Name, payload.Color, payload.Description, currentTime, foundSection.ID).Scan(&section.ID, &parentId, &section.Name, &color, &description, &section.CreatedAt, &section.UpdatedAt, &section.DeletedAt)

	// If the section does not update, return an error
	if err != nil {
		return nil, err
	}

	if parentId.Valid {
		section.ParentID = &parentId.Int64
	}
	if color.Valid {
		section.Color = &color.String
	}
	if description.Valid {
		section.Description = &description.String
	}

	// Return the section
	return &section, nil
}

func (r *SectionRepositoryPostgresDB) GetSectionsByIds(ids []int64) ([]*domain.Section, error) {
	var sections []*domain.Section
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM sections
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			section     domain.Section
			parentId    sql.NullInt64
			color       sql.NullString
			description sql.NullString
		)
		err := rows.Scan(&section.ID, &parentId, &section.Name, &color, &description, &section.CreatedAt, &section.UpdatedAt, &section.DeletedAt)
		if err != nil {
			return nil, err
		}
		if parentId.Valid {
			section.ParentID = &parentId.Int64
		}
		if color.Valid {
			section.Color = &color.String
		}
		if description.Valid {
			section.Description = &description.String
		}
		sections = append(sections, &section)
	}
	return sections, nil
}
