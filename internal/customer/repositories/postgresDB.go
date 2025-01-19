package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	sharedconstants "github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/constants"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/customer/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type CustomerRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewCustomerRepositoryPostgresDB(d *sql.DB) *CustomerRepositoryPostgresDB {
	return &CustomerRepositoryPostgresDB{
		PostgresDB: d,
	}
}

/*
FindCustomerServicesForSpecificShift returns a list of customer services for a given cycle pickup shift id and date and shift name.
It will return a list of customer services which are not already picked for the given shift and date.
It will also insert the picked service to cyclePickupShiftCustomers table.

Args:

	cyclePickupShiftID (int64): The id of the cycle pickup shift.
	date (time.Time): The date of the shift.
	shiftName (string): The name of the shift.

Returns:

	[]*domain.CustomerServices: A list of customer services.
	error: An error if there is no service for the given shift and date.
*/
func (r *CustomerRepositoryPostgresDB) FindCustomerServicesForSpecificShift(cyclePickupShiftID int64, date time.Time, shiftName string, shiftMorningStartHour int64, shiftMorningEndHour int64, shiftEveningStartHour int64, shiftEveningEndHour int64, shiftNightStartHour int64, shiftNightEndHour int64) ([]*domain.CustomerServices, error) {
	// Get inserted for this shift
	var (
		insertedCustomerServiceId int64
		currentDateTime           time.Time = time.Now()
	)
	err := r.PostgresDB.QueryRow(`
		SELECT customerServiceId FROM cyclePickupShiftCustomers
		WHERE
			cyclePickupShiftId = $1
		LIMIT 1
	`, cyclePickupShiftID).Scan(&insertedCustomerServiceId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	if insertedCustomerServiceId > 0 {
		services, err := r.QueryServices(&models.CustomersQueryServicesRequestParams{
			ID:    int(insertedCustomerServiceId),
			Page:  1,
			Limit: -1,
		})
		if err != nil {
			return nil, err
		}
		return services, nil
	}

	// Get all services for this shift
	services, err := r.QueryServices(&models.CustomersQueryServicesRequestParams{
		Page:  1,
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}
	log.Println("services length", len(services))
	if len(services) == 0 {
		return nil, errors.New("there is no service for this shift 1")
	}
	currentDayName := date.Weekday().String()
	log.Println("currentDayName", currentDayName)
	log.Println("shiftName", shiftName)
	for _, service := range services {
		// Check if service is already picked or not
		var (
			createdAt     sql.NullTime
			pickedService *domain.CustomerServices = nil
			noRows        bool
		)
		err := r.PostgresDB.QueryRow(`
			SELECT created_at FROM cyclePickupShiftCustomers
			WHERE
				customerServiceId = $1
			ORDER BY created_at DESC LIMIT 1
		`, service.ID,
		).Scan(
			&createdAt,
		)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			} else {
				noRows = true
			}
		}
		log.Println("noRows", noRows)

		if noRows {
			// Service is not picked
			// Pick service
			log.Println("repeat", *service.Repeat)
			if *service.Repeat == fmt.Sprintf("every%s", currentDayName) {
				// Set shiftName based on start and end time
				shName := shifts.MorningShift
				if service.TimeValue.Hour() >= int(shiftMorningStartHour) && service.TimeValue.Hour() < int(shiftMorningEndHour) {
					shName = shifts.MorningShift
				}
				if service.TimeValue.Hour() >= int(shiftEveningStartHour) && service.TimeValue.Hour() < int(shiftEveningEndHour) {
					shName = shifts.EveningShift
				}
				if service.TimeValue.Hour() >= int(shiftNightStartHour) || service.TimeValue.Hour() < int(shiftNightEndHour) {
					shName = shifts.NightShift
				}
				log.Println("shName shiftName", shName, shiftName)
				if shName == shiftName {
					pickedService = service
				}
			} else {
				if service.Repeat == nil || *service.Repeat == sharedconstants.SERVICE_REPEAT_DAILY || *service.Repeat == sharedconstants.SERVICE_REPEAT_WEEKLY || *service.Repeat == sharedconstants.SERVICE_REPEAT_MONTHLY {
					pickedService = service
				}
			}
			if pickedService != nil {
				// Insert service to cyclePickupShiftCustomers
				var insertedID int64
				err = r.PostgresDB.QueryRow(`
					INSERT INTO cyclePickupShiftCustomers
						(cyclePickupShiftId, customerId, customerServiceId, created_at)
					VALUES
						($1, $2, $3, $4)
					RETURNING id
				`, cyclePickupShiftID, service.CustomerID, service.ID, currentDateTime).Scan(&insertedID)
				if err != nil {
					return nil, err
				}

				// Return picked service
				return []*domain.CustomerServices{pickedService}, nil
			}
		} else {
			if createdAt.Valid {
				// Service is already picked
				createdAtDayName := createdAt.Time.Weekday().String()
				if *service.Repeat == fmt.Sprintf("every%s", currentDayName) {
					if currentDayName == createdAtDayName {
						// Set shiftName based on start and end time
						shName := shifts.MorningShift
						if service.TimeValue.Hour() >= int(shiftMorningStartHour) && service.TimeValue.Hour() < int(shiftMorningEndHour) {
							shName = shifts.MorningShift
						}
						if service.TimeValue.Hour() >= int(shiftEveningStartHour) && service.TimeValue.Hour() < int(shiftEveningEndHour) {
							shName = shifts.EveningShift
						}
						if service.TimeValue.Hour() >= int(shiftNightStartHour) || service.TimeValue.Hour() < int(shiftNightEndHour) {
							shName = shifts.NightShift
						}
						if shName == shiftName {
							var insertedID int64
							err = r.PostgresDB.QueryRow(`
								INSERT INTO cyclePickupShiftCustomers
									(cyclePickupShiftId, customerId, customerServiceId, created_at)
								VALUES
									($1, $2, $3, $4)
								RETURNING id
							`, cyclePickupShiftID, service.CustomerID, service.ID, currentDateTime).Scan(&insertedID)
							if err != nil {
								return nil, err
							}
							return []*domain.CustomerServices{service}, nil
						}
					}
				} else {
					if service.Repeat != nil {
						var insertedID int64
						if *service.Repeat == sharedconstants.SERVICE_REPEAT_DAILY {
							if currentDayName != createdAtDayName {
								err = r.PostgresDB.QueryRow(`
									INSERT INTO cyclePickupShiftCustomers
										(cyclePickupShiftId, customerId, customerServiceId, created_at)
									VALUES
										($1, $2, $3, $4)
									RETURNING id
								`, cyclePickupShiftID, service.CustomerID, service.ID, currentDateTime).Scan(&insertedID)
								if err != nil {
									return nil, err
								}
								return []*domain.CustomerServices{service}, nil
							}
						}
						if *service.Repeat == sharedconstants.SERVICE_REPEAT_WEEKLY {
							nextWeek := createdAt.Time.AddDate(0, 0, 7)
							if nextWeek.Weekday().String() == currentDayName {
								err = r.PostgresDB.QueryRow(`
									INSERT INTO cyclePickupShiftCustomers
										(cyclePickupShiftId, customerId, customerServiceId, created_at)
									VALUES
										($1, $2, $3, $4)
									RETURNING id
								`, cyclePickupShiftID, service.CustomerID, service.ID, currentDateTime).Scan(&insertedID)
								if err != nil {
									return nil, err
								}
								return []*domain.CustomerServices{service}, nil
							}
						}
						if *service.Repeat == sharedconstants.SERVICE_REPEAT_MONTHLY {
							nextMonth := createdAt.Time.AddDate(0, 1, 0)
							if nextMonth.Weekday().String() == currentDayName {
								err = r.PostgresDB.QueryRow(`
									INSERT INTO cyclePickupShiftCustomers
										(cyclePickupShiftId, customerId, customerServiceId, created_at)
									VALUES
										($1, $2, $3, $4)
									RETURNING id
								`, cyclePickupShiftID, service.CustomerID, service.ID, currentDateTime).Scan(&insertedID)
								if err != nil {
									return nil, err
								}
								return []*domain.CustomerServices{service}, nil
							}
						}
					}
				}
				continue
			}
		}
	}
	return nil, errors.New("there is no service for this shift 2")
}

