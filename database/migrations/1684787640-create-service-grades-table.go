package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type ServiceGradesMigration1684787640 struct {
}

func NewServiceGradesMigration1684787640() *ServiceGradesMigration1684787640 {
	return &ServiceGradesMigration1684787640{}
}

func (m *ServiceGradesMigration1684787640) MigrateUp() error {
	_, ee := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS servicegrades (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
            description VARCHAR(255) DEFAULT NULL,
            grade INTEGER NOT NULL DEFAULT 0,
            color VARCHAR(7) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE contractTypes ALTER COLUMN id SET DEFAULT nextval('contractTypes_id_seq'::regclass);
		ALTER TABLE contractTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE contractTypes ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('contractTypes_id_seq', (SELECT MAX(id) FROM contractTypes));
    `)
	log.Println(ee)
	return nil
}

func (m *ServiceGradesMigration1684787640) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ServiceGradesMigration1684787640';`)
	database.PostgresDB.Exec(`DROP TABLE IF EXISTS servicegrades CASCADE;`)
	return nil
}
