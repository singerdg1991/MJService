package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type CustomerStatusLogsMigration1684787679 struct {
}

func NewCustomerStatusLogsMigration1684787679() *CustomerStatusLogsMigration1684787679 {
	return &CustomerStatusLogsMigration1684787679{}
}

func (m *CustomerStatusLogsMigration1684787679) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS customerStatusLogs (
			id SERIAL PRIMARY KEY,
			customerId INT NOT NULL,
			statusValue TEXT NOT NULL,
			createdBy INT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customerStatusLogs_customerId FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE,
			CONSTRAINT fk_customerStatusLogs_createdBy FOREIGN KEY (createdBy) REFERENCES users(id) ON DELETE SET NULL
		);
    `)
	if err != nil {
		log.Println("************ Error in creating customerStatusLogs table ************", err.Error())
	}
	return nil
}

func (m *CustomerStatusLogsMigration1684787679) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'CustomerStatusLogsMigration1684787679'`)
	database.PostgresDB.Exec(`
		DROP TABLE IF EXISTS customerStatusLogs CASCADE;
    `)
	return nil
}