func makeWhereFilters(queries *models.CustomersQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf("customers.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf("customers.userId = %d", queries.UserID))
		}
		if queries.Filters.UserId.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.UserId.Op, fmt.Sprintf("%v", queries.Filters.UserId.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("customers.userId %s %s", opValue.Operator, val))
		}
		if queries.Filters.FirstName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FirstName.Op, queries.Filters.FirstName.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.firstName %s %s", opValue.Operator, val))
		}
		if queries.Filters.LastName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.LastName.Op, queries.Filters.LastName.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.lastName %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) Query(queries *models.CustomersQueryRequestParams) ([]*domain.Customer, error) {
	q := `
		SELECT customers.*,
		   users.firstName AS userFirstName,
		   users.lastName AS userLastName,
		   users.avatarUrl AS userAvatarUrl,
		   users.gender AS userGender,
		   users.email AS userEmail,
		   users.phone AS userPhone,
		   users.birthDate AS userBirthDate,
		   users.nationalCode AS userNationalCode,
		   staffs.id AS rNId,
		   u.firstName AS rNFirstName,
		   u.lastName AS rNLastName,
		   u.avatarUrl AS rNAvatarUrl
		FROM customers
		LEFT JOIN users ON customers.userId = users.id
		LEFT JOIN staffs ON customers.responsibleNurseId = staffs.id
		LEFT JOIN users u ON staffs.userId = u.id
	 `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.FirstName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" users.firstName %s", queries.Sorts.FirstName.Op))
		}
		if queries.Sorts.LastName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" users.lastName %s", queries.Sorts.LastName.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" customers.created_at %s", queries.Sorts.CreatedAt.Op))
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
	var customers []*domain.Customer
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			customer                                     domain.Customer
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
			var customerUser domain.CustomerUser
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
			customer.ResponsibleNurse = &domain.CustomerResponsibleNurse{
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
			customersRelatives []domain.CustomersRelative
		)
		for rows.Next() {
			var (
				customersRelative domain.CustomersRelative
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
			customerDiagnoses []domain.CustomerDiagnose
		)
		for rows.Next() {
			var (
				customerDiagnose domain.CustomerDiagnose
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
				customerDiagnose.Diagnose = &domain.CustomerDiagnoseDiagnose{
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
		log.Printf("%#v\n", diagnoseIDs)
		customer.DiagnoseIDs = diagnoseIDs
		customer.Diagnoses = customerDiagnoses

		customers = append(customers, &customer)
	}
	if rows != nil {
		rows.Close()
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
		if err == nil {
			for rows.Next() {
				var (
					section     domain.CustomerSection
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
						parent           domain.CustomerSection
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
						child       domain.CustomerSection
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
							parent           domain.CustomerSection
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
					customer.Sections = []domain.CustomerSection{}
				}
				customer.Sections = append(customer.Sections, section)
			}
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
		if err == nil {
			for rows.Next() {
				var (
					customerLimitation    domain.CustomerLimitation
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
				if err == nil {
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
		}
		if rows != nil {
			rows.Close()
		}
		if customer.Limitations == nil {
			customer.Limitations = []domain.CustomerLimitation{}
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
		if err == nil {
			for rows.Next() {
				var customerMotherLanguage domain.CustomerMotherLang
				err := rows.Scan(
					&customerMotherLanguage.ID,
					&customerMotherLanguage.Name,
				)
				if err == nil {
					customer.MotherLangs = append(customer.MotherLangs, customerMotherLanguage)
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if customer.MotherLangs == nil {
			customer.MotherLangs = []domain.CustomerMotherLang{}
		}
	}

	return customers, nil
}

func (r *CustomerRepositoryPostgresDB) Count(queries *models.CustomersQueryRequestParams) (int64, error) {
	q := `
		SELECT COUNT(customers.id)
		FROM customers
		LEFT JOIN users ON customers.userId = users.id
	 `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreatePersonalInfo(payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	currentTime := time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Create MotherlangIds jsonb
	motherLangIdsBytes, err := json.Marshal(payload.MotherLangIDsInt64)
	if err != nil {
		return nil, err
	}
	motherLangIdsJSON := string(motherLangIdsBytes)

	q := `
		INSERT INTO customers (
			motherLangIds,
		    status,
			statusDate,
			created_at,
			updated_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var insertedID int
	err = tx.QueryRow(
		q,
		motherLangIdsJSON,
		payload.Status,
		payload.StatusDateAsDate,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert limitations
	if payload.Limitations != nil {
		for _, limitation := range payload.Limitations {
			q := `
				INSERT INTO customerLimitations (customerId, limitationId, description)
				VALUES ($1, $2, $3)
			`
			_, err := tx.Exec(q, insertedID, limitation.LimitationID, limitation.Description)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Insert customerStatusLogs if not exists
	var (
		prevStatusValue sql.NullString
		statusValue     = payload.Status
		shouldBeCreate  = false
	)
	err = tx.QueryRow(`
		SELECT csl.statusValue FROM customerStatusLogs csl WHERE csl.customerId = $1 AND csl.statusValue = $2 ORDER BY csl.created_at DESC LIMIT 1
	`, insertedID, statusValue).Scan(&prevStatusValue)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, err
		}
	}
	if prevStatusValue.Valid {
		if prevStatusValue.String != statusValue {
			shouldBeCreate = true
		}
	} else {
		shouldBeCreate = true
	}
	if shouldBeCreate {
		_, err = tx.Exec(
			`INSERT INTO customerStatusLogs (customerId, statusValue, createdBy, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
			insertedID,
			statusValue,
			payload.AuthenticatedUser.ID,
			currentTime,
			currentTime,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer
	customer, err := r.Query(&models.CustomersQueryRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}
	if len(customer) == 0 {
		return nil, errors.New("customer not found")
	}

	// Return customer
	return customer[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdatePersonalInfo(customerId int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	currentTime := time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Create MotherlangIds jsonb
	motherLangIdsBytes, err := json.Marshal(payload.MotherLangIDsInt64)
	if err != nil {
		return nil, err
	}
	motherLangIdsJSON := string(motherLangIdsBytes)

	q := `
		UPDATE customers
		SET
			motherLangIds = $1,
		    status = $2,
			statusDate = $3,
			updated_at = $4
		WHERE id = $5
	`
	_, err = tx.Exec(
		q,
		motherLangIdsJSON,
		payload.Status,
		payload.StatusDateAsDate,
		currentTime,
		customerId,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete limitations
	q = `
		DELETE FROM customerLimitations
		WHERE customerId = $1
	`
	_, err = tx.Exec(q, customerId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert limitations
	if payload.Limitations != nil {
		for _, limitation := range payload.Limitations {
			q := `
				INSERT INTO customerLimitations (customerId, limitationId, description)
				VALUES ($1, $2, $3)
			`
			_, err := tx.Exec(q, customerId, limitation.LimitationID, limitation.Description)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Insert customerStatusLogs if not exists
	var (
		prevStatusValue sql.NullString
		statusValue     = payload.Status
		shouldBeCreate  = false
	)
	err = tx.QueryRow(`
		SELECT csl.statusValue FROM customerStatusLogs csl WHERE csl.customerId = $1 AND csl.statusValue = $2 ORDER BY csl.created_at DESC LIMIT 1
	`, customerId, statusValue).Scan(&prevStatusValue)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, err
		}
	}
	if prevStatusValue.Valid {
		if prevStatusValue.String != statusValue {
			shouldBeCreate = true
		}
	} else {
		shouldBeCreate = true
	}
	if shouldBeCreate {
		_, err = tx.Exec(
			`INSERT INTO customerStatusLogs (customerId, statusValue, createdBy, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
			customerId,
			statusValue,
			payload.AuthenticatedUser.ID,
			currentTime,
			currentTime,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer
	customer, err := r.Query(&models.CustomersQueryRequestParams{
		ID: int(customerId),
	})
	if err != nil {
		return nil, err
	}

	// Return customer
	return customer[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdateUserInformation(customerID int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	currentTime := time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	q := `
		UPDATE customers
		SET
			userId = $1,
			updated_at = $2
		WHERE id = $3
	`
	_, err = tx.Exec(
		q,
		payload.UserID,
		currentTime,
		customerID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer
	customer, err := r.Query(&models.CustomersQueryRequestParams{
		ID: int(customerID),
	})
	if err != nil {
		return nil, err
	}
	if len(customer) == 0 {
		return nil, errors.New("customer not found")
	}

	// Return customer
	return customer[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdateAdditionalInfo(payload *models.CustomersUpdateAdditionalInfoRequestBody) (*domain.Customer, error) {
	currentTime := time.Now()

	// Get customer
	customers, err := r.Query(&models.CustomersQueryRequestParams{
		ID: int(payload.CustomerID),
	})
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, errors.New("customer not found")
	}
	customerIdentity := customers[0]

	// Start transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	q := `
		UPDATE customers
		SET
			keyNo = $1,
			paymentMethod = $2,
			nurseGenderWish = $3,
			responsibleNurseId = $4,
		    limitingTheRightToSelfDeterminationStartDate = $5,
		    limitingTheRightToSelfDeterminationEndDate = $6,
		    hasLimitingTheRightToSelfDetermination = $7,
		    mobilityContract = $8,
			parkingInfo = $9,
			extraExplanation = $10,
			updated_at = $11
		WHERE id = $12
	`
	_, err = tx.Exec(
		q,
		payload.KeyNo,
		payload.PaymentMethod,
		payload.NurseGenderWish,
		payload.ResponsibleNurseID,
		payload.LimitingTheRightToSelfDeterminationStartDateAsDate,
		payload.LimitingTheRightToSelfDeterminationEndDateAsDate,
		payload.HasLimitingTheRightToSelfDetermination,
		payload.MobilityContract,
		payload.ParkingInfo,
		payload.ExtraExplanation,
		currentTime,
		payload.CustomerID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert customerContractualMobilityRestrictionLogs
	var (
		beforeValue = customerIdentity.MobilityContract
		afterValue  = payload.MobilityContract
	)
	_, err = tx.Exec(
		`INSERT INTO customerContractualMobilityRestrictionLogs (customerId, beforeValue, afterValue, createdBy, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		payload.CustomerID,
		beforeValue,
		afterValue,
		payload.AuthenticatedUser.ID,
		currentTime,
		currentTime,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete sections
	q = `
		DELETE FROM customerSections
		WHERE customerId = $1
	`
	_, err = tx.Exec(q, payload.CustomerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert sections
	if payload.SectionIDsInt64 != nil {
		for _, sectionID := range payload.SectionIDsInt64 {
			q := `
				INSERT INTO customerSections (customerId, sectionId)
				VALUES ($1, $2)
			`
			_, err := tx.Exec(q, payload.CustomerID, sectionID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Delete customersRelatives
	q = `
		DELETE FROM customersRelatives
		WHERE customerId = $1
	`
	_, err = tx.Exec(q, payload.CustomerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create customersRelatives
	if payload.RelativeIDsInt64 != nil {
		for _, relativeID := range payload.RelativeIDsInt64 {
			q := `
				INSERT INTO customersRelatives (customerId, relativeId)
				VALUES ($1, $2)
			`
			_, err := tx.Exec(q, payload.CustomerID, relativeID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Delete customersDiagnoses
	q = `
		DELETE FROM customersDiagnoses
		WHERE customerId = $1
	`
	_, err = tx.Exec(q, payload.CustomerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create customersDiagnoses
	if payload.DiagnoseIDsInt64 != nil {
		for _, diagnoseID := range payload.DiagnoseIDsInt64 {
			q := `
				INSERT INTO customersDiagnoses (customerId, diagnoseId)
				VALUES ($1, $2)
			`
			_, err := tx.Exec(q, payload.CustomerID, diagnoseID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer
	customer, err := r.Query(&models.CustomersQueryRequestParams{
		ID: int(payload.CustomerID),
	})
	if err != nil {
		return nil, err
	}
	if len(customer) == 0 {
		return nil, errors.New("customer not found")
	}

	// Return customer
	return customer[0], nil
}

func makeCreditDetailsWhereFilters(queries *models.CustomersQueryCreditDetailsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cd.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" cd.customerId = %d ", queries.CustomerID))
		}
		if queries.Filters.BankAccountNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.BankAccountNumber.Op, fmt.Sprintf("%v", queries.Filters.BankAccountNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cd.bankAccountNumber %s %s ", opValue.Operator, val))
		}
		if queries.Filters.BillingAddressName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.BillingAddressName.Op, fmt.Sprintf("%v", queries.Filters.BillingAddressName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ba.name %s %s ", opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryCreditDetails(queries *models.CustomersQueryCreditDetailsRequestParams) ([]*domain.CustomerCreditDetail, error) {
	q := `
		SELECT
			cd.id,
			cd.customerId,
			cd.billingAddressId,
			cd.bankAccountNumber,
			cd.created_at,
			cd.updated_at,
			cd.deleted_at,
			ba.id as billingAddressId,
			ba.cityId as billingAddressCityId,
			ba.street as billingAddressStreet,
			ba.name as billingAddressName,
			ba.postalCode as billingAddressPostalCode,
			ba.buildingNumber as billingAddressBuildingNumber,
			c.name as billingAddressCityName
		FROM customerCreditDetails cd
		LEFT JOIN addresses ba ON ba.id = cd.billingAddressId
		LEFT JOIN cities c ON c.id = ba.cityId
	`
	if queries != nil {
		where := makeCreditDetailsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cd.created_at %s ", queries.Sorts.CreatedAt.Op))
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

	var customerCreditDetails []*domain.CustomerCreditDetail
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			customerCreditDetail         = &domain.CustomerCreditDetail{}
			deletedAt                    sql.NullTime
			billingAddressID             sql.NullInt64
			billingAddressCityID         sql.NullInt64
			billingAddressStreet         sql.NullString
			billingAddressName           sql.NullString
			billingAddressPostalCode     sql.NullString
			billingAddressBuildingNumber sql.NullString
			billingAddressCityName       sql.NullString
		)
		err := rows.Scan(
			&customerCreditDetail.ID,
			&customerCreditDetail.CustomerID,
			&customerCreditDetail.BillingAddressID,
			&customerCreditDetail.BankAccountNumber,
			&customerCreditDetail.CreatedAt,
			&customerCreditDetail.UpdatedAt,
			&deletedAt,
			&billingAddressID,
			&billingAddressCityID,
			&billingAddressStreet,
			&billingAddressName,
			&billingAddressPostalCode,
			&billingAddressBuildingNumber,
			&billingAddressCityName,
		)
		if err != nil {
			return nil, err
		}
		customerCreditDetail.DeletedAt = exp.TerIf(deletedAt.Valid, &deletedAt.Time, nil)
		if billingAddressID.Valid {
			customerCreditDetail.BillingAddress = &domain.CustomerCreditDetailBillingAddress{
				ID: uint(billingAddressID.Int64),
			}
			if billingAddressCityID.Valid {
				customerCreditDetail.BillingAddress.City = &domain.CustomerCreditDetailBillingAddressCity{
					ID: uint(billingAddressCityID.Int64),
				}
				if billingAddressCityName.Valid {
					customerCreditDetail.BillingAddress.City.Name = billingAddressCityName.String
				}
			}
			if billingAddressStreet.Valid {
				customerCreditDetail.BillingAddress.Street = billingAddressStreet.String
			}
			if billingAddressName.Valid {
				customerCreditDetail.BillingAddress.Name = billingAddressName.String
			}
			if billingAddressPostalCode.Valid {
				customerCreditDetail.BillingAddress.PostalCode = &billingAddressPostalCode.String
			}
			if billingAddressBuildingNumber.Valid {
				customerCreditDetail.BillingAddress.BuildingNumber = billingAddressBuildingNumber.String
			}
		}
		customerCreditDetails = append(customerCreditDetails, customerCreditDetail)
	}

	return customerCreditDetails, nil
}

func (r *CustomerRepositoryPostgresDB) CountCreditDetails(queries *models.CustomersQueryCreditDetailsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cd.id)
		FROM customerCreditDetails cd
		LEFT JOIN addresses ba ON ba.id = cd.billingAddressId
		LEFT JOIN cities c ON c.id = ba.cityId
	`
	if queries != nil {
		where := makeCreditDetailsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateCreditDetails(payload *models.CustomersCreateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	currentTime := time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	q := `
		INSERT INTO customerCreditDetails (
			customerId,
			billingAddressId,
			bankAccountNumber,
			created_at,
			updated_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var insertedID int
	err = tx.QueryRow(
		q,
		payload.CustomerID,
		payload.BillingAddressID,
		payload.BankAccountNumber,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer credit details
	customerCreditDetails, err := r.QueryCreditDetails(&models.CustomersQueryCreditDetailsRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerCreditDetails) == 0 {
		return nil, errors.New("failed to get customer credit details")
	}

	return customerCreditDetails[0], nil
}

func (r *CustomerRepositoryPostgresDB) DeleteCustomerCreditDetails(payload *models.CustomersDeleteCreditDetailsRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM customerCreditDetails
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

func (r *CustomerRepositoryPostgresDB) UpdateCustomerCreditDetails(payload *models.CustomersUpdateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	currentTime := time.Now()
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	q := `
		UPDATE customerCreditDetails
		SET
			customerId = $1,
			billingAddressId = $2,
			bankAccountNumber = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id
	`
	var updatedID int
	err = tx.QueryRow(
		q,
		payload.CustomerID,
		payload.BillingAddressID,
		payload.BankAccountNumber,
		currentTime,
		payload.ID,
	).Scan(&updatedID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get customer credit details
	customerCreditDetails, err := r.QueryCreditDetails(&models.CustomersQueryCreditDetailsRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerCreditDetails) == 0 {
		return nil, errors.New("failed to get customer credit details")
	}

	return customerCreditDetails[0], nil
}

func makeAbsencesWhereFilters(queries *models.CustomersQueryAbsencesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ca.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" ca.customerId = %d ", queries.CustomerID))
		}
		if queries.Filters.StartDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartDate.Op, fmt.Sprintf("%v", queries.Filters.StartDate.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ca.start_date %s %s ", opValue.Operator, val))
		}
		if queries.Filters.EndDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndDate.Op, fmt.Sprintf("%v", queries.Filters.EndDate.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ca.end_date %s %s ", opValue.Operator, val))
		}
		if queries.Filters.Reason.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Reason.Op, fmt.Sprintf("%v", queries.Filters.Reason.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ca.reason %s %s ", opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryAbsences(queries *models.CustomersQueryAbsencesRequestParams) ([]*domain.CustomerAbsence, error) {
	q := `
		SELECT
			ca.id,
			ca.customerId,
			ca.start_date,
			ca.end_date,
			ca.reason,
			ca.attachments,
			ca.created_at,
			ca.updated_at,
			ca.deleted_at,
			c.id AS customerId,
			u.firstName AS customerFirstName,
			u.lastName AS customerLastName,
			u.avatarUrl AS customerAvatarUrl,
			u.email AS customerEmail
		FROM customerAbsences ca
		LEFT JOIN customers c ON c.id = ca.customerId
		LEFT JOIN users u ON u.id = c.userId
	`
	if queries != nil {
		where := makeAbsencesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.StartDate.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ca.start_date %s ", queries.Sorts.StartDate.Op))
		}
		if queries.Sorts.EndDate.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ca.end_date %s ", queries.Sorts.EndDate.Op))
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

	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerAbsences []*domain.CustomerAbsence
	for rows.Next() {
		var (
			customerAbsence     domain.CustomerAbsence
			reason              sql.NullString
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			deletedAt           sql.NullTime
			customerID          sql.NullInt64
			customerFirstName   sql.NullString
			customerLastName    sql.NullString
			customerAvatarUrl   sql.NullString
			customerEmail       sql.NullString
		)
		err := rows.Scan(
			&customerAbsence.ID,
			&customerAbsence.CustomerID,
			&customerAbsence.StartDate,
			&customerAbsence.EndDate,
			&reason,
			&attachments,
			&customerAbsence.CreatedAt,
			&customerAbsence.UpdatedAt,
			&deletedAt,
			&customerID,
			&customerFirstName,
			&customerLastName,
			&customerAvatarUrl,
			&customerEmail,
		)
		if err != nil {
			return nil, err
		}
		if reason.Valid {
			customerAbsence.Reason = &reason.String
		}
		if deletedAt.Valid {
			customerAbsence.DeletedAt = &deletedAt.Time
		}
		if customerID.Valid {
			customerAbsence.Customer = &domain.CustomerAbsenceCustomer{
				ID: customerID.Int64,
			}
			if customerFirstName.Valid {
				customerAbsence.Customer.FirstName = customerFirstName.String
			}
			if customerLastName.Valid {
				customerAbsence.Customer.LastName = customerLastName.String
			}
			if customerAvatarUrl.Valid {
				customerAbsence.Customer.AvatarUrl = customerAvatarUrl.String
			}
			if customerEmail.Valid {
				customerAbsence.Customer.Email = customerEmail.String
			}
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff customer: %d", err, customerAbsence.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CUSTOMER_BUCKET_NAME[len("maja."):])
				}
			}
			customerAbsence.Attachments = attachmentsMetadata
		}
		customerAbsences = append(customerAbsences, &customerAbsence)
	}

	return customerAbsences, nil
}

func (r *CustomerRepositoryPostgresDB) CountAbsences(queries *models.CustomersQueryAbsencesRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(ca.id)
		FROM customerAbsences ca
		LEFT JOIN customers c ON c.id = ca.customerId
		LEFT JOIN users u ON u.id = c.userId
	`
	if queries != nil {
		where := makeAbsencesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateAbsences(customer *domain.Customer, payload *models.CustomersCreateAbsencesRequestBody) (*domain.CustomerAbsence, error) {
	currentTime := time.Now()
	q := `
		INSERT INTO customerAbsences (
			customerId,
			start_date,
			end_date,
		    reason,
			created_at,
			updated_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var insertedID int
	err := r.PostgresDB.QueryRow(
		q,
		customer.ID,
		*payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Reason,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	// Get customer absences
	customerAbsences, err := r.QueryAbsences(&models.CustomersQueryAbsencesRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerAbsences) == 0 {
		return nil, errors.New("failed to get customer absences")
	}

	return customerAbsences[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdateAbsence(customerAbsence *domain.CustomerAbsence, payload *models.CustomersUpdateAbsenceRequestBody) (*domain.CustomerAbsence, error) {
	currentTime := time.Now()
	q := `
		UPDATE customerAbsences
		SET
			start_date = $1,
			end_date = $2,
			reason = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id
	`
	var updatedID int
	err := r.PostgresDB.QueryRow(
		q,
		*payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Reason,
		currentTime,
		customerAbsence.ID,
	).Scan(&updatedID)
	if err != nil {
		return nil, err
	}

	// Get customer absences
	customerAbsences, err := r.QueryAbsences(&models.CustomersQueryAbsencesRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerAbsences) == 0 {
		return nil, errors.New("failed to get customer absences")
	}

	return customerAbsences[0], nil
}

func (r *CustomerRepositoryPostgresDB) DeleteAbsences(payload *models.CustomersDeleteAbsencesRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				customerAbsences
			WHERE
			    id = ANY ($1) AND customerId = $2
			RETURNING id
		`,
		pq.Int64Array(payload.IDsInt64),
		payload.CustomerID,
	).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}

func makeServicesWhereFilters(queries *models.CustomersQueryServicesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cs.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" cs.customerId = %d ", queries.CustomerID))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryServices(queries *models.CustomersQueryServicesRequestParams) ([]*domain.CustomerServices, error) {
	q := `
		SELECT
			cs.id,
			cs.customerId,
			cs.serviceId,
			cs.serviceTypeId,
			cs.gradeId,
			cs.nurseWishId,
			cs.reportType,
			cs.timeValue,
			cs.repeat,
			cs.visitType,
			cs.serviceLengthMinute,
			cs.startTimeRange,
		    cs.endTimeRange,
			cs.description,
			cs.paymentMethod,
			cs.homeCareFee,
			cs.cityCouncilFee,
			cs.created_at,
			cs.updated_at,
			cs.deleted_at,
			s.id AS nurseWishId,
			u.firstName AS nurseWishFirstName,
			u.lastName AS nurseWishLastName,
			u.avatarUrl AS nurseWishAvatarUrl,
			sg.id AS gradeId,
			sg.name AS gradeName,
			sg.color AS gradeColor,
			sg.grade AS gradeGrade,
			sg.description AS gradeDescription,
			se.id AS serviceId,
			se.name AS serviceName,
			st.id AS serviceTypeId,
			st.name AS serviceTypeName,
			css.id AS customerCustomerId,
			css.keyNo AS customerKeyNo,
			css.userId AS customerUserId,
			csu.firstName AS customerFirstName,
			csu.lastName AS customerLastName
		FROM customerServices cs
		LEFT JOIN staffs s ON s.id = cs.nurseWishId
		LEFT JOIN users u ON u.id = s.userId
		LEFT JOIN customers css ON css.id = cs.customerId
		LEFT JOIN users csu ON csu.id = css.userId
		LEFT JOIN servicegrades sg ON sg.id = cs.gradeId
		LEFT JOIN services se ON se.id = cs.serviceId
		LEFT JOIN serviceTypes st ON st.id = cs.serviceTypeId
	`
	if queries != nil {
		where := makeServicesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerServices []*domain.CustomerServices
	for rows.Next() {
		customerService := &domain.CustomerServices{}
		var (
			deletedAt             sql.NullTime
			description           sql.NullString
			paymentMethod         sql.NullString
			homeCareFee           sql.NullFloat64
			cityCouncilFee        sql.NullFloat64
			nwID                  sql.NullInt64
			nurseWishFirstName    sql.NullString
			nurseWishLastName     sql.NullString
			nurseWishAvatarUrl    sql.NullString
			gradeId               sql.NullInt64
			gradeName             sql.NullString
			gradeColor            sql.NullString
			gradeGrade            sql.NullInt64
			gradeDescription      sql.NullString
			serviceId             sql.NullInt64
			serviceName           sql.NullString
			serviceTypeId         sql.NullInt64
			serviceTypeName       sql.NullString
			customerId            sql.NullInt64
			customerKeyNo         sql.NullString
			customerUserId        sql.NullInt64
			customerUserFirstName sql.NullString
			customerUserLastName  sql.NullString
		)
		err := rows.Scan(
			&customerService.ID,
			&customerService.CustomerID,
			&customerService.ServiceID,
			&customerService.ServiceTypeID,
			&customerService.GradeID,
			&customerService.NurseWishID,
			&customerService.ReportType,
			&customerService.TimeValue,
			&customerService.Repeat,
			&customerService.VisitType,
			&customerService.ServiceLengthMinute,
			&customerService.StartTimeRange,
			&customerService.EndTimeRange,
			&description,
			&paymentMethod,
			&homeCareFee,
			&cityCouncilFee,
			&customerService.CreatedAt,
			&customerService.UpdatedAt,
			&deletedAt,
			&nwID,
			&nurseWishFirstName,
			&nurseWishLastName,
			&nurseWishAvatarUrl,
			&gradeId,
			&gradeName,
			&gradeColor,
			&gradeGrade,
			&gradeDescription,
			&serviceId,
			&serviceName,
			&serviceTypeId,
			&serviceTypeName,
			&customerId,
			&customerKeyNo,
			&customerUserId,
			&customerUserFirstName,
			&customerUserLastName,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			customerService.DeletedAt = &deletedAt.Time
		}
		if description.Valid {
			customerService.Description = &description.String
		}
		if paymentMethod.Valid {
			customerService.PaymentMethod = paymentMethod.String
		}
		if homeCareFee.Valid {
			hcf := uint(homeCareFee.Float64)
			customerService.HomeCareFee = &hcf
		}
		if cityCouncilFee.Valid {
			ccf := uint(cityCouncilFee.Float64)
			customerService.CityCouncilFee = &ccf
		}
		if nwID.Valid {
			customerService.NurseWish = &domain.CustomerServiceNurseWish{
				ID: uint(nwID.Int64),
			}
			if nurseWishFirstName.Valid {
				customerService.NurseWish.FirstName = nurseWishFirstName.String
			}
			if nurseWishLastName.Valid {
				customerService.NurseWish.LastName = nurseWishLastName.String
			}
			if nurseWishAvatarUrl.Valid {
				customerService.NurseWish.AvatarUrl = nurseWishAvatarUrl.String
			}
		}
		if gradeId.Valid {
			customerService.Grade = &domain.CustomerServiceGrade{
				ID: uint(gradeId.Int64),
			}
			if gradeName.Valid {
				customerService.Grade.Name = gradeName.String
			}
			if gradeColor.Valid {
				customerService.Grade.Color = gradeColor.String
			}
			if gradeGrade.Valid {
				customerService.Grade.Grade = int(gradeGrade.Int64)
			}
			if gradeDescription.Valid {
				customerService.Grade.Description = gradeDescription.String
			}
		}
		if serviceId.Valid {
			customerService.Service = &domain.CustomerServiceService{
				ID: uint(serviceId.Int64),
			}
			if serviceName.Valid {
				customerService.Service.Name = serviceName.String
			}
		}
		if serviceTypeId.Valid {
			customerService.ServiceType = &domain.CustomerServiceServiceType{
				ID: uint(serviceTypeId.Int64),
			}
			if serviceTypeName.Valid {
				customerService.ServiceType.Name = serviceTypeName.String
			}
		}
		if customerId.Valid {
			customerService.Customer = &domain.CustomerServiceCustomer{
				ID: uint(customerId.Int64),
			}
			if customerKeyNo.Valid {
				customerService.Customer.KeyNo = customerKeyNo.String
			}
			if customerUserId.Valid {
				customerService.Customer.User = &domain.CustomerServiceCustomerUser{
					ID: uint(customerUserId.Int64),
				}
				if customerUserFirstName.Valid {
					customerService.Customer.User.FirstName = customerUserFirstName.String
				}
				if customerUserLastName.Valid {
					customerService.Customer.User.LastName = customerUserLastName.String
				}
			}
		}
		customerServices = append(customerServices, customerService)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customerServices, nil
}

func (r *CustomerRepositoryPostgresDB) CountServices(queries *models.CustomersQueryServicesRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cs.id)
		FROM customerServices cs
		LEFT JOIN staffs s ON s.id = cs.nurseWishId
		LEFT JOIN users u ON u.id = s.userId
		LEFT JOIN customers css ON css.id = cs.customerId
		LEFT JOIN users csu ON csu.id = css.userId
		LEFT JOIN servicegrades sg ON sg.id = cs.gradeId
		LEFT JOIN services se ON se.id = cs.serviceId
		LEFT JOIN serviceTypes st ON st.id = cs.serviceTypeId
	`
	if queries != nil {
		where := makeServicesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateServices(customer *domain.Customer, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	currentTime := time.Now()
	q := `
		INSERT INTO customerServices (
			customerId,
			serviceId,
			serviceTypeId,
			gradeId,
			nurseWishId,
			reportType,
			timeValue,
			repeat,
			visitType,
			serviceLengthMinute,
			startTimeRange,
		    endTimeRange,
			description,
			paymentMethod,
			homeCareFee,
			cityCouncilFee,
			created_at,
			updated_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NULL)
		RETURNING id
	`
	var insertedID int
	err := r.PostgresDB.QueryRow(
		q,
		customer.ID,
		payload.ServiceID,
		payload.ServiceTypeID,
		payload.GradeID,
		payload.NurseWishID,
		payload.ReportType,
		payload.TimeValue,
		payload.Repeat,
		payload.VisitType,
		payload.ServiceLengthMinute,
		payload.StartTimeRangeAsTime,
		payload.EndTimeRangeAsTime,
		payload.Description,
		payload.PaymentMethod,
		payload.HomeCareFee,
		payload.CityCouncilFee,
		currentTime,
		currentTime,
	).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	// Get customer services
	customerServices, err := r.QueryServices(&models.CustomersQueryServicesRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerServices) == 0 {
		return nil, errors.New("failed to get customer services")
	}

	return customerServices[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdateService(customerService *domain.CustomerServices, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	currentTime := time.Now()
	q := `
		UPDATE customerServices
		SET
			serviceId = $1,
			serviceTypeId = $2,
			gradeId = $3,
			nurseWishId = $4,
			reportType = $5,
			timeValue = $6,
			repeat = $7,
			visitType = $8,
			serviceLengthMinute = $9,
			startTimeRange = $10,
		    endTimeRange = $11,
			description = $12,
			paymentMethod = $13,
			homeCareFee = $14,
			cityCouncilFee = $15,
			updated_at = $16
		WHERE id = $17
		RETURNING id
	`
	var updatedID int
	err := r.PostgresDB.QueryRow(
		q,
		payload.ServiceID,
		payload.ServiceTypeID,
		payload.GradeID,
		payload.NurseWishID,
		payload.ReportType,
		payload.TimeValue,
		payload.Repeat,
		payload.VisitType,
		payload.ServiceLengthMinute,
		payload.StartTimeRangeAsTime,
		payload.EndTimeRangeAsTime,
		payload.Description,
		payload.PaymentMethod,
		payload.HomeCareFee,
		payload.CityCouncilFee,
		currentTime,
		customerService.ID,
	).Scan(&updatedID)
	if err != nil {
		return nil, err
	}

	// Get customer services
	customerServices, err := r.QueryServices(&models.CustomersQueryServicesRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerServices) == 0 {
		return nil, errors.New("failed to get customer services")
	}

	return customerServices[0], nil
}

func (r *CustomerRepositoryPostgresDB) DeleteServices(payload *models.CustomersDeleteServicesRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				customerServices
			WHERE
			    id = ANY ($1) AND customerId = $2
			RETURNING id
		`,
		pq.Int64Array(payload.IDsInt64),
		payload.CustomerID,
	).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}

func makeMedicinesWhereFilters(queries *models.CustomersQueryMedicinesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf("cm.id = %d", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf("cm.customerId = %d", queries.CustomerID))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryMedicines(queries *models.CustomersQueryMedicinesRequestParams) ([]*domain.CustomerMedicine, error) {
	q := `
		SELECT cm.*,
		   p.id AS pId,
		   p.title AS pTitle,
		   m.id AS mId,
		   m.name AS mName
		FROM customerMedicines cm
		LEFT JOIN prescriptions p ON cm.prescriptionId = p.id
		LEFT JOIN medicines m ON cm.medicineId = m.id
	 `
	if queries != nil {
		where := makeMedicinesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cm.created_at %s", queries.Sorts.CreatedAt.Op))
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
	var customerMedicines []*domain.CustomerMedicine
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			customerMedicine    domain.CustomerMedicine
			days                json.RawMessage
			hours               json.RawMessage
			startDate           sql.NullTime
			endDate             sql.NullTime
			warning             sql.NullString
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			deletedAt           sql.NullTime
			prescriptionId      sql.NullInt64
			prescriptionTitle   sql.NullString
			medicineId          sql.NullInt64
			medicineName        sql.NullString
		)
		err := rows.Scan(
			&customerMedicine.ID,
			&customerMedicine.CustomerID,
			&customerMedicine.PrescriptionID,
			&customerMedicine.MedicineID,
			&customerMedicine.DosageAmount,
			&customerMedicine.DosageUnit,
			&days,
			&customerMedicine.IsJustOneTime,
			&hours,
			&startDate,
			&endDate,
			&warning,
			&customerMedicine.IsUseAsNeeded,
			&attachments,
			&customerMedicine.CreatedAt,
			&customerMedicine.UpdatedAt,
			&deletedAt,
			&prescriptionId,
			&prescriptionTitle,
			&medicineId,
			&medicineName,
		)
		if err != nil {
			return nil, err
		}
		var (
			daysStr      string
			daysJSON     interface{}
			daysStrArray []string
			hoursStr     string
		)
		if err := json.Unmarshal(days, &daysStr); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(daysStr), &daysJSON); err != nil {
			return nil, err
		}
		daysArray, ok := daysJSON.([]interface{})
		if !ok {
			return nil, errors.New("failed to parse days")
		}
		for _, day := range daysArray {
			daysStrArray = append(daysStrArray, day.(string))
		}
		customerMedicine.Days = daysStrArray
		if err := json.Unmarshal(hours, &hoursStr); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(hoursStr), &customerMedicine.Hours); err != nil {
			return nil, err
		}
		if startDate.Valid {
			customerMedicine.StartDate = &startDate.Time
		}
		if endDate.Valid {
			customerMedicine.EndDate = &endDate.Time
		}
		if warning.Valid {
			customerMedicine.Warning = &warning.String
		}
		if deletedAt.Valid {
			customerMedicine.DeletedAt = &deletedAt.Time
		}
		if prescriptionId.Valid {
			customerMedicine.Prescription = &domain.CustomerMedicinePrescription{
				ID: uint(prescriptionId.Int64),
			}
			if prescriptionTitle.Valid {
				customerMedicine.Prescription.Title = prescriptionTitle.String
			}
		}
		if medicineId.Valid {
			customerMedicine.Medicine = &domain.CustomerMedicineMedicine{
				ID: uint(medicineId.Int64),
			}
			if medicineName.Valid {
				customerMedicine.Medicine.Name = medicineName.String
			}
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in customer medicine: %d", err, customerMedicine.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CUSTOMER_BUCKET_NAME[len("maja."):])
				}
			}
			customerMedicine.Attachments = attachmentsMetadata
		}
		customerMedicines = append(customerMedicines, &customerMedicine)
	}
	if rows != nil {
		rows.Close()
	}

	return customerMedicines, nil
}

func (r *CustomerRepositoryPostgresDB) CountMedicines(queries *models.CustomersQueryMedicinesRequestParams) (int64, error) {
	q := `
		SELECT COUNT(cm.id)
		FROM customerMedicines cm
		LEFT JOIN prescriptions p ON cm.prescriptionId = p.id
		LEFT JOIN medicines m ON cm.medicineId = m.id
	 `
	if queries != nil {
		where := makeMedicinesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateMedicines(customer *domain.Customer, payload *models.CustomersCreateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	q := `
		INSERT INTO customerMedicines (
			customerId,
			prescriptionId,
			medicineId,
			dosageAmount,
			dosageUnit,
			days,
			isJustOneTime,
			hours,
			start_date,
			end_date,
			warning,
			isUseAsNeeded
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12
		) RETURNING id
	`
	var (
		days       string
		hours      string
		insertedId int
	)
	if payload.Days != nil {
		daysBytes, err := json.Marshal(payload.Days)
		if err != nil {
			return nil, err
		}
		days = string(daysBytes)
	}
	if payload.HoursMetadata != nil {
		hoursBytes, err := json.Marshal(payload.Hours)
		if err != nil {
			return nil, err
		}
		hours = string(hoursBytes)
	}
	err := r.PostgresDB.QueryRow(
		q,
		customer.ID,
		payload.PrescriptionID,
		payload.MedicineID,
		payload.DosageAmount,
		payload.DosageUnit,
		days,
		payload.IsJustOneTimeAsBool,
		hours,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Warning,
		payload.IsUseAsNeededAsBool,
	).Scan(
		&insertedId,
	)
	if err != nil {
		return nil, err
	}

	// Get inserted customerMedicine
	customerMedicines, err := r.QueryMedicines(&models.CustomersQueryMedicinesRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(customerMedicines) == 0 {
		return nil, errors.New("CustomerMedicine not found")
	}
	customerMedicine := customerMedicines[0]

	return customerMedicine, nil
}

func (r *CustomerRepositoryPostgresDB) UpdateMedicine(customerMedicine *domain.CustomerMedicine, payload *models.CustomersUpdateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	q := `
		UPDATE customerMedicines SET
			prescriptionId = $1,
			medicineId = $2,
			dosageAmount = $3,
			dosageUnit = $4,
			days = $5,
			isJustOneTime = $6,
			hours = $7,
			start_date = $8,
			end_date = $9,
			warning = $10,
			isUseAsNeeded = $11,
			updated_at = NOW()
		WHERE id = $12
	`
	var (
		days  string
		hours string
	)
	if payload.DaysAsArray != nil {
		daysBytes, err := json.Marshal(payload.Days)
		if err != nil {
			return nil, err
		}
		days = string(daysBytes)
	}
	if payload.HoursMetadata != nil {
		hoursBytes, err := json.Marshal(payload.Hours)
		if err != nil {
			return nil, err
		}
		hours = string(hoursBytes)
	}
	log.Println("payload.DaysAsArray", payload.DaysAsArray)
	log.Println("payload.Days", payload.Days)
	log.Println("payload.HoursMetadata", payload.HoursMetadata)
	log.Println("payload.Hours", payload.Hours)
	log.Println("days", days)
	log.Println("hours", hours)
	_, err := r.PostgresDB.Exec(
		q,
		payload.PrescriptionID,
		payload.MedicineID,
		payload.DosageAmount,
		payload.DosageUnit,
		days,
		payload.IsJustOneTimeAsBool,
		hours,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Warning,
		payload.IsUseAsNeededAsBool,
		customerMedicine.ID,
	)
	if err != nil {
		return nil, err
	}

	// Get updated customerMedicine
	customerMedicines, err := r.QueryMedicines(&models.CustomersQueryMedicinesRequestParams{ID: int(customerMedicine.ID)})
	if err != nil {
		return nil, err
	}
	if len(customerMedicines) == 0 {
		return nil, errors.New("CustomerMedicine not found")
	}
	customerMedicine = customerMedicines[0]

	return customerMedicine, nil
}

func (r *CustomerRepositoryPostgresDB) DeleteMedicines(payload *models.CustomersDeleteMedicinesRequestBody) ([]int64, error) {
	q := `
		DELETE FROM customerMedicines
		WHERE id = ANY($1)
	`
	_, err := r.PostgresDB.Exec(
		q,
		pq.Array(payload.IDsInt64),
	)
	if err != nil {
		return nil, err
	}
	return payload.IDsInt64, nil
}

func (r *CustomerRepositoryPostgresDB) QueryOtherAttachments(payload *models.CustomersQueryOtherAttachmentsRequestParams) ([]*domain.CustomerOtherAttachment, int64, error) {
	query := `
		SELECT
		    userOtherAttachments.id,
		    userOtherAttachments.userId,
		    userOtherAttachments.title,
			userOtherAttachments.attachments,
		    userOtherAttachments.created_at,
		    userOtherAttachments.updated_at,
			userOtherAttachments.deleted_at,
			users.id AS userUserId,
		    users.firstName AS userName,
		    users.lastName AS userLastName,
		    users.avatarUrl AS userAvatarUrl
		FROM userOtherAttachments
		LEFT JOIN users ON users.id = userOtherAttachments.userId
	`
	var (
		where []string
		args  []interface{}
	)
	if payload.ID > 0 {
		where = append(where, fmt.Sprintf("userOtherAttachments.id = $%v", len(args)+1))
		args = append(args, payload.ID)
	}
	if payload.UserID > 0 {
		where = append(where, fmt.Sprintf("userOtherAttachments.userId = $%v", len(args)+1))
		args = append(args, payload.UserID)
	}
	if payload.Filters.Title.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.Title.Op, fmt.Sprintf("%v", payload.Filters.Title.Value))
		val := exp.TerIf(opValue.Value == "", "", opValue.Value)
		where = append(where, fmt.Sprintf("userOtherAttachments.title %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	var sorts []string
	if payload.Sorts.Title.Op != "" {
		sorts = append(sorts, fmt.Sprintf(" userOtherAttachments.title %s", payload.Sorts.Title.Op))
	}
	if payload.Sorts.CreatedAt.Op != "" {
		sorts = append(sorts, fmt.Sprintf(" userOtherAttachments.created_at %s", payload.Sorts.CreatedAt.Op))
	}
	if len(sorts) > 0 {
		query += " ORDER BY " + strings.Join(sorts, ",")
	}
	log.Printf("Query for customer other attachments: %s\n", query)

	var count int64
	err := r.PostgresDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count", query), args...).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	limit := exp.TerIf(payload.Limit == 0, 10, payload.Limit)
	page := exp.TerIf(payload.Page == 0, 1, payload.Page)
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT %v", limit)
	query += fmt.Sprintf(" OFFSET %v", offset)

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var (
		userOtherAttachments []*domain.CustomerOtherAttachment
	)
	for rows.Next() {
		var (
			customerOtherAttachment = &domain.CustomerOtherAttachment{}
			title                   sql.NullString
			attachments             sql.NullString
			attachmentsMetadata     []*types.UploadMetadata
			deletedAt               sql.NullTime
			userID                  sql.NullInt64
			userFirstName           sql.NullString
			userLastName            sql.NullString
			userAvatarUrl           sql.NullString
		)
		err := rows.Scan(
			&customerOtherAttachment.ID,
			&customerOtherAttachment.UserID,
			&title,
			&attachments,
			&customerOtherAttachment.CreatedAt,
			&customerOtherAttachment.UpdatedAt,
			&deletedAt,
			&userID,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, 0, err
		}

		if title.Valid {
			customerOtherAttachment.Title = title.String
		}
		if userID.Valid {
			customerOtherAttachment.User = &domain.CustomerOtherAttachmentUser{
				ID: uint(userID.Int64),
			}
			if userFirstName.Valid {
				customerOtherAttachment.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				customerOtherAttachment.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				customerOtherAttachment.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if deletedAt.Valid {
			customerOtherAttachment.DeletedAt = &deletedAt.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in customer other attachments: %d", err, customerOtherAttachment.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.CUSTOMER_BUCKET_NAME[len("maja."):])
				}
			}
			customerOtherAttachment.Attachments = attachmentsMetadata
		}

		// Find customerId
		if err := r.PostgresDB.QueryRow(`SELECT id FROM customers WHERE userId = $1`, customerOtherAttachment.UserID).Scan(&customerOtherAttachment.CustomerID); err != nil {
			log.Printf("failed to find customerId for customer other attachment: %d, %v", customerOtherAttachment.ID, err)
		}

		userOtherAttachments = append(userOtherAttachments, customerOtherAttachment)
	}

	return userOtherAttachments, count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateOtherAttachments(customer *domain.Customer, payload *models.CustomersCreateOtherAttachmentsRequestBody) (*domain.CustomerOtherAttachment, error) {
	// Insert the other attachment
	var otherAttachment domain.CustomerOtherAttachment
	err := r.PostgresDB.QueryRow(`
		INSERT INTO userOtherAttachments (userId, title, created_at, updated_at, deleted_at)
		VALUES ($1, $2, NOW(), NOW(), NULL)
		RETURNING id
	`,
		customer.UserID,
		payload.Title,
	).Scan(
		&otherAttachment.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the other attachment
	results, _, err := r.QueryOtherAttachments(&models.CustomersQueryOtherAttachmentsRequestParams{
		ID:    int(otherAttachment.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	otherAttachment = *results[0]

	// Return the other attachment
	return &otherAttachment, nil
}

func (r *CustomerRepositoryPostgresDB) UpdateCustomerOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CustomerOtherAttachment, error) {
	var otherAttachment domain.CustomerOtherAttachment

	// Current time
	currentTime := time.Now()

	// Find the other attachment by id
	results, _, err := r.QueryOtherAttachments(&models.CustomersQueryOtherAttachmentsRequestParams{
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
	foundOtherAttachment := results[0]
	if foundOtherAttachment == nil {
		return nil, errors.New("no results found")
	}

	// Marshal attachments into JSON format
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)

	// Update the other attachment
	err = r.PostgresDB.QueryRow(`
		UPDATE userOtherAttachments
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundOtherAttachment.ID,
	).Scan(
		&otherAttachment.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the other attachment
	results, _, err = r.QueryOtherAttachments(&models.CustomersQueryOtherAttachmentsRequestParams{
		ID:    int(foundOtherAttachment.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	otherAttachment = *results[0]

	// Return the other attachment
	return &otherAttachment, nil
}

func (r *CustomerRepositoryPostgresDB) UpdateCustomerOtherAttachment(customerOtherAttachment *domain.CustomerOtherAttachment, payload *models.CustomersUpdateOtherAttachmentRequestBody) (*domain.CustomerOtherAttachment, error) {
	var (
		updatedID   int
		currentTime = time.Now()
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE userOtherAttachments SET title = $1, updated_at = $2 WHERE id = $3 RETURNING id
	`,
		payload.Title,
		currentTime,
		customerOtherAttachment.ID,
	).Scan(
		&updatedID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the other attachment
	results, _, err := r.QueryOtherAttachments(&models.CustomersQueryOtherAttachmentsRequestParams{
		ID:    updatedID,
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	otherAttachment := results[0]

	return otherAttachment, nil
}

func (r *CustomerRepositoryPostgresDB) DeleteCustomerOtherAttachments(payload *models.CustomersDeleteOtherAttachmentsRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				userOtherAttachments
			WHERE
			    id = ANY ($1) AND userId = $2
			RETURNING id
		`,
		pq.Int64Array(payload.IDsInt64),
		payload.UserID,
	).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}

func makeRelativesWhereFilters(queries *models.CustomersQueryRelativesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cr.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" cr.customerId = %d ", queries.CustomerID))
		}
		if queries.Filters.CityID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CityID.Op, fmt.Sprintf("%v", queries.Filters.CityID.Value))
			val := exp.TerIf(opValue.Value == "", "", opValue.Value)
			where = append(where, fmt.Sprintf(" cr.cityId %s %s ", opValue.Operator, val))
		}
		if queries.Filters.FirstName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FirstName.Op, fmt.Sprintf("%v", queries.Filters.FirstName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.firstName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.LastName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.LastName.Op, fmt.Sprintf("%v", queries.Filters.LastName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.lastName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.PhoneNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PhoneNumber.Op, fmt.Sprintf("%v", queries.Filters.PhoneNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.phoneNumber %s %s ", opValue.Operator, val))
		}
		if queries.Filters.Relation.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Relation.Op, fmt.Sprintf("%v", queries.Filters.Relation.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.relation %s %s ", opValue.Operator, val))
		}
		if queries.Filters.AddressName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AddressName.Op, fmt.Sprintf("%v", queries.Filters.AddressName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.addressName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.AddressStreet.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AddressStreet.Op, fmt.Sprintf("%v", queries.Filters.AddressStreet.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.addressStreet %s %s ", opValue.Operator, val))
		}
		if queries.Filters.AddressBuildingNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AddressBuildingNumber.Op, fmt.Sprintf("%v", queries.Filters.AddressBuildingNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.addressBuildingNumber %s %s ", opValue.Operator, val))
		}
		if queries.Filters.AddressPostalCode.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AddressPostalCode.Op, fmt.Sprintf("%v", queries.Filters.AddressPostalCode.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cr.addressPostalCode %s %s ", opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryRelatives(queries *models.CustomersQueryRelativesRequestParams) ([]*domain.CustomerRelative, error) {
	q := `
		SELECT
			cr.id,
			cr.customerId,
			cr.addressCityId,
			cr.firstName,
			cr.lastName,
			cr.phoneNumber,
			cr.relation,
			cr.addressName,
			cr.addressStreet,
			cr.addressBuildingNumber,
			cr.addressPostalCode,
			cr.created_at,
			cr.updated_at,
			cr.deleted_at,
			c.id AS customerId,
			u.id AS userId,
			u.firstName AS customerFirstName,
			u.lastName AS customerLastName,
			u.avatarUrl AS customerAvatarUrl,
			cc.id AS customerCityId,
			cc.name AS customerCityName
		FROM customerRelatives cr
		LEFT JOIN customers c ON c.id = cr.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN cities cc ON cc.id = cr.addressCityId
	`
	if queries != nil {
		where := makeRelativesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.id %s ", queries.Sorts.ID.Op))
		}
		if queries.Sorts.FirstName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.firstName %s ", queries.Sorts.FirstName.Op))
		}
		if queries.Sorts.LastName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.lastName %s ", queries.Sorts.LastName.Op))
		}
		if queries.Sorts.PhoneNumber.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.phoneNumber %s ", queries.Sorts.PhoneNumber.Op))
		}
		if queries.Sorts.Relation.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.relation %s ", queries.Sorts.Relation.Op))
		}
		if queries.Sorts.AddressName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.addressName %s ", queries.Sorts.AddressName.Op))
		}
		if queries.Sorts.AddressStreet.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.addressStreet %s ", queries.Sorts.AddressStreet.Op))
		}
		if queries.Sorts.AddressBuildingNumber.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.addressBuildingNumber %s ", queries.Sorts.AddressBuildingNumber.Op))
		}
		if queries.Sorts.AddressPostalCode.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.addressPostalCode %s ", queries.Sorts.AddressPostalCode.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cr.created_at %s ", queries.Sorts.CreatedAt.Op))
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

	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerRelatives []*domain.CustomerRelative
	for rows.Next() {
		var (
			customerRelative domain.CustomerRelative
			deletedAt        sql.NullTime
			customerID       sql.NullInt64
			cityID           sql.NullInt64
			cCustomerID      sql.NullInt64
			userID           sql.NullInt64
			userFirstName    sql.NullString
			userLastName     sql.NullString
			userAvatarUrl    sql.NullString
			cCityID          sql.NullInt64
			addressCityName  sql.NullString
		)
		err := rows.Scan(
			&customerRelative.ID,
			&customerRelative.CustomerID,
			&cityID,
			&customerRelative.FirstName,
			&customerRelative.LastName,
			&customerRelative.PhoneNumber,
			&customerRelative.Relation,
			&customerRelative.AddressName,
			&customerRelative.AddressStreet,
			&customerRelative.AddressBuildingNumber,
			&customerRelative.AddressPostalCode,
			&customerRelative.CreatedAt,
			&customerRelative.UpdatedAt,
			&deletedAt,
			&cCustomerID,
			&userID,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&cCityID,
			&addressCityName,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			customerRelative.DeletedAt = &deletedAt.Time
		}
		if customerID.Valid {
			cid := uint(customerID.Int64)
			customerRelative.CustomerID = cid
		}
		if cityID.Valid {
			cid := uint(cityID.Int64)
			customerRelative.CityID = &cid
		}
		if cCustomerID.Valid {
			cid := uint(cCustomerID.Int64)
			customerRelative.Customer = &domain.CustomerRelativeCustomer{
				ID: cid,
			}
			if userID.Valid {
				uid := uint(userID.Int64)
				customerRelative.Customer.UserID = &uid
			}
			if userFirstName.Valid {
				customerRelative.Customer.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				customerRelative.Customer.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				customerRelative.Customer.AvatarUrl = userAvatarUrl.String
			}
		}
		if cCityID.Valid {
			cid := uint(cCityID.Int64)
			customerRelative.City = &domain.CustomerRelativeCity{
				ID: cid,
			}
			if addressCityName.Valid {
				customerRelative.City.Name = addressCityName.String
			}
		}
		customerRelatives = append(customerRelatives, &customerRelative)
	}

	return customerRelatives, nil
}

func (r *CustomerRepositoryPostgresDB) CountRelatives(queries *models.CustomersQueryRelativesRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cr.id)
		FROM customerRelatives cr
		LEFT JOIN customers c ON c.id = cr.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN cities cc ON cc.id = cr.addressCityId
	`
	if queries != nil {
		where := makeRelativesWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CustomerRepositoryPostgresDB) CreateRelatives(customer *domain.Customer, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	currentTime := time.Now()
	q := `
		INSERT INTO customerRelatives (
			customerId,
			addressCityId,
			firstName,
			lastName,
			phoneNumber,
			relation,
			addressName,
			addressStreet,
			addressBuildingNumber,
			addressPostalCode,
			created_at,
			updated_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`
	var insertedID int
	err := r.PostgresDB.QueryRow(
		q,
		customer.ID,
		payload.CityID,
		payload.FirstName,
		payload.LastName,
		payload.PhoneNumber,
		payload.Relation,
		payload.AddressName,
		payload.AddressStreet,
		payload.AddressBuildingNumber,
		payload.AddressPostalCode,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	// Get customer relatives
	customerRelatives, err := r.QueryRelatives(&models.CustomersQueryRelativesRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerRelatives) == 0 {
		return nil, errors.New("failed to get customer relatives")
	}

	return customerRelatives[0], nil
}

func (r *CustomerRepositoryPostgresDB) UpdateRelative(customerRelative *domain.CustomerRelative, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	currentTime := time.Now()
	q := `
		UPDATE customerRelatives
		SET
			addressCityId = $1,
			firstName = $2,
			lastName = $3,
			phoneNumber = $4,
			relation = $5,
			addressName = $6,
			addressStreet = $7,
			addressBuildingNumber = $8,
			addressPostalCode = $9,
			updated_at = $10
		WHERE id = $11
		RETURNING id
	`
	var updatedID int
	err := r.PostgresDB.QueryRow(
		q,
		payload.CityID,
		payload.FirstName,
		payload.LastName,
		payload.PhoneNumber,
		payload.Relation,
		payload.AddressName,
		payload.AddressStreet,
		payload.AddressBuildingNumber,
		payload.AddressPostalCode,
		currentTime,
		customerRelative.ID,
	).Scan(&updatedID)
	if err != nil {
		return nil, err
	}

	// Get customer relatives
	customerRelatives, err := r.QueryRelatives(&models.CustomersQueryRelativesRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}

	if len(customerRelatives) == 0 {
		return nil, errors.New("failed to get customer relatives")
	}

	return customerRelatives[0], nil
}

func (r *CustomerRepositoryPostgresDB) DeleteRelatives(payload *models.CustomersDeleteRelativesRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				customerRelatives
			WHERE
			    id = ANY ($1) AND customerId = $2
			RETURNING id
		`,
		pq.Int64Array(payload.IDsInt64),
		payload.CustomerID,
	).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}

func makeContractualMobilityRestrictionLogsWhereFilters(queries *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cmrl.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" cmrl.customerId = %d ", queries.CustomerID))
		}
		if queries.Filters.BeforeValue.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.BeforeValue.Op, fmt.Sprintf("%v", queries.Filters.BeforeValue.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cmrl.beforeValue %s %s ", opValue.Operator, val))
		}
		if queries.Filters.AfterValue.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.AfterValue.Op, fmt.Sprintf("%v", queries.Filters.AfterValue.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cmrl.afterValue %s %s ", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cmrl.created_at %s %s ", opValue.Operator, val))
		}
		if queries.Filters.FirstName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FirstName.Op, fmt.Sprintf("%v", queries.Filters.FirstName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.firstName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.LastName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.LastName.Op, fmt.Sprintf("%v", queries.Filters.LastName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.lastName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.FullName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FullName.Op, fmt.Sprintf("%v", queries.Filters.FullName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.firstName %s %s OR u2.lastName %s %s", opValue.Operator, val, opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryContractualMobilityRestrictionLogs(queries *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) ([]*domain.CustomerContractualMobilityRestrictionLog, error) {
	q := `
		SELECT
			cmrl.id,
			cmrl.customerId,
			cmrl.beforeValue,
			cmrl.afterValue,
			cmrl.created_at,
			cmrl.updated_at,
			cmrl.deleted_at,
			c.id AS fkCustomerId,
			u.id AS fkUserId,
			u.firstName AS fkUserFirstName,
			u.lastName AS fkUserLastName,
			u.avatarUrl AS fkUserAvatarUrl,
			u2.id AS fkUser2Id,
			u2.firstName AS fkUser2FirstName,
			u2.lastName AS fkUser2LastName,
			u2.avatarUrl AS fkUser2AvatarUrl
		FROM customerContractualMobilityRestrictionLogs cmrl
		LEFT JOIN customers c ON c.id = cmrl.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN users u2 ON u2.id = cmrl.createdBy
	`
	if queries != nil {
		where := makeContractualMobilityRestrictionLogsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.id %s ", queries.Sorts.ID.Op))
		}
		if queries.Sorts.FirstName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" u2.firstName %s ", queries.Sorts.FirstName.Op))
		}
		if queries.Sorts.LastName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" u2.lastName %s ", queries.Sorts.LastName.Op))
		}
		if queries.Sorts.BeforeValue.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.beforeValue %s ", queries.Sorts.BeforeValue.Op))
		}
		if queries.Sorts.AfterValue.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.afterValue %s ", queries.Sorts.AfterValue.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.created_at %s ", queries.Sorts.CreatedAt.Op))
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

	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerContractualMobilityRestrictionLogs []*domain.CustomerContractualMobilityRestrictionLog
	for rows.Next() {
		var (
			customerContractualMobilityRestrictionLog domain.CustomerContractualMobilityRestrictionLog
			deletedAt                                 sql.NullTime
			customerID                                sql.NullInt64
			beforeValue                               sql.NullString
			afterValue                                sql.NullString
			cCustomerID                               sql.NullInt64
			userID                                    sql.NullInt64
			userFirstName                             sql.NullString
			userLastName                              sql.NullString
			userAvatarUrl                             sql.NullString
			createdByID                               sql.NullInt64
			createdByFirstName                        sql.NullString
			createdByLastName                         sql.NullString
			createdByAvatarUrl                        sql.NullString
		)
		err := rows.Scan(
			&customerContractualMobilityRestrictionLog.ID,
			&customerID,
			&beforeValue,
			&afterValue,
			&customerContractualMobilityRestrictionLog.CreatedAt,
			&customerContractualMobilityRestrictionLog.UpdatedAt,
			&deletedAt,
			&cCustomerID,
			&userID,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&createdByID,
			&createdByFirstName,
			&createdByLastName,
			&createdByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			customerContractualMobilityRestrictionLog.DeletedAt = &deletedAt.Time
		}
		if customerID.Valid {
			customerContractualMobilityRestrictionLog.CustomerID = customerID.Int64
		}
		if beforeValue.Valid {
			customerContractualMobilityRestrictionLog.BeforeValue = &beforeValue.String
		}
		if afterValue.Valid {
			customerContractualMobilityRestrictionLog.AfterValue = &afterValue.String
		}
		if cCustomerID.Valid {
			customerContractualMobilityRestrictionLog.Customer = &domain.CustomerContractualMobilityRestrictionLogCustomer{
				ID: cCustomerID.Int64,
			}
			if userID.Valid {
				customerContractualMobilityRestrictionLog.Customer.UserID = userID.Int64
			}
			if userFirstName.Valid {
				customerContractualMobilityRestrictionLog.Customer.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				customerContractualMobilityRestrictionLog.Customer.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				customerContractualMobilityRestrictionLog.Customer.AvatarUrl = userAvatarUrl.String
			}
		}
		if createdByID.Valid {
			customerContractualMobilityRestrictionLog.CreatedBy = &domain.CustomerContractualMobilityRestrictionLogCreatedBy{
				ID: createdByID.Int64,
			}
			if createdByFirstName.Valid {
				customerContractualMobilityRestrictionLog.CreatedBy.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				customerContractualMobilityRestrictionLog.CreatedBy.LastName = createdByLastName.String
			}
			if createdByAvatarUrl.Valid {
				customerContractualMobilityRestrictionLog.CreatedBy.AvatarUrl = createdByAvatarUrl.String
			}
		}
		customerContractualMobilityRestrictionLogs = append(customerContractualMobilityRestrictionLogs, &customerContractualMobilityRestrictionLog)
	}

	return customerContractualMobilityRestrictionLogs, nil
}

func (r *CustomerRepositoryPostgresDB) CountContractualMobilityRestrictionLogs(queries *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cmrl.id)
		FROM customerContractualMobilityRestrictionLogs cmrl
		LEFT JOIN customers c ON c.id = cmrl.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN users u2 ON u2.id = cmrl.createdBy
	`
	if queries != nil {
		where := makeContractualMobilityRestrictionLogsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CustomerRepositoryPostgresDB) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerAbsence, error) {
	var absence domain.CustomerAbsence

	// Current time
	currentTime := time.Now()

	// Find the absence by id
	results, err := r.QueryAbsences(&models.CustomersQueryAbsencesRequestParams{
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
	foundAbsence := results[0]
	if foundAbsence == nil {
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

	// Update the absence
	err = r.PostgresDB.QueryRow(`
		UPDATE customerAbsences
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundAbsence.ID,
	).Scan(
		&absence.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the absence
	results, err = r.QueryAbsences(&models.CustomersQueryAbsencesRequestParams{
		ID:    int(foundAbsence.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	absence = *results[0]

	// Return the absence
	return &absence, nil
}

func (r *CustomerRepositoryPostgresDB) UpdateMedicineAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerMedicine, error) {
	var medicine domain.CustomerMedicine

	// Current time
	currentTime := time.Now()

	// Find the medicine by id
	results, err := r.QueryMedicines(&models.CustomersQueryMedicinesRequestParams{
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
	foundMedicine := results[0]
	if foundMedicine == nil {
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
		UPDATE customerMedicines
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundMedicine.ID,
	).Scan(
		&medicine.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine
	results, err = r.QueryMedicines(&models.CustomersQueryMedicinesRequestParams{
		ID:    int(foundMedicine.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	medicine = *results[0]

	// Return the medicine
	return &medicine, nil
}

func makeStatusLogsWhereFilters(queries *models.CustomersQueryStatusLogsRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cmrl.id = %d ", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" cmrl.customerId = %d ", queries.CustomerID))
		}
		if queries.Filters.StatusValue.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StatusValue.Op, fmt.Sprintf("%v", queries.Filters.StatusValue.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cmrl.statusValue %s %s ", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" cmrl.created_at %s %s ", opValue.Operator, val))
		}
		if queries.Filters.FirstName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FirstName.Op, fmt.Sprintf("%v", queries.Filters.FirstName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.firstName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.LastName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.LastName.Op, fmt.Sprintf("%v", queries.Filters.LastName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.lastName %s %s ", opValue.Operator, val))
		}
		if queries.Filters.FullName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.FullName.Op, fmt.Sprintf("%v", queries.Filters.FullName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" u2.firstName %s %s OR u2.lastName %s %s", opValue.Operator, val, opValue.Operator, val))
		}
	}
	return where
}

func (r *CustomerRepositoryPostgresDB) QueryStatusLogs(queries *models.CustomersQueryStatusLogsRequestParams) ([]*domain.CustomerStatusLog, error) {
	q := `
		SELECT
			cmrl.id,
			cmrl.customerId,
			cmrl.statusValue,
			cmrl.created_at,
			cmrl.updated_at,
			cmrl.deleted_at,
			c.id AS fkCustomerId,
			u.id AS fkUserId,
			u.firstName AS fkUserFirstName,
			u.lastName AS fkUserLastName,
			u.avatarUrl AS fkUserAvatarUrl,
			u2.id AS fkUser2Id,
			u2.firstName AS fkUser2FirstName,
			u2.lastName AS fkUser2LastName,
			u2.avatarUrl AS fkUser2AvatarUrl
		FROM customerStatusLogs cmrl
		LEFT JOIN customers c ON c.id = cmrl.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN users u2 ON u2.id = cmrl.createdBy
	`
	if queries != nil {
		where := makeStatusLogsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.id %s ", queries.Sorts.ID.Op))
		}
		if queries.Sorts.FirstName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" u2.firstName %s ", queries.Sorts.FirstName.Op))
		}
		if queries.Sorts.LastName.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" u2.lastName %s ", queries.Sorts.LastName.Op))
		}
		if queries.Sorts.StatusValue.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.statusValue %s ", queries.Sorts.StatusValue.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" cmrl.created_at %s ", queries.Sorts.CreatedAt.Op))
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

	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerStatusLogs []*domain.CustomerStatusLog
	for rows.Next() {
		var (
			customerStatusLog  domain.CustomerStatusLog
			deletedAt          sql.NullTime
			customerID         sql.NullInt64
			statusValue        sql.NullString
			cCustomerID        sql.NullInt64
			userID             sql.NullInt64
			userFirstName      sql.NullString
			userLastName       sql.NullString
			userAvatarUrl      sql.NullString
			createdByID        sql.NullInt64
			createdByFirstName sql.NullString
			createdByLastName  sql.NullString
			createdByAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&customerStatusLog.ID,
			&customerID,
			&statusValue,
			&customerStatusLog.CreatedAt,
			&customerStatusLog.UpdatedAt,
			&deletedAt,
			&cCustomerID,
			&userID,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&createdByID,
			&createdByFirstName,
			&createdByLastName,
			&createdByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			customerStatusLog.DeletedAt = &deletedAt.Time
		}
		if customerID.Valid {
			customerStatusLog.CustomerID = customerID.Int64
		}
		if statusValue.Valid {
			customerStatusLog.StatusValue = statusValue.String
		}
		if cCustomerID.Valid {
			customerStatusLog.Customer = &domain.CustomerStatusLogCustomer{
				ID: cCustomerID.Int64,
			}
			if userID.Valid {
				customerStatusLog.Customer.UserID = userID.Int64
			}
			if userFirstName.Valid {
				customerStatusLog.Customer.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				customerStatusLog.Customer.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				customerStatusLog.Customer.AvatarUrl = userAvatarUrl.String
			}
		}
		if createdByID.Valid {
			customerStatusLog.CreatedBy = &domain.CustomerStatusLogCreatedBy{
				ID: createdByID.Int64,
			}
			if createdByFirstName.Valid {
				customerStatusLog.CreatedBy.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				customerStatusLog.CreatedBy.LastName = createdByLastName.String
			}
			if createdByAvatarUrl.Valid {
				customerStatusLog.CreatedBy.AvatarUrl = createdByAvatarUrl.String
			}
		}
		customerStatusLogs = append(customerStatusLogs, &customerStatusLog)
	}

	return customerStatusLogs, nil
}

func (r *CustomerRepositoryPostgresDB) CountStatusLogs(queries *models.CustomersQueryStatusLogsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cmrl.id)
		FROM customerStatusLogs cmrl
		LEFT JOIN customers c ON c.id = cmrl.customerId
		LEFT JOIN users u ON u.id = c.userId
		LEFT JOIN users u2 ON u2.id = cmrl.createdBy
	`
	if queries != nil {
		where := makeStatusLogsWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
