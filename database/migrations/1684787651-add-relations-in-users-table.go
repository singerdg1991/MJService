package migrations

import (
	"context"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"log"
)

type UsersAddRelationToCustomerIdAndStaffIdMigration1684787651 struct {
}

func NewUsersAddRelationToCustomerIdAndStaffIdMigration1684787651() *UsersAddRelationToCustomerIdAndStaffIdMigration1684787651 {
	return &UsersAddRelationToCustomerIdAndStaffIdMigration1684787651{}
}

func (m *UsersAddRelationToCustomerIdAndStaffIdMigration1684787651) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
		DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_customer_id') THEN
            ALTER TABLE users ADD CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE SET NULL;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_staff_id') THEN
            ALTER TABLE users ADD CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE SET NULL;
        END IF;
    END $$;
    `)
	if err != nil {
		log.Println("-------------------------------------", err)
	}

	return nil
}

func (m *UsersAddRelationToCustomerIdAndStaffIdMigration1684787651) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'UsersAddRelationToCustomerIdAndStaffIdMigration1684787651'`)
	database.PostgresDB.Exec(`
		ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_customer_id;
		ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_staff_id;
    `)
	database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).Drop(context.Background())
	return nil
}
