package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Kit/exp"
	sharedconstants "github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	csDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	cmodels "github.com/hoitek/Maja-Service/internal/customer/models"
	csPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
	sgDomain "github.com/hoitek/Maja-Service/internal/servicegrade/domain"
	sgmodels "github.com/hoitek/Maja-Service/internal/servicegrade/models"
	sgPorts "github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type CycleRepositoryPostgresDB struct {
	PostgresDB          *sql.DB
	CustomerServices    csPorts.CustomerService
	ServiceGradeService sgPorts.ServiceGradeService
}

// NewCycleRepositoryPostgresDB creates a new instance of CycleRepositoryPostgresDB.
// It takes a database connection (d), a CustomerService (css), and a ServiceGradeService (sgs) as parameters.
// The function returns a pointer to a CycleRepositoryPostgresDB struct.
func NewCycleRepositoryPostgresDB(d *sql.DB, css csPorts.CustomerService, sgs sgPorts.ServiceGradeService) *CycleRepositoryPostgresDB {
	return &CycleRepositoryPostgresDB{
		PostgresDB:          d,
		CustomerServices:    css,
		ServiceGradeService: sgs,
	}
}

// makeQueryWhereFilters generates a slice of SQL WHERE conditions based on the provided query parameters.
// It handles different operators and values for filtering the cycles.
func makeQueryWhereFilters(queries *models.CyclesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" c.id = %d ", queries.ID))
		}
		if queries.SectionID != 0 {
			where = append(where, fmt.Sprintf(" c.section_id = %d ", queries.SectionID))
		}
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
		if queries.Filters.StartDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartDate.Op, fmt.Sprintf("%v", queries.Filters.StartDate.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" c.start_datetime %s %s", opValue.Operator, val))
		}
	}
	return where
}

// Query retrieves cycles based on the provided queries.
// It includes filtering, sorting, and pagination.
//
// Parameters:
// - queries: A pointer to a CyclesQueryRequestParams struct containing the query parameters.
//
// Returns:
// - A slice of pointers to Cycle structs representing the retrieved cycles.
// - An error if any occurred during the query.
func (r *CycleRepositoryPostgresDB) Query(queries *models.CyclesQueryRequestParams) ([]*domain.Cycle, error) {
	q := `SELECT * FROM cycles c`
	if queries != nil {
		where := makeQueryWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" c.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" c.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" c.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var cycles []*domain.Cycle
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cycle = &domain.Cycle{}
		err := rows.Scan(&cycle.ID, &cycle.Name, &cycle.CreatedAt, &cycle.UpdatedAt, &cycle.DeletedAt, &cycle.SectionID, &cycle.StartDate, &cycle.EndDate, &cycle.PeriodLength, &cycle.ShiftMorningStartTime, &cycle.ShiftMorningEndTime, &cycle.ShiftEveningStartTime, &cycle.ShiftEveningEndTime, &cycle.ShiftNightStartTime, &cycle.ShiftNightEndTime, &cycle.FreezePeriodDate, &cycle.WishDays)
		if err != nil {
			return nil, err
		}

		// Query staff types
		cycleStaffTypes, err := r.QueryStaffTypes(&models.CyclesQueryStaffTypesRequestParams{
			CycleID: int(cycle.ID),
		})
		if err != nil {
			return nil, err
		}
		for _, cycleStaffType := range cycleStaffTypes {
			cycle.StaffTypes = append(cycle.StaffTypes, *cycleStaffType)
		}

		// Query next staff types
		cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
			CurrentCycleID: int(cycle.ID),
		})
		if err != nil {
			return nil, err
		}
		for _, cycleNextStaffType := range cycleNextStaffTypes {
			cycle.NextStaffTypes = append(cycle.NextStaffTypes, *cycleNextStaffType)
		}

		// Set default status
		cycle.SetDefaultStatus()

		// Append to cycles
		cycles = append(cycles, cycle)
	}
	return cycles, nil
}

// Count returns the total number of cycles that match the given query parameters.
//
// Parameters:
// - queries: A pointer to a CyclesQueryRequestParams struct containing the query parameters.
//
// Returns:
// - An int64 representing the total number of cycles that match the given query parameters.
// - An error if any error occurs during the database query.
func (r *CycleRepositoryPostgresDB) Count(queries *models.CyclesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(c.id) FROM cycles c `
	if queries != nil {
		where := makeQueryWhereFilters(queries)
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

// makeStaffTypesWhereFilters generates the WHERE filters based on the provided CyclesQueryStaffTypesRequestParams.
//
// Parameters:
// - queries *models.CyclesQueryStaffTypesRequestParams: The query parameters used to filter staff types.
// Return type(s):
// []string: A slice of strings representing the WHERE filters.
func makeStaffTypesWhereFilters(queries *models.CyclesQueryStaffTypesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" st.id = %d ", queries.ID))
		}
		if queries.CycleID != 0 {
			where = append(where, fmt.Sprintf(" st.cycleId = %d ", queries.CycleID))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" st.shiftName %s %s", opValue.Operator, val))
		}
		if queries.Filters.DateRangeFrom.Op != "" && queries.Filters.DateRangeTo.Op != "" {
			opValueFrom := utils.GetDBOperatorAndValue(queries.Filters.DateRangeFrom.Op, fmt.Sprintf("%v", queries.Filters.DateRangeFrom.Value))
			opValueTo := utils.GetDBOperatorAndValue(queries.Filters.DateRangeTo.Op, fmt.Sprintf("%v", queries.Filters.DateRangeTo.Value))
			valFrom := exp.TerIf(opValueFrom.Value == "", "", fmt.Sprintf("'%s'", opValueFrom.Value))
			valTo := exp.TerIf(opValueTo.Value == "", "", fmt.Sprintf("'%s'", opValueTo.Value))
			where = append(where, fmt.Sprintf(" st.datetime >= %s AND st.datetime <= %s ", valFrom, valTo))
		} else {
			if queries.Filters.DateTime.Op != "" {
				opValue := utils.GetDBOperatorAndValue(queries.Filters.DateTime.Op, fmt.Sprintf("%v", queries.Filters.DateTime.Value))
				val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
				where = append(where, fmt.Sprintf(" st.datetime %s %s", opValue.Operator, val))
			}
		}
		if queries.Filters.RoleID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleID.Op, fmt.Sprintf("%v", queries.Filters.RoleID.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" st.roleId %s %s", opValue.Operator, val))
		}
		if len(queries.RoleIDsInt64) > 0 {
			where = append(where, fmt.Sprintf(" st.roleId = ANY ('{%s}') ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(queries.RoleIDsInt64)), ","), "[]")))
		}
		if queries.Filters.RoleName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleName.Op, fmt.Sprintf("%v", queries.Filters.RoleName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.StartHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartHour.Op, fmt.Sprintf("%v", queries.Filters.StartHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.startHour %s %s", opValue.Operator, val))
		}
		if queries.Filters.EndHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndHour.Op, fmt.Sprintf("%v", queries.Filters.EndHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.endHour %s %s", opValue.Operator, val))
		}
	}
	return where
}

// CalculateStaffTypeNeededAndRemindStaffCount calculates the needed and reminded staff count for each staff type in the given slice.
//
// Parameters:
// - staffTypes: A slice of pointers to CycleStaffType objects representing staff types.
//
// Returns:
// - A slice of pointers to CycleStaffType objects with updated reminder counts.
// - An error if any occurred during the calculation.
func (r *CycleRepositoryPostgresDB) CalculateStaffTypeNeededAndRemindStaffCount(staffTypes []*domain.CycleStaffType) ([]*domain.CycleStaffType, error) {
	for _, staffType := range staffTypes {
		count, err := r.CountPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
			CycleID:                int(staffType.CycleID),
			CycleStaffTypeIDs:      fmt.Sprintf("[%d]", int(staffType.ID)),
			CycleStaffTypeIDsInt64: []int64{int64(staffType.ID)},
		})
		if err != nil {
			log.Printf("Error while query staff types: %s", err.Error())
			return nil, err
		}
		staffType.UsedStaffCount = uint(count)
		if staffType.NeededStaffCount < staffType.UsedStaffCount {
			staffType.RemindStaffCount = staffType.NeededStaffCount
		} else {
			staffType.RemindStaffCount = staffType.NeededStaffCount - staffType.UsedStaffCount
		}
		log.Printf("staffType: %#v\n", staffType.RemindStaffCount)
	}
	return staffTypes, nil
}

// CalculateNextStaffTypeNeededAndRemindStaffCount calculates the next staff type needed and remind staff count.
//
// It takes a slice of next staff types as input and returns the updated slice with the used staff count and remind staff count for each next staff type.
// It returns an error if there is an issue while querying the next staff types.
//
// Parameters:
// - nextStaffTypes (slice of *domain.CycleNextStaffType): A slice of next staff types.
//
// Return:
// - A slice of *domain.CycleNextStaffType: The updated slice of next staff types.
// - error: An error if there is an issue while querying the next staff types.
func (r *CycleRepositoryPostgresDB) CalculateNextStaffTypeNeededAndRemindStaffCount(nextStaffTypes []*domain.CycleNextStaffType) ([]*domain.CycleNextStaffType, error) {
	for _, nextStaffType := range nextStaffTypes {
		count, err := r.CountIncomingCyclePickupShifts(&models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
			CycleID:                    int(nextStaffType.CurrentCycleID),
			CycleNextStaffTypeIDs:      fmt.Sprintf("[%d]", int(nextStaffType.ID)),
			CycleNextStaffTypeIDsInt64: []int64{int64(nextStaffType.ID)},
		})
		if err != nil {
			log.Printf("Error while query next staff types: %s", err.Error())
			return nil, err
		}
		nextStaffType.UsedStaffCount = uint(count)
		if nextStaffType.NeededStaffCount < nextStaffType.UsedStaffCount {
			nextStaffType.RemindStaffCount = nextStaffType.NeededStaffCount
		} else {
			nextStaffType.RemindStaffCount = nextStaffType.NeededStaffCount - nextStaffType.UsedStaffCount
		}
	}
	return nextStaffTypes, nil
}

// FindAllStaffTypesByCycleID finds all staff types for a given cycle ID.
//
// Parameters:
// - cycleID: The ID of the cycle for which to retrieve staff types.
//
// Returns:
// - A slice of pointers to CycleStaffType objects, each representing a staff type for the given cycle.
// - An error if any occurred during the query.
func (r *CycleRepositoryPostgresDB) FindAllStaffTypesByCycleID(cycleID int64) ([]*domain.CycleStaffType, error) {
	q := `
		SELECT
			st.id,
			st.cycleId,
			st.roleId,
			st.datetime,
			st.shiftName,
			st.neededStaffCount,
			st.startHour,
			st.endHour,
			st.isUnplanned,
			st.created_at,
			st.updated_at,
			st.deleted_at,
			r.id as roleId,
			r.name as roleName
		FROM cycleStaffTypes st
		JOIN _roles r ON r.id = st.roleId
		WHERE st.cycleId = $1;
	`

	var staffTypes []*domain.CycleStaffType
	rows, err := r.PostgresDB.Query(q, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			staffType domain.CycleStaffType
			role      domain.CycleStaffTypeRole
		)
		err := rows.Scan(
			&staffType.ID,
			&staffType.CycleID,
			&staffType.RoleID,
			&staffType.DateTime,
			&staffType.ShiftName,
			&staffType.NeededStaffCount,
			&staffType.StartHour,
			&staffType.EndHour,
			&staffType.IsUnplanned,
			&staffType.CreatedAt,
			&staffType.UpdatedAt,
			&staffType.DeletedAt,
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		staffType.Role = &role
		staffTypes = append(staffTypes, &staffType)
	}
	st, err := r.CalculateStaffTypeNeededAndRemindStaffCount(staffTypes)
	if err != nil {
		return nil, err
	}
	staffTypes = st
	return staffTypes, nil
}

// QueryStaffTypes retrieves a list of staff types based on the provided query parameters.
//
// queries: A pointer to a CyclesQueryStaffTypesRequestParams struct containing the query parameters.
// Returns a slice of pointers to CycleStaffType structs and an error if any.
func (r *CycleRepositoryPostgresDB) QueryStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) ([]*domain.CycleStaffType, error) {
	q := `
		SELECT
			st.id,
			st.cycleId,
			st.roleId,
			st.datetime,
			st.shiftName,
			st.neededStaffCount,
			st.startHour,
			st.endHour,
			st.isUnplanned,
			st.created_at,
			st.updated_at,
			st.deleted_at,
			r.id as roleId,
			r.name as roleName
		FROM cycleStaffTypes st
		JOIN _roles r ON r.id = st.roleId
	`
	if queries != nil {
		where := makeStaffTypesWhereFilters(queries)
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

	var staffTypes []*domain.CycleStaffType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			staffType domain.CycleStaffType
			role      domain.CycleStaffTypeRole
		)
		err := rows.Scan(
			&staffType.ID,
			&staffType.CycleID,
			&staffType.RoleID,
			&staffType.DateTime,
			&staffType.ShiftName,
			&staffType.NeededStaffCount,
			&staffType.StartHour,
			&staffType.EndHour,
			&staffType.IsUnplanned,
			&staffType.CreatedAt,
			&staffType.UpdatedAt,
			&staffType.DeletedAt,
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		staffType.Role = &role
		staffTypes = append(staffTypes, &staffType)
	}
	st, err := r.CalculateStaffTypeNeededAndRemindStaffCount(staffTypes)
	if err != nil {
		return nil, err
	}
	staffTypes = st
	return staffTypes, nil
}

// CountStaffTypes returns the count of staff types based on the provided queries.
//
// Parameters:
// - queries: A pointer to a CyclesQueryStaffTypesRequestParams struct containing the filters for the query.
//
// Returns:
// - int64: The count of staff types.
// - error: An error if the query execution fails.
func (r *CycleRepositoryPostgresDB) CountStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) (int64, error) {
	q := `SELECT COUNT(st.id) FROM cycleStaffTypes st `
	if queries != nil {
		where := makeStaffTypesWhereFilters(queries)
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

// Create creates a new cycle based on the provided request body.
//
// payload: a pointer to a CyclesCreateRequestBody struct containing the cycle's name, start date, end date, period length, shift times, freeze period date, and wish days.
// Returns a pointer to a Cycle struct and an error if any.
func (r *CycleRepositoryPostgresDB) Create(payload *models.CyclesCreateRequestBody) (*domain.Cycle, error) {
	var cycle domain.Cycle

	// Current time
	currentTime := time.Now()

	// Calculate cycle name based on count of cycles
	var count int64
	err := r.PostgresDB.QueryRow(`
		SELECT COUNT(*) FROM cycles
	`).Scan(&count)
	if err != nil {
		return nil, err
	}
	payload.Name = fmt.Sprintf("Cycle %d", count+1)

	// Insert the cycle
	err = r.PostgresDB.QueryRow(`
		INSERT INTO cycles (name, created_at, updated_at, deleted_at, sectionId, start_datetime, end_datetime, periodLength, shiftMorningStartTime, shiftMorningEndTime, shiftEveningStartTime, shiftEveningEndTime, shiftNightStartTime, shiftNightEndTime, freeze_period_datetime, wishDays)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, name, created_at, updated_at, deleted_at, sectionId, start_datetime, end_datetime, periodLength, shiftMorningStartTime, shiftMorningEndTime, shiftEveningStartTime, shiftEveningEndTime, shiftNightStartTime, shiftNightEndTime, freeze_period_datetime, wishDays
	`,
		payload.Name,
		currentTime,
		currentTime,
		nil,
		payload.SectionID,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.PeriodLength,
		payload.ShiftMorningStartTime,
		payload.ShiftMorningEndTime,
		payload.ShiftEveningStartTime,
		payload.ShiftEveningEndTime,
		payload.ShiftNightStartTime,
		payload.ShiftNightEndTime,
		payload.FreezePeriodDateAsDate,
		payload.WishDays,
	).Scan(
		&cycle.ID,
		&cycle.Name,
		&cycle.CreatedAt,
		&cycle.UpdatedAt,
		&cycle.DeletedAt,
		&cycle.SectionID,
		&cycle.StartDate,
		&cycle.EndDate,
		&cycle.PeriodLength,
		&cycle.ShiftMorningStartTime,
		&cycle.ShiftMorningEndTime,
		&cycle.ShiftEveningStartTime,
		&cycle.ShiftEveningEndTime,
		&cycle.ShiftNightStartTime,
		&cycle.ShiftNightEndTime,
		&cycle.FreezePeriodDate,
		&cycle.WishDays,
	)
	if err != nil {
		return nil, err
	}

	// Return the cycle
	return &cycle, nil
}

// Delete deletes cycles from the database based on the provided payload.
//
// Parameters:
// - payload: a pointer to a models.CyclesDeleteRequestBody struct containing the IDs of the cycles to be deleted.
//
// Returns:
// - []int64: an array of the deleted cycle IDs.
// - error: an error if the deletion fails.
func (r *CycleRepositoryPostgresDB) Delete(payload *models.CyclesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM cycles
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

// Update updates a cycle by id.
//
// payload: a pointer to a CyclesCreateRequestBody struct containing the cycle's updated parameters.
// id: the id of the cycle to be updated.
// Returns a pointer to a Cycle struct and an error if any.
func (r *CycleRepositoryPostgresDB) Update(payload *models.CyclesCreateRequestBody, id int64) (*domain.Cycle, error) {
	var cycle domain.Cycle

	// Current time
	currentTime := time.Now()

	// Find the cycle by id
	var foundCycleID int64
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM cycles
		WHERE id = $1
	`, id).Scan(&foundCycleID)
	if err != nil {
		return nil, err
	}

	// Update the cycle
	err = r.PostgresDB.QueryRow(`
		UPDATE cycles
		SET updated_at = $1, sectionId = $2, start_datetime = $3, end_datetime = $4, periodLength = $5, shiftMorningStartTime = $6, shiftMorningEndTime = $7, shiftEveningStartTime = $8, shiftEveningEndTime = $9, shiftNightStartTime = $10, shiftNightEndTime = $11, freeze_period_datetime = $12, wishDays = $13
		WHERE id = $14
		RETURNING id, name, created_at, updated_at, deleted_at, sectionId, start_datetime, end_datetime, periodLength, shiftMorningStartTime, shiftMorningEndTime, shiftEveningStartTime, shiftEveningEndTime, shiftNightStartTime, shiftNightEndTime, freeze_period_datetime, wishDays
	`,
		currentTime,
		payload.SectionID,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.PeriodLength,
		payload.ShiftMorningStartTime,
		payload.ShiftMorningEndTime,
		payload.ShiftEveningStartTime,
		payload.ShiftEveningEndTime,
		payload.ShiftNightStartTime,
		payload.ShiftNightEndTime,
		payload.FreezePeriodDateAsDate,
		payload.WishDays,
		foundCycleID,
	).Scan(
		&cycle.ID,
		&cycle.SectionID,
		&cycle.Name,
		&cycle.CreatedAt,
		&cycle.UpdatedAt,
		&cycle.DeletedAt,
		&cycle.SectionID,
		&cycle.StartDate,
		&cycle.EndDate,
		&cycle.PeriodLength,
		&cycle.ShiftMorningStartTime,
		&cycle.ShiftMorningEndTime,
		&cycle.ShiftEveningStartTime,
		&cycle.ShiftEveningEndTime,
		&cycle.ShiftNightStartTime,
		&cycle.ShiftNightEndTime,
		&cycle.FreezePeriodDate,
		&cycle.WishDays,
	)
	if err != nil {
		return nil, err
	}

	// Return the cycle
	return &cycle, nil
}

