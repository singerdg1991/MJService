package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type ContractTypesMigration1684787639 struct {
}

func NewContractTypesMigration1684787639() *ContractTypesMigration1684787639 {
	return &ContractTypesMigration1684787639{}
}

func (m *ContractTypesMigration1684787639) MigrateUp() error {
	_, ee := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS contracttypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
            description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE contractTypes ALTER COLUMN id SET DEFAULT nextval('contractTypes_id_seq'::regclass);
		ALTER TABLE contractTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE contractTypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO contractTypes
			(id, name, description, created_at, updated_at, deleted_at)
		VALUES
			(1, 'permanent', '', '2020-01-01', '2020-01-01', NULL),
			(2, 'osa-aikainen', '', '2020-01-01', '2020-01-01', NULL),
			(3, 'määräaikainen', '', '2020-01-01', '2020-01-01', NULL),
			(4, 'zero contract', '', '2020-01-01', '2020-01-01', NULL)
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('contractTypes_id_seq', (SELECT MAX(id) FROM contractTypes));
    `)
	log.Println(ee)
	return nil
}

func (m *ContractTypesMigration1684787639) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'ContractTypesMigration1684787639';`)
	database.PostgresDB.Exec(`DROP TABLE IF EXISTS contracttypes CASCADE;`)
	return nil
}
