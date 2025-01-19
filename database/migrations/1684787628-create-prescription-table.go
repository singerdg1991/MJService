package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type PrescriptionsMigration1684787628 struct {
}

func NewPrescriptionsMigration1684787628() *PrescriptionsMigration1684787628 {
	return &PrescriptionsMigration1684787628{}
}

func (m *PrescriptionsMigration1684787628) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS prescriptions (
			id SERIAL PRIMARY KEY,
			customerId INT DEFAULT NULL,
			title VARCHAR(255) DEFAULT NULL,
			datetime TIMESTAMP DEFAULT NULL,
			doctorFullName VARCHAR(255) DEFAULT NULL,
            start_date TIMESTAMP DEFAULT NULL,
            end_date TIMESTAMP DEFAULT NULL,
			status VARCHAR(255) DEFAULT NULL,
			attachments JSONB DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_customer_id FOREIGN KEY (customerId) REFERENCES customers(id) ON DELETE CASCADE
		);
		ALTER TABLE prescriptions ALTER COLUMN id SET DEFAULT nextval('prescriptions_id_seq'::regclass);
		ALTER TABLE prescriptions ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE prescriptions ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('prescriptions_id_seq', (SELECT MAX(id) FROM prescriptions));
    `)
	if err != nil {
		log.Println("-------------------------------------", err)
	}
	return nil
}

func (m *PrescriptionsMigration1684787628) MigrateDown() error {
	database.PostgresDB.Exec("DELETE FROM _migrations WHERE name = 'PrescriptionsMigration1684787628'")
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS prescriptions CASCADE;
    `)
	return nil
}