// UpdateStaffType updates the staff type for a cycle.
//
// It takes a CyclesUpdateStaffTypeRequestBody payload, an ID, and a boolean indicating
// whether the staff type is unplanned. It returns a Cycle and an error if any.
func (r *CycleRepositoryPostgresDB) UpdateStaffType(payload *models.CyclesUpdateStaffTypeRequestBody, id int64, isUnplanned bool) (*domain.Cycle, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Find cycle staff type by cycle id and datetime and shift name
	var (
		cycleStaffType domain.CycleStaffType
		isUpdateMode   = true
	)
	err = r.PostgresDB.QueryRow(`
		SELECT *
		FROM cycleStaffTypes
		WHERE cycleId = $1 AND datetime = $2 AND shiftName = $3 AND roleId = $4
	`,
		cycle.ID,
		payload.DateTimeAsDate,
		payload.ShiftName,
		payload.RoleID,
	).Scan(
		&cycleStaffType.ID,
		&cycleStaffType.CycleID,
		&cycleStaffType.RoleID,
		&cycleStaffType.DateTime,
		&cycleStaffType.ShiftName,
		&cycleStaffType.NeededStaffCount,
		&cycleStaffType.StartHour,
		&cycleStaffType.EndHour,
		&cycleStaffType.IsUnplanned,
		&cycleStaffType.CreatedAt,
		&cycleStaffType.UpdatedAt,
		&cycleStaffType.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			isUpdateMode = false
		} else {
			return nil, err
		}
	}
	if cycleStaffType.ID == 0 {
		isUpdateMode = false
	}
	currentTime := time.Now()
	if isUpdateMode {
		// Update the cycle staff type
		err = r.PostgresDB.QueryRow(`
			UPDATE cycleStaffTypes
			SET updated_at = $1, roleId = $2, neededStaffCount = $3
			WHERE id = $4
			RETURNING id, cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, isUnplanned, created_at, updated_at, deleted_at
		`,
			currentTime,
			payload.RoleID,
			payload.NeededStaffCount,
			cycleStaffType.ID,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.IsUnplanned,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Insert the cycle staff type
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleStaffTypes (cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, isUnplanned, created_at, updated_at, deleted_at
		`,
			cycle.ID,
			payload.RoleID,
			payload.DateTimeAsDate,
			payload.ShiftName,
			payload.NeededStaffCount,
			payload.StartHourAsTime,
			payload.EndHourAsTime,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.IsUnplanned,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	cycleStaffType.Role = payload.Role

	// Append the cycle staff type to the cycle
	cycle.StaffTypes = append(cycle.StaffTypes, cycleStaffType)

	// Query next staff types
	cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
		CurrentCycleID: int(cycle.ID),
	})
	if err != nil {
		return nil, err
	}
	for _, cycleNextStaffType := range cycleNextStaffTypes {
		cycle.NextStaffTypes = append(cycle.NextStaffTypes, *cycleNextStaffType)
	}

	// Return the cycle
	return cycle, nil
}

// UpdateStaffTypes updates the staff types of a cycle.
//
// It takes a CyclesUpdateStaffTypesRequestBody payload and a cycle ID as arguments.
// Returns the updated cycle and an error if any occurs during the process.
func (r *CycleRepositoryPostgresDB) UpdateStaffTypes(payload *models.CyclesUpdateStaffTypesRequestBody, id int64) (*domain.Cycle, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]
	cycle.StaffTypes = []domain.CycleStaffType{}

	// Create a transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Delete all cycle staff types by cycle id
	_, err = tx.Exec(`
		DELETE FROM cycleStaffTypes
		WHERE cycleId = $1
	`, cycle.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert cycle staff types
	currentTime := time.Now()
	for _, cycleStaffType := range payload.StaffTypes {
		var insertedCycleStaffType domain.CycleStaffType
		err = tx.QueryRow(`
			INSERT INTO cycleStaffTypes (cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, isUnplanned, created_at, updated_at, deleted_at
		`,
			cycle.ID,
			cycleStaffType.RoleID,
			cycleStaffType.DateTimeAsDate,
			cycleStaffType.ShiftName,
			cycleStaffType.NeededStaffCount,
			cycleStaffType.StartHourAsTime,
			cycleStaffType.EndHourAsTime,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&insertedCycleStaffType.ID,
			&insertedCycleStaffType.CycleID,
			&insertedCycleStaffType.RoleID,
			&insertedCycleStaffType.DateTime,
			&insertedCycleStaffType.ShiftName,
			&insertedCycleStaffType.NeededStaffCount,
			&insertedCycleStaffType.StartHour,
			&insertedCycleStaffType.EndHour,
			&insertedCycleStaffType.IsUnplanned,
			&insertedCycleStaffType.CreatedAt,
			&insertedCycleStaffType.UpdatedAt,
			&insertedCycleStaffType.DeletedAt,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		insertedCycleStaffType.Role = cycleStaffType.Role

		// Append the cycle staff type to the cycle
		cycle.StaffTypes = append(cycle.StaffTypes, insertedCycleStaffType)

		// Query next staff types
		cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
			CurrentCycleID: int(cycle.ID),
		})
		if err != nil {
			return nil, err
		}
		for _, cycleNextStaffType := range cycleNextStaffTypes {
			cycle.NextStaffTypes = append(cycle.NextStaffTypes, *cycleNextStaffType)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Return the cycle
	return cycle, nil
}

func (r *CycleRepositoryPostgresDB) UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload *models.CyclesUpdateStaffTypeRequestBody, migratedCycleID int64, currentCycleID int64) (*domain.Cycle, error) {
	// Find last cycle staff type
	var (
		lastStaffType   *domain.CycleNextStaffType
		lastStaffTypeID = payload.ID
	)
	log.Printf("lastStaffTypeID: %d\n", lastStaffTypeID)
	if lastStaffTypeID != 0 {
		cycleStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
			ID: int(lastStaffTypeID),
		})
		if err != nil {
			return nil, err
		}
		if len(cycleStaffTypes) == 0 {
			return nil, errors.New("cycle staff type not found")
		}
		log.Printf("len(cycleStaffTypes): %d\n", len(cycleStaffTypes))
		if len(cycleStaffTypes) > 0 {
			lastStaffType = cycleStaffTypes[0]
		}
	}
	if lastStaffType == nil {
		return nil, errors.New("last staff type not found")
	}

	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(migratedCycleID),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Find cycle staff type by cycle id and datetime and shift name
	var (
		cycleStaffType domain.CycleStaffType
		isUpdateMode   = true
	)
	err = r.PostgresDB.QueryRow(`
		SELECT *
		FROM cycleStaffTypes
		WHERE cycleId = $1 AND datetime = $2 AND shiftName = $3 AND roleId = $4
	`,
		cycle.ID,
		payload.DateTimeAsDate,
		payload.ShiftName,
		payload.RoleID,
	).Scan(
		&cycleStaffType.ID,
		&cycleStaffType.CycleID,
		&cycleStaffType.RoleID,
		&cycleStaffType.DateTime,
		&cycleStaffType.ShiftName,
		&cycleStaffType.NeededStaffCount,
		&cycleStaffType.StartHour,
		&cycleStaffType.EndHour,
		&cycleStaffType.IsUnplanned,
		&cycleStaffType.CreatedAt,
		&cycleStaffType.UpdatedAt,
		&cycleStaffType.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			isUpdateMode = false
		} else {
			return nil, err
		}
	}
	if cycleStaffType.ID == 0 {
		isUpdateMode = false
	}
	currentTime := time.Now()
	if isUpdateMode {
		// Update the cycle staff type
		err = r.PostgresDB.QueryRow(`
			UPDATE cycleStaffTypes
			SET updated_at = $1, roleId = $2, neededStaffCount = $3
			WHERE id = $4
			RETURNING id, cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, isUnplanned, created_at, updated_at, deleted_at
		`,
			currentTime,
			payload.RoleID,
			payload.NeededStaffCount,
			cycleStaffType.ID,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.IsUnplanned,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Insert the cycle staff type
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleStaffTypes (cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, cycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, isUnplanned, created_at, updated_at, deleted_at
		`,
			cycle.ID,
			payload.RoleID,
			payload.DateTimeAsDate,
			payload.ShiftName,
			payload.NeededStaffCount,
			payload.StartHourAsTime,
			payload.EndHourAsTime,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.IsUnplanned,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	cycleStaffType.Role = payload.Role

	// Append the cycle staff type to the cycle
	cycle.StaffTypes = append(cycle.StaffTypes, cycleStaffType)

	// Query next staff types
	cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
		CurrentCycleID: int(cycle.ID),
	})
	if err != nil {
		return nil, err
	}
	for _, cycleNextStaffType := range cycleNextStaffTypes {
		cycle.NextStaffTypes = append(cycle.NextStaffTypes, *cycleNextStaffType)
	}

	// Migrate pickup shift
	pickedUpShifts, err := r.QueryIncomingCyclePickupShifts(&models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		CycleID:               int(currentCycleID),
		CycleNextStaffTypeIDs: fmt.Sprintf("[%d]", int(cycleStaffType.ID)),
	})
	if err != nil {
		return nil, err
	}
	for _, pickedUpShift := range pickedUpShifts {
		cstids := []interface{}{
			float64(cycleStaffType.ID),
		}
		cstidsInt64, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(cstids)
		if err != nil {
			return nil, err
		}
		_, err = r.PickupShift(&models.CyclesCreatePickupShiftRequestBody{
			CycleID:           int(migratedCycleID),
			StaffID:           int(pickedUpShift.Staff.ID),
			CycleStaffTypeIDs: cstids,
			DateTime:          payload.DateTime,
			Staff: &domain.CyclePickupShiftStaff{
				ID:        pickedUpShift.Staff.ID,
				UserID:    pickedUpShift.Staff.UserID,
				FirstName: pickedUpShift.Staff.FirstName,
				LastName:  pickedUpShift.Staff.LastName,
			},
			CycleStaffTypeIDsInt64: cstidsInt64,
			DateTimeAsDate:         payload.DateTimeAsDate,
			ShiftName:              payload.ShiftName,
		})
		if err != nil {
			return nil, err
		}
	}

	// Return the cycle
	return cycle, nil
}

// Duplicate creates a new cycle by duplicating an existing one.
//
// It takes a payload of type models.CyclesDuplicateRequestBody, which contains the cycle to be duplicated.
// Returns a pointer to the newly created domain.Cycle and an error.
func (r *CycleRepositoryPostgresDB) Duplicate(payload *models.CyclesDuplicateRequestBody) (*domain.Cycle, error) {
	// Calculate the start date and end date
	currentTimeStart := time.Now()
	currentTimeEnd := currentTimeStart

	// Find the difference between start date and end date
	diffDate := payload.Cycle.EndDate.Sub(payload.Cycle.StartDate)
	currentTimeEnd = currentTimeEnd.Add(diffDate)

	// Update the start date to the current start time
	payload.Cycle.StartDate = currentTimeStart

	// Update the end date to the current end time
	payload.Cycle.EndDate = &currentTimeEnd

	var (
		startDateStr        = payload.Cycle.StartDate.Format(time.RFC3339)
		endDateStr          = payload.Cycle.EndDate.Format(time.RFC3339)
		freezePeriodDateStr = payload.Cycle.FreezePeriodDate.Format(time.RFC3339)
	)
	createdCycle, err := r.Create(&models.CyclesCreateRequestBody{
		StartDate:              startDateStr,
		EndDate:                &endDateStr,
		PeriodLength:           payload.Cycle.PeriodLength,
		ShiftMorningStartTime:  payload.Cycle.ShiftMorningStartTime,
		ShiftMorningEndTime:    payload.Cycle.ShiftMorningEndTime,
		ShiftEveningStartTime:  payload.Cycle.ShiftEveningStartTime,
		ShiftEveningEndTime:    payload.Cycle.ShiftEveningEndTime,
		ShiftNightStartTime:    payload.Cycle.ShiftNightStartTime,
		ShiftNightEndTime:      payload.Cycle.ShiftNightEndTime,
		FreezePeriodDate:       freezePeriodDateStr,
		WishDays:               payload.Cycle.WishDays,
		Name:                   payload.Cycle.Name,
		SectionID:              int(payload.Cycle.SectionID),
		StartDateAsDate:        &payload.Cycle.StartDate,
		EndDateAsDate:          payload.Cycle.EndDate,
		FreezePeriodDateAsDate: &payload.Cycle.FreezePeriodDate,
	})
	if err != nil {
		return nil, err
	}

	// Create staff types
	var staffTypesData []*models.CyclesUpdateStaffTypeRequestBody
	for _, staffType := range payload.CycleStaffTypes {
		var newDateTime = staffType.DateTime.Add(currentTimeStart.Sub(payload.Cycle.StartDate))
		var dateTimeStr = newDateTime.Format(time.RFC3339)
		var startHourStr = staffType.StartHour.Format("15:04")
		var endHourStr = staffType.EndHour.Format("15:04")
		staffTypesData = append(staffTypesData, &models.CyclesUpdateStaffTypeRequestBody{
			ShiftName:        staffType.ShiftName,
			DateTime:         dateTimeStr,
			NeededStaffCount: int(staffType.NeededStaffCount),
			StartHour:        startHourStr,
			EndHour:          endHourStr,
			RoleID:           int(staffType.RoleID),
			DateTimeAsDate:   &newDateTime,
			Role:             staffType.Role,
		})
	}
	createdCycle, err = r.UpdateStaffTypes(&models.CyclesUpdateStaffTypesRequestBody{
		StaffTypes: staffTypesData,
	}, int64(createdCycle.ID))
	if err != nil {
		r.Delete(&models.CyclesDeleteRequestBody{
			IDs:      nil,
			IDsInt64: []int64{int64(createdCycle.ID)},
		})
		return nil, err
	}

	return createdCycle, nil
}

