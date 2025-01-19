package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/notification/constants"
	"github.com/hoitek/Maja-Service/internal/notification/domain"
	"github.com/hoitek/Maja-Service/internal/notification/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type NotificationRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewNotificationRepositoryPostgresDB(d *sql.DB) *NotificationRepositoryPostgresDB {
	return &NotificationRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.NotificationsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" n.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" n.userId = %d", queries.UserID))
		}
		if queries.Type == constants.NOTIFICATION_TYPE_NOTIFICATION {
			where = append(where, " n.status IS NULL")
		}
		if queries.Type == constants.NOTIFICATION_TYPE_REQUEST {
			where = append(where, " n.status IS NOT NULL")
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" n.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" n.description %s %s", opValue.Operator, val))
		}
		if queries.Filters.ReadAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ReadAt.Op, fmt.Sprintf("%v", queries.Filters.ReadAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" n.read_at %s %s", opValue.Operator, val))
		}
		if queries.Filters.IsForSystem.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.IsForSystem.Op, fmt.Sprintf("%v", queries.Filters.IsForSystem.Value))
			val := exp.TerIf(opValue.Value == "", "", opValue.Value)
			where = append(where, fmt.Sprintf(" n.isForSystem %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *NotificationRepositoryPostgresDB) Query(queries *models.NotificationsQueryRequestParams) ([]*domain.Notification, error) {
	q := `
		SELECT
			n.id,
			n.userId,
			n.title,
			n.description,
			n.read_at,
			n.readBy,
			n.extra,
			n.isForSystem,
			n.status,
			n.status_at,
			n.created_at,
			n.updated_at,
			n.deleted_at,
			u.id AS userId,
			u.firstName AS userFirstName,
			u.lastName AS userLastName,
			u.avatarUrl AS userAvatarUrl,
			u2.id AS readyById,
			u2.firstName AS readyByFirstName,
			u2.lastName AS readyByLastName,
			u2.avatarUrl AS readyByAvatarUrl
		FROM notifications n
		LEFT JOIN users u ON u.id = n.userId
		LEFT JOIN users u2 ON u2.id = n.readBy
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" n.name %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" n.description %s", queries.Sorts.Description.Op))
		}
		if queries.Sorts.ReadAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" n.read_at %s", queries.Sorts.ReadAt.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" n.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var notifications []*domain.Notification
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			notification    domain.Notification
			readAt          sql.NullTime
			readBy          sql.NullInt64
			extra           json.RawMessage
			isForSystem     sql.NullBool
			deletedAt       sql.NullTime
			userId          sql.NullInt64
			userFirstName   sql.NullString
			userLastName    sql.NullString
			userAvatarUrl   sql.NullString
			readById        sql.NullInt64
			readByFirstName sql.NullString
			readByLastName  sql.NullString
			readByAvatarUrl sql.NullString
			status          sql.NullString
			statusAt        sql.NullTime
		)
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Description,
			&readAt,
			&readBy,
			&extra,
			&isForSystem,
			&status,
			&statusAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&readById,
			&readByFirstName,
			&readByLastName,
			&readByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if userId.Valid {
			uid := uint(userId.Int64)
			notification.UserID = &uid
			notification.User = &domain.NotificationUser{
				ID: uid,
			}
			if userFirstName.Valid {
				notification.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				notification.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				notification.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if readAt.Valid {
			notification.ReadAt = &readAt.Time
		}
		if readBy.Valid {
			rb := uint(readBy.Int64)
			notification.ReadBy = &rb
		}
		if extra != nil {
			var extraData interface{}
			err := json.Unmarshal(extra, &extraData)
			if err != nil {
				return nil, err
			}
			notification.Extra = extraData
		}
		if isForSystem.Valid {
			notification.IsForSystem = isForSystem.Bool
		}
		if deletedAt.Valid {
			notification.DeletedAt = &deletedAt.Time
		}
		if readById.Valid {
			notification.ReadByUser = &domain.NotificationUser{
				ID: uint(readById.Int64),
			}
			if readByFirstName.Valid {
				notification.ReadByUser.FirstName = readByFirstName.String
			}
			if readByLastName.Valid {
				notification.ReadByUser.LastName = readByLastName.String
			}
			if readByAvatarUrl.Valid {
				notification.ReadByUser.AvatarUrl = readByAvatarUrl.String
			}
		}
		if status.Valid {
			notification.Status = &status.String
		}
		if statusAt.Valid {
			notification.StatusAt = &statusAt.Time
		}
		notifications = append(notifications, &notification)
	}
	return notifications, nil
}

func (r *NotificationRepositoryPostgresDB) Count(queries *models.NotificationsQueryRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(n.id)
		FROM notifications n
		LEFT JOIN users u ON u.id = n.readBy
	`
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

func (r *NotificationRepositoryPostgresDB) Delete(payload *models.NotificationsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM notifications
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

func (r *NotificationRepositoryPostgresDB) GetNotificationsByIds(ids []int64) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	rows, err := r.PostgresDB.Query(`
		SELECT
			n.id,
			n.userId,
			n.title,
			n.description,
			n.read_at,
			n.readBy,
			n.extra,
			n.isForSystem,
			n.status,
			n.status_at,
			n.created_at,
			n.updated_at,
			n.deleted_at,
			u.id AS userId,
			u.firstName AS userFirstName,
			u.lastName AS userLastName,
			u.avatarUrl AS userAvatarUrl,
			u2.id AS readyById,
			u2.firstName AS readyByFirstName,
			u2.lastName AS readyByLastName,
			u2.avatarUrl AS readyByAvatarUrl
		FROM notifications n
		LEFT JOIN users u ON u.id = n.userId
		LEFT JOIN users u2 ON u2.id = n.readBy
		WHERE n.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			notification    domain.Notification
			readAt          sql.NullTime
			readBy          sql.NullInt64
			extra           json.RawMessage
			isForSystem     sql.NullBool
			deletedAt       sql.NullTime
			userId          sql.NullInt64
			userFirstName   sql.NullString
			userLastName    sql.NullString
			userAvatarUrl   sql.NullString
			readById        sql.NullInt64
			readByFirstName sql.NullString
			readByLastName  sql.NullString
			readByAvatarUrl sql.NullString
			status          sql.NullString
			statusAt        sql.NullTime
		)
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Description,
			&readAt,
			&readBy,
			&extra,
			&isForSystem,
			&status,
			&statusAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&readById,
			&readByFirstName,
			&readByLastName,
			&readByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if userId.Valid {
			notification.User = &domain.NotificationUser{
				ID: uint(userId.Int64),
			}
			if userFirstName.Valid {
				notification.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				notification.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				notification.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if readAt.Valid {
			notification.ReadAt = &readAt.Time
		}
		if readBy.Valid {
			rb := uint(readBy.Int64)
			notification.ReadBy = &rb
		}
		if extra != nil {
			var extraData interface{}
			err := json.Unmarshal(extra, &extraData)
			if err != nil {
				return nil, err
			}
		}
		if isForSystem.Valid {
			notification.IsForSystem = isForSystem.Bool
		}
		if deletedAt.Valid {
			notification.DeletedAt = &deletedAt.Time
		}
		if readById.Valid {
			notification.ReadByUser = &domain.NotificationUser{
				ID: uint(readById.Int64),
			}
			if readByFirstName.Valid {
				notification.ReadByUser.FirstName = readByFirstName.String
			}
			if readByLastName.Valid {
				notification.ReadByUser.LastName = readByLastName.String
			}
			if readByAvatarUrl.Valid {
				notification.ReadByUser.AvatarUrl = readByAvatarUrl.String
			}
		}
		if status.Valid {
			notification.Status = &status.String
		}
		if statusAt.Valid {
			notification.StatusAt = &statusAt.Time
		}
		notifications = append(notifications, &notification)
	}
	return notifications, nil
}

func (r *NotificationRepositoryPostgresDB) CreateNotification(userId int64, title, description string, status string, isForSystem bool, extra interface{}) (*domain.Notification, error) {
	var notificationID int
	currentTime := time.Now()
	extraBytes, err := json.Marshal(map[string]interface{}{
		"data": extra,
	})
	if err != nil {
		return nil, err
	}
	extraJSON := string(extraBytes)
	err = r.PostgresDB.QueryRow(`
		INSERT INTO notifications (userId, title, description, extra, isForSystem, status, status_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userId, title, description, extraJSON, isForSystem, status, currentTime).Scan(&notificationID)
	if err != nil {
		return nil, err
	}
	notifications, err := r.Query(&models.NotificationsQueryRequestParams{
		ID: notificationID,
	})
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return nil, errors.New("no notifications found")
	}
	notification := notifications[0]
	return notification, nil
}

func (r *NotificationRepositoryPostgresDB) UpdateStatus(payload *models.NotificationsUpdateStatusRequestBody, notificationID int64) (*domain.Notification, error) {
	var (
		updatedID   int64
		currentTime = time.Now()
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE notifications
		SET status = $1, status_at = $2
		WHERE id = $3
		RETURNING id
	`, payload.Status, currentTime, notificationID).Scan(
		&updatedID,
	)
	if err != nil {
		return nil, err
	}
	notifications, err := r.Query(&models.NotificationsQueryRequestParams{
		ID: int(updatedID),
	})
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return nil, errors.New("no notifications found")
	}
	notification := notifications[0]
	return notification, nil
}
