package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type AddressesMigration1684787626 struct {
}

func NewAddressesMigration1684787626() *AddressesMigration1684787626 {
	return &AddressesMigration1684787626{}
}

func (m *AddressesMigration1684787626) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
			staffId INT DEFAULT NULL,
			customerId INT DEFAULT NULL,
			cityId INT NOT NULL,
			street VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			postalCode VARCHAR(255) DEFAULT NULL,
		    buildingNumber VARCHAR(255) DEFAULT NULL,
            isDeliveryAddress BOOLEAN DEFAULT FALSE,
            isMainAddress BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
		    CONSTRAINT fk_city_id FOREIGN KEY (cityId) REFERENCES cities(id) ON DELETE CASCADE
		);
		ALTER TABLE addresses ALTER COLUMN id SET DEFAULT nextval('addresses_id_seq'::regclass);
		ALTER TABLE addresses ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE addresses ALTER COLUMN updated_at SET DEFAULT now();
    `)

	return nil
}

func (m *AddressesMigration1684787626) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'AddressesMigration1684787626';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS addresses CASCADE;
    `)
	return nil
}