// makeNextStaffTypesWhereFilters generates WHERE filters for next staff types based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryNextStaffTypesRequestParams as an argument.
// Returns a slice of strings representing the WHERE conditions.
func makeNextStaffTypesWhereFilters(queries *models.CyclesQueryNextStaffTypesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" st.id = %d ", queries.ID))
		}
		if queries.CurrentCycleID != 0 {
			where = append(where, fmt.Sprintf(" st.currentCycleId = %d ", queries.CurrentCycleID))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" st.shiftName %s %s", opValue.Operator, val))
		}
		if queries.Filters.DateTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.DateTime.Op, fmt.Sprintf("%v", queries.Filters.DateTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" st.datetime %s %s", opValue.Operator, val))
		}
		if queries.Filters.RoleID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleID.Op, fmt.Sprintf("%v", queries.Filters.RoleID.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" st.roleId %s %s", opValue.Operator, val))
		}
		if len(queries.RoleIDsInt64) > 0 {
			where = append(where, fmt.Sprintf(" st.roleId = ANY ('{%s}') ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(queries.RoleIDsInt64)), ","), "[]")))
		}
		if queries.Filters.RoleName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleName.Op, fmt.Sprintf("%v", queries.Filters.RoleName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.StartHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartHour.Op, fmt.Sprintf("%v", queries.Filters.StartHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.startHour %s %s", opValue.Operator, val))
		}
		if queries.Filters.EndHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndHour.Op, fmt.Sprintf("%v", queries.Filters.EndHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.endHour %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryNextStaffTypes retrieves next staff types based on the provided queries.
//
// Parameters:
// - queries: A pointer to a CyclesQueryNextStaffTypesRequestParams struct containing the query parameters.
//
// Returns:
// - A slice of pointers to CycleNextStaffType structs representing the retrieved next staff types.
// - An error if any occurred during the query.
func (r *CycleRepositoryPostgresDB) QueryNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) ([]*domain.CycleNextStaffType, error) {
	q := `
		SELECT
			st.id,
			st.currentCycleId,
			st.roleId,
			st.datetime,
			st.shiftName,
			st.neededStaffCount,
			st.startHour,
			st.endHour,
			st.created_at,
			st.updated_at,
			st.deleted_at,
			r.id as roleId,
			r.name as roleName
		FROM cycleNextStaffTypes st
		JOIN _roles r ON r.id = st.roleId
	`
	if queries != nil {
		where := makeNextStaffTypesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit != -1 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"

	var staffTypes []*domain.CycleNextStaffType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			staffType domain.CycleNextStaffType
			role      domain.CycleNextStaffTypeRole
		)
		err := rows.Scan(
			&staffType.ID,
			&staffType.CurrentCycleID,
			&staffType.RoleID,
			&staffType.DateTime,
			&staffType.ShiftName,
			&staffType.NeededStaffCount,
			&staffType.StartHour,
			&staffType.EndHour,
			&staffType.CreatedAt,
			&staffType.UpdatedAt,
			&staffType.DeletedAt,
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		staffType.Role = &role
		staffTypes = append(staffTypes, &staffType)
	}
	st, err := r.CalculateNextStaffTypeNeededAndRemindStaffCount(staffTypes)
	if err != nil {
		return nil, err
	}
	staffTypes = st
	return staffTypes, nil
}

// CountNextStaffTypes counts the number of next staff types based on the provided query parameters.
//
// Parameters:
// - queries: A pointer to a CyclesQueryNextStaffTypesRequestParams struct containing the query parameters.
//
// Returns:
// - An integer representing the count of next staff types.
// - An error if any occurred during the query.
func (r *CycleRepositoryPostgresDB) CountNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) (int64, error) {
	q := `SELECT COUNT(st.id) FROM cycleNextStaffTypes st `
	if queries != nil {
		where := makeNextStaffTypesWhereFilters(queries)
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

// UpdateNextStaffType updates the next staff type for a cycle.
//
// Parameters:
// - payload: A pointer to a CyclesUpdateNextStaffTypeRequestBody struct containing the update parameters.
// - id: The ID of the cycle to update.
//
// Returns:
// - A slice of pointers to CycleNextStaffType structs representing the updated cycle next staff types.
// - An error if any occurred during the update.
func (r *CycleRepositoryPostgresDB) UpdateNextStaffType(payload *models.CyclesUpdateNextStaffTypeRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Find cycle staff type by cycle id and datetime and shift name
	var (
		cycleStaffType domain.CycleNextStaffType
		isUpdateMode   = true
	)
	err = r.PostgresDB.QueryRow(`
		SELECT *
		FROM cycleNextStaffTypes
		WHERE currentCycleId = $1 AND datetime = $2 AND shiftName = $3
	`,
		cycle.ID,
		payload.DateTimeAsDate,
		payload.ShiftName,
	).Scan(
		&cycleStaffType.ID,
		&cycleStaffType.CurrentCycleID,
		&cycleStaffType.RoleID,
		&cycleStaffType.DateTime,
		&cycleStaffType.ShiftName,
		&cycleStaffType.NeededStaffCount,
		&cycleStaffType.StartHour,
		&cycleStaffType.EndHour,
		&cycleStaffType.CreatedAt,
		&cycleStaffType.UpdatedAt,
		&cycleStaffType.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			isUpdateMode = false
		} else {
			return nil, err
		}
	}
	if cycleStaffType.ID == 0 {
		isUpdateMode = false
	}
	currentTime := time.Now()
	if isUpdateMode {
		// Update the cycle staff type
		err = r.PostgresDB.QueryRow(`
			UPDATE cycleNextStaffTypes
			SET updated_at = $1, roleId = $2, neededStaffCount = $3
			WHERE id = $4
			RETURNING id, currentCycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at
		`,
			currentTime,
			payload.RoleID,
			payload.NeededStaffCount,
			cycleStaffType.ID,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CurrentCycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Insert the cycle staff type
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleNextStaffTypes (currentCycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, currentCycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at
		`,
			cycle.ID,
			payload.RoleID,
			payload.DateTimeAsDate,
			payload.ShiftName,
			payload.NeededStaffCount,
			payload.StartHourAsTime,
			payload.EndHourAsTime,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&cycleStaffType.ID,
			&cycleStaffType.CurrentCycleID,
			&cycleStaffType.RoleID,
			&cycleStaffType.DateTime,
			&cycleStaffType.ShiftName,
			&cycleStaffType.NeededStaffCount,
			&cycleStaffType.StartHour,
			&cycleStaffType.EndHour,
			&cycleStaffType.CreatedAt,
			&cycleStaffType.UpdatedAt,
			&cycleStaffType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	cycleStaffType.Role = payload.Role

	// Query cycle next staff types
	cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
		CurrentCycleID: int(cycle.ID),
	})
	if err != nil {
		return nil, err
	}

	return cycleNextStaffTypes, nil
}

// UpdateNextStaffTypes updates the next staff types for a cycle based on the provided payload and cycle ID.
//
// Parameters:
// - payload: A pointer to a CyclesUpdateNextStaffTypesRequestBody struct containing the payload for updating staff types.
// - id: An int64 representing the ID of the cycle to update.
//
// Returns:
// - A slice of pointers to CycleNextStaffType structs representing the updated next staff types.
// - An error if any occurred during the update process.
func (r *CycleRepositoryPostgresDB) UpdateNextStaffTypes(payload *models.CyclesUpdateNextStaffTypesRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Create a transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Delete all cycle staff types by cycle id
	_, err = tx.Exec(`
		DELETE FROM cycleNextStaffTypes
		WHERE currentCycleId = $1
	`, cycle.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert cycle staff types
	currentTime := time.Now()
	for _, cycleStaffType := range payload.StaffTypes {
		var insertedCycleStaffType domain.CycleNextStaffType
		err = tx.QueryRow(`
			INSERT INTO cycleNextStaffTypes (currentCycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, currentCycleId, roleId, datetime, shiftName, neededStaffCount, startHour, endHour, created_at, updated_at, deleted_at
		`,
			cycle.ID,
			cycleStaffType.RoleID,
			cycleStaffType.DateTimeAsDate,
			cycleStaffType.ShiftName,
			cycleStaffType.NeededStaffCount,
			cycleStaffType.StartHourAsTime,
			cycleStaffType.EndHourAsTime,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&insertedCycleStaffType.ID,
			&insertedCycleStaffType.CurrentCycleID,
			&insertedCycleStaffType.RoleID,
			&insertedCycleStaffType.DateTime,
			&insertedCycleStaffType.ShiftName,
			&insertedCycleStaffType.NeededStaffCount,
			&insertedCycleStaffType.StartHour,
			&insertedCycleStaffType.EndHour,
			&insertedCycleStaffType.CreatedAt,
			&insertedCycleStaffType.UpdatedAt,
			&insertedCycleStaffType.DeletedAt,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		insertedCycleStaffType.Role = cycleStaffType.Role
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Query cycle next staff types
	cycleNextStaffTypes, err := r.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
		CurrentCycleID: int(cycle.ID),
	})
	if err != nil {
		return nil, err
	}

	return cycleNextStaffTypes, nil
}

// makeCycleShiftsWhereFilters generates WHERE filters for cycle shifts based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryCycleShiftsRequestParams as an argument.
// Returns a slice of strings representing the WHERE conditions.
func makeCycleShiftsWhereFilters(queries *models.CyclesQueryCycleShiftsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cs.id = %d ", queries.ID))
		}
		if queries.CycleID != 0 {
			where = append(where, fmt.Sprintf(" cs.cycleId = %d ", queries.CycleID))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cs.shiftName %s %s", opValue.Operator, val))
		}
		if queries.Filters.DateTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.DateTime.Op, fmt.Sprintf("%v", queries.Filters.DateTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cs.datetime %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cs.status %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryCycleShifts queries cycle shifts based on the given query parameters.
//
// Parameters:
// - queries: A pointer to a CyclesQueryCycleShiftsRequestParams struct containing the query parameters.
//
// Returns:
// - A slice of pointers to domain.CycleShift structs representing the queried cycle shifts.
// - An error if any error occurs during the database query.
func (r *CycleRepositoryPostgresDB) QueryCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) ([]*domain.CycleShift, error) {
	q := `
		SELECT
			cs.id,
			cs.exchangeKey,
			cs.cycleId,
			cs.staffTypeIds,
			cs.shiftName,
			cs.vehicleType,
			cs.startLocation,
			cs.datetime,
			cs.status,
			cs.created_at,
			cs.updated_at,
			cs.deleted_at
		FROM cycleShifts cs
	`
	if queries != nil {
		where := makeCycleShiftsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit > -1 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"

	var shifts []*domain.CycleShift
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			shift                domain.CycleShift
			staffTypeIDs         sql.NullString
			staffTypeIDsMetadata []uint
			vehicleType          sql.NullString
			startLocation        sql.NullString
		)
		err := rows.Scan(
			&shift.ID,
			&shift.ExchangeKey,
			&shift.CycleID,
			&staffTypeIDs,
			&shift.ShiftName,
			&vehicleType,
			&startLocation,
			&shift.DateTime,
			&shift.Status,
			&shift.CreatedAt,
			&shift.UpdatedAt,
			&shift.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		if staffTypeIDs.Valid {
			err := json.Unmarshal([]byte(staffTypeIDs.String), &staffTypeIDsMetadata)
			if err != nil {
				log.Printf("Error while unmarshalling staffType ids in query cycle shift: %v\n", err.Error())
				return nil, err
			} else {
				shift.StaffTypeIDs = staffTypeIDsMetadata
			}
		}
		if len(shift.StaffTypeIDs) > 0 {
			for _, staffTypeID := range shift.StaffTypeIDs {
				cycleStaffType, err := r.QueryStaffTypes(&models.CyclesQueryStaffTypesRequestParams{
					ID: int(staffTypeID),
				})
				if err != nil {
					log.Printf("Error while querying staff type in query cycle shift: %v\n", err.Error())
					return nil, err
				}
				if len(cycleStaffType) > 0 {
					shift.StaffTypes = append(shift.StaffTypes, cycleStaffType[0])
				}
			}
		}
		if vehicleType.Valid {
			shift.VehicleType = &vehicleType.String
		}
		if startLocation.Valid {
			shift.StartLocation = &startLocation.String
		}
		shifts = append(shifts, &shift)
	}
	return shifts, nil
}

// CountCycleShifts counts the number of cycle shifts based on the provided query parameters.
//
// Parameters:
// - queries: A pointer to a CyclesQueryCycleShiftsRequestParams struct containing the query parameters.
//
// Returns:
// - An int64 representing the total number of cycle shifts that match the given query parameters.
// - An error if any error occurs during the database query.
func (r *CycleRepositoryPostgresDB) CountCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cs.id)
		FROM cycleShifts cs
	`
	if queries != nil {
		where := makeCycleShiftsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	log.Printf("count cycle shifts query: %v\n", q)

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// makePickupShiftsWhereFilters generates SQL WHERE filters based on the provided query parameters for pickup shifts.
//
// It takes a pointer to models.CyclesQueryPickupShiftsRequestParams as an argument.
// Returns a slice of strings representing the SQL WHERE conditions.
func makePickupShiftsWhereFilters(queries *models.CyclesQueryPickupShiftsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ps.id = %d ", queries.ID))
		}
		if queries.CycleID != 0 {
			where = append(where, fmt.Sprintf(" ps.cycleId = %d ", queries.CycleID))
		}
		if queries.StaffID != 0 {
			where = append(where, fmt.Sprintf(" ps.staffId = %d ", queries.StaffID))
		}
		if len(queries.CycleStaffTypeIDsInt64) > 0 {
			var w []string
			for _, cycleStaffTypeID := range queries.CycleStaffTypeIDsInt64 {
				w = append(w, fmt.Sprintf(" ps.cycleStaffTypeId = %d ", cycleStaffTypeID))
			}
			where = append(where, fmt.Sprintf("(%s)", strings.Join(w, " OR ")))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.shiftName %s %s", opValue.Operator, val))
		}
		if queries.Filters.DateTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.DateTime.Op, fmt.Sprintf("%v", queries.Filters.DateTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ps.datetime %s %s", opValue.Operator, val))
		}
		if queries.RangeDateTimeStartAsTime != nil {
			where = append(where, fmt.Sprintf(" ps.datetime >= '%s' ", queries.RangeDateTimeStartAsTime.Format(time.RFC3339)))
		}
		if queries.RangeDateTimeEndAsTime != nil {
			where = append(where, fmt.Sprintf(" ps.datetime <= '%s' ", queries.RangeDateTimeEndAsTime.Format(time.RFC3339)))
		}
		if len(queries.ShiftNamesAsArray) > 0 {
			where = append(where, fmt.Sprintf(" cst.shiftName = ANY ('{%s}') ", strings.Trim(strings.Join(queries.ShiftNamesAsArray, ","), "[]")))
		}
		if queries.Filters.StartHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartHour.Op, fmt.Sprintf("%v", queries.Filters.StartHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.startHour %s %s", opValue.Operator, val))
		}
		if queries.Filters.EndHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndHour.Op, fmt.Sprintf("%v", queries.Filters.EndHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.endHour %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ps.status %s %s", opValue.Operator, val))
		}
	}
	log.Println(where)
	return where
}

// QueryPickupShifts retrieves a list of pickup shifts based on the provided query parameters.
//
// The queries parameter is a struct containing various query parameters such as cycle ID, staff ID, shift ID, etc.
// It returns a slice of CyclePickupShift objects and an error.
func (r *CycleRepositoryPostgresDB) QueryPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) ([]*domain.CyclePickupShift, error) {
	q := `
		SELECT
			ps.id,
			ps.cycleId,
			ps.staffId,
			ps.shiftId,
			ps.cycleStaffTypeId,
			ps.datetime,
			ps.status,
			ps.prevStatus,
			ps.startKilometer,
			ps.reasonOfTheCancellation,
			ps.reasonOfTheReactivation,
			ps.reasonOfTheResume,
			ps.reasonOfThePause,
			ps.isUnplanned,
			ps.created_at,
			ps.updated_at,
			ps.deleted_at,
			ps.started_at,
			ps.ended_at,
			ps.cancelled_at,
			ps.delayed_at,
			ps.paused_at,
			ps.resumed_at,
			ps.reactivated_at,
			s.id as staffStaffID,
			s.userId as staffUserID,
			u.firstName as staffFirstName,
			u.lastName as staffLastName,
			u.avatarUrl as staffAvatarUrl,
			s.vehicleTypes as staffVehicleTypes,
			s.vehicleLicenseTypes as staffVehicleLicenseTypes,
			cst.id as cycleStaffTypeCycleStaffTypeId,
            cst.roleId as cycleStaffTypeCycleStaffTypeRoleId,
            r.id as cycleStaffTypeCycleStaffTypeRoleRoleID,
            r.name as cycleStaffTypeCycleStaffTypeRoleRoleName,
            cst.datetime as cycleStaffTypeCycleStaffTypeDateTime,
			cst.shiftName as cycleStaffTypeCycleStaffTypeShiftName,
            cst.startHour as cycleStaffTypeCycleStaffTypeStartHour,
            cst.endHour as cycleStaffTypeCycleStaffTypeEndHour,
			cst.isUnplanned as cycleStaffTypeCycleStaffTypeIsUnplanned
		FROM cyclePickupShifts ps
		JOIN cycleStaffTypes cst ON ps.cycleStaffTypeId = cst.id
		JOIN _roles r ON cst.roleId = r.id
		JOIN staffs s ON ps.staffId = s.id
		JOIN users u ON s.userId = u.id
	`
	if queries != nil {
		where := makePickupShiftsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit > -1 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"

	var pickupShifts []*domain.CyclePickupShift
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			pickupShift                              domain.CyclePickupShift
			sID                                      sql.NullInt64
			shiftID                                  sql.NullInt64
			cstID                                    sql.NullInt64
			staffID                                  sql.NullInt64
			staffUserID                              sql.NullInt64
			staffFirstName                           sql.NullString
			staffLastName                            sql.NullString
			staffAvatarUrl                           sql.NullString
			staffVehicleTypes                        *json.RawMessage
			staffVehicleLicenseTypes                 *json.RawMessage
			startKilometer                           sql.NullString
			reasonOfTheCancellation                  sql.NullString
			reasonOfTheReactivation                  sql.NullString
			reasonOfTheResume                        sql.NullString
			reasonOfThePause                         sql.NullString
			deletedAt                                sql.NullTime
			startedAt                                sql.NullTime
			endedAt                                  sql.NullTime
			cancelledAt                              sql.NullTime
			delayedAt                                sql.NullTime
			pausedAt                                 sql.NullTime
			resumedAt                                sql.NullTime
			reactivatedAt                            sql.NullTime
			cycleStaffTypeCycleStaffTypeId           sql.NullInt64
			cycleStaffTypeCycleStaffTypeRoleId       sql.NullInt64
			cycleStaffTypeCycleStaffTypeRoleRoleID   sql.NullInt64
			cycleStaffTypeCycleStaffTypeRoleRoleName sql.NullString
			cycleStaffTypeCycleStaffTypeDateTime     sql.NullTime
			cycleStaffTypeCycleStaffTypeShiftName    sql.NullString
			cycleStaffTypeCycleStaffTypeStartHour    sql.NullTime
			cycleStaffTypeCycleStaffTypeEndHour      sql.NullTime
			cycleStaffTypeCycleStaffTypeIsUnplanned  sql.NullBool
		)
		err := rows.Scan(
			&pickupShift.ID,
			&pickupShift.CycleID,
			&sID,
			&shiftID,
			&cstID,
			&pickupShift.DateTime,
			&pickupShift.Status,
			&pickupShift.PrevStatus,
			&startKilometer,
			&reasonOfTheCancellation,
			&reasonOfTheReactivation,
			&reasonOfTheResume,
			&reasonOfThePause,
			&pickupShift.IsUnplanned,
			&pickupShift.CreatedAt,
			&pickupShift.UpdatedAt,
			&deletedAt,
			&startedAt,
			&endedAt,
			&cancelledAt,
			&delayedAt,
			&pausedAt,
			&resumedAt,
			&reactivatedAt,
			&staffID,
			&staffUserID,
			&staffFirstName,
			&staffLastName,
			&staffAvatarUrl,
			&staffVehicleTypes,
			&staffVehicleLicenseTypes,
			&cycleStaffTypeCycleStaffTypeId,
			&cycleStaffTypeCycleStaffTypeRoleId,
			&cycleStaffTypeCycleStaffTypeRoleRoleID,
			&cycleStaffTypeCycleStaffTypeRoleRoleName,
			&cycleStaffTypeCycleStaffTypeDateTime,
			&cycleStaffTypeCycleStaffTypeShiftName,
			&cycleStaffTypeCycleStaffTypeStartHour,
			&cycleStaffTypeCycleStaffTypeEndHour,
			&cycleStaffTypeCycleStaffTypeIsUnplanned,
		)
		if err != nil {
			return nil, err
		}
		if startKilometer.Valid {
			pickupShift.StartKilometer = &startKilometer.String
		}
		if reasonOfTheCancellation.Valid {
			pickupShift.ReasonOfTheCancellation = &reasonOfTheCancellation.String
		}
		if reasonOfTheReactivation.Valid {
			pickupShift.ReasonOfTheReactivation = &reasonOfTheReactivation.String
		}
		if reasonOfTheResume.Valid {
			pickupShift.ReasonOfTheResume = &reasonOfTheResume.String
		}
		if reasonOfThePause.Valid {
			pickupShift.ReasonOfThePause = &reasonOfThePause.String
		}
		if deletedAt.Valid {
			pickupShift.DeletedAt = &deletedAt.Time
		}
		if startedAt.Valid {
			pickupShift.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			pickupShift.EndedAt = &endedAt.Time
		}
		if cancelledAt.Valid {
			pickupShift.CancelledAt = &cancelledAt.Time
		}
		if delayedAt.Valid {
			pickupShift.DelayedAt = &delayedAt.Time
		}
		if pausedAt.Valid {
			pickupShift.PausedAt = &pausedAt.Time
		}
		if resumedAt.Valid {
			pickupShift.ResumedAt = &resumedAt.Time
		}
		if reactivatedAt.Valid {
			pickupShift.ReactivatedAt = &reactivatedAt.Time
		}
		if staffID.Valid {
			pickupShift.Staff = &domain.CyclePickupShiftStaff{
				ID: uint(staffID.Int64),
			}
			if staffUserID.Valid {
				pickupShift.Staff.UserID = uint(staffUserID.Int64)
			}
			if staffFirstName.Valid {
				pickupShift.Staff.FirstName = staffFirstName.String
			}
			if staffLastName.Valid {
				pickupShift.Staff.LastName = staffLastName.String
			}
			if staffAvatarUrl.Valid {
				pickupShift.Staff.AvatarUrl = staffAvatarUrl.String
			}
			if staffVehicleTypes != nil {
				var vt interface{}
				err = json.Unmarshal(*staffVehicleTypes, &vt)
				if err != nil {
					return nil, err
				}
				if vt != nil {
					pickupShift.Staff.VehicleTypes = vt
					var vehicleTypes []string
					err = json.Unmarshal(*staffVehicleTypes, &vehicleTypes)
					if err != nil {
						return nil, err
					}
					// FIXME: This is a temporary solution
					if len(vehicleTypes) > 0 {
						pickupShift.Staff.SelectedVehicleForCycle = "own"
						pickupShift.Staff.SelectedVehicleTypeForCycle = vehicleTypes[0]
					} else {
						pickupShift.Staff.SelectedVehicleTypeForCycle = "public_transportation"
					}
				}
			}
			if staffVehicleLicenseTypes != nil {
				var vlt interface{}
				err = json.Unmarshal(*staffVehicleLicenseTypes, &vlt)
				if err != nil {
					return nil, err
				}
				if vlt != nil {
					pickupShift.Staff.VehicleLicenseTypes = vlt
				}
			}
		}
		if shiftID.Valid {
			pickupShift.Shift = &domain.CycleShift{
				ID: uint(shiftID.Int64),
			}
			shifts, err := r.QueryCycleShifts(&models.CyclesQueryCycleShiftsRequestParams{
				ID: int(shiftID.Int64),
			})
			if err != nil {
				log.Printf("Error while fetching shift in pickup shift: %v\n", err.Error())
			} else {
				pickupShift.Shift = shifts[0]
			}
		}
		if cycleStaffTypeCycleStaffTypeId.Valid {
			pickupShift.CycleStaffType = &domain.CyclePickupShiftCycleStaffType{
				ID: uint(cycleStaffTypeCycleStaffTypeId.Int64),
			}
			if cycleStaffTypeCycleStaffTypeRoleRoleID.Valid {
				pickupShift.CycleStaffType.Role = &domain.CycleStaffTypeRole{
					ID: uint(cycleStaffTypeCycleStaffTypeRoleRoleID.Int64),
				}
				if cycleStaffTypeCycleStaffTypeRoleRoleName.Valid {
					pickupShift.CycleStaffType.Role.Name = cycleStaffTypeCycleStaffTypeRoleRoleName.String
				}
			}
			if cycleStaffTypeCycleStaffTypeDateTime.Valid {
				pickupShift.CycleStaffType.DateTime = cycleStaffTypeCycleStaffTypeDateTime.Time
			}
			if cycleStaffTypeCycleStaffTypeShiftName.Valid {
				pickupShift.CycleStaffType.ShiftName = cycleStaffTypeCycleStaffTypeShiftName.String
			}
			if cycleStaffTypeCycleStaffTypeStartHour.Valid {
				pickupShift.CycleStaffType.StartHour = cycleStaffTypeCycleStaffTypeStartHour.Time
			}
			if cycleStaffTypeCycleStaffTypeEndHour.Valid {
				pickupShift.CycleStaffType.EndHour = cycleStaffTypeCycleStaffTypeEndHour.Time
			}
			if cycleStaffTypeCycleStaffTypeIsUnplanned.Valid {
				pickupShift.CycleStaffType.IsUnplanned = cycleStaffTypeCycleStaffTypeIsUnplanned.Bool
			}
		}
		pickupShifts = append(pickupShifts, &pickupShift)
	}

	// Find customer services by cycle pickup shifts id
	for _, pickupShift := range pickupShifts {
		var (
			shiftMorningStartTime string = "08:00"
			shiftMorningEndTime   string = "16:00"
			shiftEveningStartTime string = "16:00"
			shiftEveningEndTime   string = "00:00"
			shiftNightStartTime   string = "00:00"
			shiftNightEndTime     string = "08:00"
			shiftMorningStartHour int64
			shiftMorningEndHour   int64
			shiftEveningStartHour int64
			shiftEveningEndHour   int64
			shiftNightStartHour   int64
			shiftNightEndHour     int64
		)
		err := r.PostgresDB.QueryRow(`
			SELECT shiftMorningStartTime, shiftMorningEndTime, shiftEveningStartTime, shiftEveningEndTime, shiftNightStartTime, shiftNightEndTime
			FROM cycles
			WHERE id = $1
		`, pickupShift.Shift.ID).Scan(&shiftMorningStartTime, &shiftMorningEndTime, &shiftEveningStartTime, &shiftEveningEndTime, &shiftNightStartTime, &shiftNightEndTime)
		if err != nil {
			log.Printf("error finding shift: %v\n", err.Error())
		}
		if shiftMorningStartTime != "" {
			shiftMorningStartHour, err = strconv.ParseInt(shiftMorningStartTime[:2], 10, 64)
			if err != nil {
				shiftMorningStartHour = 8
			}
		}
		if shiftMorningEndTime != "" {
			shiftMorningEndHour, err = strconv.ParseInt(shiftMorningEndTime[:2], 10, 64)
			if err != nil {
				shiftMorningEndHour = 16
			}
		}
		if shiftEveningStartTime != "" {
			shiftEveningStartHour, err = strconv.ParseInt(shiftEveningStartTime[:2], 10, 64)
			if err != nil {
				shiftEveningStartHour = 16
			}
		}
		if shiftEveningEndTime != "" {
			shiftEveningEndHour, err = strconv.ParseInt(shiftEveningEndTime[:2], 10, 64)
			if err != nil {
				shiftEveningEndHour = 0
			}
		}
		if shiftNightStartTime != "" {
			shiftNightStartHour, err = strconv.ParseInt(shiftNightStartTime[:2], 10, 64)
			if err != nil {
				shiftNightStartHour = 0
			}
		}
		if shiftNightEndTime != "" {
			shiftNightEndHour, err = strconv.ParseInt(shiftNightEndTime[:2], 10, 64)
			if err != nil {
				shiftNightEndHour = 8
			}
		}
		customerServices, err := r.CustomerServices.FindCustomerServicesForSpecificShift(
			int64(pickupShift.ID),
			pickupShift.DateTime,
			pickupShift.CycleStaffType.ShiftName,
			shiftMorningStartHour,
			shiftMorningEndHour,
			shiftEveningStartHour,
			shiftEveningEndHour,
			shiftNightStartHour,
			shiftNightEndHour,
		)
		if err != nil {
			log.Printf("error finding customer services for specific shift: %v\n", err.Error())
		}
		if customerServices == nil {
			continue
		}
		pickupShift.CustomerServices = customerServices
	}

	// Find user roles
	for _, pickupShift := range pickupShifts {
		if pickupShift.Staff != nil {
			var (
				role domain.CyclePickupShiftStaffRole
			)
			err := r.PostgresDB.QueryRow(`
				SELECT usersRoles.roleId, _roles.name
				FROM usersRoles
				INNER JOIN _roles ON _roles.id = usersRoles.roleId
				WHERE usersRoles.userId = $1
			`, pickupShift.Staff.UserID).Scan(&role.ID, &role.Name)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Printf("error finding user roles: %v\n", err.Error())
				}
				continue
			}
			pickupShift.Staff.Roles = append(pickupShift.Staff.Roles, &role)
		}
	}

	// Find user permissions
	for _, pickupShift := range pickupShifts {
		if pickupShift.Staff != nil {
			var userRoles = pickupShift.Staff.Roles
			for index, role := range pickupShift.Staff.Roles {
				var (
					permission domain.CyclePickupShiftStaffRolePermission
				)
				err := r.PostgresDB.QueryRow(`
					SELECT _rolesPermissions.permissionId, _permissions.name, _permissions.title
					FROM _rolesPermissions
					INNER JOIN _permissions ON _permissions.id = _rolesPermissions.permissionId
					WHERE _rolesPermissions.roleId = $1
				`, role.ID).Scan(&permission.ID, &permission.Name, &permission.Title)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						log.Printf("error finding user permissions: %v\n", err.Error())
					}
					continue
				}
				userRoles[index].Permissions = append(userRoles[index].Permissions, permission)
			}
			pickupShift.Staff.Roles = userRoles
		}
	}

	return pickupShifts, nil
}

// CountPickupShifts counts the number of pickup shifts based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryPickupShiftsRequestParams as an argument.
// Returns the count as an int64 and an error.
func (r *CycleRepositoryPostgresDB) CountPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(ps.id)
		FROM cyclePickupShifts ps
		JOIN cycleStaffTypes cst ON ps.cycleStaffTypeId = cst.id
		JOIN _roles r ON cst.roleId = r.id
		JOIN staffs s ON ps.staffId = s.id
		JOIN users u ON s.userId = u.id
	`
	if queries != nil {
		where := makePickupShiftsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	log.Printf("count pickup shifts query: %v\n", q)

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CreateShiftExchangeKey generates a unique exchange key for a given set of staff type IDs.
//
// It takes a slice of int64 staffTypeIds as an argument.
// Returns a string representing the exchange key and an error if any.
func (r *CycleRepositoryPostgresDB) CreateShiftExchangeKey(staffTypeIds []int64) (string, error) {
	if len(staffTypeIds) == 0 {
		return "", fmt.Errorf("staffTypeIds cannot be empty")
	}

	// Sort the staffTypeIds to ensure the key is consistent
	sort.Slice(staffTypeIds, func(i, j int) bool { return staffTypeIds[i] < staffTypeIds[j] })

	// Convert the staffTypeIds to a string and concatenate them
	var sb strings.Builder
	for _, id := range staffTypeIds {
		sb.WriteString(fmt.Sprintf("%d-", id))
	}
	concatenatedString := sb.String()

	// Generate a SHA-256 hash of the concatenated string
	hash := sha256.New()
	hash.Write([]byte(concatenatedString))
	hashedBytes := hash.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString, nil
}

// CreateShiftIfNotExist creates a new shift if it does not exist in the database.
//
// Parameters:
//   - cycleId: the ID of the cycle
//   - staffTypeIds: a list of staff type IDs
//   - shiftName: the name of the shift
//   - datetime: the date and time of the shift
//   - status: the status of the shift
//
// Return type(s):
//   - *int64: the ID of the created shift
//   - error: any error that occurred during the creation process
func (r *CycleRepositoryPostgresDB) CreateShiftIfNotExist(cycleId int64, staffTypeIds []int64, shiftName string, datetime *time.Time, status string) (*int64, error) {
	var (
		currentTime = time.Now()
		shiftID     int64
	)
	var vIDs sql.NullString
	err := r.PostgresDB.QueryRow(`
		SELECT id, staffTypeIds
		FROM cycleShifts
		WHERE shiftName = $1 AND datetime = $2
	`,
		shiftName,
		datetime,
	).Scan(&shiftID, &vIDs)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	exchangeKey, err := r.CreateShiftExchangeKey(staffTypeIds)
	if err != nil {
		return nil, err
	}
	if shiftID == 0 {
		b, err := json.Marshal(staffTypeIds)
		if err != nil {
			return nil, err
		}
		staffTypeIDsJSON := string(b)
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleShifts (exchangeKey, cycleId, staffTypeIds, shiftName, datetime, status, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`,
			exchangeKey,
			cycleId,
			staffTypeIDsJSON,
			shiftName,
			datetime,
			status,
			currentTime,
			currentTime,
			nil,
		).Scan(&shiftID)
		if err != nil {
			return nil, err
		}
	} else {
		var laststaffTypeIDs []uint
		if vIDs.Valid {
			err := json.Unmarshal([]byte(vIDs.String), &laststaffTypeIDs)
			if err != nil {
				return nil, err
			}
		}
		for _, laststaffTypeID := range laststaffTypeIDs {
			var found = false
			for _, staffTypeId := range staffTypeIds {
				if uint(staffTypeId) == laststaffTypeID {
					found = true
					break
				}
			}
			if found {
				continue
			}
			staffTypeIds = append(staffTypeIds, int64(laststaffTypeID))
		}
		exchangeKey, err = r.CreateShiftExchangeKey(staffTypeIds)
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(staffTypeIds)
		if err != nil {
			return nil, err
		}
		staffTypeIDsJSON := string(b)
		err = r.PostgresDB.QueryRow(`
            UPDATE cycleShifts
            SET exchangeKey = $1, staffTypeIds = $2, status = $3, updated_at = $4
            WHERE id = $5
            RETURNING id
        `,
			exchangeKey,
			staffTypeIDsJSON,
			status,
			currentTime,
			shiftID,
		).Scan(&shiftID)
		if err != nil {
			return nil, err
		}
	}
	return &shiftID, nil
}

