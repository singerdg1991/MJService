package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/utils"

	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
)

type UserRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewUserRepositoryPostgresDB(d *sql.DB) *UserRepositoryPostgresDB {
	return &UserRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.UsersQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf("users.id = %d", queries.ID))
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
		if queries.Filters.Email.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Email.Op, queries.Filters.Email.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.email %s %s", opValue.Operator, val))
		}
		if queries.Filters.WorkPhoneNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.WorkPhoneNumber.Op, queries.Filters.WorkPhoneNumber.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.workPhoneNumber %s %s", opValue.Operator, val))
		}
		if queries.Filters.Username.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Username.Op, queries.Filters.Username.Value)
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf("users.username %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *UserRepositoryPostgresDB) ReAttachStaffIDAndCustomerID(user *domain.User) *domain.User {
	var (
		staffID    sql.NullInt64
		customerID sql.NullInt64
	)
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM staffs
		WHERE userId = $1
	`, user.ID).Scan(&staffID)
	if err != nil {
		return user
	}
	if staffID.Valid {
		sid := uint(staffID.Int64)
		user.StaffID = &sid
		res, err := r.PostgresDB.Exec(`
			UPDATE users
			SET staffId = $1
			WHERE id = $2
		`, staffID.Int64, user.ID)
		if err != nil {
			log.Println("Error updating staffId: ", err)
			return user
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Println("Error getting rows affected: ", err)
			return user
		}
		if rowsAffected == 0 {
			log.Println("No rows affected")
			return user
		}
	}
	err = r.PostgresDB.QueryRow(`
		SELECT id
		FROM customers
		WHERE userId = $1
	`, user.ID).Scan(&customerID)
	if err != nil {
		return user
	}
	if customerID.Valid {
		cid := uint(customerID.Int64)
		user.CustomerID = &cid
		res, err := r.PostgresDB.Exec(`
			UPDATE users
			SET customerId = $1
			WHERE id = $2
		`, customerID.Int64, user.ID)
		if err != nil {
			log.Println("Error updating customerId: ", err)
			return user
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Println("Error getting rows affected: ", err)
			return user
		}
		if rowsAffected == 0 {
			log.Println("No rows affected")
			return user
		}
	}
	return user
}

func (r *UserRepositoryPostgresDB) Query(queries *models.UsersQueryRequestParams) ([]*domain.User, error) {
	// Build query with filters
	q := `SELECT * FROM users `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	log.Println("Query: ", q)

	// Find all users with r.RawDB.Query
	var users []*domain.User
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			user                    domain.User
			customerID              sql.NullInt64
			staffID                 sql.NullInt64
			firstName               sql.NullString
			lastName                sql.NullString
			username                sql.NullString
			password                sql.NullString
			email                   sql.NullString
			phone                   sql.NullString
			telephone               sql.NullString
			registrationNumber      sql.NullString
			workPhoneNumber         sql.NullString
			gender                  sql.NullString
			accountNumber           sql.NullString
			nationalCode            sql.NullString
			birthDate               sql.NullTime
			avatarUrl               sql.NullString
			forcedChangePassword    sql.NullBool
			suspendedAt             sql.NullTime
			privacyPolicyAcceptedAt sql.NullTime
			createdAt               sql.NullTime
			updatedAt               sql.NullTime
			deletedAt               sql.NullTime
		)
		err := rows.Scan(
			&user.ID,
			&customerID,
			&staffID,
			&firstName,
			&lastName,
			&username,
			&password,
			&email,
			&phone,
			&telephone,
			&registrationNumber,
			&workPhoneNumber,
			&gender,
			&accountNumber,
			&nationalCode,
			&birthDate,
			&avatarUrl,
			&forcedChangePassword,
			&suspendedAt,
			&privacyPolicyAcceptedAt,
			&user.UserType,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if customerID.Valid {
			cid := uint(customerID.Int64)
			user.CustomerID = &cid
		}
		if staffID.Valid {
			sid := uint(staffID.Int64)
			user.StaffID = &sid
		}
		if firstName.Valid {
			user.FirstName = firstName.String
		}
		if lastName.Valid {
			user.LastName = lastName.String
		}
		if username.Valid {
			user.Username = username.String
		}
		if password.Valid {
			user.Password = password.String
		}
		if email.Valid {
			user.Email = email.String
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if telephone.Valid {
			user.Telephone = &telephone.String
		}
		if registrationNumber.Valid {
			user.RegistrationNumber = &registrationNumber.String
		}
		if workPhoneNumber.Valid {
			user.WorkPhoneNumber = &workPhoneNumber.String
		}
		if gender.Valid {
			user.Gender = &gender.String
		}
		if accountNumber.Valid {
			user.AccountNumber = &accountNumber.String
		}
		if nationalCode.Valid {
			user.NationalCode = nationalCode.String
		}
		if birthDate.Valid {
			user.BirthDate = birthDate.Time
		}
		if avatarUrl.Valid {
			user.AvatarUrl = avatarUrl.String
		}
		if forcedChangePassword.Valid {
			user.ForcedChangePassword = forcedChangePassword.Bool
		}
		if suspendedAt.Valid {
			user.SuspendedAt = &suspendedAt.Time
		}
		if privacyPolicyAcceptedAt.Valid {
			user.PrivacyPolicyAcceptedAt = &privacyPolicyAcceptedAt.Time
		}
		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}
		if user.StaffID == nil || user.CustomerID == nil {
			user = *r.ReAttachStaffIDAndCustomerID(&user)
		}
		users = append(users, &user)
	}

	// Find user language skills
	for _, user := range users {
		var (
			languageSkill domain.UserLanguageSkillRes
		)
		err := r.PostgresDB.QueryRow(`
			SELECT userLanguageSkills.languageSkillId, languageSkills.name
			FROM userLanguageSkills
			INNER JOIN languageSkills ON languageSkills.id = userLanguageSkills.languageSkillId
			WHERE userLanguageSkills.userId = $1
		`, user.ID).Scan(&languageSkill.ID, &languageSkill.Name)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			continue
		}
		user.LanguageSkills = append(user.LanguageSkills, languageSkill)
	}

	// Find user roles
	for _, user := range users {
		var (
			role domain.UserRole
		)
		err := r.PostgresDB.QueryRow(`
			SELECT usersRoles.roleId, _roles.name
			FROM usersRoles
			LEFT JOIN _roles ON _roles.id = usersRoles.roleId
			WHERE usersRoles.userId = $1
		`, user.ID).Scan(&role.ID, &role.Name)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			continue
		}
		user.RoleIDs = append(user.RoleIDs, domain.UserRoleID{
			ID: role.ID,
		})
		user.Roles = append(user.Roles, role)
	}

	// Find user permissions
	for _, user := range users {
		if user.Roles != nil {
			for index, role := range user.Roles {
				rows, err := r.PostgresDB.Query(`
					SELECT _rolesPermissions.permissionId, _permissions.name, _permissions.title
					FROM _rolesPermissions
					LEFT JOIN _permissions ON _permissions.id = _rolesPermissions.permissionId
					WHERE _rolesPermissions.roleId = $1
				`, role.ID)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, err
					}
					continue
				}
				defer rows.Close()
				var (
					permissions []domain.UserRolePermission
				)
				for rows.Next() {
					var permission domain.UserRolePermission
					err := rows.Scan(&permission.ID, &permission.Name, &permission.Title)
					if err != nil {
						if !errors.Is(err, sql.ErrNoRows) {
							return nil, err
						}
						continue
					}
					permissions = append(permissions, permission)
				}
				user.Roles[index].Permissions = permissions
			}
		}
	}

	return users, nil
}

func (r *UserRepositoryPostgresDB) Count(queries *models.UsersQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM users `
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

func (r *UserRepositoryPostgresDB) Create(payload *models.UsersCreateRequestBody) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	bDate, err := time.Parse(time.RFC3339, payload.BirthDate)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	var accountNumber *string
	if payload.AccountNumber != "" {
		accountNumber = &payload.AccountNumber
	}
	user := &domain.User{
		FirstName:            payload.FirstName,
		LastName:             payload.LastName,
		Username:             payload.Username,
		Password:             payload.Password,
		Email:                payload.Email,
		Phone:                payload.Phone,
		Gender:               &payload.Gender,
		WorkPhoneNumber:      &payload.WorkPhoneNumber,
		AccountNumber:        accountNumber,
		RegistrationNumber:   &payload.RegistrationNumber,
		ForcedChangePassword: payload.ForcedChangePassword,
		NationalCode:         payload.NationalCode,
		BirthDate:            bDate,
		AvatarUrl:            payload.AvatarUrl,
		SuspendedAt:          nil,
	}

	err = tx.QueryRow(`
			INSERT INTO users
				(firstName, lastName, username, password, gender, workPhoneNumber, accountNumber, registrationNumber, forcedChangePassword, email, phone, nationalCode, birthDate, avatarUrl, suspended_at, created_at, updated_at, deleted_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
			RETURNING id
		`,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Gender,
		user.WorkPhoneNumber,
		user.AccountNumber,
		user.RegistrationNumber,
		user.ForcedChangePassword,
		user.Email,
		user.Phone,
		user.NationalCode,
		user.BirthDate,
		user.AvatarUrl,
		user.SuspendedAt,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	).Scan(&user.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	// Insert user language skills
	for _, languageSkillId := range payload.LanguageSkillIds {
		_, err = tx.Exec(`
			INSERT INTO userLanguageSkills
			    (userId, languageSkillId)
			VALUES
			    ($1, $2)
		`,
			user.ID,
			languageSkillId,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Find user by id
	users, err := r.Query(&models.UsersQueryRequestParams{ID: int(user.ID)})
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

func (r *UserRepositoryPostgresDB) Update(payload *models.UsersUpdateRequestBody, id int) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	users, err := r.Query(&models.UsersQueryRequestParams{ID: id})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	user := users[0]

	// Parse birthdate
	bDate, err := utils.TryParseToDateTime(payload.BirthDate)
	if err != nil {
		return nil, err
	}

	// Parse updated at
	var currentTime = time.Now()

	// Update user
	var accountNumber *string
	if payload.AccountNumber != "" {
		accountNumber = &payload.AccountNumber
	}
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Username = payload.Username
	user.Email = payload.Email
	user.Phone = payload.Phone
	user.AccountNumber = accountNumber
	user.ForcedChangePassword = payload.ForcedChangePassword
	user.NationalCode = payload.NationalCode
	user.RegistrationNumber = &payload.RegistrationNumber
	user.WorkPhoneNumber = &payload.WorkPhoneNumber
	user.NationalCode = payload.NationalCode
	user.BirthDate = bDate
	user.Gender = &payload.Gender
	user.AvatarUrl = payload.AvatarUrl
	user.UpdatedAt = currentTime
	if payload.Password != "" {
		user.Password = payload.Password
	}

	_, err = tx.Exec(`
		UPDATE users
		SET
		    firstName = $1,
		    lastName = $2,
		    username = $3,
		    password = $4,
		    email = $5,
		    phone = $6,
		    accountNumber = $7,
		    forcedChangePassword = $8,
		    nationalCode = $9,
		    registrationNumber = $10,
		    workPhoneNumber = $11,
		    birthDate = $12,
		    gender = $13,
		    avatarUrl = $14,
		    updated_at = $15
		WHERE id = $16
	`,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
		user.AccountNumber,
		user.ForcedChangePassword,
		user.NationalCode,
		user.RegistrationNumber,
		user.WorkPhoneNumber,
		user.BirthDate,
		user.Gender,
		user.AvatarUrl,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	// Update user language skills
	_, err = tx.Exec(`
		DELETE FROM userLanguageSkills
		WHERE userId = $1
	`, user.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	for _, languageSkillId := range payload.LanguageSkillIds {
		_, err = tx.Exec(`
			INSERT INTO userLanguageSkills
			    (userId, languageSkillId)
			VALUES
			    ($1, $2)
		`,
			user.ID,
			languageSkillId,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Find user by id
	users, err = r.Query(&models.UsersQueryRequestParams{ID: int(user.ID)})
	if err != nil {
		return nil, err
	}
	log.Printf("user %#v\n", users[0])

	return users[0], nil
}

func (r *UserRepositoryPostgresDB) Delete(data *models.UsersDeleteRequestBody) ([]uint, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	var ids []uint
	idsStr := strings.Split(data.IDs, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		ids = append(ids, uint(id))
	}
	var rowsAffected int64
	err = tx.QueryRow(`DELETE FROM users WHERE id IN ($1) RETURNING 1`, strings.Join(idsStr, ",")).Scan(&rowsAffected)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("no rows affected")
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *UserRepositoryPostgresDB) UpdateAcceptPolicy(acceptPolicy bool, id int) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	users, err := r.Query(&models.UsersQueryRequestParams{ID: id})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(users) == 0 {
		tx.Rollback()
		return nil, errors.New("user not found")
	}
	user := users[0]

	currentTime := time.Now()
	user.PrivacyPolicyAcceptedAt = nil
	if acceptPolicy {
		user.PrivacyPolicyAcceptedAt = &currentTime
	}
	user.UpdatedAt = currentTime

	res, err := tx.Exec(`
		UPDATE users
		SET
		    privacy_policy_accepted_at = $1,
		    updated_at = $2
		WHERE id = $3
	`, user.PrivacyPolicyAcceptedAt, user.UpdatedAt, user.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("no rows affected")
	}

	user.ID = uint(id)

	return user, nil
}

func (r *UserRepositoryPostgresDB) UpdatePassword(newPassword string, id int) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	users, err := r.Query(&models.UsersQueryRequestParams{ID: id})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(users) == 0 {
		tx.Rollback()
		return nil, errors.New("user not found")
	}
	user := users[0]

	currentTime := time.Now()
	user.Password = newPassword
	user.UpdatedAt = currentTime
	user.ForcedChangePassword = false

	res, err := tx.Exec(`
		UPDATE users
		SET
		    password = $1,
		    forcedChangePassword = $2,
		    updated_at = $3
		WHERE id = $4
	`, user.Password, user.ForcedChangePassword, user.UpdatedAt, user.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("no rows affected")
	}

	user.ID = uint(id)

	return user, nil
}

func (r *UserRepositoryPostgresDB) CreateUserForCustomer(payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	bDate, err := time.Parse(time.RFC3339, payload.DateOfBirth)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	currentTime := time.Now()
	customerId := uint(payload.CustomerID)
	user := &domain.User{
		CustomerID:           &customerId,
		FirstName:            payload.FirstName,
		LastName:             payload.LastName,
		Password:             payload.Password,
		ForcedChangePassword: payload.ForcedChangePassword,
		Gender:               &payload.Gender,
		NationalCode:         payload.NationalCode,
		BirthDate:            bDate,
		Email:                payload.Email,
		Phone:                payload.PhoneNumber,
		SuspendedAt:          nil,
		UserType:             constants.USER_TYPE_CUSTOMER,
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		DeletedAt:            nil,
	}

	err = tx.QueryRow(`
			INSERT INTO users
				(customerId, firstName, lastName, password, forcedChangePassword, gender, nationalCode, birthDate, email, phone, userType, suspended_at, created_at, updated_at, deleted_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id
		`,
		user.CustomerID,
		user.FirstName,
		user.LastName,
		user.Password,
		user.ForcedChangePassword,
		user.Gender,
		user.NationalCode,
		user.BirthDate,
		user.Email,
		user.Phone,
		user.UserType,
		user.SuspendedAt,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	).Scan(&user.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Find user by id
	users, err := r.Query(&models.UsersQueryRequestParams{ID: int(user.ID)})
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

func (r *UserRepositoryPostgresDB) UpdateUserForCustomer(userID int64, payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	bDate, err := time.Parse(time.RFC3339, payload.DateOfBirth)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	currentTime := time.Now()
	customerId := uint(payload.CustomerID)
	user := &domain.User{
		CustomerID:   &customerId,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Gender:       &payload.Gender,
		NationalCode: payload.NationalCode,
		BirthDate:    bDate,
		Email:        payload.Email,
		Phone:        payload.PhoneNumber,
		SuspendedAt:  nil,
		UserType:     constants.USER_TYPE_CUSTOMER,
		UpdatedAt:    currentTime,
	}

	res, err := tx.Exec(`
		UPDATE users
		SET
			customerId = $1,
			firstName = $2,
			lastName = $3,
			gender = $4,
			nationalCode = $5,
			birthDate = $6,
			email = $7,
			phone = $8,
			updated_at = $9,
			userType = $10
		WHERE id = $11
	`,
		user.CustomerID,
		user.FirstName,
		user.LastName,
		user.Gender,
		user.NationalCode,
		user.BirthDate,
		user.Email,
		user.Phone,
		user.UpdatedAt,
		userID,
		user.UserType,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	// Find user by id
	users, err := r.Query(&models.UsersQueryRequestParams{ID: int(userID)})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (r *UserRepositoryPostgresDB) DeleteUserForCustomer(userID int64) error {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, userID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPostgresDB) UpdateCustomerIDForUser(userID int64, customerID int64) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	users, err := r.Query(&models.UsersQueryRequestParams{ID: int(userID)})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	user := users[0]

	currentTime := time.Now()
	newCustomerID := uint(customerID)
	user.CustomerID = &newCustomerID
	user.UpdatedAt = currentTime

	res, err := tx.Exec(`
		UPDATE users
		SET
		    customerId = $1,
		    updated_at = $2
		WHERE id = $3
	`, user.CustomerID, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}

	user.ID = uint(userID)

	return user, nil
}

func (r *UserRepositoryPostgresDB) UpdateUserAdditionalInfoForCustomer(payload *sharedmodels.UpdateUserAdditionalInfoForCustomer) (*domain.User, error) {
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	users, err := r.Query(&models.UsersQueryRequestParams{ID: int(payload.UserID)})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	user := users[0]

	currentTime := time.Now()
	user.UpdatedAt = currentTime

	res, err := tx.Exec(`
		UPDATE users
		SET
		    updated_at = $1
		WHERE id = $2
	`, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}

	user.ID = uint(payload.UserID)

	return user, nil
}
