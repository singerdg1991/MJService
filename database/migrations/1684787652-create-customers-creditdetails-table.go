package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type CustomersCreditDetailsMigration1684787652 struct {
}

func NewCustomersCreditDetailsMigration1684787652() *CustomersCreditDetailsMigration1684787652 {
	return &CustomersCreditDetailsMigration1684787652{}
}

func (m *CustomersCreditDetailsMigration1684787652) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerCreditDetails (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			billingAddressId INT NOT NULL,
			bankAccountNumber VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT now(),
            updated_at TIMESTAMP DEFAULT now(),
            deleted_at TIMESTAMP DEFAULT NULL,
            CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
            CONSTRAINT fk_billing_address_id FOREIGN KEY (billingAddressId) REFERENCES addresses(id) ON DELETE CASCADE
		);
    `)

	return nil
}

func (m *CustomersCreditDetailsMigration1684787652) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersCreditDetailsMigration1684787652'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerCreditDetails;
    `)
	return nil
}