// PickupShift creates a new cycle pickup shift.
//
// It takes a payload of type *models.CyclesCreatePickupShiftRequestBody, which contains the necessary information to create a new cycle pickup shift.
// Returns a slice of *domain.CyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) PickupShift(payload *models.CyclesCreatePickupShiftRequestBody) ([]*domain.CyclePickupShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: payload.CycleID,
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Create shift if not exists in cycleShifts
	var shiftStatus = constants.SHIFT_STATUS_NOT_STARTED
	shiftID, err := r.CreateShiftIfNotExist(int64(cycle.ID), payload.CycleStaffTypeIDsInt64, payload.ShiftName, payload.DateTimeAsDate, shiftStatus)
	if err != nil {
		return nil, err
	}
	if shiftID == nil {
		return nil, errors.New("shift not found")
	}

	// Insert the cycle pickup shift
	var (
		currentTime = time.Now()
		status      = constants.VISIT_NOT_STARTED
	)
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	for _, cycleStaffTypeIDInt64 := range payload.CycleStaffTypeIDsInt64 {
		var insertedID int64
		err = tx.QueryRow(`
			INSERT INTO cyclePickupShifts (cycleId, staffId, shiftId, cycleStaffTypeId, datetime, status, created_at, updated_at, deleted_at, started_at, ended_at, cancelled_at, delayed_at, paused_at, resumed_at, reactivated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			RETURNING id
		`,
			cycle.ID,
			payload.StaffID,
			shiftID,
			cycleStaffTypeIDInt64,
			payload.DateTimeAsDate,
			status,
			currentTime,
			currentTime,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		).Scan(
			&insertedID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Query cycle pickup shifts
	cyclePickupShiftStaffs, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		CycleID: int(cycle.ID),
		StaffID: payload.StaffID,
		Filters: models.CycleQueryPickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: payload.DateTimeAsDate.Format("2006-01-02T15:04:05"),
			},
		},
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShiftStaffs) == 0 {
		return nil, errors.New("there are no cycle pickup shifts")
	}

	// Create cycleChats for online visits
	for _, cyclePickupShiftStaff := range cyclePickupShiftStaffs {
		if len(cyclePickupShiftStaff.CustomerServices) <= 0 {
			continue
		}
		var isOnline = false
		for _, customerService := range cyclePickupShiftStaff.CustomerServices {
			if customerService.VisitType != nil {
				if *customerService.VisitType == sharedconstants.VISIT_TYPE_ONLINE {
					isOnline = true
					break
				}
			}
		}
		if !isOnline {
			continue
		}
		var (
			insertedCycleChatID int64
			foundCycleChatID    int64
		)
		// Check if the cycle chat already exists with this cycleId and cyclePickupShiftId and senderUserId and recipientUserId
		err = r.PostgresDB.QueryRow(`
			SELECT id
			FROM cycleChats
			WHERE cycleId = $1
			AND cyclePickupShiftId = $2
			AND senderUserId = $3
			AND recipientUserId = $4
		`,
			cycle.ID,
			cyclePickupShiftStaff.ID,
			payload.Staff.UserID,
			cyclePickupShiftStaff.CustomerServices[0].CustomerID,
		).Scan(&foundCycleChatID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		if foundCycleChatID > 0 {
			continue
		}

		// Insert the cycle chat
		var (
			senderUserId    = payload.Staff.UserID
			recipientUserId = cyclePickupShiftStaff.CustomerServices[0].Customer.User.ID
		)

		// Insert the cycle chat
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleChats (cycleId, cyclePickupShiftId, senderUserId, recipientUserId, isSystem, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
			cycle.ID,
			cyclePickupShiftStaff.ID,
			senderUserId,
			recipientUserId,
			false,
			currentTime,
			currentTime,
		).Scan(&insertedCycleChatID)
		if err != nil {
			return nil, err
		}
		// FIXME: Remove following line
		break
	}

	return cyclePickupShiftStaffs, nil
}

