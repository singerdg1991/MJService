package migrations

import (
	"log"

	"github.com/hoitek/Maja-Service/database"
)

type PushesMigration1687787691 struct {
}

func NewPushesMigration1687787691() *PushesMigration1687787691 {
	return &PushesMigration1687787691{}
}

func (m *PushesMigration1687787691) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS pushes (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
			endpoint TEXT NOT NULL,
			keysAuth TEXT NOT NULL,
			keysP256dh TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_pushes_userId FOREIGN KEY (userId) REFERENCES users(id)
		);
		ALTER TABLE pushes ALTER COLUMN id SET DEFAULT nextval('pushes_id_seq'::regclass);
		ALTER TABLE pushes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE pushes ALTER COLUMN updated_at SET DEFAULT now();
		SELECT setval('pushes_id_seq', (SELECT MAX(id) FROM pushes));
    `)
	log.Println(err)
	return nil
}

func (m *PushesMigration1687787691) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'PushesMigration1687787691';`)
	if err != nil {
		log.Printf("-------------------- Failed to delete migration PushesMigration1687787691: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS pushes;
    `)
	if err != nil {
		log.Printf("-------------------- Failed to delete pushes: %v\n", err)
	}
	return nil
}
