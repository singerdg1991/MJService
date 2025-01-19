package migrations

import (
	"context"
	"log"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/database"
	sharedconstant "github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/security"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
	"github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/internal/user/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeedUsers1684787622 struct {
}

func NewSeedUsers1684787622() *SeedUsers1684787622 {
	return &SeedUsers1684787622{}
}

func (m *SeedUsers1684787622) MigrateUp() error {
	var (
		genderMale = "male"
		password   = security.HashPassword("111111")
	)
	var users = []*domain.User{
		{
			ID: 1,
			RoleIDs: []domain.UserRoleID{
				{
					ID: RoleOwner,
				},
			},
			FirstName:               "Owner",
			LastName:                "",
			Username:                "owner",
			Password:                password,
			Email:                   "owner@yahoo.com",
			Gender:                  &genderMale,
			Phone:                   "0987654321",
			NationalCode:            "1234567890",
			BirthDate:               time.Now(),
			UserType:                sharedconstant.USER_TYPE_STAFF,
			AvatarUrl:               "https://i.pravatar.cc/150?img=35",
			ForcedChangePassword:    false,
			SuspendedAt:             nil,
			PrivacyPolicyAcceptedAt: nil,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
			DeletedAt:               nil,
		},
	}
	tx, err := database.PostgresDB.Begin()
	if err != nil {
		log.Printf("error in creating transaction: %s", err.Error())
		return err
	}
	for _, user := range users {
		_, err := tx.Exec(`
			INSERT INTO users
			    (id, firstName, lastName, username, password, email, phone, nationalCode, birthDate, avatarUrl, forcedChangePassword, userType, suspended_at, privacy_policy_accepted_at, created_at, updated_at, deleted_at)
			VALUES
			    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			ON CONFLICT(id) DO NOTHING;
		`, user.ID, user.FirstName, user.LastName, user.Username, user.Password, user.Email, user.Phone, user.NationalCode, user.BirthDate, user.AvatarUrl, user.ForcedChangePassword, user.UserType, user.SuspendedAt, user.PrivacyPolicyAcceptedAt, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		if err != nil {
			tx.Rollback()
			log.Println("Error in seeding users: ", err)
		} else {
			log.Printf("User with id %d created successfully", user.ID)
			// Assign roles to user
			for _, role := range user.RoleIDs {
				_, err := tx.Exec(`
					INSERT INTO usersRoles
					    (userId, roleId)
					VALUES
					    ($1, $2)
					ON CONFLICT(userId, roleId) DO NOTHING;
				`, user.ID, role.ID)
				if err != nil {
					tx.Rollback()
					log.Println("Error in seeding user roles: ", err)
				} else {
					log.Printf("User role with id %d assigned to user with id %d successfully", role.ID, user.ID)
				}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("error in commiting transaction: %s", err.Error())
		return err
	}
	database.PostgresDB.Exec(`SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));`)
	database.PostgresDB.Exec(`SELECT setval('usersRoles_id_seq', (SELECT MAX(id) FROM usersRoles));`)

	for _, user := range users {
		// Check if user already exists
		filter := bson.M{
			"id": user.ID,
		}
		count, err := database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).CountDocuments(context.Background(), filter)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		userRepositoryPostgresDB := exp.TerIf[ports.UserRepositoryPostgresDB](
			config.AppConfig.Environment == constants.ENVIRONMENT_TESTING,
			repositories.NewUserRepositoryStub(),
			repositories.NewUserRepositoryPostgresDB(database.PostgresDB),
		)
		u, err := userRepositoryPostgresDB.Query(&models.UsersQueryRequestParams{
			ID: int(user.ID),
		})
		if err != nil {
			return err
		}
		if len(u) == 0 {
			continue
		}
		data, err := u[0].ToJson()
		if err != nil {
			return err
		}

		// Convert data to BSON document
		var sectionData bson.M
		err = bson.UnmarshalExtJSON([]byte(data), true, &sectionData)
		if err != nil {
			return err
		}

		// Create user
		_, err = database.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).InsertOne(context.Background(), sectionData, options.InsertOne())
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SeedUsers1684787622) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'SeedUsers1684787622';`)
	database.PostgresDB.Exec(`
		DELETE FROM users WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
    `)
	return nil
}