// makeIncomingCyclePickupShiftsWhereFilters makes the where filters for the incoming cycle pickup shifts query
func makeIncomingCyclePickupShiftsWhereFilters(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ps.id = %d ", queries.ID))
		}
		if queries.CycleID != 0 {
			where = append(where, fmt.Sprintf(" ps.cycleId = %d ", queries.CycleID))
		}
		if queries.StaffID != 0 {
			where = append(where, fmt.Sprintf(" ps.staffId = %d ", queries.StaffID))
		}
		if len(queries.CycleNextStaffTypeIDsInt64) > 0 {
			var w []string
			for _, cycleNextStaffTypeID := range queries.CycleNextStaffTypeIDsInt64 {
				w = append(w, fmt.Sprintf(" ps.cycleNextStaffTypeId = %d ", cycleNextStaffTypeID))
			}
			where = append(where, fmt.Sprintf("(%s)", strings.Join(w, " OR ")))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.shiftName %s %s", opValue.Operator, val))
		}
		if queries.Filters.DateTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.DateTime.Op, fmt.Sprintf("%v", queries.Filters.DateTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ps.datetime %s %s", opValue.Operator, val))
		}
		if queries.RangeDateTimeStartAsTime != nil {
			where = append(where, fmt.Sprintf(" ps.datetime >= '%s' ", queries.RangeDateTimeStartAsTime.Format(time.RFC3339)))
		}
		if queries.RangeDateTimeEndAsTime != nil {
			where = append(where, fmt.Sprintf(" ps.datetime <= '%s' ", queries.RangeDateTimeEndAsTime.Format(time.RFC3339)))
		}
		if len(queries.ShiftNamesAsArray) > 0 {
			where = append(where, fmt.Sprintf(" cst.shiftName = ANY ('{%s}') ", strings.Trim(strings.Join(queries.ShiftNamesAsArray, ","), "[]")))
		}
		if queries.Filters.StartHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartHour.Op, fmt.Sprintf("%v", queries.Filters.StartHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.startHour %s %s", opValue.Operator, val))
		}
		if queries.Filters.EndHour.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndHour.Op, fmt.Sprintf("%v", queries.Filters.EndHour.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cst.endHour %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryIncomingCyclePickupShifts retrieves the incoming cycle pickup shifts from the PostgresDB based on the provided queries.
//
// Parameters:
// - queries: The CyclesQueryIncomingCyclePickupShiftsRequestParams containing the query parameters for filtering the results.
//
// Returns:
// - []*domain.CycleIncomingCyclePickupShift: The list of incoming cycle pickup shifts.
// - error: An error if the query execution fails.
func (r *CycleRepositoryPostgresDB) QueryIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) ([]*domain.CycleIncomingCyclePickupShift, error) {
	q := `
		SELECT
			ps.id,
			ps.cycleId,
			ps.staffId,
			ps.cycleNextStaffTypeId,
			ps.datetime,
			ps.created_at,
			ps.updated_at,
			ps.deleted_at,
			s.id as staffStaffID,
			s.userId as staffUserID,
			u.firstName as staffFirstName,
			u.lastName as staffLastName,
			u.avatarUrl as staffAvatarUrl,
			cst.id as cycleNextStaffTypeCycleNextStaffTypeId,
            cst.roleId as cycleNextStaffTypeCycleNextStaffTypeRoleId,
            r.id as cycleNextStaffTypeCycleNextStaffTypeRoleRoleID,
            r.name as cycleNextStaffTypeCycleNextStaffTypeRoleRoleName,
            cst.datetime as cycleNextStaffTypeCycleNextStaffTypeDateTime,
			cst.shiftName as cycleNextStaffTypeCycleNextStaffTypeShiftName,
            cst.startHour as cycleNextStaffTypeCycleNextStaffTypeStartHour,
            cst.endHour as cycleNextStaffTypeCycleNextStaffTypeEndHour
		FROM incomingCyclePickupShifts ps
		JOIN cycleNextStaffTypes cst ON ps.cycleNextStaffTypeId = cst.id
		JOIN _roles r ON cst.roleId = r.id
		JOIN staffs s ON ps.staffId = s.id
		JOIN users u ON s.userId = u.id
	`
	if queries != nil {
		where := makeIncomingCyclePickupShiftsWhereFilters(queries)
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

	var pickupShifts []*domain.CycleIncomingCyclePickupShift
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			pickupShift                                      domain.CycleIncomingCyclePickupShift
			sID                                              sql.NullInt64
			cstID                                            sql.NullInt64
			staffID                                          sql.NullInt64
			staffUserID                                      sql.NullInt64
			staffFirstName                                   sql.NullString
			staffLastName                                    sql.NullString
			staffAvatarUrl                                   sql.NullString
			deletedAt                                        sql.NullTime
			cycleNextStaffTypeCycleNextStaffTypeId           sql.NullInt64
			cycleNextStaffTypeCycleNextStaffTypeRoleId       sql.NullInt64
			cycleNextStaffTypeCycleNextStaffTypeRoleRoleID   sql.NullInt64
			cycleNextStaffTypeCycleNextStaffTypeRoleRoleName sql.NullString
			cycleNextStaffTypeCycleNextStaffTypeDateTime     sql.NullTime
			cycleNextStaffTypeCycleNextStaffTypeShiftName    sql.NullString
			cycleNextStaffTypeCycleNextStaffTypeStartHour    sql.NullTime
			cycleNextStaffTypeCycleNextStaffTypeEndHour      sql.NullTime
		)
		err := rows.Scan(
			&pickupShift.ID,
			&pickupShift.CycleID,
			&sID,
			&cstID,
			&pickupShift.DateTime,
			&pickupShift.CreatedAt,
			&pickupShift.UpdatedAt,
			&deletedAt,
			&staffID,
			&staffUserID,
			&staffFirstName,
			&staffLastName,
			&staffAvatarUrl,
			&cycleNextStaffTypeCycleNextStaffTypeId,
			&cycleNextStaffTypeCycleNextStaffTypeRoleId,
			&cycleNextStaffTypeCycleNextStaffTypeRoleRoleID,
			&cycleNextStaffTypeCycleNextStaffTypeRoleRoleName,
			&cycleNextStaffTypeCycleNextStaffTypeDateTime,
			&cycleNextStaffTypeCycleNextStaffTypeShiftName,
			&cycleNextStaffTypeCycleNextStaffTypeStartHour,
			&cycleNextStaffTypeCycleNextStaffTypeEndHour,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			pickupShift.DeletedAt = &deletedAt.Time
		}
		if staffID.Valid {
			pickupShift.Staff = &domain.CycleIncomingCyclePickupShiftStaff{
				ID: uint(staffID.Int64),
			}
			if staffUserID.Valid {
				pickupShift.Staff.UserID = uint(staffUserID.Int64)
			}
			if staffFirstName.Valid {
				pickupShift.Staff.FirstName = staffFirstName.String
			}
			if staffLastName.Valid {
				pickupShift.Staff.LastName = staffLastName.String
			}
			if staffAvatarUrl.Valid {
				pickupShift.Staff.AvatarUrl = staffAvatarUrl.String
			}
		}
		if cycleNextStaffTypeCycleNextStaffTypeId.Valid {
			pickupShift.CycleNextStaffType = &domain.CycleIncomingCyclePickupShiftCycleNextStaffType{
				ID: uint(cycleNextStaffTypeCycleNextStaffTypeId.Int64),
			}
			if cycleNextStaffTypeCycleNextStaffTypeRoleRoleID.Valid {
				pickupShift.CycleNextStaffType.Role = &domain.CycleNextStaffTypeRole{
					ID: uint(cycleNextStaffTypeCycleNextStaffTypeRoleRoleID.Int64),
				}
				if cycleNextStaffTypeCycleNextStaffTypeRoleRoleName.Valid {
					pickupShift.CycleNextStaffType.Role.Name = cycleNextStaffTypeCycleNextStaffTypeRoleRoleName.String
				}
			}
			if cycleNextStaffTypeCycleNextStaffTypeDateTime.Valid {
				pickupShift.CycleNextStaffType.DateTime = cycleNextStaffTypeCycleNextStaffTypeDateTime.Time
			}
			if cycleNextStaffTypeCycleNextStaffTypeShiftName.Valid {
				pickupShift.CycleNextStaffType.ShiftName = cycleNextStaffTypeCycleNextStaffTypeShiftName.String
			}
			if cycleNextStaffTypeCycleNextStaffTypeStartHour.Valid {
				pickupShift.CycleNextStaffType.StartHour = cycleNextStaffTypeCycleNextStaffTypeStartHour.Time
			}
			if cycleNextStaffTypeCycleNextStaffTypeEndHour.Valid {
				pickupShift.CycleNextStaffType.EndHour = cycleNextStaffTypeCycleNextStaffTypeEndHour.Time
			}
		}
		pickupShifts = append(pickupShifts, &pickupShift)
	}
	return pickupShifts, nil
}

