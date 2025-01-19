package migrations

import (
	"github.com/hoitek/Maja-Service/database"
)

type PaymentTypesMigration1684787621 struct {
}

func NewPaymentTypesMigration1684787621() *PaymentTypesMigration1684787621 {
	return &PaymentTypesMigration1684787621{}
}

func (m *PaymentTypesMigration1684787621) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS paymentTypes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE paymentTypes ALTER COLUMN id SET DEFAULT nextval('paymentTypes_id_seq'::regclass);
		ALTER TABLE paymentTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE paymentTypes ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO paymentTypes
			(id, name, created_at, updated_at, deleted_at)
		VALUES
			(1, 'per visit', '2020-01-01', '2020-01-01', '2020-01-01'),
			(2, 'per hour', '2020-01-01', '2020-01-01', '2020-01-01'),
			(3, 'per month', '2020-01-01', '2020-01-01', '2020-01-01'),
			(4, 'per year', '2020-01-01', '2020-01-01', '2020-01-01'),
			(5, 'per week', '2020-01-01', '2020-01-01', '2020-01-01')
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('paymentTypes_id_seq', (SELECT MAX(id) FROM paymentTypes));
    `)
	return nil
}

func (m *PaymentTypesMigration1684787621) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'PaymentTypesMigration1684787621';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS paymentTypes CASCADE;
    `)
	return nil
}
