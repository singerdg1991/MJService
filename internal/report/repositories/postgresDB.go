package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	customerDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	cycleModels "github.com/hoitek/Maja-Service/internal/cycle/models"
	crPorts "github.com/hoitek/Maja-Service/internal/cycle/ports"
	"github.com/hoitek/Maja-Service/internal/report/domain"
	"github.com/hoitek/Maja-Service/internal/report/models"
	"github.com/lib/pq"
)

type ReportRepositoryPostgresDB struct {
	PostgresDB      *sql.DB
	CycleRepository crPorts.CycleRepositoryPostgresDB
}

func NewReportRepositoryPostgresDB(d *sql.DB, cycleRepository crPorts.CycleRepositoryPostgresDB) *ReportRepositoryPostgresDB {
	return &ReportRepositoryPostgresDB{
		PostgresDB:      d,
		CycleRepository: cycleRepository,
	}
}

func QueryArrangementsTableWhereFilters(q *models.ReportsQueryArrangementsTableRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add ID filter
		if q.ID > 0 {
			where = append(where, fmt.Sprintf(" cps.id = $%d", paramCount))
			args = append(args, q.ID)
			paramCount++
		}

		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			// query += fmt.Sprintf(" cps.id = ANY($%d)", paramCount)
			// args = append(args, q.FilterIDsAsInt64Slice)
			// paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					// Convert string value to integer for numeric fields
					if strings.HasSuffix(agg.Field, ".id") || strings.HasSuffix(agg.Field, "Hour") || strings.HasSuffix(agg.Field, "Kilometer") {
						where = append(where, fmt.Sprintf(" CAST(%s AS INTEGER) < $%d", agg.Field, paramCount))
					} else {
						where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					}
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) CountArrangementsTable(queries *models.ReportsQueryArrangementsTableRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cps.id) as total_count
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		LEFT JOIN cycleStaffTypes cst ON cps.cycleStaffTypeId = cst.id
	`
	var args []interface{}
	if queries != nil {
		where, argsValues, _ := QueryArrangementsTableWhereFilters(queries)
		if len(where) > 0 {
			args = append(args, argsValues...)
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepositoryPostgresDB) QueryArrangementsTable(q *models.ReportsQueryArrangementsTableRequestParams) ([]*domain.ReportArrangementTable, error) {
	query := `
		SELECT 
			cps.id,
			c.sectionId as sectionId,
			cps.cycleId,
			s.id as staffId,
			s.userId as staffUserId,
			u.firstName as staffFirstName,
			u.lastName as staffLastName,
			u.avatarUrl as staffAvatarUrl,
			cst.id as shiftId,
			cst.shiftName,
			cst.startHour,
			cst.endHour,
			cst.isUnplanned as shiftIsUnplanned,
			cst.datetime as shiftDatetime,
			cps.status,
			cps.prevStatus,
			cps.startKilometer,
			cps.reasonOfTheCancellation,
			cps.reasonOfTheReactivation,
			cps.reasonOfTheResume,
			cps.reasonOfThePause,
			cps.isUnplanned,
			cps.datetime,
			cps.created_at,
			cps.updated_at,
			cps.deleted_at,
			cps.started_at,
			cps.ended_at,
			cps.cancelled_at,
			cps.delayed_at
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		LEFT JOIN cycleStaffTypes cst ON cps.cycleStaffTypeId = cst.id
	`

	// Add filters
	where, args, _ := QueryArrangementsTableWhereFilters(q)
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Add sorting
	query += " ORDER BY cps.datetime DESC"

	// Add pagination
	limit := exp.TerIf(q.Limit == 0, 10, q.Limit)
	q.Page = exp.TerIf(q.Page == 0, 1, q.Page)
	offset := (q.Page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	log.Printf("Query: %s", query)
	log.Printf("Args: %v", args)

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.ReportArrangementTable

	for rows.Next() {
		var item domain.ReportArrangementTable
		var staffID, staffUserID, staffFirstName, staffLastName, staffAvatarURL sql.NullString
		var shiftID sql.NullInt64
		var shiftName, startHour, endHour, shiftDateTime sql.NullString
		var shiftIsUnplanned sql.NullBool
		var startKilometer, reasonOfTheCancellation, reasonOfTheReactivation, reasonOfTheResume, reasonOfThePause sql.NullString
		var startedAt, endedAt, cancelledAt, delayedAt, deletedAt sql.NullTime
		var sectionId sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&sectionId,
			&item.CycleID,
			&staffID,
			&staffUserID,
			&staffFirstName,
			&staffLastName,
			&staffAvatarURL,
			&shiftID,
			&shiftName,
			&startHour,
			&endHour,
			&shiftIsUnplanned,
			&shiftDateTime,
			&item.Status,
			&item.PrevStatus,
			&startKilometer,
			&reasonOfTheCancellation,
			&reasonOfTheReactivation,
			&reasonOfTheResume,
			&reasonOfThePause,
			&item.IsUnplanned,
			&item.DateTime,
			&item.CreatedAt,
			&item.UpdatedAt,
			&deletedAt,
			&startedAt,
			&endedAt,
			&cancelledAt,
			&delayedAt,
		)
		if err != nil {
			return nil, err
		}

		if staffID.Valid {
			userID, _ := strconv.ParseUint(staffUserID.String, 10, 64)
			item.Staff = &domain.ReportArrangementTableStaff{
				ID:        uint(userID),
				UserID:    uint(userID),
				FirstName: staffFirstName.String,
				LastName:  staffLastName.String,
				AvatarUrl: staffAvatarURL.String,
			}
		}

		if shiftID.Valid {
			startHourTime, _ := time.Parse("15:04", startHour.String)
			endHourTime, _ := time.Parse("15:04", endHour.String)
			shiftDateTimeTime, _ := time.Parse(time.RFC3339, shiftDateTime.String)

			item.Shift = &domain.ReportArrangementTableShift{
				ID:          uint(shiftID.Int64),
				ShiftName:   shiftName.String,
				StartHour:   startHourTime,
				EndHour:     endHourTime,
				IsUnplanned: shiftIsUnplanned.Bool,
				DateTime:    shiftDateTimeTime,
			}
		}

		if startKilometer.Valid {
			item.StartKilometer = &startKilometer.String
		}

		if reasonOfTheCancellation.Valid {
			item.ReasonOfTheCancellation = &reasonOfTheCancellation.String
		}

		if reasonOfTheReactivation.Valid {
			item.ReasonOfTheReactivation = &reasonOfTheReactivation.String
		}

		if reasonOfTheResume.Valid {
			item.ReasonOfTheResume = &reasonOfTheResume.String
		}

		if reasonOfThePause.Valid {
			item.ReasonOfThePause = &reasonOfThePause.String
		}

		if startedAt.Valid {
			item.StartedAt = &startedAt.Time
		}

		if endedAt.Valid {
			item.EndedAt = &endedAt.Time
		}

		if cancelledAt.Valid {
			item.CancelledAt = &cancelledAt.Time
		}

		if delayedAt.Valid {
			item.DelayedAt = &delayedAt.Time
		}

		if deletedAt.Valid {
			item.DeletedAt = &deletedAt.Time
		}

		items = append(items, &item)
	}

	return items, nil
}

func QueryShiftsTableWhereFilters(q *models.ReportsQueryShiftsTableRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		if q.ID > 0 {
			where = append(where, fmt.Sprintf(" s.id = $%d", paramCount))
			args = append(args, q.ID)
			paramCount++
		}

		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) CountShiftsTable(queries *models.ReportsQueryShiftsTableRequestParams) (int64, error) {
	q := `
		SELECT 
			COUNT(cs.id) as total_count
		FROM cycleShifts cs
	`
	var args []interface{}
	if queries != nil {
		where, argsValues, _ := QueryShiftsTableWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
			args = append(args, argsValues...)
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting shifts: %v", err)
	}

	return count, nil
}

func (r *ReportRepositoryPostgresDB) QueryShiftsTable(q *models.ReportsQueryShiftsTableRequestParams) ([]*domain.ReportShiftTable, error) {
	query := `
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

	var args []interface{}
	if q != nil {
		where, argsValues, _ := QueryShiftsTableWhereFilters(q)
		if len(where) > 0 {
			query += " WHERE " + strings.Join(where, " AND ")
			args = append(args, argsValues...)
		}

		// Add sorting
		query += " ORDER BY cs.datetime DESC"

		// Add pagination
		limit := exp.TerIf(q.Limit == 0, 10, q.Limit)
		q.Page = exp.TerIf(q.Page == 0, 1, q.Page)
		offset := (q.Page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying shifts: %v", err)
	}
	defer rows.Close()

	var items []*domain.ReportShiftTable
	for rows.Next() {
		item := &domain.ReportShiftTable{}
		var (
			staffTypeIDs         sql.NullString
			staffTypeIDsMetadata []uint
			vehicleType          sql.NullString
			startLocation        sql.NullString
		)
		err := rows.Scan(
			&item.ID,
			&item.ExchangeKey,
			&item.CycleID,
			&staffTypeIDs,
			&item.ShiftName,
			&vehicleType,
			&startLocation,
			&item.DateTime,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
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
				item.StaffTypeIDs = staffTypeIDsMetadata
			}
		}
		if len(item.StaffTypeIDs) > 0 {
			for _, staffTypeID := range item.StaffTypeIDs {
				cycleStaffType, err := r.CycleRepository.QueryStaffTypes(&cycleModels.CyclesQueryStaffTypesRequestParams{
					ID: int(staffTypeID),
				})
				if err != nil {
					log.Printf("Error while querying staff type in query cycle shift: %v\n", err.Error())
					return nil, err
				}
				if len(cycleStaffType) > 0 {
					item.StaffTypes = append(item.StaffTypes, cycleStaffType[0])
				}
			}
		}
		if vehicleType.Valid {
			item.VehicleType = &vehicleType.String
		}
		if startLocation.Valid {
			item.StartLocation = &startLocation.String
		}
		items = append(items, item)
	}

	return items, nil
}