// CountIncomingCyclePickupShifts counts the number of incoming cycle pickup shifts based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryIncomingCyclePickupShiftsRequestParams as an argument.
// Returns the count as an int64 and an error.
func (r *CycleRepositoryPostgresDB) CountIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(ps.id)
		FROM incomingCyclePickupShifts ps
		JOIN cycleNextStaffTypes cst ON ps.cycleNextStaffTypeId = cst.id
		JOIN _roles r ON cst.roleId = r.id
		JOIN staffs s ON ps.staffId = s.id
		JOIN users u ON s.userId = u.id
	`
	if queries != nil {
		where := makeIncomingCyclePickupShiftsWhereFilters(queries)
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

// PickupShiftIncomingCycle inserts a new cycle pickup shift for the given cycle ID, staff ID, and date.
//
// It takes a payload of type models.CyclesCreateIncomingCyclePickupShiftRequestBody.
// Returns a slice of domain.CycleIncomingCyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) PickupShiftIncomingCycle(payload *models.CyclesCreateIncomingCyclePickupShiftRequestBody) ([]*domain.CycleIncomingCyclePickupShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: payload.CycleID,
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Insert the cycle pickup shift
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	for _, cycleNextStaffTypeIDInt64 := range payload.CycleNextStaffTypeIDsInt64 {
		var insertedID int64
		err = tx.QueryRow(`
			INSERT INTO incomingCyclePickupShifts (cycleId, staffId, cycleNextStaffTypeId, datetime, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
			cycle.ID,
			payload.StaffID,
			cycleNextStaffTypeIDInt64,
			payload.DateTimeAsDate,
			currentTime,
			currentTime,
			nil,
		).Scan(
			&insertedID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Query cycle pickup shifts
	cycleIncomingCyclePickupShiftStaffs, err := r.QueryIncomingCyclePickupShifts(&models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		CycleID: int(cycle.ID),
		StaffID: payload.StaffID,
		Filters: models.CycleQueryIncomingCyclePickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: payload.DateTimeAsDate.Format("2006-01-02T15:04:05"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(cycleIncomingCyclePickupShiftStaffs) == 0 {
		return nil, errors.New("there are no cycle pickup shifts")
	}

	return cycleIncomingCyclePickupShiftStaffs, nil
}

// FindAllNextStaffTypesByCycleID retrieves all next staff types for a given cycle ID.
//
// The cycleID parameter is the ID of the cycle for which to retrieve next staff types.
// Returns a slice of CycleNextStaffType objects and an error if any.
func (r *CycleRepositoryPostgresDB) FindAllNextStaffTypesByCycleID(cycleID int64) ([]*domain.CycleNextStaffType, error) {
	q := `
		SELECT
			st.id,
			st.currentCycleId,
			st.roleId,
			st.datetime,
			st.shiftName,
			st.neededStaffCount,
			st.startHour,
			st.endHour,
			st.created_at,
			st.updated_at,
			st.deleted_at,
			r.id as roleId,
			r.name as roleName
		FROM cycleNextStaffTypes st
		JOIN _roles r ON r.id = st.roleId
		WHERE st.currentCycleId = $1;
	`

	var staffTypes []*domain.CycleNextStaffType
	rows, err := r.PostgresDB.Query(q, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			staffType domain.CycleNextStaffType
			role      domain.CycleNextStaffTypeRole
		)
		err := rows.Scan(
			&staffType.ID,
			&staffType.CurrentCycleID,
			&staffType.RoleID,
			&staffType.DateTime,
			&staffType.ShiftName,
			&staffType.NeededStaffCount,
			&staffType.StartHour,
			&staffType.EndHour,
			&staffType.CreatedAt,
			&staffType.UpdatedAt,
			&staffType.DeletedAt,
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		staffType.Role = &role
		staffTypes = append(staffTypes, &staffType)
	}
	st, err := r.CalculateNextStaffTypeNeededAndRemindStaffCount(staffTypes)
	if err != nil {
		return nil, err
	}
	staffTypes = st
	return staffTypes, nil
}

// FindVisitsForStaffInSpecificShift finds picked up shifts for a staff member in a specific shift.
//
// Parameters:
// - cycleID: the ID of the cycle to search for
// - staffID: the ID of the staff member to find visits for
// - datetime: the date and time for the search
// - shiftName: the name of the shift to filter by
// Returns:
// - []*domain.CyclePickupShift: a slice of picked up shifts for the staff member
// - error: an error if the operation fails
func (r *CycleRepositoryPostgresDB) FindVisitsForStaffInSpecificShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CyclePickupShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(cycleID),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		CycleID: int(cycle.ID),
		StaffID: int(staffID),
		Filters: models.CycleQueryPickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: datetime.Format("2006-01-02T15:04:05"),
			},
			ShiftName: filters.FilterValue[string]{
				Op:    "equals",
				Value: shiftName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return pickupShifts, nil
}

// AssignShiftsToStaff assigns shifts to a staff member.
//
// It takes a payload of type models.CyclesCreateShiftAssignToMeRequestBody and a target staff ID.
// It returns a slice of domain.CyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) AssignShiftsToStaff(payload *models.CyclesCreateShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CyclePickupShift, error) {
	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err := r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), int64(payload.StaffID), payload.DateAsDate, payload.ShiftName)
	if err != nil {
		return nil, err
	}

	// If there are no pickup shifts, return an error
	if len(pickupShifts) == 0 {
		return nil, errors.New("there are no pickup shifts")
	}

	// Update staffId in cyclePickupShifts to the current user
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range pickupShifts {
		_, err = tx.Exec(`
			UPDATE cyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			targetStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err = r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), targetStaffID, payload.DateAsDate, payload.ShiftName)
	if err != nil {
		return nil, err
	}

	return pickupShifts, nil
}

// SwapShifts swaps the shifts between two staff members.
//
// It takes a payload of type models.CyclesCreateShiftSwapRequestBody,
// which contains the cycle ID, source and target staff IDs, dates, and shift names.
// It returns two slices of domain.CyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) SwapShifts(payload *models.CyclesCreateShiftSwapRequestBody) ([]*domain.CyclePickupShift, []*domain.CyclePickupShift, error) {
	// Find pickedUp shifts for source staff based on CycleID StaffID Date ShiftName
	sourcePickupShifts, err := r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), int64(payload.SourceStaffID), payload.SourceDateAsDate, payload.SourceShiftName)
	if err != nil {
		return nil, nil, err
	}

	// If there are no pickup shifts, return an error
	if len(sourcePickupShifts) == 0 {
		return nil, nil, errors.New("there are no pickup shifts for source staff")
	}

	// Find pickedUp shifts for target staff based on CycleID StaffID Date ShiftName
	targetPickupShifts, err := r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), int64(payload.TargetStaffID), payload.TargetDateAsDate, payload.TargetShiftName)
	if err != nil {
		return nil, nil, err
	}

	// If there are no pickup shifts, return an error
	if len(targetPickupShifts) == 0 {
		return nil, nil, errors.New("there are no pickup shifts for target staff")
	}

	// Swap the shifts between the source and target staff
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, nil, err
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range sourcePickupShifts {
		_, err = tx.Exec(`
			UPDATE cyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			payload.TargetStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, nil, err
			}
			return nil, nil, err
		}
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range targetPickupShifts {
		_, err = tx.Exec(`
			UPDATE cyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			payload.SourceStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, nil, err
			}
			return nil, nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	// Find pickedUp shifts for source staff based on CycleID StaffID Date ShiftName
	sourcePickupShifts, err = r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), int64(payload.SourceStaffID), payload.SourceDateAsDate, payload.SourceShiftName)
	if err != nil {
		return nil, nil, err
	}

	// Find pickedUp shifts for target staff based on CycleID StaffID Date ShiftName
	targetPickupShifts, err = r.FindVisitsForStaffInSpecificShift(int64(payload.CycleID), int64(payload.TargetStaffID), payload.TargetDateAsDate, payload.TargetShiftName)
	if err != nil {
		return nil, nil, err
	}

	return sourcePickupShifts, targetPickupShifts, nil
}

// ShiftStart starts a cycle shift.
//
// payload: a pointer to a CyclesCreateShiftStartRequestBody struct containing the cycle ID, shift ID, vehicle type, start location, and other parameters needed to start the shift.
// Returns a pointer to a CycleShift struct and an error if any.
func (r *CycleRepositoryPostgresDB) ShiftStart(payload *models.CyclesCreateShiftStartRequestBody) (*domain.CycleShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: payload.CycleID,
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}

	// Find cycle shift by id
	var (
		shift                domain.CycleShift
		staffTypeIDs         sql.NullString
		staffTypeIDsMetadata []uint
		vehicleType          sql.NullString
		startLocation        sql.NullString
	)
	err = r.PostgresDB.QueryRow(`
		SELECT
			s.id,
			s.exchangeKey,
			s.cycleId,
			s.staffTypeIds,
			s.shiftName,
			s.vehicleType,
			s.startLocation,
			s.datetime,
			s.status,
			s.created_at,
			s.updated_at,
			s.deleted_at
		FROM cycleShifts s
		WHERE s.id = $1
	`,
		payload.ShiftID,
	).Scan(
		&shift.ID,
		&shift.ExchangeKey,
		&shift.CycleID,
		&staffTypeIDs,
		&shift.ShiftName,
		&vehicleType,
		&startLocation,
		&shift.DateTime,
		&shift.Status,
		&shift.CreatedAt,
		&shift.UpdatedAt,
		&shift.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("shift not found")
		}
		return nil, err
	}
	if staffTypeIDs.Valid {
		err := json.Unmarshal([]byte(staffTypeIDs.String), &staffTypeIDsMetadata)
		if err != nil {
			log.Printf("Error while unmarshalling staffType ids in query cycle shift: %v\n", err.Error())
			return nil, err
		} else {
			shift.StaffTypeIDs = staffTypeIDsMetadata
		}
	}
	if vehicleType.Valid {
		shift.VehicleType = &vehicleType.String
	}
	if startLocation.Valid {
		shift.StartLocation = &startLocation.String
	}

	// Update the status of the shift
	var (
		currentTime = time.Now()
		newStatus   = constants.SHIFT_STATUS_STARTED
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cycleShifts
		SET status = $1, vehicleType = $2, startLocation = $3, updated_at = $4
		WHERE id = $5
	`,
		newStatus,
		payload.VehicleType,
		payload.StartLocation,
		currentTime,
		shift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Find cycle shifts
	shift.Status = newStatus

	return &shift, nil
}

// ShiftEnd updates the status of a cycle shift to ended and returns the updated cycle shift.
//
// payload: a pointer to a CyclesCreateShiftEndRequestBody struct containing the cycle ID, shift ID,
// vehicle type, start location, and other parameters needed to update the shift.
// Returns a pointer to a CycleShift struct and an error if any.
func (r *CycleRepositoryPostgresDB) ShiftEnd(payload *models.CyclesCreateShiftEndRequestBody) (*domain.CycleShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: payload.CycleID,
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}

	// Find cycle shift by id
	var (
		shift                domain.CycleShift
		staffTypeIDs         sql.NullString
		staffTypeIDsMetadata []uint
		vehicleType          sql.NullString
		startLocation        sql.NullString
	)
	err = r.PostgresDB.QueryRow(`
		SELECT
			s.id,
			s.exchangeKey,
			s.cycleId,
			s.staffTypeIds,
			s.shiftName,
			s.vehicleType,
			s.startLocation,
			s.datetime,
			s.status,
			s.created_at,
			s.updated_at,
			s.deleted_at
		FROM cycleShifts s
		WHERE s.id = $1
	`,
		payload.ShiftID,
	).Scan(
		&shift.ID,
		&shift.ExchangeKey,
		&shift.CycleID,
		&staffTypeIDs,
		&shift.ShiftName,
		&vehicleType,
		&startLocation,
		&shift.DateTime,
		&shift.Status,
		&shift.CreatedAt,
		&shift.UpdatedAt,
		&shift.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("shift not found")
		}
		return nil, err
	}
	if staffTypeIDs.Valid {
		err := json.Unmarshal([]byte(staffTypeIDs.String), &staffTypeIDsMetadata)
		if err != nil {
			log.Printf("Error while unmarshalling staffType ids in query cycle shift: %v\n", err.Error())
			return nil, err
		} else {
			shift.StaffTypeIDs = staffTypeIDsMetadata
		}
	}
	if vehicleType.Valid {
		shift.VehicleType = &vehicleType.String
	}
	if startLocation.Valid {
		shift.StartLocation = &startLocation.String
	}

	// Validate vehicleType and startLocation should be equals to payload
	if shift.VehicleType != nil && *shift.VehicleType != payload.VehicleType {
		return nil, errors.New("vehicleType should be same as what was started")
	}
	if shift.StartLocation != nil && *shift.StartLocation != *payload.StartLocation {
		return nil, errors.New("startLocation should be same as what was started")
	}

	// Update the status of the shift
	var (
		currentTime = time.Now()
		newStatus   = constants.SHIFT_STATUS_ENDED
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cycleShifts
		SET status = $1, updated_at = $2
		WHERE id = $3
	`,
		newStatus,
		currentTime,
		shift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Find cycle shifts
	shift.Status = newStatus

	return &shift, nil
}

// makeShiftCustomerHomeKeysWhereFilters generates a list of WHERE clause filters for querying shift customer home keys.
//
// It takes a pointer to models.CyclesQueryShiftCustomerHomeKeysRequestParams as a parameter, which contains the query parameters.
// Returns a slice of strings representing the WHERE clause filters.
func makeShiftCustomerHomeKeysWhereFilters(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" s.id = %d ", queries.ID))
		}
		if queries.ShiftID != 0 {
			where = append(where, fmt.Sprintf(" s.shiftId = %d ", queries.ShiftID))
		}
		if queries.Filters.KeyNo.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.KeyNo.Op, fmt.Sprintf("%v", queries.Filters.KeyNo.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" s.keyNo %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" s.status %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryShiftCustomerHomeKeys retrieves a list of shift customer home keys based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryShiftCustomerHomeKeysRequestParams as an argument.
// Returns a list of pointers to domain.CycleShiftCustomerHomeKey and an error.
func (r *CycleRepositoryPostgresDB) QueryShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) ([]*domain.CycleShiftCustomerHomeKey, error) {
	q := `
		SELECT
			s.id,
			s.shiftId,
			s.keyNo,
			s.status,
			s.reason,
			s.created_at,
			s.created_by,
			s.updated_at,
			s.deleted_at,
			u.id as userId,
			u.firstName as userFirstName,
			u.lastName as userLastName,
			u.avatarUrl as userAvatarUrl
		FROM cycleShiftCustomerHomeKeys s
		LEFT JOIN users u ON s.created_by = u.id
	`
	if queries != nil {
		where := makeShiftCustomerHomeKeysWhereFilters(queries)
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

	var shiftCustomerHomeKeys []*domain.CycleShiftCustomerHomeKey
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			shiftCustomerHomeKey domain.CycleShiftCustomerHomeKey
			reason               sql.NullString
			createdBy            sql.NullInt64
			deletedAt            sql.NullTime
			userId               sql.NullInt64
			userFirstName        sql.NullString
			userLastName         sql.NullString
			userAvatarUrl        sql.NullString
		)
		err := rows.Scan(
			&shiftCustomerHomeKey.ID,
			&shiftCustomerHomeKey.ShiftID,
			&shiftCustomerHomeKey.KeyNo,
			&shiftCustomerHomeKey.Status,
			&reason,
			&shiftCustomerHomeKey.CreatedAt,
			&createdBy,
			&shiftCustomerHomeKey.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if reason.Valid {
			shiftCustomerHomeKey.Reason = &reason.String
		}
		if userId.Valid {
			shiftCustomerHomeKey.CreatedBy = &domain.CycleShiftCustomerHomeKeyCreatedBy{
				ID: userId.Int64,
			}
			if userFirstName.Valid {
				shiftCustomerHomeKey.CreatedBy.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				shiftCustomerHomeKey.CreatedBy.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				shiftCustomerHomeKey.CreatedBy.AvatarUrl = userAvatarUrl.String
			}
		}
		if deletedAt.Valid {
			shiftCustomerHomeKey.DeletedAt = &deletedAt.Time
		}
		shiftCustomerHomeKeys = append(shiftCustomerHomeKeys, &shiftCustomerHomeKey)
	}
	return shiftCustomerHomeKeys, nil
}

// CountShiftCustomerHomeKeys counts the number of shift customer home keys based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryShiftCustomerHomeKeysRequestParams as an argument.
// Returns the count as an int64 and an error.
func (r *CycleRepositoryPostgresDB) CountShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(s.id)
		FROM cycleShiftCustomerHomeKeys s
		LEFT JOIN users u ON s.created_by = u.id
	`
	if queries != nil {
		where := makeShiftCustomerHomeKeysWhereFilters(queries)
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

// CreateCycleShiftCustomerHomeKey creates a new cycle shift customer home key or updates an existing one.
//
// It takes a payload of type *models.CyclesShiftCustomerHomeKeyRequestBody as an argument.
// Returns a *domain.CycleShiftCustomerHomeKey and an error.
func (r *CycleRepositoryPostgresDB) CreateCycleShiftCustomerHomeKey(payload *models.CyclesShiftCustomerHomeKeyRequestBody) (*domain.CycleShiftCustomerHomeKey, error) {
	// Find shift by id
	shifts, err := r.QueryCycleShifts(&models.CyclesQueryCycleShiftsRequestParams{
		ID: payload.ShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(shifts) == 0 {
		return nil, errors.New("shift not found")
	}

	// Find cycle shift customer home key by id
	shiftCustomerHomeKeys, err := r.QueryShiftCustomerHomeKeys(&models.CyclesQueryShiftCustomerHomeKeysRequestParams{
		ShiftID: payload.ShiftID,
	})
	if err != nil {
		return nil, err
	}
	var (
		foundKeyNo                    = false
		foundedShiftCustomerHomeKeyID int64
	)
	if len(shiftCustomerHomeKeys) > 0 {
		for _, shiftCustomerHomeKey := range shiftCustomerHomeKeys {
			if shiftCustomerHomeKey.KeyNo == payload.KeyNo {
				foundKeyNo = true
				foundedShiftCustomerHomeKeyID = int64(shiftCustomerHomeKey.ID)
				break
			}
		}
	}

	// Update if found keyNo, else create new
	var (
		id          int64
		currentTime = time.Now()
	)
	if foundKeyNo {
		// Update the cycle shift customer home key
		_, err = r.PostgresDB.Exec(`
			UPDATE cycleShiftCustomerHomeKeys
			SET status = $1, reason = $2, updated_at = $3
			WHERE id = $4
		`,
			payload.Status,
			payload.Reason,
			currentTime,
			foundedShiftCustomerHomeKeyID,
		)
		if err != nil {
			return nil, err
		}
		id = foundedShiftCustomerHomeKeyID
	} else {
		// Insert the cycle shift customer home key
		err = r.PostgresDB.QueryRow(`
			INSERT INTO cycleShiftCustomerHomeKeys (shiftId, keyNo, status, reason, created_at, created_by, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`,
			payload.ShiftID,
			payload.KeyNo,
			payload.Status,
			payload.Reason,
			currentTime,
			payload.AuthenticatedUser.ID,
			currentTime,
			nil,
		).Scan(
			&id,
		)
		if err != nil {
			return nil, err
		}
	}

	// Find cycle shift customer home key by id
	shiftCustomerHomeKeys, err = r.QueryShiftCustomerHomeKeys(&models.CyclesQueryShiftCustomerHomeKeysRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(shiftCustomerHomeKeys) == 0 {
		return nil, errors.New("shift customer home key not found")
	}

	return shiftCustomerHomeKeys[0], nil
}

// VisitStart updates the status of a cycle pickup shift to started.
//
// It takes a payload of type *models.CyclesCreateVisitStartRequestBody,
// which contains the ID of the cycle pickup shift to be updated.
// Returns the updated cycle pickup shift and an error if any.
func (r *CycleRepositoryPostgresDB) VisitStart(payload *models.CyclesCreateVisitStartRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_STARTED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, started_at = $2, startKilometer = $3, updated_at = $4, prevStatus = $5
		WHERE id = $6
	`,
		newStatus,
		currentTime,
		payload.StartKilometer,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.StartKilometer = &payload.StartKilometer
	cyclePickupShift.StartedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitEnd updates the status of a cycle pickup shift to ended.
//
// It takes a payload of type *models.CyclesCreateVisitEndRequestBody,
// which contains the ID of the cycle pickup shift to be updated.
// Returns the updated cycle pickup shift and an error if any.
func (r *CycleRepositoryPostgresDB) VisitEnd(payload *models.CyclesCreateVisitEndRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_ENDED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, ended_at = $2, updated_at = $3, prevStatus = $4
		WHERE id = $5
	`,
		newStatus,
		currentTime,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.EndedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitCancel updates the status of a cycle pickup shift to cancelled.
//
// It takes a payload of type models.CyclesCreateVisitCancelRequestBody,
// which contains the ID of the cycle pickup shift to be updated along with the reason for cancellation.
// Returns the updated cycle pickup shift and an error if any.
func (r *CycleRepositoryPostgresDB) VisitCancel(payload *models.CyclesCreateVisitCancelRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_CANCELLED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, cancelled_at = $2, reasonOfTheCancellation = $3, updated_at = $4, prevStatus = $5
		WHERE id = $6
	`,
		newStatus,
		currentTime,
		payload.Reason,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.ReasonOfTheCancellation = &payload.Reason
	cyclePickupShift.CancelledAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitDelay updates the status of a cycle pickup shift to delayed.
//
// It takes a payload of type models.CyclesCreateVisitDelayRequestBody,
// which contains the ID of the cycle pickup shift to be updated.
// Returns the updated cycle pickup shift and an error if any.
func (r *CycleRepositoryPostgresDB) VisitDelay(payload *models.CyclesCreateVisitDelayRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_DELAYED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, delayed_at = $2, updated_at = $3, prevStatus = $4
		WHERE id = $5
	`,
		newStatus,
		currentTime,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.DelayedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitPause pauses a visit in the cycle pickup shift.
//
// It takes a CyclesCreateVisitPauseRequestBody payload as an argument, which contains the ID of the cycle pickup shift to pause and the reason for the pause.
// It returns the updated CyclePickupShift and an error if any occurs during the operation.
func (r *CycleRepositoryPostgresDB) VisitPause(payload *models.CyclesCreateVisitPauseRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_PAUSED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, paused_at = $2, reasonOfThePause = $3, updated_at = $4, prevStatus = $5
		WHERE id = $6
	`,
		newStatus,
		currentTime,
		payload.Reason,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.ReasonOfThePause = &payload.Reason
	cyclePickupShift.PausedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitResume resumes a visit for a cycle pickup shift.
//
// It takes a payload of type models.CyclesCreateVisitResumeRequestBody, which contains the ID of the cycle pickup shift to resume and the reason for the resume.
// It returns the updated cycle pickup shift and an error if any occurs during the process.
func (r *CycleRepositoryPostgresDB) VisitResume(payload *models.CyclesCreateVisitResumeRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = constants.VISIT_RESUMED
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, resumed_at = $2, reasonOfTheResume = $3, updated_at = $4, prevStatus = $5
		WHERE id = $6
	`,
		newStatus,
		currentTime,
		payload.Reason,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.ReasonOfTheResume = &payload.Reason
	cyclePickupShift.ResumedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// VisitReactive reactivates a previously ended or delayed pickup shift.
//
// It takes a CyclesCreateVisitReactiveRequestBody payload as an argument, which contains the ID of the pickup shift to reactivate and the reason for the reactivation.
// Returns the updated CyclePickupShift and an error if any occurs during the process.
func (r *CycleRepositoryPostgresDB) VisitReactive(payload *models.CyclesCreateVisitReactiveRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the status of the pickup shift
	var (
		currentTime = time.Now()
		newStatus   = cyclePickupShift.PrevStatus
		prevStatus  = cyclePickupShift.Status
	)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET status = $1, reactivated_at = $2, reasonOfTheReactivation = $3, updated_at = $4, prevStatus = $5
		WHERE id = $6
	`,
		newStatus,
		currentTime,
		payload.Reason,
		currentTime,
		prevStatus,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the status of the pickup shift
	cyclePickupShift.Status = newStatus
	cyclePickupShift.PrevStatus = prevStatus
	cyclePickupShift.ReasonOfTheReactivation = &payload.Reason
	cyclePickupShift.ReactivatedAt = &currentTime
	cyclePickupShift.UpdatedAt = currentTime

	return cyclePickupShift, nil
}

// AssignVisitToStaff assigns a visit to a staff member.
//
// It takes a payload of type models.CyclesCreateVisitAssignRequestBody, which contains the ID of the cycle pickup shift and the ID of the staff member to assign the visit to.
// It returns the updated cycle pickup shift and an error if the operation fails.
func (r *CycleRepositoryPostgresDB) AssignVisitToStaff(payload *models.CyclesCreateVisitAssignRequestBody) (*domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift := cyclePickupShifts[0]

	// Update the staffId in the cyclePickupShifts
	var currentTime = time.Now()
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET staffId = $1, updated_at = $2
		WHERE id = $3
	`,
		payload.StaffID,
		currentTime,
		cyclePickupShift.ID,
	)
	if err != nil {
		return nil, err
	}

	// Set the staffId of the pickup shift
	cyclePickupShifts, err = r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}
	cyclePickupShift = cyclePickupShifts[0]

	return cyclePickupShift, nil
}

// SwapVisits swaps the staffId of two cyclePickupShifts.
//
// It takes a payload of type models.CyclesCreateVisitSwapRequestBody, which contains the IDs of the source and target cyclePickupShifts to be swapped.
// It returns the updated source and target cyclePickupShifts, along with an error if the operation fails.
func (r *CycleRepositoryPostgresDB) SwapVisits(payload *models.CyclesCreateVisitSwapRequestBody) (*domain.CyclePickupShift, *domain.CyclePickupShift, error) {
	// Find cyclePickupShift by id
	sourceCyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.SourceCyclePickupShiftID,
	})
	if err != nil {
		return nil, nil, err
	}
	if len(sourceCyclePickupShifts) == 0 {
		return nil, nil, errors.New("source visit not found")
	}
	sourceCyclePickupShift := sourceCyclePickupShifts[0]

	// Find cyclePickupShift by id
	targetCyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.TargetCyclePickupShiftID,
	})
	if err != nil {
		return nil, nil, err
	}
	if len(targetCyclePickupShifts) == 0 {
		return nil, nil, errors.New("target visit not found")
	}
	targetCyclePickupShift := targetCyclePickupShifts[0]

	// Swap the staffId in the cyclePickupShifts
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, nil, err
	}

	// Update the staffId in the cyclePickupShifts
	_, err = tx.Exec(`
		UPDATE cyclePickupShifts
		SET staffId = $1, updated_at = $2
		WHERE id = $3
	`,
		targetCyclePickupShift.Staff.ID,
		currentTime,
		sourceCyclePickupShift.ID,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, nil, err
		}
		return nil, nil, err
	}

	// Update the staffId in the cyclePickupShifts
	_, err = tx.Exec(`
		UPDATE cyclePickupShifts
		SET staffId = $1, updated_at = $2
		WHERE id = $3
	`,
		sourceCyclePickupShift.Staff.ID,
		currentTime,
		targetCyclePickupShift.ID,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, nil, err
		}
		return nil, nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	// Find cyclePickupShift by id
	sourceCyclePickupShifts, err = r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.SourceCyclePickupShiftID,
	})
	if err != nil {
		return nil, nil, err
	}

	// Find cyclePickupShift by id
	targetCyclePickupShifts, err = r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.TargetCyclePickupShiftID,
	})
	if err != nil {
		return nil, nil, err
	}

	return sourceCyclePickupShifts[0], targetCyclePickupShifts[0], nil
}

