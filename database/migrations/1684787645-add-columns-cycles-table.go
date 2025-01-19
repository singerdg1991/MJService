package migrations

import "github.com/hoitek/Maja-Service/database"

type CyclesAddColumnsMigration1684787645 struct {
}

func NewCyclesAddColumnsMigration1684787645() *CyclesAddColumnsMigration1684787645 {
	return &CyclesAddColumnsMigration1684787645{}
}

func (m *CyclesAddColumnsMigration1684787645) MigrateUp() error {
	database.PostgresDB.Exec(`
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS sectionId INT NOT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS start_datetime TIMESTAMP DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS end_datetime TIMESTAMP DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS periodLength VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftMorningStartTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftMorningEndTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftEveningStartTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftEveningEndTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftNightStartTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS shiftNightEndTime VARCHAR(255) DEFAULT NULL;
		ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS freeze_period_datetime TIMESTAMP DEFAULT NULL;
        ALTER TABLE IF EXISTS cycles ADD COLUMN IF NOT EXISTS wishDays INT NOT NULL DEFAULT 0;
    `)
	return nil
}

func (m *CyclesAddColumnsMigration1684787645) MigrateDown() error {
	database.PostgresDB.Exec(`
		DELETE FROM _migrations WHERE name = 'CyclesAddColumnsMigration1684787645';
	`)
	database.PostgresDB.Exec(`
		ALTER TABLE cycles DROP COLUMN start_datetime;
		ALTER TABLE cycles DROP COLUMN end_datetime;
		ALTER TABLE cycles DROP COLUMN periodLength;
		ALTER TABLE cycles DROP COLUMN shiftMorningStartTime;
		ALTER TABLE cycles DROP COLUMN shiftMorningEndTime;
		ALTER TABLE cycles DROP COLUMN shiftEveningStartTime;
		ALTER TABLE cycles DROP COLUMN shiftEveningEndTime;
		ALTER TABLE cycles DROP COLUMN shiftNightStartTime;
		ALTER TABLE cycles DROP COLUMN shiftNightEndTime;
		ALTER TABLE cycles DROP COLUMN freeze_period_datetime;
        ALTER TABLE cycles DROP COLUMN wishDays;
    `)
	return nil
}
