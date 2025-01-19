package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomerRelativesMigration1684787675 struct {
}

func NewCustomerRelativesMigration1684787675() *CustomerRelativesMigration1684787675 {
	return &CustomerRelativesMigration1684787675{}
}

func (m *CustomerRelativesMigration1684787675) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerRelatives (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			addressCityId INT DEFAULT NULL,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			phoneNumber VARCHAR(255) NOT NULL,
			relation VARCHAR(255) NOT NULL,
			addressName VARCHAR(255) NOT NULL,
			addressStreet VARCHAR(255) NOT NULL,
			addressBuildingNumber VARCHAR(255) NOT NULL,
			addressPostalCode VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customerRelatives_customerId FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_customerRelatives_addressCityId FOREIGN KEY (addressCityId) REFERENCES cities(id) ON DELETE SET NULL
		);
    `)
	if err != nil {
		log.Println("************ Error in creating customerRelatives table ************", err.Error())
	}
	return nil
}

func (m *CustomerRelativesMigration1684787675) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomerRelativesMigration1684787675'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerRelatives CASCADE;
    `)
	return nil
}