func QueryVisitsTableWhereFilters(q *models.ReportsQueryVisitsTableRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		if q.ID > 0 {
			where = append(where, fmt.Sprintf(" v.id = $%d", paramCount))
			args = append(args, q.ID)
			paramCount++
		}

		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) CountVisitsTable(queries *models.ReportsQueryVisitsTableRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cps.id) as total_count
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
	`
	var args []interface{}
	if queries != nil {
		where, argsValues, _ := QueryVisitsTableWhereFilters(queries)
		if len(where) > 0 {
			args = append(args, argsValues...)
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepositoryPostgresDB) QueryVisitsTable(q *models.ReportsQueryVisitsTableRequestParams) ([]*domain.ReportVisitTable, error) {
	query := `
		SELECT
			cps.id,
			cps.cycleId,
			s.id as staffId,
			s.userId as staffUserId,
			u.firstName as staffFirstName,
			u.lastName as staffLastName,
			u.avatarUrl as staffAvatarUrl,
			cps.status,
			cps.prevStatus,
			cps.startKilometer,
			cps.reasonOfTheCancellation,
			cps.reasonOfTheReactivation,
			cps.reasonOfTheResume,
			cps.reasonOfThePause,
			cps.isUnplanned,
			cps.datetime,
			cps.created_at,
			cps.updated_at,
			cps.deleted_at,
			cps.started_at,
			cps.ended_at,
			cps.cancelled_at,
			cps.delayed_at
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
	`

	// Add filters
	where, args, _ := QueryVisitsTableWhereFilters(q)
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Add sorting
	query += " ORDER BY cps.datetime DESC"

	// Add pagination
	limit := exp.TerIf(q.Limit == 0, 10, q.Limit)
	q.Page = exp.TerIf(q.Page == 0, 1, q.Page)
	offset := (q.Page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	log.Printf("Visit's Query: %s\n", query)

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.ReportVisitTable

	for rows.Next() {
		var item domain.ReportVisitTable
		var staffID, staffUserID, staffFirstName, staffLastName, staffAvatarURL sql.NullString
		var startKilometer, reasonOfTheCancellation, reasonOfTheReactivation, reasonOfTheResume, reasonOfThePause sql.NullString
		var startedAt, endedAt, cancelledAt, delayedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&item.ID,
			&item.CycleID,
			&staffID,
			&staffUserID,
			&staffFirstName,
			&staffLastName,
			&staffAvatarURL,
			&item.Status,
			&item.PrevStatus,
			&startKilometer,
			&reasonOfTheCancellation,
			&reasonOfTheReactivation,
			&reasonOfTheResume,
			&reasonOfThePause,
			&item.IsUnplanned,
			&item.DateTime,
			&item.CreatedAt,
			&item.UpdatedAt,
			&deletedAt,
			&startedAt,
			&endedAt,
			&cancelledAt,
			&delayedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if startKilometer.Valid {
			item.StartKilometer = &startKilometer.String
		}
		if reasonOfTheCancellation.Valid {
			item.ReasonOfTheCancellation = &reasonOfTheCancellation.String
		}
		if reasonOfTheReactivation.Valid {
			item.ReasonOfTheReactivation = &reasonOfTheReactivation.String
		}
		if reasonOfTheResume.Valid {
			item.ReasonOfTheResume = &reasonOfTheResume.String
		}
		if reasonOfThePause.Valid {
			item.ReasonOfThePause = &reasonOfThePause.String
		}
		if deletedAt.Valid {
			item.DeletedAt = &deletedAt.Time
		}
		if startedAt.Valid {
			item.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			item.EndedAt = &endedAt.Time
		}
		if cancelledAt.Valid {
			item.CancelledAt = &cancelledAt.Time
		}
		if delayedAt.Valid {
			item.DelayedAt = &delayedAt.Time
		}

		// Handle staff information
		if staffID.Valid {
			item.Staff = &domain.ReportVisitTableStaff{
				ID:        uint(stringToUint(staffID.String)),
				UserID:    uint(stringToUint(staffUserID.String)),
				FirstName: staffFirstName.String,
				LastName:  staffLastName.String,
				AvatarUrl: staffAvatarURL.String,
			}
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *ReportRepositoryPostgresDB) CountCustomersTable(queries *models.ReportsQueryCustomersTableRequestParams) (int64, error) {
	q := `
		SELECT 
			COUNT(c.id) AS total_count
		FROM 
			customers c
		LEFT JOIN 
			users u1 ON c.userId = u1.id
		LEFT JOIN 
			staffs s ON c.responsibleNurseId = s.id
		LEFT JOIN 
			users u2 ON s.userId = u2.id
	`
	var args []interface{}
	if queries != nil {
		where, argsValues, _ := QueryCustomersTableWhereFilters(queries)
		if len(where) > 0 {
			args = append(args, argsValues...)
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepositoryPostgresDB) QueryCustomersTable(q *models.ReportsQueryCustomersTableRequestParams) ([]*domain.ReportCustomerTable, error) {
	query := `
		SELECT 
			c.*,
			u1.firstName AS userFirstName,
			u1.lastName AS userLastName,
			u1.avatarUrl AS userAvatarUrl,
			u1.gender AS userGender,
			u1.email AS userEmail,
			u1.phone AS userPhone,
			u1.birthDate AS userBirthDate,
			u1.nationalCode AS userNationalCode,
			s.id AS rNId,
			u2.firstName AS rNFirstName,
			u2.lastName AS rNLastName,
			u2.avatarUrl AS rNAvatarUrl
		FROM 
			customers c
		LEFT JOIN 
			users u1 ON c.userId = u1.id
		LEFT JOIN 
			staffs s ON c.responsibleNurseId = s.id
		LEFT JOIN 
			users u2 ON s.userId = u2.id
	`

	// Add filters
	where, args, _ := QueryCustomersTableWhereFilters(q)
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Add sorting
	query += " ORDER BY c.created_at DESC"

	// Add pagination
	limit := exp.TerIf(q.Limit == 0, 10, q.Limit)
	q.Page = exp.TerIf(q.Page == 0, 1, q.Page)
	offset := (q.Page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	log.Printf("Q: %s\n", query)

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*domain.ReportCustomerTable
	for rows.Next() {
		var (
			customer                                     domain.ReportCustomerTable
			userID                                       sql.NullInt64
			responsibleNurseID                           sql.NullInt64
			motherLangIDs                                json.RawMessage
			parkingInfo                                  sql.NullString
			extraExplanation                             sql.NullString
			limitingTheRightToSelfDeterminationStartDate sql.NullTime
			limitingTheRightToSelfDeterminationEndDate   sql.NullTime
			mobilityContract                             sql.NullString
			keyNo                                        sql.NullString
			paymentMethod                                sql.NullString
			userFirstName                                sql.NullString
			userLastName                                 sql.NullString
			userAvatarUrl                                sql.NullString
			userGender                                   sql.NullString
			userEmail                                    sql.NullString
			userPhone                                    sql.NullString
			userBirthDate                                sql.NullTime
			userNationalCode                             sql.NullString
			rNId                                         sql.NullInt64
			rNFirstName                                  sql.NullString
			rNLastName                                   sql.NullString
			rNAvatarUrl                                  sql.NullString
		)
		err := rows.Scan(
			&customer.ID,
			&userID,
			&responsibleNurseID,
			&motherLangIDs,
			&customer.NurseGenderWish,
			&customer.Status,
			&customer.StatusDate,
			&parkingInfo,
			&extraExplanation,
			&customer.HasLimitingTheRightToSelfDetermination,
			&limitingTheRightToSelfDeterminationStartDate,
			&limitingTheRightToSelfDeterminationEndDate,
			&mobilityContract,
			&keyNo,
			&paymentMethod,
			&customer.CreatedAt,
			&customer.UpdatedAt,
			&customer.DeletedAt,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&userGender,
			&userEmail,
			&userPhone,
			&userBirthDate,
			&userNationalCode,
			&rNId,
			&rNFirstName,
			&rNLastName,
			&rNAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if userID.Valid {
			var customerUser customerDomain.CustomerUser
			customer.UserID = &userID.Int64
			if userFirstName.Valid {
				customerUser.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				customerUser.LastName = userLastName.String
			}
			if userGender.Valid {
				customerUser.Gender = userGender.String
			}
			if userAvatarUrl.Valid {
				customerUser.AvatarUrl = userAvatarUrl.String
			}
			if userEmail.Valid {
				customerUser.Email = userEmail.String
			}
			if userPhone.Valid {
				customerUser.Phone = userPhone.String
			}
			if userBirthDate.Valid {
				customerUser.BirthDate = &userBirthDate.Time
			}
			if userNationalCode.Valid {
				customerUser.NationalCode = userNationalCode.String
			}
			customerUser.ID = *customer.UserID
			customer.User = &customerUser
		} else {
			customer.UserID = nil
		}
		if responsibleNurseID.Valid {
			customer.ResponsibleNurseID = &responsibleNurseID.Int64
		}
		if motherLangIDs != nil {
			var motherLangIDsArray []int64
			err := json.Unmarshal(motherLangIDs, &motherLangIDsArray)
			if err != nil {
				return nil, err
			}
			customer.MotherLangIDs = motherLangIDsArray
		}
		if parkingInfo.Valid {
			customer.ParkingInfo = &parkingInfo.String
		}
		if extraExplanation.Valid {
			customer.ExtraExplanation = &extraExplanation.String
		}
		if limitingTheRightToSelfDeterminationStartDate.Valid {
			customer.LimitingTheRightToSelfDeterminationStartDate = &limitingTheRightToSelfDeterminationStartDate.Time
		}
		if limitingTheRightToSelfDeterminationEndDate.Valid {
			customer.LimitingTheRightToSelfDeterminationEndDate = &limitingTheRightToSelfDeterminationEndDate.Time
		}
		if mobilityContract.Valid {
			customer.MobilityContract = &mobilityContract.String
		}
		if keyNo.Valid {
			customer.KeyNo = &keyNo.String
		}
		if paymentMethod.Valid {
			customer.PaymentMethod = &paymentMethod.String
		}
		if rNId.Valid {
			customer.ResponsibleNurse = &customerDomain.CustomerResponsibleNurse{
				ID: uint(rNId.Int64),
			}
			if rNFirstName.Valid {
				customer.ResponsibleNurse.FirstName = rNFirstName.String
			}
			if rNLastName.Valid {
				customer.ResponsibleNurse.LastName = rNLastName.String
			}
			if rNAvatarUrl.Valid {
				customer.ResponsibleNurse.AvatarUrl = rNAvatarUrl.String
			}
		}

		// Get customersRelatives
		q := `
			SELECT
			    cr.id,
				cr.firstName,
				cr.lastName,
				cr.relation
			FROM customersRelatives crs
			LEFT JOIN customerRelatives cr ON crs.relativeId = cr.id
			WHERE crs.customerId = $1
		`
		rows, err := r.PostgresDB.Query(q, customer.ID)
		if err != nil {
			return nil, err
		}
		var (
			relativeIDs        []int64
			customersRelatives []customerDomain.CustomersRelative
		)
		for rows.Next() {
			var (
				customersRelative customerDomain.CustomersRelative
			)
			err := rows.Scan(
				&customersRelative.ID,
				&customersRelative.FirstName,
				&customersRelative.LastName,
				&customersRelative.Relation,
			)
			if err != nil {
				return nil, err
			}
			relativeIDs = append(relativeIDs, int64(customersRelative.ID))
			customersRelatives = append(customersRelatives, customersRelative)
		}
		if rows != nil {
			rows.Close()
		}
		customer.RelativeIDs = relativeIDs
		customer.Relatives = customersRelatives

		// Get customersDiagnoses
		q = `
			SELECT
			    cds.id,
				cds.customerId,
				cds.diagnoseId,
				d.id AS dId,
				d.title AS dTitle
			FROM customersDiagnoses cds
			LEFT JOIN diagnoses d ON cds.diagnoseId = d.id
			WHERE cds.customerId = $1
		`
		rows, err = r.PostgresDB.Query(q, customer.ID)
		if err != nil {
			return nil, err
		}
		var (
			diagnoseIDs       []int64
			customerDiagnoses []customerDomain.CustomerDiagnose
		)
		for rows.Next() {
			var (
				customerDiagnose customerDomain.CustomerDiagnose
				diagnoseID       sql.NullInt64
				diagnoseTitle    sql.NullString
			)
			err := rows.Scan(
				&customerDiagnose.ID,
				&customerDiagnose.CustomerID,
				&customerDiagnose.DiagnoseID,
				&diagnoseID,
				&diagnoseTitle,
			)
			if err != nil {
				return nil, err
			}
			if diagnoseID.Valid {
				customerDiagnose.Diagnose = &customerDomain.CustomerDiagnoseDiagnose{
					ID: uint(diagnoseID.Int64),
				}
				if diagnoseTitle.Valid {
					customerDiagnose.Diagnose.Title = diagnoseTitle.String
				}
			}
			diagnoseIDs = append(diagnoseIDs, int64(customerDiagnose.DiagnoseID))
			customerDiagnoses = append(customerDiagnoses, customerDiagnose)
		}
		if rows != nil {
			rows.Close()
		}
		customer.DiagnoseIDs = diagnoseIDs
		customer.Diagnoses = customerDiagnoses
		customers = append(customers, &customer)
	}

	// Get sections for each customer
	for _, customer := range customers {
		q := `
			SELECT
			    sections.*
			FROM customerSections
			LEFT JOIN sections ON customerSections.sectionId = sections.id
			WHERE customerSections.customerId = $1
		`
		rows, err := r.PostgresDB.Query(q, customer.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var (
				section     customerDomain.CustomerSection
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
					parent           customerDomain.CustomerSection
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
					child       customerDomain.CustomerSection
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
						parent           customerDomain.CustomerSection
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
			if rows != nil {
				rows.Close()
			}
			if customer.Sections == nil {
				customer.Sections = []customerDomain.CustomerSection{}
			}
			customer.Sections = append(customer.Sections, section)
		}
	}

	// Get limitations for each customer
	for _, customer := range customers {
		q := `
			SELECT
			    customerLimitations.id AS customerLimitationId,
			    customerLimitations.customerId AS customerLimitationCustomerId,
			    customerLimitations.description AS customerLimitationDescription,
			    limitations.id AS limitationId,
			    limitations.name AS limitationName,
			    limitations.description AS limitationDescription
			FROM customerLimitations
			LEFT JOIN limitations ON customerLimitations.limitationId = limitations.id
			WHERE customerLimitations.customerId = $1
		`
		rows, err := r.PostgresDB.Query(q, customer.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var (
				customerLimitation    customerDomain.CustomerLimitation
				limitationName        string
				description           sql.NullString
				limitationDescription sql.NullString
			)
			err := rows.Scan(
				&customerLimitation.ID,
				&customerLimitation.CustomerID,
				&description,
				&customerLimitation.LimitationID,
				&limitationName,
				&limitationDescription,
			)
			if err != nil {
				return nil, err
			}
			if description.Valid {
				customerLimitation.Description = &description.String
			}
			if limitationDescription.Valid {
				customerLimitation.Limitation.Description = &limitationDescription.String
			}
			customerLimitation.Limitation.ID = customerLimitation.LimitationID
			customerLimitation.Limitation.Name = limitationName
			customer.Limitations = append(customer.Limitations, customerLimitation)
		}
	}

	// Get mother languages for each customer
	for _, customer := range customers {
		q := `
			SELECT
				ls.id, ls.name
			FROM languageskills ls
			WHERE ls.id = ANY ($1)
		`
		rows, err := r.PostgresDB.Query(q, pq.Int64Array(customer.MotherLangIDs))
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var customerMotherLanguage customerDomain.CustomerMotherLang
			err := rows.Scan(
				&customerMotherLanguage.ID,
				&customerMotherLanguage.Name,
			)
			if err != nil {
				return nil, err
			}
			customer.MotherLangs = append(customer.MotherLangs, customerMotherLanguage)
		}
	}

	return customers, nil
}

func QueryCustomersTableWhereFilters(q *models.ReportsQueryCustomersTableRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		if q.ID > 0 {
			where = append(where, fmt.Sprintf(" c.id = $%d", paramCount))
			args = append(args, q.ID)
			paramCount++
		}

		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

// GetShiftSchedulingChart returns data for the shift scheduling range bar chart
func (r *ReportRepositoryPostgresDB) GetShiftSchedulingChart(q *models.ReportsQueryArrangementsTableRequestParams) ([]domain.ShiftSchedulingChartData, error) {
	query := `
		SELECT
			s.id as staffId,
			u.firstName as staffFirstName,
			u.lastName as staffLastName,
			cst.shiftName,
			cst.startHour,
			cst.endHour,
			cst.datetime,
			cps.status,
			cps.isUnplanned
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cycleStaffTypes cst ON cps.cycleStaffTypeId = cst.id
	`

	// Add filters
	where, args, _ := QueryArrangementsTableWhereFilters(q)
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Order by staff name and shift start time
	query += " ORDER BY u.firstName, u.lastName, cst.startHour"

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chartData []domain.ShiftSchedulingChartData

	for rows.Next() {
		var (
			data                                     domain.ShiftSchedulingChartData
			staffID                                  sql.NullInt64
			staffFirstName, staffLastName, shiftName sql.NullString
			startHour, endHour, shiftDate            sql.NullTime
			status                                   sql.NullString
			isUnplanned                              sql.NullBool
		)

		err := rows.Scan(
			&staffID,
			&staffFirstName,
			&staffLastName,
			&shiftName,
			&startHour,
			&endHour,
			&shiftDate,
			&status,
			&isUnplanned,
		)
		if err != nil {
			return nil, err
		}

		// Format staff name
		var staffName string
		if staffFirstName.Valid && staffLastName.Valid {
			staffName = staffFirstName.String + " " + staffLastName.String
		}

		// Create chart data entry
		data = domain.ShiftSchedulingChartData{
			StaffID:     uint(staffID.Int64),
			StaffName:   staffName,
			ShiftName:   shiftName.String,
			StartHour:   startHour.Time,
			EndHour:     endHour.Time,
			Date:        shiftDate.Time,
			Status:      status.String,
			IsUnplanned: isUnplanned.Bool,
		}

		chartData = append(chartData, data)
	}

	return chartData, nil
}

func QueryShiftDurationByCustomersWhereFilters(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			// query += fmt.Sprintf(" cps.id = ANY($%d)", paramCount)
			// args = append(args, q.FilterIDsAsInt64Slice)
			// paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					// Convert string value to integer for numeric fields
					if strings.HasSuffix(agg.Field, ".id") || strings.HasSuffix(agg.Field, "Hour") || strings.HasSuffix(agg.Field, "Kilometer") {
						where = append(where, fmt.Sprintf(" CAST(%s AS INTEGER) < $%d", agg.Field, paramCount))
					} else {
						where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					}
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

// GetShiftDistributionByCustomer returns data for the shift distribution pie chart
func (r *ReportRepositoryPostgresDB) GetShiftDistributionByCustomer(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams, period string) ([]domain.ShiftDistributionByCustomer, error) {
	query := `
		WITH customer_shifts AS (
			SELECT
				cu.id as customerId,
				u2.firstName || ' ' || u2.lastName as customerName,
				COUNT(DISTINCT cps.id) as shift_count
			FROM cyclePickupShifts cps
			LEFT JOIN cyclePickupShiftCustomers cpsc ON cps.id = cpsc.cyclePickupShiftId
			LEFT JOIN customerServices cs ON cpsc.customerServiceId = cs.id
			LEFT JOIN customers cu ON cs.customerId = cu.id
			LEFT JOIN users u2 ON cu.userId = u2.id
			LEFT JOIN cycles c ON cps.cycleId = c.id
			LEFT JOIN staffs s ON cps.staffId = s.id
			LEFT JOIN users u ON s.userId = u.id
			WHERE cu.id IS NOT NULL
	`

	// Add time period filter
	switch period {
	case "daily":
		query += " AND DATE(cps.datetime) = CURRENT_DATE"
	case "weekly":
		query += " AND DATE(cps.datetime) BETWEEN CURRENT_DATE - INTERVAL '7 days' AND CURRENT_DATE"
	}

	// Add other filters
	where, args, _ := QueryShiftDurationByCustomersWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Complete the CTE and calculate percentages
	query += `
			GROUP BY cu.id, u2.firstName, u2.lastName
		)
		SELECT
			cs.customerId,
			cs.customerName,
			cs.shift_count as numberOfShifts,
			ROUND((cs.shift_count::numeric * 100.0) / NULLIF((SELECT SUM(shift_count) FROM customer_shifts), 0), 2) as percentage
		FROM customer_shifts cs
		ORDER BY cs.shift_count DESC
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributionData []domain.ShiftDistributionByCustomer

	for rows.Next() {
		var data domain.ShiftDistributionByCustomer
		var customerID sql.NullInt64
		var customerName sql.NullString
		var numberOfShifts sql.NullInt64
		var percentage sql.NullFloat64

		err := rows.Scan(
			&customerID,
			&customerName,
			&numberOfShifts,
			&percentage,
		)
		if err != nil {
			return nil, err
		}

		if customerID.Valid && customerName.Valid && numberOfShifts.Valid && percentage.Valid {
			data = domain.ShiftDistributionByCustomer{
				CustomerID:     uint(customerID.Int64),
				CustomerName:   customerName.String,
				NumberOfShifts: int(numberOfShifts.Int64),
				Percentage:     percentage.Float64,
			}
			distributionData = append(distributionData, data)
		}
	}

	return distributionData, nil
}

