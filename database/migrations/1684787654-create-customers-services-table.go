package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomersServicesMigration1684787654 struct {
}

func NewCustomersServicesMigration1684787654() *CustomersServicesMigration1684787654 {
	return &CustomersServicesMigration1684787654{}
}

func (m *CustomersServicesMigration1684787654) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerServices (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			serviceId INT NOT NULL,
			serviceTypeId INT NOT NULL,
			gradeId INT NOT NULL,
			nurseWishId INT NOT NULL,
			reportType VARCHAR(255) NOT NULL,
            timeValue TIME NOT NULL,
            repeat VARCHAR(255) NOT NULL,
			visitType VARCHAR(255) NOT NULL,
            serviceLengthMinute INT NOT NULL,
            startTimeRange TIME NOT NULL,
            endTimeRange TIME NOT NULL,
            description TEXT DEFAULT NULL,
            paymentMethod VARCHAR(255) DEFAULT 'own',
			homeCareFee INT DEFAULT NULL,
			cityCouncilFee INT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
            CONSTRAINT fk_customerServices_customerId FOREIGN KEY (customerId) REFERENCES customers(id),
            CONSTRAINT fk_customerServices_serviceId FOREIGN KEY (serviceId) REFERENCES services(id),
            CONSTRAINT fk_customerServices_serviceTypeId FOREIGN KEY (serviceTypeId) REFERENCES serviceTypes(id),
            CONSTRAINT fk_customerServices_gradeId FOREIGN KEY (gradeId) REFERENCES servicegrades(id),
            CONSTRAINT fk_customerServices_nurseWishId FOREIGN KEY (nurseWishId) REFERENCES staffs(id)
		);
    `)
	if err != nil {
		log.Println("----------------", err.Error())
	}
	return nil
}

func (m *CustomersServicesMigration1684787654) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomersServicesMigration1684787654'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerServices CASCADE;
    `)
	return nil
}
