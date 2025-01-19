package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	evaluationConstants "github.com/hoitek/Maja-Service/internal/evaluation/constants"
	"github.com/hoitek/Maja-Service/internal/staff/constants"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"github.com/hoitek/Maja-Service/internal/staff/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type StaffRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewStaffRepositoryPostgresDB(d *sql.DB) *StaffRepositoryPostgresDB {
	return &StaffRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.StaffsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf("staffs.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf("staffs.userId = %d", queries.UserID))
		}
		if queries.Filters.UserId.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.UserId.Op, fmt.Sprintf("%v", queries.Filters.UserId.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("staffs.userId %s %s", opValue.Operator, val))
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
		if queries.Filters.PhoneNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PhoneNumber.Op, queries.Filters.PhoneNumber.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.phone %s %s", opValue.Operator, val))
		}
		if queries.Filters.Team.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Team.Op, queries.Filters.Team.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("teams.name %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *StaffRepositoryPostgresDB) Query(queries *models.StaffsQueryRequestParams) ([]*domain.Staff, error) {
	q := `
		SELECT staffs.*,
		   users.username AS userName,
		   users.firstName AS userFirstName,
		   users.lastName AS userLastName,
		   users.phone AS userPhoneNumber,
		   users.email AS userEmail,
		   users.avatarUrl AS userAvatarUrl,
		   users.workPhoneNumber AS userWorkPhoneNumber,
		   users.gender AS userGender,
		   users.accountNumber AS userAccountNumber,
		   users.telephone AS userTelephone,
           users.registrationNumber AS userRegistrationNumber,
           users.nationalCode AS userNationalCode,
           users.birthDate AS userBirthDate,
		   paymentTypes.id AS paymentTypeId,
		   paymentTypes.name AS paymentTypeName
		FROM staffs
		LEFT JOIN users ON staffs.userId = users.id
		LEFT JOIN paymentTypes ON staffs.paymentTypeId = paymentTypes.id
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
		if queries.Sorts.PhoneNumber.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" users.phone %s", queries.Sorts.PhoneNumber.Op))
		}
		if queries.Sorts.Team.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" teams.name %s", queries.Sorts.Team.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" staffs.created_at %s", queries.Sorts.CreatedAt.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		if limit > -1 {
			q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	}
	q += ";"
	log.Println(q)
	var staffs []*domain.Staff
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var staff domain.Staff
		var (
			userName            sql.NullString
			workPhoneNumber     sql.NullString
			gender              sql.NullString
			accountNumber       sql.NullString
			telephone           sql.NullString
			registrationNumber  sql.NullString
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			userFirstName       sql.NullString
			userLastName        sql.NullString
			userPhoneNumber     sql.NullString
			userEmail           sql.NullString
			userAvatarUrl       sql.NullString
			paymentTypeID       sql.NullInt64
			paymentTypeName     sql.NullString
			jobTitle            sql.NullString
			certificateCode     sql.NullString
			nationalCode        sql.NullString
		)
		var (
			birthDate           sql.NullTime
			vehicleTypes        *json.RawMessage
			vehicleLicenseTypes *json.RawMessage
		)
		err := rows.Scan(
			&staff.ID,
			&staff.UserID,
			&staff.PaymentTypeID,
			&jobTitle,
			&certificateCode,
			&staff.ExperienceAmount,
			&staff.ExperienceAmountUnit,
			&staff.IsSubcontractor,
			&staff.CompanyRegistrationNumber,
			&staff.OrganizationNumber,
			&staff.PercentLengthInContract,
			&staff.HourLengthInContract,
			&staff.Salary,
			&vehicleTypes,
			&vehicleLicenseTypes,
			&staff.TrialTime,
			&attachments,
			&staff.JoinedAt,
			&staff.ContractStartedAt,
			&staff.ContractExpiresAt,
			&staff.CreatedAt,
			&staff.UpdatedAt,
			&staff.DeletedAt,
			&userName,
			&userFirstName,
			&userLastName,
			&userPhoneNumber,
			&userEmail,
			&userAvatarUrl,
			&workPhoneNumber,
			&gender,
			&accountNumber,
			&telephone,
			&registrationNumber,
			&nationalCode,
			&birthDate,
			&paymentTypeID,
			&paymentTypeName,
		)
		if err != nil {
			return nil, err
		}
		if jobTitle.Valid {
			staff.JobTitle = &jobTitle.String
		}
		if certificateCode.Valid {
			staff.CertificateCode = &certificateCode.String
		}
		if paymentTypeID.Valid {
			pid := uint(paymentTypeID.Int64)
			staff.PaymentType = &domain.StaffPaymentTypeRes{
				ID: pid,
			}
			if paymentTypeName.Valid {
				staff.PaymentType.Name = paymentTypeName.String
			}
		}
		if vehicleTypes != nil {
			var vt interface{}
			err = json.Unmarshal(*vehicleTypes, &vt)
			if err != nil {
				return nil, err
			}
			if vt != nil {
				staff.VehicleTypes = vt
			}
		}
		if vehicleLicenseTypes != nil {
			var vlt interface{}
			err = json.Unmarshal(*vehicleLicenseTypes, &vlt)
			if err != nil {
				return nil, err
			}
			if vlt != nil {
				staff.VehicleLicenseTypes = vlt
			}
		}
		staff.User = domain.StaffUser{
			ID: staff.UserID,
		}
		if userFirstName.Valid {
			staff.User.FirstName = userFirstName.String
		}
		if userLastName.Valid {
			staff.User.LastName = userLastName.String
		}
		if userEmail.Valid {
			staff.User.Email = userEmail.String
		}
		if userPhoneNumber.Valid {
			staff.User.Phone = userPhoneNumber.String
		}
		if userName.Valid {
			staff.User.Username = userName.String
		}
		if userAvatarUrl.Valid {
			staff.User.AvatarUrl = userAvatarUrl.String
		}
		if workPhoneNumber.Valid {
			staff.User.WorkPhoneNumber = workPhoneNumber.String
		}
		if gender.Valid {
			staff.User.Gender = gender.String
		}
		if accountNumber.Valid {
			staff.User.AccountNumber = accountNumber.String
		}
		if telephone.Valid {
			staff.User.Telephone = telephone.String
		}
		if registrationNumber.Valid {
			staff.User.RegistrationNumber = registrationNumber.String
		}
		if nationalCode.Valid {
			staff.User.NationalCode = nationalCode.String
		}
		if birthDate.Valid {
			staff.User.BirthDate = &birthDate.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff: %d", err, staff.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			staff.Attachments = attachmentsMetadata
		}
		staffs = append(staffs, &staff)
	}
	if rows != nil {
		rows.Close()
	}

	// Get user roles
	for _, staff := range staffs {
		q := `
			SELECT _roles.id, _roles.name
			FROM usersRoles
			LEFT JOIN _roles ON usersRoles.roleId = _roles.id
			WHERE usersRoles.userId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.UserID)
		if err == nil {
			for rows.Next() {
				var (
					id   uint
					name string
				)
				err := rows.Scan(&id, &name)
				if err == nil {
					staff.User.Roles = append(staff.User.Roles, domain.StaffUserRole{
						ID:   id,
						Name: name,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.User.Roles == nil {
			staff.User.Roles = []domain.StaffUserRole{}
		}
	}

	// Get user language skills
	for _, staff := range staffs {
		q := `
			SELECT languageSkills.id, languageSkills.name
			FROM userLanguageSkills
			LEFT JOIN languageSkills ON userLanguageSkills.languageSkillId = languageSkills.id
			WHERE userLanguageSkills.userId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.UserID)
		if err == nil {
			for rows.Next() {
				var (
					id            uint
					languageSkill string
				)
				err := rows.Scan(&id, &languageSkill)
				if err == nil {
					staff.User.LanguageSkills = append(staff.User.LanguageSkills, domain.StaffUserLanguageSkill{
						ID:   id,
						Name: languageSkill,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.User.LanguageSkills == nil {
			staff.User.LanguageSkills = []domain.StaffUserLanguageSkill{}
		}
	}

	// Get staff limitations
	for _, staff := range staffs {
		q := `
			SELECT limitations.id AS lid, limitations.name AS lname, staffLimitations.description AS sdescription
			FROM staffLimitations
			LEFT JOIN limitations ON staffLimitations.limitationId = limitations.id
			WHERE staffLimitations.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id          uint
					name        string
					description string
				)
				err := rows.Scan(&id, &name, &description)
				if err == nil {
					staff.Limitations = append(staff.Limitations, sharedmodels.SharedLimitation{
						ID:          id,
						Name:        name,
						Description: description,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.Limitations == nil {
			staff.Limitations = []sharedmodels.SharedLimitation{}
		}
	}

	// Get Addresses for each staff
	for _, staff := range staffs {
		q := `
			SELECT id, name, city, zipCode, state
			FROM addresses
			WHERE addresses.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id      int
					name    string
					city    string
					zipCode string
					state   string
				)
				err := rows.Scan(&id, name, &city, &zipCode, &state)
				if err == nil {
					staff.Addresses = append(staff.Addresses, domain.StaffAddress{
						ID:      uint(id),
						Name:    name,
						City:    city,
						ZipCode: zipCode,
						State:   state,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.Addresses == nil {
			staff.Addresses = []domain.StaffAddress{}
		}
	}

	// Get sections for each staff
	for _, staff := range staffs {
		q := `
			SELECT sections.id, sections.name
			FROM staffSections
			LEFT JOIN sections ON staffSections.sectionId = sections.id
			WHERE staffSections.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id   int
					name string
				)
				err := rows.Scan(&id, &name)
				if err == nil {
					staff.Sections = append(staff.Sections, domain.StaffSectionRes{
						ID:   uint(id),
						Name: name,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.Sections == nil {
			staff.Sections = []domain.StaffSectionRes{}
		}
	}

	// Get licenses for each staff
	for _, staff := range staffs {
		q := `
			SELECT licenses.id, licenses.name, staffLicenses.id, staffLicenses.expire_date
			FROM staffLicenses
			LEFT JOIN licenses ON staffLicenses.licenseId = licenses.id
			WHERE staffLicenses.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					licenseId      int
					staffLicenseId int
					name           string
					expireDate     *sql.NullTime
				)
				err := rows.Scan(&licenseId, &name, &staffLicenseId, &expireDate)
				if err == nil {
					newLicense := domain.StaffLicensesRes{
						ID:        uint(staffLicenseId),
						StaffID:   staff.ID,
						LicenseID: uint(licenseId),
						License: domain.StaffLicensesResLicense{
							ID:   uint(licenseId),
							Name: name,
						},
					}
					if expireDate.Valid {
						newLicense.ExpireDate = &expireDate.Time
					}
					staff.Licenses = append(staff.Licenses, newLicense)
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.Licenses == nil {
			staff.Licenses = []domain.StaffLicensesRes{}
		}
	}

	// Get contract types for each staff
	for _, staff := range staffs {
		q := `
			SELECT contractTypes.id, contractTypes.name
			FROM staffContractTypes
			LEFT JOIN contractTypes ON staffContractTypes.contractTypeId = contractTypes.id
			WHERE staffContractTypes.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id   int
					name string
				)
				err := rows.Scan(&id, &name)
				if err == nil {
					staff.ContractTypes = append(staff.ContractTypes, domain.StaffContractTypeRes{
						ID:   uint(id),
						Name: name,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.ContractTypes == nil {
			staff.ContractTypes = []domain.StaffContractTypeRes{}
		}
	}

	// Get shift types for each staff
	for _, staff := range staffs {
		q := `
			SELECT shiftTypes.id, shiftTypes.name
			FROM staffShiftTypes
			LEFT JOIN shiftTypes ON staffShiftTypes.shiftTypeId = shiftTypes.id
			WHERE staffShiftTypes.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id   int
					name string
				)
				err := rows.Scan(&id, &name)
				if err == nil {
					staff.ShiftTypes = append(staff.ShiftTypes, domain.StaffShiftTypeRes{
						ID:   uint(id),
						Name: name,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.ShiftTypes == nil {
			staff.ShiftTypes = []domain.StaffShiftTypeRes{}
		}
	}

	// Get absences for each staff
	for _, staff := range staffs {
		q := `
			SELECT
				staffAbsences.id,
				staffAbsences.staffId,
				staffAbsences.start_date,
				staffAbsences.end_date,
				staffAbsences.reason,
				staffAbsences.status,
				staffAbsences.statusBy,
				staffAbsences.status_at,
				users.firstName AS statusByName,
				users.lastName AS statusByLastName,
				users.avatarUrl AS statusByAvatarUrl
			FROM staffAbsences
			LEFT JOIN users ON users.id = staffAbsences.statusBy
			WHERE staffAbsences.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					staffAbsence  = &domain.StaffAbsenceRes{}
					reason        sql.NullString
					status        sql.NullString
					statusBy      sql.NullInt64
					statusAt      sql.NullTime
					userFirstName sql.NullString
					userLastName  sql.NullString
					userAvatarUrl sql.NullString
				)
				err := rows.Scan(
					&staffAbsence.ID,
					&staffAbsence.StaffID,
					&staffAbsence.StartDate,
					&staffAbsence.EndDate,
					&reason,
					&status,
					&statusBy,
					&statusAt,
					&userFirstName,
					&userLastName,
					&userAvatarUrl,
				)
				if err != nil {
					return nil, err
				}
				if reason.Valid {
					staffAbsence.Reason = &reason.String
				}
				if status.Valid {
					staffAbsence.Status = &status.String
				}
				if statusBy.Valid {
					staffAbsence.StatusBy = &domain.StaffAbsencesQueryResStatusBy{
						ID: uint(statusBy.Int64),
					}
					if userFirstName.Valid {
						staffAbsence.StatusBy.FirstName = userFirstName.String
					}
					if userLastName.Valid {
						staffAbsence.StatusBy.LastName = userLastName.String
					}
					if userAvatarUrl.Valid {
						staffAbsence.StatusBy.AvatarUrl = userAvatarUrl.String
					}
				}
				if statusAt.Valid {
					staffAbsence.StatusAt = &statusAt.Time
				}
				staff.Absences = append(staff.Absences, *staffAbsence)
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.Absences == nil {
			staff.Absences = []domain.StaffAbsenceRes{}
		}
	}

	// Get staffTypes for each staff
	for _, staff := range staffs {
		q := `
			SELECT staffTypes.id, staffTypes.name
			FROM staffStaffTypes
			LEFT JOIN staffTypes ON staffStaffTypes.staffTypeId = staffTypes.id
			WHERE staffStaffTypes.staffId = $1
		`
		rows, err := r.PostgresDB.Query(q, staff.ID)
		if err == nil {
			for rows.Next() {
				var (
					id   int
					name string
				)
				err := rows.Scan(&id, &name)
				if err == nil {
					staff.StaffTypes = append(staff.StaffTypes, domain.StaffTypesRes{
						ID:   uint(id),
						Name: name,
					})
				}
			}
		}
		if rows != nil {
			rows.Close()
		}
		if staff.StaffTypes == nil {
			staff.StaffTypes = []domain.StaffTypesRes{}
		}
	}

	// Get grace count for each staff
	for _, staff := range staffs {
		q := `
			SELECT COUNT(*)
			FROM evaluations
			WHERE staffId = $1 and evaluationType = $2
		`
		err := r.PostgresDB.QueryRow(q, staff.ID, evaluationConstants.EVALUATION_TYPE_GRACE).Scan(&staff.Grace)
		if err != nil {
			log.Printf("Error getting grace count for staff %d: %s", staff.ID, err)
			continue
		}
	}

	// Get warning count for each staff
	for _, staff := range staffs {
		q := `
			SELECT COUNT(*)
			FROM evaluations
			WHERE staffId = $1 and evaluationType = $2
		`
		err := r.PostgresDB.QueryRow(q, staff.ID, evaluationConstants.EVALUATION_TYPE_WARNING).Scan(&staff.Warning)
		if err != nil {
			log.Printf("Error getting warning count for staff %d: %s", staff.ID, err)
			continue
		}
	}

	// Get attention count for each staff
	for _, staff := range staffs {
		q := `
			SELECT COUNT(*)
			FROM evaluations
			WHERE staffId = $1 and evaluationType = $2
		`
		err := r.PostgresDB.QueryRow(q, staff.ID, evaluationConstants.EVALUATION_TYPE_ATTENTION).Scan(&staff.Attention)
		if err != nil {
			log.Printf("Error getting attention count for staff %d: %s", staff.ID, err)
			continue
		}
	}

	return staffs, nil
}

func (r *StaffRepositoryPostgresDB) Count(queries *models.StaffsQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(staffs.id)
		FROM staffs
		LEFT JOIN users ON staffs.userId = users.id
		LEFT JOIN paymentTypes ON staffs.paymentTypeId = paymentTypes.id
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
		if queries.Sorts.PhoneNumber.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" users.phone %s", queries.Sorts.PhoneNumber.Op))
		}
		if queries.Sorts.Team.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" teams.name %s", queries.Sorts.Team.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
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

func (r *StaffRepositoryPostgresDB) Delete(payload *models.StaffsDeleteRequestBody) ([]int64, error) {
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
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM staffs
			WHERE id = ANY ($1)
			RETURNING id
		`, pq.Int64Array(ids),
	).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return ids, nil
}

func (r *StaffRepositoryPostgresDB) DeleteLicenses(payload *models.StaffsDeleteLicensesRequestBody) ([]int64, *domain.Staff, error) {
	var ids []int64
	idsStr := strings.Split(payload.IDs, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, int64(id))
	}
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				staffLicenses
			WHERE
			    id = ANY ($1) AND staffId = $2
			RETURNING id
		`,
		pq.Int64Array(ids),
		payload.StaffID,
	).Scan(&rowsAffected)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil, err
		}
	}

	// Get the staff
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			ID: payload.StaffID,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	if len(staffs) == 0 {
		return ids, nil, nil
	}

	// Set the staff
	staff := staffs[0]

	return ids, staff, nil
}

func (r *StaffRepositoryPostgresDB) CreateOrUpdateContract(payload *models.StaffsCreateOrUpdateContractRequestBody) (*domain.Staff, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	var staff domain.Staff

	// Current time
	currentTime := time.Now()

	var (
		joinedAt          *time.Time
		contractStartedAt *time.Time
		contractExpiresAt *time.Time
		trialTime         *time.Time
	)

	// Parse joinedAt
	if payload.JoinedAt != "" {
		ja, err := utils.TryParseToDateTime(payload.JoinedAt)
		if err != nil {
			return nil, err
		}
		joinedAt = &ja
	}

	// Parse contractStartedAt
	if payload.ContractStartedAt != "" {
		csa, err := utils.TryParseToDateTime(payload.ContractStartedAt)
		if err != nil {
			return nil, err
		}
		contractStartedAt = &csa
	}

	// Parse contractExpiresAt
	if payload.ContractExpiresAt != "" {
		cea, err := utils.TryParseToDateTime(payload.ContractExpiresAt)
		if err != nil {
			return nil, err
		}
		contractExpiresAt = &cea
	}

	// Parse trialTime
	if payload.TrialTime != "" {
		qt, err := utils.TryParseToDateTime(payload.TrialTime)
		if err != nil {
			return nil, err
		}
		trialTime = &qt
	}

	// Get the staff by staff id or create a new one
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			UserID: payload.UserID,
		},
	)
	if err != nil {
		return nil, err
	}
	var isSubcontractor bool
	if payload.CompanyRegistrationNumber != "" {
		isSubcontractor = true
	}
	var paymentTypeID *int
	if isSubcontractor {
		paymentTypeID = &payload.PaymentTypeAsMetadata.ID
	}
	if len(staffs) == 0 {
		err := r.PostgresDB.QueryRow(`
			INSERT INTO staffs (
				userId,
				companyRegistrationNumber,
				isSubcontractor,
				paymentTypeId,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, payload.UserID, payload.CompanyRegistrationNumber, isSubcontractor, paymentTypeID, currentTime, currentTime).Scan(&staff.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		payload.ID = int(staff.ID)
	} else {
		staff.ID = staffs[0].ID
		payload.ID = int(staff.ID)
	}

	// Update the staff
	err = tx.QueryRow(
		`
		UPDATE staffs
		SET
		    joined_at = $1,
			contract_started_at = $2,
			contract_expires_at = $3,
		    trial_time = $4,
		    organizationNumber = $5,
		    jobTitle = $6,
		    certificateCode = $7,
		    percentLengthInContract = $8,
		    hourLengthInContract = $9,
		    salary = $10,
		    experienceAmount = $11,
		    experienceAmountUnit = $12,
		    updated_at = $13,
			companyRegistrationNumber = $14,
			isSubcontractor = $15,
			paymentTypeId = $16
		WHERE id = $17
		RETURNING id
	`, joinedAt,
		contractStartedAt,
		contractExpiresAt,
		trialTime,
		payload.OrganizationNumber,
		payload.JobTitle,
		payload.CertificateCode,
		payload.PercentLengthInContract,
		payload.HourLengthInContract,
		payload.Salary,
		payload.ExperienceAmount,
		payload.ExperienceAmountUnit,
		currentTime,
		payload.CompanyRegistrationNumber,
		isSubcontractor,
		paymentTypeID,
		staff.ID,
	).Scan(
		&staff.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	staff.IsSubcontractor = isSubcontractor

	// Delete staff section
	_, err = tx.Exec(
		`
		DELETE FROM staffSections WHERE staffId = $1
	`, payload.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert staff section
	for _, sectionID := range payload.SectionIDs {
		_, err = tx.Exec(
			`
			INSERT INTO staffSections (staffId, sectionId) VALUES ($1, $2)
		`, payload.ID, sectionID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete staff shift type
	_, err = tx.Exec(
		`
		DELETE FROM staffShiftTypes WHERE staffId = $1 AND shiftTypeId = ANY ($2)
	`, payload.ID, pq.Int64Array(payload.ShiftTypeIDs),
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert staff shift type
	for _, shiftTypeID := range payload.ShiftTypeIDs {
		_, err = tx.Exec(
			`
			INSERT INTO staffShiftTypes (staffId, shiftTypeId) VALUES ($1, $2)
		`, payload.ID, shiftTypeID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete staff contract type
	_, err = tx.Exec(
		`
		DELETE FROM staffContractTypes WHERE staffId = $1;
	`, payload.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert staff contract type
	for _, contractTypeID := range payload.ContractTypeIDs {
		_, err = tx.Exec(
			`
			INSERT INTO staffContractTypes (staffId, contractTypeId) VALUES ($1, $2);
		`, payload.ID, contractTypeID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete staff types
	_, err = tx.Exec(
		`
		DELETE FROM staffStaffTypes WHERE staffId = $1;
	`, payload.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert staff types
	for _, staffTypeID := range payload.StaffTypeIDs {
		_, err = tx.Exec(
			`
			INSERT INTO staffStaffTypes (staffId, staffTypeId) VALUES ($1, $2)
		`, payload.ID, staffTypeID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete userRoles
	_, err = tx.Exec(
		`
		DELETE FROM usersRoles WHERE userId = $1 AND roleId = ANY ($2)
	`, payload.UserID, pq.Int64Array(payload.RoleIDs))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert userRoles
	for _, roleID := range payload.RoleIDs {
		_, err = tx.Exec(
			`
			INSERT INTO usersRoles (userId, roleId) VALUES ($1, $2)
		`, payload.UserID, roleID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create staff chats if not exists
	if err := r.CreateBaseChatForAllStaffsIfNotExists(); err != nil {
		log.Printf("failed to create base chat for all staffs: %v\n", err)
	}

	// Get the staff
	staffs, err = r.Query(
		&models.StaffsQueryRequestParams{
			ID: int(staff.ID),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(staffs) == 0 {
		return nil, errors.New("no staff found")
	}

	// Set the staff
	staff = *staffs[0]

	// Return the staff
	return &staff, nil
}

func (r *StaffRepositoryPostgresDB) HasLicense(staffID uint, licenseID uint) (bool, error) {
	var count int
	err := r.PostgresDB.QueryRow(
		`
		SELECT COUNT(*) FROM staffLicenses WHERE staffId = $1 AND licenseId = $2
	`, staffID, licenseID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *StaffRepositoryPostgresDB) CreateLicenses(staff *domain.Staff, payload *models.StaffsCreateLicensesRequestBody) (*domain.Staff, *domain.StaffLicensesRes, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// Get current time
	currentTime := time.Now()

	var expireDate *time.Time
	if payload.ExpireDate != "" {
		exDate, err := utils.TryParseToDateTime(payload.ExpireDate)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
		expireDate = &exDate
	}

	// Create staff license
	var (
		insertedId uint
	)
	err = tx.QueryRow(
		`
		INSERT INTO staffLicenses (staffId, licenseId, expire_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, staff.ID, payload.License.ID, expireDate, currentTime, currentTime).Scan(&insertedId)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	var (
		staffLicenseRes domain.StaffLicensesRes
		expireDateRes   sql.NullTime
	)
	err = tx.QueryRow(
		`
		SELECT np.id, np.expire_date, p.id, p.name FROM staffLicenses np
		INNER JOIN licenses p ON p.id = np.licenseId
		WHERE np.id = $1
	`, insertedId).Scan(
		&staffLicenseRes.ID,
		&expireDateRes,
		&staffLicenseRes.License.ID,
		&staffLicenseRes.License.Name,
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	if expireDateRes.Valid {
		staffLicenseRes.ExpireDate = &expireDateRes.Time
	}
	staffLicenseRes.StaffID = staff.ID
	staffLicenseRes.LicenseID = staffLicenseRes.License.ID

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// Get the staff
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			ID: int(staff.ID),
		},
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	if len(staffs) == 0 {
		tx.Rollback()
		return nil, nil, errors.New("no staff found")
	}

	// Set the staff
	staff = staffs[0]

	// Return the staff
	return staff, &staffLicenseRes, nil
}

func (r *StaffRepositoryPostgresDB) CreateAbsences(staff *domain.Staff, payload *models.StaffsCreateAbsencesRequestBody) (*domain.StaffAbsencesQueryRes, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Get current time
	currentTime := time.Now()

	var (
		startDate *time.Time
		endDate   *time.Time
		reason    *string
		status    = constants.ABSENCE_STATUS_PENDING
	)

	if payload.StartDate != "" {
		sDate, err := utils.TryParseToDateTime(payload.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &sDate
	}
	if payload.EndDate != "" {
		eDate, err := utils.TryParseToDateTime(payload.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &eDate
	}
	if payload.Reason != "" {
		reason = &payload.Reason
	}

	// Create staff absence
	var insertedId int64
	err = tx.QueryRow(
		`
		INSERT INTO staffAbsences (staffId, start_date, end_date, reason, status, statusBy, status_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, staff.ID, startDate, endDate, reason, status, nil, nil, currentTime, currentTime).Scan(&insertedId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get the staff
	absences, _, err := r.QueryAbsences(
		&models.StaffsQueryAbsencesRequestParams{
			ID: int(insertedId),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(absences) == 0 {
		return nil, errors.New("no absence found")
	}

	// Set the absence
	absence := absences[0]

	// Return the staff
	return absence, nil
}

func (r *StaffRepositoryPostgresDB) CreateEmptyStaffForUser(userId int) (*domain.Staff, error) {
	// Get the staff
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			UserID: userId,
		},
	)
	if err != nil {
		return nil, err
	}

	// Return the staff if exists
	if len(staffs) > 0 {
		return staffs[0], nil
	}

	// Create staff
	_, err = r.PostgresDB.Exec(`
		INSERT INTO staffs (userId) VALUES ($1)
	`, userId)
	if err != nil {
		return nil, err
	}

	// Create staff chats if not exists
	if err := r.CreateBaseChatForAllStaffsIfNotExists(); err != nil {
		log.Printf("failed to create base chat for all staffs: %v\n", err)
	}

	// Get the staff
	staffs, err = r.Query(
		&models.StaffsQueryRequestParams{
			UserID: userId,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(staffs) == 0 {
		return nil, errors.New("no staff found")
	}
	return staffs[0], nil
}

func (r *StaffRepositoryPostgresDB) GenerateStaffOrganizationNumber() (string, error) {
	var (
		min = 1000
		max = 9999
	)
	randomNumber := utils.GenerateRandomNumber(min, max)
	year := fmt.Sprintf("%v", time.Now().Year())[2:]

	organizationNumber := fmt.Sprintf("%v%v", year, randomNumber)
	for {
		var count int
		err := r.PostgresDB.QueryRow(`
			SELECT COUNT(*) FROM staffs WHERE organizationNumber = $1
		`, organizationNumber).Scan(&count)
		if err != nil {
			return "", err
		}
		if count == 0 {
			break
		}
		randomNumber = utils.GenerateRandomNumber(min, max)
		organizationNumber = fmt.Sprintf("%v%v", time.Now().Year(), randomNumber)
	}
	return organizationNumber, nil
}

func (r *StaffRepositoryPostgresDB) QueryLicenses(payload *models.StaffsQueryLicensesRequestParams) ([]*domain.StaffLicensesRes, int64, error) {
	query := `
		SELECT staffLicenses.id, staffLicenses.licenseId, staffLicenses.staffId, staffLicenses.expire_date, staffLicenses.attachments, licenses.name
		FROM staffLicenses
		LEFT JOIN licenses ON licenses.id = staffLicenses.licenseId
	`
	var (
		where []string
		args  []interface{}
	)
	if payload.ID > 0 {
		where = append(where, fmt.Sprintf("staffLicenses.id = $%v", len(args)+1))
		args = append(args, payload.ID)
	}
	if payload.Filters.StaffID.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.StaffID.Op, fmt.Sprintf("%v", payload.Filters.StaffID.Value))
		val := exp.TerIf(opValue.Value == "", "", opValue.Value)
		where = append(where, fmt.Sprintf("staffLicenses.staffId %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if payload.Filters.Name.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.Name.Op, fmt.Sprintf("%v", payload.Filters.Name.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("licenses.name %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if payload.Filters.CreatedAt.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.CreatedAt.Op, fmt.Sprintf("%v", payload.Filters.CreatedAt.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("staffLicenses.created_at %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	var sorts []string
	if payload.Sorts.CreatedAt.Op != "" {
		sorts = append(sorts, fmt.Sprintf(" staffLicenses.created_at %s", payload.Sorts.CreatedAt.Op))
	}
	if len(sorts) > 0 {
		query += " ORDER BY " + strings.Join(sorts, ",")
	}
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

	var staffLicenses []*domain.StaffLicensesRes
	for rows.Next() {
		var (
			expireDate          sql.NullTime
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
		)
		var staffLicense domain.StaffLicensesRes
		err := rows.Scan(
			&staffLicense.ID,
			&staffLicense.LicenseID,
			&staffLicense.StaffID,
			&expireDate,
			&attachments,
			&staffLicense.License.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		staffLicense.License.ID = staffLicense.LicenseID
		if expireDate.Valid {
			staffLicense.ExpireDate = &expireDate.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff licences: %d", err, staffLicense.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			staffLicense.Attachments = attachmentsMetadata
		}
		staffLicenses = append(staffLicenses, &staffLicense)
	}

	return staffLicenses, count, nil
}

func (r *StaffRepositoryPostgresDB) QueryAbsences(payload *models.StaffsQueryAbsencesRequestParams) ([]*domain.StaffAbsencesQueryRes, int64, error) {
	query := `
		SELECT
		    staffAbsences.id,
		    staffAbsences.staffId,
		    staffAbsences.start_date,
		    staffAbsences.end_date,
		    staffAbsences.reason,
			staffAbsences.attachments,
		    staffAbsences.status,
		    staffAbsences.statusBy,
		    staffAbsences.status_at,
		    users.firstName AS statusByName,
		    users.lastName AS statusByLastName,
		    users.avatarUrl AS statusByAvatarUrl
		FROM staffAbsences
		LEFT JOIN users ON users.id = staffAbsences.statusBy
	`
	var (
		where []string
		args  []interface{}
	)
	if payload.ID > 0 {
		where = append(where, fmt.Sprintf("staffAbsences.id = $%v", len(args)+1))
		args = append(args, payload.ID)
	}
	if payload.StaffID > 0 {
		where = append(where, fmt.Sprintf("staffAbsences.staffId = $%v", len(args)+1))
		args = append(args, payload.StaffID)
	}
	if payload.Filters.StartDate.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.StartDate.Op, fmt.Sprintf("%v", payload.Filters.StartDate.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("staffAbsences.start_date %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if payload.Filters.EndDate.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.EndDate.Op, fmt.Sprintf("%v", payload.Filters.EndDate.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("staffAbsences.end_date %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if payload.Filters.Reason.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.Reason.Op, fmt.Sprintf("%v", payload.Filters.Reason.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("staffAbsences.reason %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if payload.Filters.Status.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.Status.Op, fmt.Sprintf("%v", payload.Filters.Status.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("staffAbsences.status %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

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
		staffAbsences []*domain.StaffAbsencesQueryRes
	)
	for rows.Next() {
		var (
			staffAbsence        = &domain.StaffAbsencesQueryRes{}
			reason              sql.NullString
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			status              sql.NullString
			statusBy            sql.NullInt64
			statusAt            sql.NullTime
			userFirstName       sql.NullString
			userLastName        sql.NullString
			userAvatarUrl       sql.NullString
		)
		err := rows.Scan(
			&staffAbsence.ID,
			&staffAbsence.StaffID,
			&staffAbsence.StartDate,
			&staffAbsence.EndDate,
			&reason,
			&attachments,
			&status,
			&statusBy,
			&statusAt,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, 0, err
		}
		if reason.Valid {
			staffAbsence.Reason = &reason.String
		}
		if status.Valid {
			staffAbsence.Status = &status.String
		}
		if statusBy.Valid {
			staffAbsence.StatusBy = &domain.StaffAbsencesQueryResStatusBy{
				ID: uint(statusBy.Int64),
			}
			if userFirstName.Valid {
				staffAbsence.StatusBy.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				staffAbsence.StatusBy.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				staffAbsence.StatusBy.AvatarUrl = userAvatarUrl.String
			}
		}
		if statusAt.Valid {
			staffAbsence.StatusAt = &statusAt.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff absences: %d", err, staffAbsence.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			staffAbsence.Attachments = attachmentsMetadata
		}
		staffAbsences = append(staffAbsences, staffAbsence)
	}

	return staffAbsences, count, nil
}

func (r *StaffRepositoryPostgresDB) FindStaffLicenseByID(id int) (*domain.StaffLicensesRes, error) {
	var (
		expireDate   sql.NullTime
		staffLicense domain.StaffLicensesRes
	)
	err := r.PostgresDB.QueryRow(`
		SELECT staffLicenses.id, staffLicenses.staffId, staffLicenses.licenseId, licenses.name, staffLicenses.expire_date
		FROM staffLicenses
		LEFT JOIN licenses ON licenses.id = staffLicenses.licenseId
		WHERE staffLicenses.id = $1
	`, id).Scan(
		&staffLicense.ID,
		&staffLicense.StaffID,
		&staffLicense.LicenseID,
		&staffLicense.License.Name,
		&expireDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	staffLicense.License.ID = staffLicense.LicenseID
	if expireDate.Valid {
		staffLicense.ExpireDate = &expireDate.Time
	}
	return &staffLicense, nil
}

func (r *StaffRepositoryPostgresDB) HasLicenseExcept(staffID uint, licenseID uint, staffLicenseID uint) (bool, error) {
	var count int
	err := r.PostgresDB.QueryRow(`
		SELECT COUNT(*) FROM staffLicenses WHERE staffId = $1 AND licenseId = $2 AND id != $3
	`, staffID, licenseID, staffLicenseID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *StaffRepositoryPostgresDB) UpdateLicense(staffLicense *domain.StaffLicensesRes, payload *models.StaffsUpdateLicenseRequestBody) (*domain.Staff, *domain.StaffLicensesRes, error) {
	var expireDate *time.Time
	if payload.ExpireDate != "" {
		exDate, err := utils.TryParseToDateTime(payload.ExpireDate)
		if err != nil {
			return nil, nil, err
		}
		expireDate = &exDate
	}

	var (
		id int
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE staffLicenses SET licenseId = $1, expire_date = $2 WHERE id = $3 RETURNING id
	`,
		payload.License.ID,
		expireDate,
		staffLicense.ID,
	).Scan(
		&id,
	)
	if err != nil {
		return nil, nil, err
	}

	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			ID: int(staffLicense.StaffID),
		},
	)
	if err != nil {
		return nil, nil, err
	}
	if len(staffs) == 0 {
		return nil, nil, errors.New("no staff found")
	}
	staffLicense = &domain.StaffLicensesRes{
		ID:        staffLicense.ID,
		StaffID:   staffLicense.StaffID,
		LicenseID: uint(payload.License.ID),
		License: domain.StaffLicensesResLicense{
			ID:   uint(payload.License.ID),
			Name: staffLicense.License.Name,
		},
		ExpireDate: expireDate,
	}
	return staffs[0], staffLicense, nil
}

func (r *StaffRepositoryPostgresDB) FindStaffAbsenceByID(id int) (*domain.StaffAbsenceRes, error) {
	var (
		staffAbsence domain.StaffAbsenceRes
		endDate      sql.NullTime
		reason       sql.NullString
	)
	err := r.PostgresDB.QueryRow(`
		SELECT staffAbsences.id, staffAbsences.staffId, staffAbsences.start_date, staffAbsences.end_date, staffAbsences.reason
		FROM staffAbsences
		WHERE staffAbsences.id = $1
	`, id).Scan(
		&staffAbsence.ID,
		&staffAbsence.StaffID,
		&staffAbsence.StartDate,
		&endDate,
		&reason,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if endDate.Valid {
		staffAbsence.EndDate = &endDate.Time
	}
	if reason.Valid {
		staffAbsence.Reason = &reason.String
	}
	return &staffAbsence, nil
}

func (r *StaffRepositoryPostgresDB) UpdateAbsence(staffAbsence *domain.StaffAbsenceRes, payload *models.StaffsUpdateAbsenceRequestBody) (*domain.Staff, *domain.StaffAbsenceRes, error) {
	var (
		endDate *time.Time
		reason  *string
	)
	startDate, err := utils.TryParseToDateTime(payload.StartDate)
	if err != nil {
		return nil, nil, err
	}
	if payload.EndDate != "" {
		ed, err := utils.TryParseToDateTime(payload.EndDate)
		if err != nil {
			return nil, nil, err
		}
		endDate = &ed
	}
	if payload.Reason != "" {
		reason = &payload.Reason
	}

	var (
		id       int
		status   = constants.ABSENCE_STATUS_PENDING
		statusAt = time.Now()
	)
	err = r.PostgresDB.QueryRow(`
		UPDATE staffAbsences SET start_date = $1, end_date = $2, reason = $3, status = $4, status_at = $5 WHERE id = $6 RETURNING id
	`,
		startDate,
		endDate,
		reason,
		status,
		statusAt,
		staffAbsence.ID,
	).Scan(
		&id,
	)
	if err != nil {
		return nil, nil, err
	}

	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			ID: int(staffAbsence.StaffID),
		},
	)
	if err != nil {
		return nil, nil, err
	}
	if len(staffs) == 0 {
		return nil, nil, errors.New("no staff found")
	}
	staffAbsence = &domain.StaffAbsenceRes{
		ID:        uint(id),
		StaffID:   staffAbsence.StaffID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    reason,
		Status:    &status,
		StatusBy:  nil,
		StatusAt:  &statusAt,
	}
	return staffs[0], staffAbsence, nil
}

func (r *StaffRepositoryPostgresDB) DeleteAbsences(payload *models.StaffsDeleteAbsencesRequestBody) ([]int64, *domain.Staff, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				staffAbsences
			WHERE
			    id = ANY ($1) AND staffId = $2
			RETURNING id
		`,
		pq.Int64Array(payload.IDsInt64),
		payload.StaffID,
	).Scan(&rowsAffected)
	if err != nil {
		return nil, nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, nil, errors.New("no rows affected")
	}

	// Get the staff
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			ID: payload.StaffID,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	if len(staffs) == 0 {
		return nil, nil, errors.New("no staff found")
	}

	// Set the staff
	staff := staffs[0]

	return payload.IDsInt64, staff, nil
}

func (r *StaffRepositoryPostgresDB) QueryProfile(query *models.StaffsQueryProfileRequestParams) (*domain.StaffProfile, error) {
	var (
		staffProfile domain.StaffProfile
	)
	return &staffProfile, nil
}

func (r *StaffRepositoryPostgresDB) CreateStaff(userId int64, payload *models.StaffsCreatePersonalInfoRequestBody) (*domain.Staff, error) {
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			UserID: int(userId),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(staffs) > 0 {
		return staffs[0], nil
	}
	// Create vehicleTypes jsonb
	vehicleTypesBytes, err := json.Marshal(payload.VehicleTypesAsStringSlice)
	if err != nil {
		return nil, err
	}
	vehicleTypesJSON := string(vehicleTypesBytes)

	// Create vehicleLicenseTypes jsonb
	vehicleLicenseTypesBytes, err := json.Marshal(payload.VehicleLicenseTypesAsStringSlice)
	if err != nil {
		return nil, err
	}
	vehicleLicenseTypesJSON := string(vehicleLicenseTypesBytes)

	// Insert into staffs
	var insertedId int64
	err = r.PostgresDB.QueryRow(`
		INSERT INTO staffs (userId, vehicleTypes, vehicleLicenseTypes) VALUES ($1, $2, $3) RETURNING id
	`, userId, vehicleTypesJSON, vehicleLicenseTypesJSON).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Create staff chats if not exists
	if err := r.CreateBaseChatForAllStaffsIfNotExists(); err != nil {
		log.Printf("failed to create base chat for all staffs: %v\n", err)
	}

	// Insert limitations
	for _, limitation := range payload.Limitations {
		_, err = r.PostgresDB.Exec(`
			INSERT INTO staffLimitations (staffId, limitationId, description, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, NOW(), NOW(), NULL)
		`, insertedId, limitation.LimitationID, limitation.Description)
		if err != nil {
			return nil, err
		}
	}

	staffs, err = r.Query(
		&models.StaffsQueryRequestParams{
			UserID: int(userId),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(staffs) == 0 {
		return nil, errors.New("no staff found")
	}

	return staffs[0], nil
}

func (r *StaffRepositoryPostgresDB) UpdateStaff(userId int64, payload *models.StaffsUpdatePersonalInfoRequestBody) (*domain.Staff, error) {
	staffs, err := r.Query(
		&models.StaffsQueryRequestParams{
			UserID: int(userId),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(staffs) == 0 {
		return nil, errors.New("no staff found")
	}
	staff := staffs[0]

	// Create vehicleTypes jsonb
	vehicleTypesBytes, err := json.Marshal(payload.VehicleTypesAsStringSlice)
	if err != nil {
		return nil, err
	}
	vehicleTypesJSON := string(vehicleTypesBytes)

	// Create vehicleLicenseTypes jsonb
	vehicleLicenseTypesBytes, err := json.Marshal(payload.VehicleLicenseTypesAsStringSlice)
	if err != nil {
		return nil, err
	}
	vehicleLicenseTypesJSON := string(vehicleLicenseTypesBytes)

	// Update staffs
	_, err = r.PostgresDB.Exec(`
		UPDATE staffs SET vehicleTypes = $1, vehicleLicenseTypes = $2 WHERE id = $3
	`, vehicleTypesJSON, vehicleLicenseTypesJSON, staff.ID)
	if err != nil {
		return nil, err
	}

	// Delete limitations
	_, err = r.PostgresDB.Exec(`
		DELETE FROM staffLimitations WHERE staffId = $1
	`, staff.ID)
	if err != nil {
		return nil, err
	}

	// Insert limitations
	for _, limitation := range payload.Limitations {
		_, err = r.PostgresDB.Exec(`
			INSERT INTO staffLimitations (staffId, limitationId, description, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, NOW(), NOW(), NULL)
		`, staff.ID, limitation.LimitationID, limitation.Description)
		if err != nil {
			return nil, err
		}
	}

	return staff, nil
}

func (r *StaffRepositoryPostgresDB) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffAbsencesQueryRes, error) {
	var absence domain.StaffAbsencesQueryRes

	// Current time
	currentTime := time.Now()

	// Find the absence by id
	results, _, err := r.QueryAbsences(&models.StaffsQueryAbsencesRequestParams{
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
		UPDATE staffAbsences
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
	results, _, err = r.QueryAbsences(&models.StaffsQueryAbsencesRequestParams{
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

func (r *StaffRepositoryPostgresDB) UpdateLicenseAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffLicensesRes, error) {
	var license domain.StaffLicensesRes

	// Current time
	currentTime := time.Now()

	// Find the license by id
	results, _, err := r.QueryLicenses(&models.StaffsQueryLicensesRequestParams{
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
	foundLicense := results[0]
	if foundLicense == nil {
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

	// Update the license
	err = r.PostgresDB.QueryRow(`
		UPDATE staffLicenses
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundLicense.ID,
	).Scan(
		&license.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the license
	results, _, err = r.QueryLicenses(&models.StaffsQueryLicensesRequestParams{
		ID:    int(foundLicense.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	license = *results[0]

	// Return the license
	return &license, nil
}

func (r *StaffRepositoryPostgresDB) QueryOtherAttachments(payload *models.StaffsQueryOtherAttachmentsRequestParams) ([]*domain.StaffOtherAttachment, int64, error) {
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
	log.Printf("Query for staff other attachments: %s\n", query)

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
		userOtherAttachments []*domain.StaffOtherAttachment
	)
	for rows.Next() {
		var (
			staffOtherAttachment = &domain.StaffOtherAttachment{}
			title                sql.NullString
			attachments          sql.NullString
			attachmentsMetadata  []*types.UploadMetadata
			deletedAt            sql.NullTime
			userID               sql.NullInt64
			userFirstName        sql.NullString
			userLastName         sql.NullString
			userAvatarUrl        sql.NullString
		)
		err := rows.Scan(
			&staffOtherAttachment.ID,
			&staffOtherAttachment.UserID,
			&title,
			&attachments,
			&staffOtherAttachment.CreatedAt,
			&staffOtherAttachment.UpdatedAt,
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
			staffOtherAttachment.Title = title.String
		}
		if userID.Valid {
			staffOtherAttachment.User = &domain.StaffOtherAttachmentUser{
				ID: uint(userID.Int64),
			}
			if userFirstName.Valid {
				staffOtherAttachment.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				staffOtherAttachment.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				staffOtherAttachment.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if deletedAt.Valid {
			staffOtherAttachment.DeletedAt = &deletedAt.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff other attachments: %d", err, staffOtherAttachment.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			staffOtherAttachment.Attachments = attachmentsMetadata
		}

		// Find staffId
		if err := r.PostgresDB.QueryRow(`SELECT id FROM staffs WHERE userId = $1`, staffOtherAttachment.UserID).Scan(&staffOtherAttachment.StaffID); err != nil {
			log.Printf("failed to find staffId for staff other attachment: %d, %v", staffOtherAttachment.ID, err)
		}

		userOtherAttachments = append(userOtherAttachments, staffOtherAttachment)
	}

	return userOtherAttachments, count, nil
}

func (r *StaffRepositoryPostgresDB) CreateOtherAttachments(staff *domain.Staff, payload *models.StaffsCreateOtherAttachmentsRequestBody) (*domain.StaffOtherAttachment, error) {
	// Insert the other attachment
	var otherAttachment domain.StaffOtherAttachment
	err := r.PostgresDB.QueryRow(`
		INSERT INTO userOtherAttachments (userId, title, created_at, updated_at, deleted_at)
		VALUES ($1, $2, NOW(), NOW(), NULL)
		RETURNING id
	`,
		staff.UserID,
		payload.Title,
	).Scan(
		&otherAttachment.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the other attachment
	results, _, err := r.QueryOtherAttachments(&models.StaffsQueryOtherAttachmentsRequestParams{
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

func (r *StaffRepositoryPostgresDB) UpdateStaffOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.StaffOtherAttachment, error) {
	var otherAttachment domain.StaffOtherAttachment

	// Current time
	currentTime := time.Now()

	// Find the other attachment by id
	results, _, err := r.QueryOtherAttachments(&models.StaffsQueryOtherAttachmentsRequestParams{
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
	results, _, err = r.QueryOtherAttachments(&models.StaffsQueryOtherAttachmentsRequestParams{
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

func (r *StaffRepositoryPostgresDB) UpdateStaffOtherAttachment(staffOtherAttachment *domain.StaffOtherAttachment, payload *models.StaffsUpdateOtherAttachmentRequestBody) (*domain.StaffOtherAttachment, error) {
	var (
		updatedID   int
		currentTime = time.Now()
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE userOtherAttachments SET title = $1, updated_at = $2 WHERE id = $3 RETURNING id
	`,
		payload.Title,
		currentTime,
		staffOtherAttachment.ID,
	).Scan(
		&updatedID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the other attachment
	results, _, err := r.QueryOtherAttachments(&models.StaffsQueryOtherAttachmentsRequestParams{
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

func (r *StaffRepositoryPostgresDB) DeleteStaffOtherAttachments(payload *models.StaffsDeleteOtherAttachmentsRequestBody) ([]int64, error) {
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

func (r *StaffRepositoryPostgresDB) QueryLibraries(payload *models.StaffsQueryLibrariesRequestParams) ([]*domain.StaffLibrary, int64, error) {
	query := `
		SELECT
		    userLibraries.id,
		    userLibraries.userId,
		    userLibraries.title,
			userLibraries.attachments,
		    userLibraries.created_at,
		    userLibraries.updated_at,
			userLibraries.deleted_at,
			users.id AS userUserId,
		    users.firstName AS userName,
		    users.lastName AS userLastName,
		    users.avatarUrl AS userAvatarUrl
		FROM userLibraries
		LEFT JOIN users ON users.id = userLibraries.userId
	`
	var (
		where []string
		args  []interface{}
	)
	if payload.ID > 0 {
		where = append(where, fmt.Sprintf("userLibraries.id = $%v", len(args)+1))
		args = append(args, payload.ID)
	}
	if payload.UserID > 0 {
		where = append(where, fmt.Sprintf("userLibraries.userId = $%v", len(args)+1))
		args = append(args, payload.UserID)
	}
	if payload.Filters.Title.Op != "" {
		opValue := utils.GetDBOperatorAndValue(payload.Filters.Title.Op, fmt.Sprintf("%v", payload.Filters.Title.Value))
		val := exp.TerIf(opValue.Value == "", "", opValue.Value)
		where = append(where, fmt.Sprintf("userLibraries.title %v $%v", opValue.Operator, len(args)+1))
		args = append(args, val)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	var sorts []string
	if payload.Sorts.Title.Op != "" {
		sorts = append(sorts, fmt.Sprintf(" userLibraries.title %s", payload.Sorts.Title.Op))
	}
	if payload.Sorts.CreatedAt.Op != "" {
		sorts = append(sorts, fmt.Sprintf(" userLibraries.created_at %s", payload.Sorts.CreatedAt.Op))
	}
	if len(sorts) > 0 {
		query += " ORDER BY " + strings.Join(sorts, ",")
	}

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
		userLibraries []*domain.StaffLibrary
	)
	for rows.Next() {
		var (
			staffLibrary        = &domain.StaffLibrary{}
			title               sql.NullString
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			deletedAt           sql.NullTime
			userID              sql.NullInt64
			userFirstName       sql.NullString
			userLastName        sql.NullString
			userAvatarUrl       sql.NullString
		)
		err := rows.Scan(
			&staffLibrary.ID,
			&staffLibrary.UserID,
			&title,
			&attachments,
			&staffLibrary.CreatedAt,
			&staffLibrary.UpdatedAt,
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
			staffLibrary.Title = title.String
		}
		if userID.Valid {
			staffLibrary.User = &domain.StaffLibraryUser{
				ID: uint(userID.Int64),
			}
			if userFirstName.Valid {
				staffLibrary.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				staffLibrary.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				staffLibrary.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if deletedAt.Valid {
			staffLibrary.DeletedAt = &deletedAt.Time
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in staff libraries: %d", err, staffLibrary.ID)
			} else {
				for _, library := range attachmentsMetadata {
					library.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			staffLibrary.Attachments = attachmentsMetadata
		}

		// Find staffId
		if err := r.PostgresDB.QueryRow(`SELECT id FROM staffs WHERE userId = $1`, staffLibrary.UserID).Scan(&staffLibrary.StaffID); err != nil {
			log.Printf("failed to find staffId for staff library: %d, %v", staffLibrary.ID, err)
		}

		userLibraries = append(userLibraries, staffLibrary)
	}

	return userLibraries, count, nil
}

func (r *StaffRepositoryPostgresDB) CreateLibraries(staff *domain.Staff, payload *models.StaffsCreateLibrariesRequestBody) (*domain.StaffLibrary, error) {
	// Insert the library
	var otherAttachment domain.StaffLibrary
	err := r.PostgresDB.QueryRow(`
		INSERT INTO userLibraries (userId, title, created_at, updated_at, deleted_at)
		VALUES ($1, $2, NOW(), NOW(), NULL)
		RETURNING id
	`,
		staff.UserID,
		payload.Title,
	).Scan(
		&otherAttachment.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the library
	results, _, err := r.QueryLibraries(&models.StaffsQueryLibrariesRequestParams{
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

	// Return the library
	return &otherAttachment, nil
}

func (r *StaffRepositoryPostgresDB) UpdateStaffLibraries(attachments []*types.UploadMetadata, id int64) (*domain.StaffLibrary, error) {
	var otherAttachment domain.StaffLibrary

	// Current time
	currentTime := time.Now()

	// Find the library by id
	results, _, err := r.QueryLibraries(&models.StaffsQueryLibrariesRequestParams{
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
	foundLibrary := results[0]
	if foundLibrary == nil {
		return nil, errors.New("no results found")
	}

	// Marshal attachments into JSON format
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)

	// Update the library
	err = r.PostgresDB.QueryRow(`
		UPDATE userLibraries
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundLibrary.ID,
	).Scan(
		&otherAttachment.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the library
	results, _, err = r.QueryLibraries(&models.StaffsQueryLibrariesRequestParams{
		ID:    int(foundLibrary.ID),
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

	// Return the library
	return &otherAttachment, nil
}

func (r *StaffRepositoryPostgresDB) UpdateStaffLibrary(staffLibrary *domain.StaffLibrary, payload *models.StaffsUpdateLibraryRequestBody) (*domain.StaffLibrary, error) {
	var (
		updatedID   int
		currentTime = time.Now()
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE userLibraries SET title = $1, updated_at = $2 WHERE id = $3 RETURNING id
	`,
		payload.Title,
		currentTime,
		staffLibrary.ID,
	).Scan(
		&updatedID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the library
	results, _, err := r.QueryLibraries(&models.StaffsQueryLibrariesRequestParams{
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

func (r *StaffRepositoryPostgresDB) DeleteStaffLibraries(payload *models.StaffsDeleteLibrariesRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(
		`
			DELETE FROM
				userLibraries
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

func (r *StaffRepositoryPostgresDB) UpdateStaffAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Staff, error) {
	var staff domain.Staff

	// Current time
	currentTime := time.Now()

	// Find the staff by id
	results, err := r.Query(&models.StaffsQueryRequestParams{
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
	foundStaff := results[0]
	if foundStaff == nil {
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

	// Update the staff
	err = r.PostgresDB.QueryRow(`
		UPDATE staffs
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundStaff.ID,
	).Scan(
		&staff.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the staff
	results, err = r.Query(&models.StaffsQueryRequestParams{
		ID:    int(foundStaff.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	staff = *results[0]

	// Return the staff
	return &staff, nil
}

// makeStaffChatsWhereFilters generates a slice of SQL WHERE conditions based on the provided query parameters for staff chats.
//
// It takes a pointer to models.StaffsQueryChatsRequestParams as an argument.
// Returns a slice of strings representing the WHERE conditions.
func makeStaffChatsWhereFilters(queries *models.StaffsQueryChatsRequestParams) []string {
	var where []string
	log.Printf("queries: %#v\n", queries)
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cc.id = %d ", queries.ID))
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

// QueryChats retrieves a list of staff chats based on the provided query parameters.
//
// The `queries` parameter is a pointer to `models.StaffsQueryChatsRequestParams` that contains filters for the query.
// It returns a slice of pointers to `domain.StaffChat` and an error.
func (r *StaffRepositoryPostgresDB) QueryChats(queries *models.StaffsQueryChatsRequestParams) ([]*domain.StaffChat, error) {
	q := `
		SELECT
			cc.id,
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
		FROM staffChats cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeStaffChatsWhereFilters(queries)
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

	var chats []*domain.StaffChat
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chat                   domain.StaffChat
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
			chat.SenderUser = &domain.StaffChatUser{
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
			chat.RecipientUser = &domain.StaffChatUser{
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
				log.Printf("failed to unmarshal attachments metadata: %v in staff chats: %d", err, chat.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			chat.Attachments = attachmentsMetadata
		}
		chats = append(chats, &chat)
	}

	// Fetch last message for each staffChat from chatMessages
	for _, chat := range chats {
		var (
			id          int64
			message     sql.NullString
			messageType string
			staffChatId int64
		)
		err := r.PostgresDB.QueryRow(`
			SELECT id, staffChatId, message, messageType FROM staffChatMessages WHERE staffChatId = $1
		`, chat.ID).Scan(&id, &staffChatId, &message, &messageType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			log.Printf("failed to fetch last message: %v in staff chat: %d", err, chat.ID)
			continue
		}
		_ = staffChatId
		if message.Valid && messageType == constants.CHAT_MESSAGE_TYPE_TEXT {
			chat.Message = &message.String
		}
	}
	return chats, nil
}

// CountChats returns the number of staff chats based on the provided query parameters.
//
// It takes a pointer to models.StaffsQueryChatsRequestParams as an argument.
// Returns the count of staff chats as int64 and an error if any.
func (r *StaffRepositoryPostgresDB) CountChats(queries *models.StaffsQueryChatsRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cc.id)
		FROM staffChats cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeStaffChatsWhereFilters(queries)
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

func makeStaffChatMessagesWhereFilters(queries *models.StaffsQueryChatMessagesRequestParams) []string {
	var where []string
	log.Printf("queries: %#v\n", queries)
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" cc.id = %d ", queries.ID))
		}
		if queries.StaffChatID != 0 {
			where = append(where, fmt.Sprintf(" cc.staffChatId = %d ", queries.StaffChatID))
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

// QueryChatMessages retrieves staff chat messages based on the provided query parameters.
//
// It takes a pointer to models.StaffsQueryChatMessagesRequestParams as an argument.
// Returns a slice of pointers to domain.StaffChatMessage and an error if any.
func (r *StaffRepositoryPostgresDB) QueryChatMessages(queries *models.StaffsQueryChatMessagesRequestParams) ([]*domain.StaffChatMessage, error) {
	q := `
		SELECT
			cc.id,
			cc.staffChatId,
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
		FROM staffChatMessages cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeStaffChatMessagesWhereFilters(queries)
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

	var chatMessages []*domain.StaffChatMessage
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chatMessage            domain.StaffChatMessage
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
			&chatMessage.StaffChatID,
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
			chatMessage.SenderUser = &domain.StaffChatMessageUser{
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
			chatMessage.RecipientUser = &domain.StaffChatMessageUser{
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
				log.Printf("failed to unmarshal attachments metadata: %v in staff chat messages: %d", err, chatMessage.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.STAFF_BUCKET_NAME[len("maja."):])
				}
			}
			chatMessage.Attachments = attachmentsMetadata
		}
		chatMessages = append(chatMessages, &chatMessage)
	}

	return chatMessages, nil
}

// CountChatMessages returns the number of staff chat messages based on the provided query parameters.
//
// It takes a pointer to models.StaffsQueryChatMessagesRequestParams as an argument.
// Returns the count of staff chat messages as int64 and an error if any.
func (r *StaffRepositoryPostgresDB) CountChatMessages(queries *models.StaffsQueryChatMessagesRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(cc.id)
		FROM staffChatMessages cc
		LEFT JOIN users u ON cc.senderUserId = u.id
		LEFT JOIN users u2 ON cc.recipientUserId = u2.id
	`
	if queries != nil {
		where := makeStaffChatMessagesWhereFilters(queries)
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
// It takes a StaffsCreateChatMessageRequestBody payload as input, which contains the staff chat ID, sender user ID, recipient user ID, message, and attachments.
// Returns a pointer to a StaffChatMessage and an error.
func (r *StaffRepositoryPostgresDB) CreateChatMessage(payload *models.StaffsCreateChatMessageRequestBody) (*domain.StaffChatMessage, error) {
	// Get current time
	var (
		currentTime           = time.Now()
		isSystem              = false
		insertedChatMessageID int64
	)
	err := r.PostgresDB.QueryRow(`
		INSERT INTO staffChatMessages (staffChatId, senderUserId, recipientUserId, isSystem, message, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.StaffChatID,
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
	chatMessages, err := r.QueryChatMessages(&models.StaffsQueryChatMessagesRequestParams{
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
func (r *StaffRepositoryPostgresDB) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffChatMessage, error) {
	var chatMessage domain.StaffChatMessage

	// Current time
	currentTime := time.Now()

	// Find the medicine by id
	results, err := r.QueryChatMessages(&models.StaffsQueryChatMessagesRequestParams{
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
		UPDATE staffChatMessages
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
	results, err = r.QueryChatMessages(&models.StaffsQueryChatMessagesRequestParams{
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

func (s *StaffRepositoryPostgresDB) CreateBaseChatForAllStaffsIfNotExists() error {
	// Get all staffs
	staffs, err := s.Query(&models.StaffsQueryRequestParams{
		Limit: -1,
	})
	if err != nil {
		return err
	}

	// Create a chat for each staff if it doesn't exist, each staff for each staff
	for _, sender := range staffs {
		for _, recipient := range staffs {
			if sender.ID == recipient.ID {
				continue
			}

			// Check if chat exists
			var exists bool
			err := s.PostgresDB.QueryRow(`
				SELECT EXISTS (
					SELECT 1 FROM staffChats WHERE (senderUserId = $1 OR recipientUserId = $2) AND (senderUserId = $3 OR recipientUserId = $4)
				)
			`, sender.UserID, sender.UserID, recipient.UserID, recipient.UserID).Scan(&exists)
			if err != nil {
				return err
			}
			if exists {
				continue
			}

			// Create chat if it doesn't exist
			var currentTime = time.Now()
			_, err = s.PostgresDB.Exec(`
				INSERT INTO staffChats (senderUserId, recipientUserId, isSystem, created_at, updated_at, deleted_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`,
				sender.UserID,
				recipient.UserID,
				false,
				currentTime,
				currentTime,
				nil,
			)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