// CreateUnplannedVisit creates an unplanned visit for a cycle.
//
// It takes a CyclesCreateVisitUnplannedRequestBody as a parameter, which includes the cycle ID, staff ID, date, time, length, and other relevant information.
// It returns a CyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) CreateUnplannedVisit(payload *models.CyclesCreateVisitUnplannedRequestBody) (*domain.CyclePickupShift, error) {
	// Get service grade
	qr, err := r.ServiceGradeService.Query(&sgmodels.ServiceGradesQueryRequestParams{
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	serviceGrades, ok := qr.Items.([]*sgDomain.ServiceGrade)
	if !ok {
		return nil, errors.New("service grade not found")
	}
	if len(serviceGrades) == 0 {
		return nil, errors.New("service grade not found")
	}
	serviceGrade := serviceGrades[0]

	// Validate staff roles
	if len(payload.Staff.Roles) <= 0 {
		return nil, errors.New("staff roles not found")
	}

	// Create customer service
	var (
		role               = payload.Staff.Roles[0]
		startTimeRange     = payload.Time
		startTimeRangeTime = payload.TimeAsTime
		endTimeRangeTime   = payload.TimeAsTime.Add(time.Minute * time.Duration(payload.Length))
		endTimeRange       = endTimeRangeTime.Format("15:04")
		reportType         = "medical"
		paymentMethod      = "own"
		homeCareFee        = 0
		repeat             = ""
		gradeID            = serviceGrade.ID
		shiftName          = shifts.MorningShift
	)

	// Set shiftName based on start and end time
	// TODO: Read this later from settings i used also in FindCustomerServicesForSpecificShift method
	if payload.TimeAsTime.Hour() >= 6 && payload.TimeAsTime.Hour() < 14 {
		shiftName = shifts.MorningShift
	}
	if payload.TimeAsTime.Hour() >= 14 && payload.TimeAsTime.Hour() < 22 {
		shiftName = shifts.EveningShift
	}
	if payload.TimeAsTime.Hour() >= 22 || payload.TimeAsTime.Hour() < 6 {
		shiftName = shifts.NightShift
	}

	// Create customer service with staff
	customerService, err := r.CustomerServices.CreateServices(payload.Customer, &cmodels.CustomersCreateServicesRequestBody{
		CustomerID:           payload.CustomerID,
		ServiceID:            payload.ServiceID,
		ServiceTypeID:        payload.ServiceTypeID,
		GradeID:              int(gradeID),
		NurseWishID:          payload.StaffID,
		ReportType:           reportType,
		TimeValue:            payload.Time,
		Repeat:               repeat,
		VisitType:            sharedconstants.VISIT_TYPE_ONSITE,
		ServiceLengthMinute:  payload.Length,
		StartTimeRange:       startTimeRange,
		EndTimeRange:         &endTimeRange,
		Description:          nil,
		PaymentMethod:        paymentMethod,
		HomeCareFee:          &homeCareFee,
		CityCouncilFee:       nil,
		TimeValueAsTime:      payload.TimeAsTime,
		StartTimeRangeAsTime: startTimeRangeTime,
		EndTimeRangeAsTime:   &endTimeRangeTime,
	})
	if err != nil {
		return nil, err
	}

	// Create cycle staff type
	cycleStaffType, err := r.UpdateStaffType(&models.CyclesUpdateStaffTypeRequestBody{
		ShiftName:        shiftName,
		DateTime:         payload.Date,
		StartHour:        startTimeRange,
		EndHour:          endTimeRange,
		NeededStaffCount: 1,
		RoleID:           int(role.ID),
		DateTimeAsDate:   payload.DateAsDate,
		Role: &domain.CycleStaffTypeRole{
			ID:   role.ID,
			Name: role.Name,
		},
		StartHourAsTime: startTimeRangeTime,
		EndHourAsTime:   &endTimeRangeTime,
	}, int64(payload.CycleID), true)
	if err != nil {
		_, _ = r.CustomerServices.DeleteServices(&cmodels.CustomersDeleteServicesRequestBody{
			IDs:        []int64{int64(customerService.ID)},
			IDsInt64:   []int64{int64(customerService.ID)},
			CustomerID: payload.CustomerID,
		})
		return nil, err
	}

	// Update isUnplanned to true
	_, err = r.PostgresDB.Exec(`
		UPDATE cycleStaffTypes
		SET isUnplanned = $1
		WHERE id = $2
	`,
		true,
		cycleStaffType.ID,
	)
	if err != nil {
		return nil, err
	}

	// Create cycle pickup shift
	var cycleStaffTypeIDsInt64 = []int64{
		int64(cycleStaffType.ID),
	}
	cyclePickupShifts, err := r.PickupShift(&models.CyclesCreatePickupShiftRequestBody{
		CycleID:                payload.CycleID,
		StaffID:                payload.StaffID,
		CycleStaffTypeIDs:      cycleStaffTypeIDsInt64,
		DateTime:               payload.Date,
		Staff:                  payload.Staff,
		CycleStaffTypeIDsInt64: cycleStaffTypeIDsInt64,
		DateTimeAsDate:         payload.DateAsDate,
		ShiftName:              shiftName,
	})
	if err != nil {
		_, _ = r.CustomerServices.DeleteServices(&cmodels.CustomersDeleteServicesRequestBody{
			IDs:        []int64{int64(customerService.ID)},
			IDsInt64:   []int64{int64(customerService.ID)},
			CustomerID: payload.CustomerID,
		})

		// Drop staff type by id
		_, _ = r.PostgresDB.Exec(`DELETE FROM cycleStaffTypes WHERE id = $1`, cycleStaffType.ID)
		return nil, err
	}
	cyclePickupShift := cyclePickupShifts[0]
	cyclePickupShift.CustomerServices = []*csDomain.CustomerServices{customerService}

	// Update isUnplanned to true in cyclePickupShifts
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShifts
		SET isUnplanned = $1
		WHERE id = $2
	`,
		true,
		cyclePickupShift.ID,
	)
	if err != nil {
		_, _ = r.CustomerServices.DeleteServices(&cmodels.CustomersDeleteServicesRequestBody{
			IDs:        []int64{int64(customerService.ID)},
			IDsInt64:   []int64{int64(customerService.ID)},
			CustomerID: payload.CustomerID,
		})

		// Drop staff type by id
		_, _ = r.PostgresDB.Exec(`DELETE FROM cycleStaffTypes WHERE id = $1`, cycleStaffType.ID)

		// Drop cycle pickup shift by id
		_, _ = r.PostgresDB.Exec(`DELETE FROM cyclePickupShifts WHERE id = $1`, cyclePickupShift.ID)
		return nil, err
	}

	// Set isUnplanned to true
	cyclePickupShift.IsUnplanned = true

	return cyclePickupShift, nil
}

// makeVisitTodosWhereFilters generates a list of WHERE filters for a database query based on the provided CyclesQueryVisitsTodosRequestParams.
//
// It takes a pointer to models.CyclesQueryVisitsTodosRequestParams as an input parameter.
// Returns a slice of strings representing the WHERE filters.
func makeVisitTodosWhereFilters(queries *models.CyclesQueryVisitsTodosRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" s.id = %d ", queries.ID))
		}
		if queries.CyclePickupShiftID != 0 {
			where = append(where, fmt.Sprintf(" s.cyclePickupShiftId = %d ", queries.CyclePickupShiftID))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" s.status %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryVisitTodos retrieves visit todos from the database based on the given queries.
//
// Parameters:
// - queries: A pointer to a CyclesQueryVisitsTodosRequestParams struct containing the query parameters.
//
// Returns:
// - A slice of CyclePickupShiftTodo pointers representing the retrieved visit todos.
// - An error if the query execution fails.
func (r *CycleRepositoryPostgresDB) QueryVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) ([]*domain.CyclePickupShiftTodo, error) {
	q := `
		SELECT
			s.id,
			s.cyclePickupShiftId,
			s.title,
			s.timeValue,
			s.dateValue,
			s.description,
			s.attachments,
			s.notDoneReason,
			s.status,
			s.done_at,
			s.not_done_at,
			s.created_at,
			s.updated_at,
			s.deleted_at,
			u.id as userId,
			u.firstName as userFirstName,
			u.lastName as userLastName,
			u.avatarUrl as userAvatarUrl
		FROM cyclePickupShiftTodos s
		LEFT JOIN users u ON s.createdBy = u.id
	`
	if queries != nil {
		where := makeVisitTodosWhereFilters(queries)
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

	var visitTodos []*domain.CyclePickupShiftTodo
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			visitTodo           domain.CyclePickupShiftTodo
			description         sql.NullString
			attachments         json.RawMessage
			attachmentsMetadata []*types.UploadMetadata
			notDoneReason       sql.NullString
			doneAt              sql.NullTime
			notDoneAt           sql.NullTime
			deletedAt           sql.NullTime
			userId              sql.NullInt64
			userFirstName       sql.NullString
			userLastName        sql.NullString
			userAvatarUrl       sql.NullString
		)
		err := rows.Scan(
			&visitTodo.ID,
			&visitTodo.CyclePickupShiftID,
			&visitTodo.Title,
			&visitTodo.TimeValue,
			&visitTodo.DateValue,
			&description,
			&attachments,
			&notDoneReason,
			&visitTodo.Status,
			&doneAt,
			&notDoneAt,
			&visitTodo.CreatedAt,
			&visitTodo.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			visitTodo.Description = &description.String
		}
		if attachments != nil {
			err = json.Unmarshal(attachments, &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in visit todo: %d", err, visitTodo.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CYCLE_BUCKET_NAME[len("maja."):])
				}
			}
		}
		visitTodo.Attachments = attachmentsMetadata
		if notDoneReason.Valid {
			visitTodo.NotDoneReason = &notDoneReason.String
		}
		if doneAt.Valid {
			visitTodo.DoneAt = &doneAt.Time
		}
		if notDoneAt.Valid {
			visitTodo.NotDoneAt = &notDoneAt.Time
		}
		if deletedAt.Valid {
			visitTodo.DeletedAt = &deletedAt.Time
		}
		if userId.Valid {
			visitTodo.CreatedBy = &domain.CyclePickupShiftTodoCreatedBy{
				ID: userId.Int64,
			}
			if userFirstName.Valid {
				visitTodo.CreatedBy.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				visitTodo.CreatedBy.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				visitTodo.CreatedBy.AvatarUrl = userAvatarUrl.String
			}
		}
		visitTodos = append(visitTodos, &visitTodo)
	}
	return visitTodos, nil
}

// CountVisitTodos returns the total number of visit todos that match the given query parameters.
//
// Parameters:
// - queries: A pointer to a CyclesQueryVisitsTodosRequestParams struct containing the query parameters.
//
// Returns:
// - An int64 representing the total number of visit todos that match the given query parameters.
// - An error if any error occurs during the database query.
func (r *CycleRepositoryPostgresDB) CountVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(s.id)
		FROM cyclePickupShiftTodos s
		LEFT JOIN users u ON s.createdBy = u.id
	`
	if queries != nil {
		where := makeVisitTodosWhereFilters(queries)
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

// CreateVisitTodo creates a new visit todo for a cycle pickup shift.
//
// It takes a CyclesCreateVisitTodoRequestBody as a parameter, which contains the necessary information to create a new visit todo.
// It returns a CyclePickupShiftTodo and an error.
func (r *CycleRepositoryPostgresDB) CreateVisitTodo(payload *models.CyclesCreateVisitTodoRequestBody) (*domain.CyclePickupShiftTodo, error) {
	// Find cyclePickupShift by id
	cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: payload.CyclePickupShiftID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShifts) == 0 {
		return nil, errors.New("visit not found")
	}

	// Create cycle pickup shift todo
	var (
		currentTime = time.Now()
		status      = constants.TOOD_STATUS_PENDING
		insertedID  int
	)
	err = r.PostgresDB.QueryRow(`
		INSERT INTO cyclePickupShiftTodos (cyclePickupShiftId, title, timeValue, dateValue, description, status, created_at, updated_at, deleted_at, createdBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		payload.CyclePickupShiftID,
		payload.Title,
		payload.TimeValueAsTime,
		payload.DateValueAsDate,
		payload.Description,
		status,
		currentTime,
		currentTime,
		nil,
		payload.AuthenticatedUser.ID,
	).Scan(
		&insertedID,
	)
	if err != nil {
		return nil, err
	}

	// Find cycle pickup shift todo by id
	cyclePickupShiftTodos, err := r.QueryVisitTodos(&models.CyclesQueryVisitsTodosRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShiftTodos) == 0 {
		return nil, errors.New("visit todo not found")
	}

	return cyclePickupShiftTodos[0], nil
}

// UpdateVisitTodoStatus updates the status of a visit todo.
//
// It takes a payload of type models.CyclesUpdateVisitTodoStatusRequestBody and an id of type int64 as parameters.
// It returns a pointer to domain.CyclePickupShiftTodo and an error.
func (r *CycleRepositoryPostgresDB) UpdateVisitTodoStatus(payload *models.CyclesUpdateVisitTodoStatusRequestBody, id int64) (*domain.CyclePickupShiftTodo, error) {
	// Find cyclePickupShiftTodo by id
	cyclePickupShiftTodos, err := r.QueryVisitTodos(&models.CyclesQueryVisitsTodosRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShiftTodos) == 0 {
		return nil, errors.New("visit todo not found")
	}
	cyclePickupShiftTodo := cyclePickupShiftTodos[0]

	// Update the status of the pickup shift todo
	var (
		currentTime = time.Now()
	)
	if payload.Status == constants.TODO_STATUS_DONE {
		_, err = r.PostgresDB.Exec(`
			UPDATE cyclePickupShiftTodos
			SET status = $1, updated_at = $2, done_at = $3
			WHERE id = $4
		`,
			payload.Status,
			currentTime,
			currentTime,
			id,
		)
		if err != nil {
			return nil, err
		}

		// Set the new data
		cyclePickupShiftTodo.Status = payload.Status
		cyclePickupShiftTodo.UpdatedAt = currentTime
		cyclePickupShiftTodo.DoneAt = &currentTime
	} else {
		_, err = r.PostgresDB.Exec(`
			UPDATE cyclePickupShiftTodos
			SET status = $1, updated_at = $2, not_done_at = $3, notDoneReason = $4
			WHERE id = $5
		`,
			payload.Status,
			currentTime,
			currentTime,
			payload.NotDoneReason,
			id,
		)
		if err != nil {
			return nil, err
		}

		// Set the new data
		cyclePickupShiftTodo.Status = payload.Status
		cyclePickupShiftTodo.UpdatedAt = currentTime
		cyclePickupShiftTodo.NotDoneAt = &currentTime
		cyclePickupShiftTodo.NotDoneReason = payload.NotDoneReason
	}

	return cyclePickupShiftTodo, nil
}

// UpdateVisitTodoAttachments updates the attachments of a visit todo in the database.
//
// Parameters:
// - attachments: a slice of *types.UploadMetadata representing the new attachments to be set.
// - id: an int64 representing the ID of the visit todo to be updated.
//
// Returns:
// - *domain.CyclePickupShiftTodo: the updated visit todo with the new attachments.
// - error: an error if the update operation fails.
func (r *CycleRepositoryPostgresDB) UpdateVisitTodoAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CyclePickupShiftTodo, error) {
	// Find cyclePickupShiftTodo by id
	cyclePickupShiftTodos, err := r.QueryVisitTodos(&models.CyclesQueryVisitsTodosRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cyclePickupShiftTodos) == 0 {
		return nil, errors.New("visit todo not found")
	}
	cyclePickupShiftTodo := cyclePickupShiftTodos[0]

	// Update the attachments of the pickup shift todo
	var (
		currentTime = time.Now()
	)
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)
	_, err = r.PostgresDB.Exec(`
		UPDATE cyclePickupShiftTodos
		SET attachments = $1, updated_at = $2
		WHERE id = $3
	`,
		attachmentsJSON,
		currentTime,
		id,
	)
	if err != nil {
		return nil, err
	}

	// Set the attachments of the pickup shift todo
	cyclePickupShiftTodo.Attachments = attachments
	cyclePickupShiftTodo.UpdatedAt = currentTime

	return cyclePickupShiftTodo, nil
}

// FindVisitsForStaffInSpecificIncomingShift finds picked up shifts for a staff member in a specific incoming shift.
//
// Parameters:
// - cycleID: the ID of the cycle to search for
// - staffID: the ID of the staff member to find visits for
// - datetime: the date and time for the search
// - shiftName: the name of the shift to filter by
// Returns:
// - []*domain.CycleIncomingCyclePickupShift: a slice of picked up shifts for the staff member
// - error: an error if the operation fails
func (r *CycleRepositoryPostgresDB) FindVisitsForStaffInSpecificIncomingShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CycleIncomingCyclePickupShift, error) {
	// Find cycle by id
	cycles, err := r.Query(&models.CyclesQueryRequestParams{
		ID: int(cycleID),
	})
	if err != nil {
		return nil, err
	}
	if len(cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	cycle := cycles[0]

	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err := r.QueryIncomingCyclePickupShifts(&models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		CycleID: int(cycle.ID),
		StaffID: int(staffID),
		Filters: models.CycleQueryIncomingCyclePickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: datetime.Format("2006-01-02T15:04:05"),
			},
			ShiftName: filters.FilterValue[string]{
				Op:    "equals",
				Value: shiftName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return pickupShifts, nil
}

// AssignIncomingShiftsToStaff assigns incoming shifts to a staff member.
//
// It takes a payload of type models.CyclesCreateIncomingShiftAssignToMeRequestBody and a target staff ID.
// It returns a slice of domain.CycleIncomingCyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) AssignIncomingShiftsToStaff(payload *models.CyclesCreateIncomingShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CycleIncomingCyclePickupShift, error) {
	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err := r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), int64(payload.StaffID), payload.DateAsDate, payload.ShiftName)
	if err != nil {
		return nil, err
	}

	// If there are no pickup shifts, return an error
	if len(pickupShifts) == 0 {
		return nil, errors.New("there are no pickup shifts")
	}

	// Update staffId in cyclePickupShifts to the current user
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range pickupShifts {
		_, err = tx.Exec(`
			UPDATE incomingCyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			targetStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Find pickedUp shifts for staff based on CycleID StaffID Date ShiftName
	pickupShifts, err = r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), targetStaffID, payload.DateAsDate, payload.ShiftName)
	if err != nil {
		return nil, err
	}

	return pickupShifts, nil
}

// SwapIncomingShifts swaps the incoming shifts between two staff members.
//
// It takes a payload of type models.CyclesCreateIncomingShiftSwapRequestBody,
// which contains the cycle ID, source and target staff IDs, dates, and shift names.
// It returns two slices of domain.CycleIncomingCyclePickupShift and an error.
func (r *CycleRepositoryPostgresDB) SwapIncomingShifts(payload *models.CyclesCreateIncomingShiftSwapRequestBody) ([]*domain.CycleIncomingCyclePickupShift, []*domain.CycleIncomingCyclePickupShift, error) {
	// Find pickedUp shifts for source staff based on CycleID StaffID Date ShiftName
	sourcePickupShifts, err := r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), int64(payload.SourceStaffID), payload.SourceDateAsDate, payload.SourceShiftName)
	if err != nil {
		return nil, nil, err
	}

	// If there are no pickup shifts, return an error
	if len(sourcePickupShifts) == 0 {
		return nil, nil, errors.New("there are no pickup shifts for source staff")
	}

	// Find pickedUp shifts for target staff based on CycleID StaffID Date ShiftName
	targetPickupShifts, err := r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), int64(payload.TargetStaffID), payload.TargetDateAsDate, payload.TargetShiftName)
	if err != nil {
		return nil, nil, err
	}

	// If there are no pickup shifts, return an error
	if len(targetPickupShifts) == 0 {
		return nil, nil, errors.New("there are no pickup shifts for target staff")
	}

	// Swap the shifts between the source and target staff
	var currentTime = time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, nil, err
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range sourcePickupShifts {
		_, err = tx.Exec(`
			UPDATE incomingCyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			payload.TargetStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, nil, err
			}
			return nil, nil, err
		}
	}

	// Update the staffId in the cyclePickupShifts
	for _, pickupShift := range targetPickupShifts {
		_, err = tx.Exec(`
			UPDATE incomingCyclePickupShifts
			SET staffId = $1, updated_at = $2
			WHERE id = $3
		`,
			payload.SourceStaffID,
			currentTime,
			pickupShift.ID,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, nil, err
			}
			return nil, nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	// Find pickedUp shifts for source staff based on CycleID StaffID Date ShiftName
	sourcePickupShifts, err = r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), int64(payload.SourceStaffID), payload.SourceDateAsDate, payload.SourceShiftName)
	if err != nil {
		return nil, nil, err
	}

	// Find pickedUp shifts for target staff based on CycleID StaffID Date ShiftName
	targetPickupShifts, err = r.FindVisitsForStaffInSpecificIncomingShift(int64(payload.CycleID), int64(payload.TargetStaffID), payload.TargetDateAsDate, payload.TargetShiftName)
	if err != nil {
		return nil, nil, err
	}

	return sourcePickupShifts, targetPickupShifts, nil
}

// makeCycleChatsWhereFilters generates a slice of SQL WHERE conditions based on the provided query parameters for cycle chats.
//
// It takes a pointer to models.CyclesQueryChatsRequestParams as an argument.
// Returns a slice of strings representing the WHERE conditions.
func makeCycleChatsWhereFilters(queries *models.CyclesQueryChatsRequestParams) []string {
	var where []string
	log.Printf("queries: %#v\n", queries)
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cc.id = %d ", queries.ID))
		}
		if queries.CycleID != 0 {
			where = append(where, fmt.Sprintf(" cps.cycleId = %d ", queries.CycleID))
		}
		if queries.CyclePickupShiftID != 0 {
			where = append(where, fmt.Sprintf(" cc.cyclePickupShiftId = %d ", queries.CyclePickupShiftID))
		}
		if queries.SenderUserID != 0 {
			where = append(where, fmt.Sprintf(" cc.senderUserId = %d ", queries.SenderUserID))
		}
		if queries.RecipientUserID != 0 {
			where = append(where, fmt.Sprintf(" cc.recipientUserId = %d ", queries.RecipientUserID))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cc.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryChats retrieves a list of cycle chats based on the provided query parameters.
//
// The `queries` parameter is a pointer to `models.CyclesQueryChatsRequestParams` that contains filters for the query.
// It returns a slice of pointers to `domain.CycleChat` and an error.
func (r *CycleRepositoryPostgresDB) QueryChats(queries *models.CyclesQueryChatsRequestParams) ([]*domain.CycleChat, error) {
	q := `
		SELECT
			cc.id,
			cc.cyclePickupShiftId,
			cc.senderUserId,
			cc.recipientUserId,
			cc.isSystem,
			cc.message,
			cc.attachments,
			cc.created_at,
			cc.updated_at,
			cc.deleted_at,
			u.id as senderUserId,
			u.firstName as senderUserFirstName,
			u.lastName as senderUserLastName,
			u.avatarUrl as senderUserAvatarUrl,
			u2.id as recipientUserId,
			u2.firstName as recipientUserFirstName,
			u2.lastName as recipientUserLastName,
			u2.avatarUrl as recipientUserAvatarUrl
		FROM cycleChats cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
		LEFT JOIN cyclePickupShifts cps ON cc.cyclePickupShiftId = cps.id
	`
	if queries != nil {
		where := makeCycleChatsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cc.id %s ", queries.Sorts.ID.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cc.created_at %s ", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit > -1 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"

	var chats []*domain.CycleChat
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chat                   domain.CycleChat
			message                sql.NullString
			attachments            json.RawMessage
			attachmentsMetadata    []*types.UploadMetadata
			deletedAt              sql.NullTime
			senderUserId           sql.NullInt64
			senderUserFirstName    sql.NullString
			senderUserLastName     sql.NullString
			senderUserAvatarUrl    sql.NullString
			recipientUserId        sql.NullInt64
			recipientUserFirstName sql.NullString
			recipientUserLastName  sql.NullString
			recipientUserAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&chat.ID,
			&chat.CyclePickupShiftID,
			&chat.SenderUserID,
			&chat.RecipientUserID,
			&chat.IsSystem,
			&message,
			&attachments,
			&chat.CreatedAt,
			&chat.UpdatedAt,
			&deletedAt,
			&senderUserId,
			&senderUserFirstName,
			&senderUserLastName,
			&senderUserAvatarUrl,
			&recipientUserId,
			&recipientUserFirstName,
			&recipientUserLastName,
			&recipientUserAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if message.Valid {
			chat.Message = &message.String
		}
		if deletedAt.Valid {
			chat.DeletedAt = &deletedAt.Time
		}
		if senderUserId.Valid {
			chat.SenderUser = &domain.CycleChatUser{
				ID: uint(senderUserId.Int64),
			}
			if senderUserFirstName.Valid {
				chat.SenderUser.FirstName = senderUserFirstName.String
			}
			if senderUserLastName.Valid {
				chat.SenderUser.LastName = senderUserLastName.String
			}
			if senderUserAvatarUrl.Valid {
				chat.SenderUser.AvatarUrl = senderUserAvatarUrl.String
			}
		}
		if recipientUserId.Valid {
			chat.RecipientUser = &domain.CycleChatUser{
				ID: uint(recipientUserId.Int64),
			}
			if recipientUserFirstName.Valid {
				chat.RecipientUser.FirstName = recipientUserFirstName.String
			}
			if recipientUserLastName.Valid {
				chat.RecipientUser.LastName = recipientUserLastName.String
			}
			if recipientUserAvatarUrl.Valid {
				chat.RecipientUser.AvatarUrl = recipientUserAvatarUrl.String
			}
		}
		if attachments != nil {
			err = json.Unmarshal(attachments, &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in cycle chats: %d", err, chat.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CYCLE_BUCKET_NAME[len("maja."):])
				}
			}
			chat.Attachments = attachmentsMetadata
		}
		chats = append(chats, &chat)
	}

	// Fetch cyclePickupShift
	for _, chat := range chats {
		cyclePickupShifts, err := r.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
			ID:    int(chat.CyclePickupShiftID),
			Page:  1,
			Limit: 1,
		})
		if err != nil {
			return nil, err
		}
		if len(cyclePickupShifts) > 0 {
			chat.CyclePickupShift = cyclePickupShifts[0]
		}
	}

	// Fetch last message for each cycleChat from chatMessages
	for _, chat := range chats {
		var (
			id          int64
			message     sql.NullString
			messageType string
			cycleChatId int64
		)
		err := r.PostgresDB.QueryRow(`
			SELECT id, cycleChatId, message, messageType FROM cycleChatMessages WHERE cycleChatId = $1
		`, chat.ID).Scan(&id, &cycleChatId, &message, &messageType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			log.Printf("failed to fetch last message: %v in cycle chat: %d", err, chat.ID)
			continue
		}
		_ = cycleChatId
		if message.Valid && messageType == constants.CHAT_MESSAGE_TYPE_TEXT {
			chat.Message = &message.String
		}
	}
	return chats, nil
}

// CountChats returns the number of cycle chats based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryChatsRequestParams as an argument.
// Returns the count of cycle chats as int64 and an error if any.
func (r *CycleRepositoryPostgresDB) CountChats(queries *models.CyclesQueryChatsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cc.id)
		FROM cycleChats cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
		LEFT JOIN cyclePickupShifts cps ON cc.cyclePickupShiftId = cps.id
	`
	if queries != nil {
		where := makeCycleChatsWhereFilters(queries)
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

func makeCycleChatMessagesWhereFilters(queries *models.CyclesQueryChatMessagesRequestParams) []string {
	var where []string
	log.Printf("queries: %#v\n", queries)
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cc.id = %d ", queries.ID))
		}
		if queries.CycleChatID != 0 {
			where = append(where, fmt.Sprintf(" cc.cycleChatId = %d ", queries.CycleChatID))
		}
		if queries.SenderUserID != 0 {
			where = append(where, fmt.Sprintf(" cc.senderUserId = %d ", queries.SenderUserID))
		}
		if queries.RecipientUserID != 0 {
			where = append(where, fmt.Sprintf(" cc.recipientUserId = %d ", queries.RecipientUserID))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cc.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

// QueryChatMessages retrieves cycle chat messages based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryChatMessagesRequestParams as an argument.
// Returns a slice of pointers to domain.CycleChatMessage and an error if any.
func (r *CycleRepositoryPostgresDB) QueryChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) ([]*domain.CycleChatMessage, error) {
	q := `
		SELECT
			cc.id,
			cc.cycleChatId,
			cc.senderUserId,
			cc.recipientUserId,
			cc.isSystem,
			cc.message,
			cc.messageType,
			cc.attachments,
			cc.created_at,
			cc.updated_at,
			cc.deleted_at,
			u.id AS senderUserId,
			u.firstName AS senderUserFirstName,
			u.lastName AS senderUserLastName,
			u.avatarUrl AS senderUserAvatarUrl,
			u2.id AS recipientUserId,
			u2.firstName AS recipientUserFirstName,
			u2.lastName AS recipientUserLastName,
			u2.avatarUrl AS recipientUserAvatarUrl
		FROM cycleChatMessages cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeCycleChatMessagesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cc.id %s ", queries.Sorts.ID.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cc.created_at %s ", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		if queries.Limit > -1 {
			limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
			queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
			offset := (queries.Page - 1) * limit
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"

	var chatMessages []*domain.CycleChatMessage
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chatMessage            domain.CycleChatMessage
			message                sql.NullString
			attachments            json.RawMessage
			attachmentsMetadata    []*types.UploadMetadata
			deletedAt              sql.NullTime
			senderUserId           sql.NullInt64
			senderUserFirstName    sql.NullString
			senderUserLastName     sql.NullString
			senderUserAvatarUrl    sql.NullString
			recipientUserId        sql.NullInt64
			recipientUserFirstName sql.NullString
			recipientUserLastName  sql.NullString
			recipientUserAvatarUrl sql.NullString
		)
		err = rows.Scan(
			&chatMessage.ID,
			&chatMessage.CycleChatID,
			&chatMessage.SenderUserID,
			&chatMessage.RecipientUserID,
			&chatMessage.IsSystem,
			&message,
			&chatMessage.MessageType,
			&attachments,
			&chatMessage.CreatedAt,
			&chatMessage.UpdatedAt,
			&deletedAt,
			&senderUserId,
			&senderUserFirstName,
			&senderUserLastName,
			&senderUserAvatarUrl,
			&recipientUserId,
			&recipientUserFirstName,
			&recipientUserLastName,
			&recipientUserAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if message.Valid {
			chatMessage.Message = &message.String
		}
		if deletedAt.Valid {
			chatMessage.DeletedAt = &deletedAt.Time
		}
		if senderUserId.Valid {
			chatMessage.SenderUser = &domain.CycleChatMessageUser{
				ID: uint(senderUserId.Int64),
			}
			if senderUserFirstName.Valid {
				chatMessage.SenderUser.FirstName = senderUserFirstName.String
			}
			if senderUserLastName.Valid {
				chatMessage.SenderUser.LastName = senderUserLastName.String
			}
			if senderUserAvatarUrl.Valid {
				chatMessage.SenderUser.AvatarUrl = senderUserAvatarUrl.String
			}
		}
		if recipientUserId.Valid {
			chatMessage.RecipientUser = &domain.CycleChatMessageUser{
				ID: uint(recipientUserId.Int64),
			}
			if recipientUserFirstName.Valid {
				chatMessage.RecipientUser.FirstName = recipientUserFirstName.String
			}
			if recipientUserLastName.Valid {
				chatMessage.RecipientUser.LastName = recipientUserLastName.String
			}
			if recipientUserAvatarUrl.Valid {
				chatMessage.RecipientUser.AvatarUrl = recipientUserAvatarUrl.String
			}
		}
		if attachments != nil {
			err = json.Unmarshal(attachments, &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in cycle chat messages: %d", err, chatMessage.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CYCLE_BUCKET_NAME[len("maja."):])
				}
			}
			chatMessage.Attachments = attachmentsMetadata
		}
		chatMessages = append(chatMessages, &chatMessage)
	}

	return chatMessages, nil
}

// CountChatMessages returns the number of cycle chat messages based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryChatMessagesRequestParams as an argument.
// Returns the count of cycle chat messages as int64 and an error if any.
func (r *CycleRepositoryPostgresDB) CountChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cc.id)
		FROM cycleChatMessages cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeCycleChatMessagesWhereFilters(queries)
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

// CreateChatMessage creates a new chat message in the database.
//
// It takes a CyclesCreateChatMessageRequestBody payload as input, which contains the cycle chat ID, sender user ID, recipient user ID, message, and attachments.
// Returns a pointer to a CycleChatMessage and an error.
func (r *CycleRepositoryPostgresDB) CreateChatMessage(payload *models.CyclesCreateChatMessageRequestBody) (*domain.CycleChatMessage, error) {
	// Get current time
	var (
		currentTime           = time.Now()
		isSystem              = false
		insertedChatMessageID int64
	)
	err := r.PostgresDB.QueryRow(`
		INSERT INTO cycleChatMessages (cycleChatId, senderUserId, recipientUserId, isSystem, message, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.CycleChatID,
		payload.SenderUserID,
		payload.RecipientUserID,
		isSystem,
		payload.Message,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedChatMessageID)
	if err != nil {
		return nil, err
	}

	// Get chat message
	chatMessages, err := r.QueryChatMessages(&models.CyclesQueryChatMessagesRequestParams{
		ID:    int(insertedChatMessageID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(chatMessages) == 0 {
		return nil, errors.New("failed to create chat message")
	}
	chatMessage := chatMessages[0]

	return chatMessage, nil
}

// UpdateChatMessageAttachments updates the attachments of a chat message.
//
// It takes in the previous attachments, new attachments, and the ID of the chat message.
// It returns the updated chat message and an error if any.
func (r *CycleRepositoryPostgresDB) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CycleChatMessage, error) {
	var chatMessage domain.CycleChatMessage

	// Current time
	currentTime := time.Now()

	// Find the medicine by id
	results, err := r.QueryChatMessages(&models.CyclesQueryChatMessagesRequestParams{
		ID:    int(id),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	foundChatMessage := results[0]
	if foundChatMessage == nil {
		return nil, errors.New("no results found")
	}

	// Marshal attachments into JSON format
	for _, attachment := range previousAttachments {
		// Check if the attachment already exists with fileName
		var exists bool
		for _, a := range attachments {
			if a.FileName == attachment.FileName {
				exists = true
				break
			}
		}
		if !exists {
			attachments = append(attachments, &attachment)
		}
	}
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)

	// Update the medicine
	err = r.PostgresDB.QueryRow(`
		UPDATE cycleChatMessages
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundChatMessage.ID,
	).Scan(
		&chatMessage.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine
	results, err = r.QueryChatMessages(&models.CyclesQueryChatMessagesRequestParams{
		ID:    int(foundChatMessage.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	chatMessage = *results[0]

	// Return the medicine
	return &chatMessage, nil
}
