package migrations

import (
	"context"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/internal/user/domain"
)

type UsersMigration1684787622 struct {
}

func NewUsersMigration1684787622() *UsersMigration1684787622 {
	return &UsersMigration1684787622{}
}

func (m *UsersMigration1684787622) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			customerId INTEGER DEFAULT NULL,
			staffId INTEGER DEFAULT NULL,
			firstName VARCHAR(255) DEFAULT NULL,
		    lastName VARCHAR(255) DEFAULT NULL,
		    username VARCHAR(255) DEFAULT NULL,
		    password VARCHAR(255) DEFAULT NULL,
			email VARCHAR(255) DEFAULT NULL,
			phone VARCHAR(255) DEFAULT NULL,
		    telephone VARCHAR(255) DEFAULT NULL,
		    registrationNumber VARCHAR(255) DEFAULT NULL,
		    workPhoneNumber VARCHAR(255) DEFAULT NULL,
		    gender VARCHAR(255) DEFAULT NULL,
		    accountNumber VARCHAR(255) DEFAULT NULL,
		    nationalCode VARCHAR(255) DEFAULT NULL,
		    birthDate TIMESTAMP DEFAULT NULL,
		    avatarUrl VARCHAR(255) DEFAULT NULL,
		    forcedChangePassword BOOLEAN NOT NULL DEFAULT FALSE,
		    suspended_at TIMESTAMP DEFAULT NULL,
		    privacy_policy_accepted_at TIMESTAMP DEFAULT NULL,
            userType VARCHAR(255) DEFAULT 'staff',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
		    deleted_at TIMESTAMP DEFAULT NULL
		);
		CREATE TABLE IF NOT EXISTS userLanguageSkills (
			id SERIAL PRIMARY KEY,
			userId INTEGER NOT NULL,
			languageSkillId INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT userLanguageSkills_userId_fkey FOREIGN KEY (userId) REFERENCES users (id),
		    CONSTRAINT userLanguageSkills_languageSkillId_fkey FOREIGN KEY (languageSkillId) REFERENCES languageSkills (id)
		);
		ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);
		ALTER TABLE users ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE users ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE users ALTER COLUMN deleted_at SET DEFAULT NULL;
		ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
		ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
		ALTER TABLE userLanguageSkills ALTER COLUMN id SET DEFAULT nextval('userLanguageSkills_id_seq'::regclass);
		ALTER TABLE userLanguageSkills ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE userLanguageSkills ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE userLanguageSkills ALTER COLUMN deleted_at SET DEFAULT NULL;
    `)

	return nil
}

func (m *UsersMigration1684787622) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'UsersMigration1684787622';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS userLanguageSkills CASCADE;
    `)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS users CASCADE;
    `)
	database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).Drop(context.Background())
	return nil
}