func QueryShiftDurationAnalysisWhereFilters(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			// query += fmt.Sprintf(" cps.id = ANY($%d)", paramCount)
			// args = append(args, q.FilterIDsAsInt64Slice)
			// paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					// Convert string value to integer for numeric fields
					if strings.HasSuffix(agg.Field, ".id") || strings.HasSuffix(agg.Field, "Hour") || strings.HasSuffix(agg.Field, "Kilometer") {
						where = append(where, fmt.Sprintf(" CAST(%s AS INTEGER) < $%d", agg.Field, paramCount))
					} else {
						where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					}
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

// GetShiftDurationAnalysis returns data for the shift duration analysis scatter chart
func (r *ReportRepositoryPostgresDB) GetShiftDurationAnalysis(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) ([]domain.ShiftDurationAnalysis, error) {
	query := `
		SELECT
			c.id as customerId,
			u2.firstName as customerFirstName,
			u2.lastName as customerLastName,
			CASE 
				WHEN cst.endHour < cst.startHour THEN
					EXTRACT(EPOCH FROM (cst.endHour::time - cst.startHour::time + INTERVAL '24 hours'))/3600
				ELSE
					EXTRACT(EPOCH FROM (cst.endHour::time - cst.startHour::time))/3600
			END as durationHours,
			COUNT(*) as numberOfShifts
		FROM cyclePickupShifts cps
		LEFT JOIN cycleStaffTypes cst ON cps.cycleStaffTypeId = cst.id
		LEFT JOIN cyclePickupShiftCustomers cpsc ON cps.id = cpsc.cyclePickupShiftId
		LEFT JOIN customerServices cs ON cpsc.customerServiceId = cs.id
		LEFT JOIN customers cu ON cs.customerId = cu.id
		LEFT JOIN users u2 ON cu.userId = u2.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		WHERE cu.id IS NOT NULL
	`

	// Add filters
	where, args, _ := QueryShiftDurationAnalysisWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Group by customer and duration
	query += " GROUP BY c.id, cu.id, u2.firstName, u2.lastName, durationHours"

	// Order by customer name and duration
	query += " ORDER BY u2.firstName, u2.lastName, durationHours"

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analysisData []domain.ShiftDurationAnalysis

	for rows.Next() {
		var data domain.ShiftDurationAnalysis
		var customerID sql.NullInt64
		var customerFirstName sql.NullString
		var customerLastName sql.NullString
		var durationHours sql.NullFloat64
		var numberOfShifts sql.NullInt64

		err := rows.Scan(
			&customerID,
			&customerFirstName,
			&customerLastName,
			&durationHours,
			&numberOfShifts,
		)
		if err != nil {
			return nil, err
		}

		if customerID.Valid && customerFirstName.Valid && customerLastName.Valid && durationHours.Valid && numberOfShifts.Valid {
			data = domain.ShiftDurationAnalysis{
				CustomerID:     uint(customerID.Int64),
				CustomerName:   customerFirstName.String + " " + customerLastName.String,
				DurationHours:  durationHours.Float64,
				NumberOfShifts: int(numberOfShifts.Int64),
			}
			analysisData = append(analysisData, data)
		}
	}

	return analysisData, nil
}

func QueryShiftsByStaffAndCustomerWhereFilters(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			// query += fmt.Sprintf(" cps.id = ANY($%d)", paramCount)
			// args = append(args, q.FilterIDsAsInt64Slice)
			// paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					if strings.HasSuffix(agg.Field, ".id") || strings.HasSuffix(agg.Field, "Hour") || strings.HasSuffix(agg.Field, "Kilometer") {
						where = append(where, fmt.Sprintf(" CAST(%s AS INTEGER) < $%d", agg.Field, paramCount))
					} else {
						where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					}
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

// GetShiftsByStaffAndCustomer returns data for the staff-customer shifts stacked bar chart
func (r *ReportRepositoryPostgresDB) GetShiftsByStaffAndCustomer(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) ([]domain.ShiftsByStaffAndCustomer, error) {
	query := `
		SELECT 
			s.id as staffId,
			u.firstName || ' ' || u.lastName as staffName,
			cu.id as customerId,
			u2.firstName || ' ' || u2.lastName as customerName,
			COUNT(DISTINCT cps.id) as shiftCount
		FROM cyclePickupShifts cps
		LEFT JOIN staffs s ON cps.staffId = s.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN cyclePickupShiftCustomers cpsc ON cps.id = cpsc.cyclePickupShiftId
		LEFT JOIN customers cu ON cpsc.customerId = cu.id
		LEFT JOIN users u2 ON cu.userId = u2.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		WHERE s.id IS NOT NULL AND cu.id IS NOT NULL
	`

	// Add filters
	where, args, _ := QueryShiftsByStaffAndCustomerWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Group by staff and customer
	query += " GROUP BY s.id, u.firstName, u.lastName, cu.id, u2.firstName, u2.lastName"

	// Order by staff name and customer name
	query += " ORDER BY staffName, customerName"

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.ShiftsByStaffAndCustomer
	for rows.Next() {
		var item domain.ShiftsByStaffAndCustomer
		err := rows.Scan(
			&item.StaffId,
			&item.StaffName,
			&item.CustomerId,
			&item.CustomerName,
			&item.ShiftCount,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

// Helper function to convert string to uint
func stringToUint(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

func QueryCustomerLifecycleWhereFilters(q *models.ReportsQueryCustomerLifecycleChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter using customerSections join
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" cs.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.id = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.FilterIDsAsInt64Slice))
			paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) GetCustomerLifecycle(q *models.ReportsQueryCustomerLifecycleChartRequestParams) ([]domain.CustomerLifecycle, error) {
	query := `
		WITH lifecycle_stages AS (
			SELECT 
				CASE 
					WHEN status = 'active' THEN 'Active'
					WHEN status = 'inactive' THEN 'Inactive'
					WHEN status = 'dead' THEN 'Dead'
					ELSE 'Unknown'
				END as stage,
				CASE 
					WHEN status = 'active' THEN 1
					WHEN status = 'inactive' THEN 2
					WHEN status = 'dead' THEN 3
					ELSE 4
				END as stage_order
			FROM customers c
			LEFT JOIN customerSections cs ON c.id = cs.customerId
			WHERE c.deleted_at IS NULL
	`

	// Add filters
	where, args, _ := QueryCustomerLifecycleWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Close the CTE and add the main query
	query += `
		)
		SELECT 
			stage,
			COUNT(*) as count,
			stage_order
		FROM lifecycle_stages
		WHERE stage != 'Unknown'
		GROUP BY stage, stage_order
		ORDER BY stage_order
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.CustomerLifecycle
	for rows.Next() {
		var item domain.CustomerLifecycle
		err := rows.Scan(
			&item.Stage,
			&item.Count,
			&item.StageOrder,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func QueryShiftCountPerDayWhereFilters(q *models.ReportsQueryShiftCountPerDayChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			// query += fmt.Sprintf(" cps.id = ANY($%d)", paramCount)
			// args = append(args, q.FilterIDsAsInt64Slice)
			// paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) GetShiftCountPerDay(q *models.ReportsQueryShiftCountPerDayChartRequestParams) ([]domain.ShiftCountPerDay, error) {
	query := `
		WITH days AS (
			SELECT 
				day_name,
				day_order
			FROM (
				VALUES 
					('Monday', 1),
					('Tuesday', 2),
					('Wednesday', 3),
					('Thursday', 4),
					('Friday', 5),
					('Saturday', 6),
					('Sunday', 7)
			) AS d(day_name, day_order)
		)
		SELECT 
			days.day_name as dayOfWeek,
			COALESCE(COUNT(DISTINCT cps.id), 0) as shiftCount,
			days.day_order as dayOrder
		FROM days
		LEFT JOIN cycleStaffTypes cst ON EXTRACT(DOW FROM cst.datetime) = days.day_order - 1
		LEFT JOIN cyclePickupShifts cps ON cps.cycleStaffTypeId = cst.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		WHERE cps.id IS NOT NULL
	`

	// Add filters
	where, args, _ := QueryShiftCountPerDayWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Group by day and order by day order
	query += " GROUP BY days.day_name, days.day_order ORDER BY days.day_order"

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.ShiftCountPerDay
	for rows.Next() {
		var item domain.ShiftCountPerDay
		err := rows.Scan(
			&item.DayOfWeek,
			&item.ShiftCount,
			&item.DayOrder,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func QueryShiftHeatmapWhereFilters(q *models.ReportsQueryShiftHeatmapChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" cps.id = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.FilterIDsAsInt64Slice))
			paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) GetShiftHeatmap(q *models.ReportsQueryShiftHeatmapChartRequestParams) ([]domain.ShiftHeatmap, error) {
	query := `
		WITH days AS (
			SELECT 
				day_name,
				day_order
			FROM (
				VALUES 
					('Monday', 1),
					('Tuesday', 2),
					('Wednesday', 3),
					('Thursday', 4),
					('Friday', 5),
					('Saturday', 6),
					('Sunday', 7)
			) AS d(day_name, day_order)
		),
		hours AS (
			SELECT generate_series(0, 23) AS hour
		)
		SELECT 
			days.day_name as dayOfWeek,
			days.day_order as dayOrder,
			hours.hour as hour,
			COALESCE(COUNT(DISTINCT cps.id), 0) as shiftCount
		FROM days
		CROSS JOIN hours
		LEFT JOIN cycleStaffTypes cst ON 
			EXTRACT(DOW FROM cst.datetime) = days.day_order - 1 AND
			EXTRACT(HOUR FROM cst.datetime) = hours.hour
		LEFT JOIN cyclePickupShifts cps ON cps.cycleStaffTypeId = cst.id
		LEFT JOIN cycles c ON cps.cycleId = c.id
		WHERE cps.id IS NOT NULL
	`

	// Add filters
	where, args, _ := QueryShiftHeatmapWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Group by day and hour, order by day and hour
	query += " GROUP BY days.day_name, days.day_order, hours.hour ORDER BY days.day_order, hours.hour"

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.ShiftHeatmap
	for rows.Next() {
		var item domain.ShiftHeatmap
		err := rows.Scan(
			&item.DayOfWeek,
			&item.DayOrder,
			&item.Hour,
			&item.ShiftCount,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func QueryCustomerLocationWhereFilters(q *models.ReportsQueryCustomerLocationChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter using customerSections join
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" cs.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.id = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.FilterIDsAsInt64Slice))
			paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) GetCustomerLocation(q *models.ReportsQueryCustomerLocationChartRequestParams) ([]domain.CustomerLocation, error) {
	query := `
		WITH customer_cities AS (
			SELECT DISTINCT
				ci.id as cityId,
				ci.name as cityName,
				c.id as customerId
			FROM customers c
			LEFT JOIN customerSections cs ON c.id = cs.customerId
			LEFT JOIN addresses a ON c.id = a.customerId
			LEFT JOIN cities ci ON a.cityId = ci.id
			WHERE c.deleted_at IS NULL
			AND a.isMainAddress = true
			AND ci.id IS NOT NULL
	`

	// Add filters
	where, args, _ := QueryCustomerLocationWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Close the CTE and add the main query
	query += `
		)
		SELECT 
			cityId,
			cityName,
			COUNT(DISTINCT customerId) as customerCount
		FROM customer_cities
		GROUP BY cityId, cityName
		ORDER BY customerCount DESC
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.CustomerLocation
	for rows.Next() {
		var item domain.CustomerLocation
		err := rows.Scan(
			&item.CityID,
			&item.CityName,
			&item.CustomerCount,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func QueryCustomerAgeGroupWhereFilters(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) ([]string, []interface{}, int) {
	var (
		where      []string
		args       []interface{} = []interface{}{}
		paramCount int           = 1
	)
	if q != nil {
		// Add section IDs filter using customerSections join
		if len(q.SectionIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" cs.sectionId = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.SectionIDsAsInt64Slice))
			paramCount++
		}

		// Add filter IDs
		if len(q.FilterIDsAsInt64Slice) > 0 {
			where = append(where, fmt.Sprintf(" c.id = ANY($%d)", paramCount))
			args = append(args, pq.Int64Array(q.FilterIDsAsInt64Slice))
			paramCount++
		}

		// Add aggregations
		if len(q.AggregationsValue) > 0 {
			for _, agg := range q.AggregationsValue {
				switch agg.OperationType {
				case "equal":
					where = append(where, fmt.Sprintf(" %s = $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "notEqual":
					where = append(where, fmt.Sprintf(" %s != $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "greaterThan":
					where = append(where, fmt.Sprintf(" %s > $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "lessThan":
					where = append(where, fmt.Sprintf(" %s < $%d", agg.Field, paramCount))
					args = append(args, agg.Value)
					paramCount++
				case "contains":
					where = append(where, fmt.Sprintf(" %s ILIKE $%d", agg.Field, paramCount))
					args = append(args, "%"+agg.Value+"%")
					paramCount++
				}
			}
		}
	}
	return where, args, paramCount
}

func (r *ReportRepositoryPostgresDB) GetCustomerAgeGroup(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) ([]domain.CustomerAgeGroup, error) {
	query := `
		WITH customer_ages AS (
			SELECT 
				c.id as customerId,
				CASE 
					WHEN u.birthDate IS NULL THEN 'Unknown'
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 21 THEN '0-20'
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 41 THEN '21-40'
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 61 THEN '41-60'
					ELSE '61+'
				END as age_group,
				CASE 
					WHEN u.birthDate IS NULL THEN 0
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 21 THEN 1
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 41 THEN 2
					WHEN EXTRACT(YEAR FROM AGE(NOW(), u.birthDate)) < 61 THEN 3
					ELSE 4
				END as group_order
			FROM customers c
			LEFT JOIN customerSections cs ON c.id = cs.customerId
			LEFT JOIN users u ON c.userId = u.id
			WHERE c.deleted_at IS NULL
	`

	// Add filters
	where, args, _ := QueryCustomerAgeGroupWhereFilters(q)
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// Close the CTE and add the main query
	query += `
		)
		SELECT 
			age_group as ageGroup,
			COUNT(DISTINCT customerId) as customerCount,
			group_order as groupOrder
		FROM customer_ages
		GROUP BY age_group, group_order
		ORDER BY group_order
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.CustomerAgeGroup
	for rows.Next() {
		var item domain.CustomerAgeGroup
		err := rows.Scan(
			&item.AgeGroup,
			&item.CustomerCount,
			&item.GroupOrder,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}
