package migrations

import (
	"context"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"log"
)

type CustomersMigration1684787650 struct {
}

func NewCustomersMigration1684787650() *CustomersMigration1684787650 {
	return &CustomersMigration1684787650{}
}

func (m *CustomersMigration1684787650) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			userId INT DEFAULT NULL,
			responsibleNurseId INT DEFAULT NULL,
			motherLangIds JSONB DEFAULT NULL,
			nurseGenderWish VARCHAR(255) DEFAULT NULL,
			status VARCHAR(255) DEFAULT NULL,
            statusDate TIMESTAMP DEFAULT NULL,
            parkingInfo VARCHAR(255) DEFAULT NULL,
            extraExplanation TEXT DEFAULT NULL,
            hasLimitingTheRightToSelfDetermination BOOLEAN DEFAULT false,
            limitingTheRightToSelfDeterminationStartDate TIMESTAMP DEFAULT NULL,
            limitingTheRightToSelfDeterminationEndDate TIMESTAMP DEFAULT NULL,
            mobilityContract TEXT DEFAULT NULL,
            keyNo VARCHAR(255) DEFAULT NULL,
            paymentMethod VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_user_id FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_responsible_nurse_id FOREIGN KEY (responsibleNurseId) REFERENCES staffs(id) ON DELETE SET NULL
		);
		CREATE TABLE IF NOT EXISTS customerSections (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			sectionId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_section_id FOREIGN KEY (sectionId) REFERENCES sections(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS customerLimitations (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			limitationId INT NOT NULL,
			description TEXT DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_limitation_id FOREIGN KEY (limitationId) REFERENCES limitations(id) ON DELETE CASCADE
		);
		ALTER TABLE customers ALTER COLUMN id SET DEFAULT nextval('customers_id_seq'::regclass);
		ALTER TABLE customers ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE customers ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE customers ALTER COLUMN deleted_at SET DEFAULT NULL;
		ALTER TABLE customerSections ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE customerSections ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE customerSections ALTER COLUMN deleted_at SET DEFAULT NULL;
		ALTER TABLE customerLimitations ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE customerLimitations ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE customerLimitations ALTER COLUMN deleted_at SET DEFAULT NULL;
    `)
	if err != nil {
		log.Println("-------------------------------------", err)
	}
	return nil
}

func (m *CustomersMigration1684787650) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersMigration1684787650'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerLimitations CASCADE;
        DROP TABLE IF EXISTS customerSections CASCADE;
		DROP TABLE IF EXISTS customers CASCADE;
    `)
	database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).Drop(context.Background())
	return nil
}
