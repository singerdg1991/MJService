package migrations

import "github.com/hoitek/Maja-Service/database"

type LanguageSkillsAddDescriptionColumnMigration1684787635 struct {
}

func NewLanguageSkillsAddDescriptionColumnMigration1684787635() *LanguageSkillsAddDescriptionColumnMigration1684787635 {
	return &LanguageSkillsAddDescriptionColumnMigration1684787635{}
}

func (m *LanguageSkillsAddDescriptionColumnMigration1684787635) MigrateUp() error {
	database.PostgresDB.Exec(`
		ALTER TABLE IF EXISTS languageskills ADD COLUMN IF NOT EXISTS description VARCHAR(255) DEFAULT NULL;
    `)
	return nil
}

func (m *LanguageSkillsAddDescriptionColumnMigration1684787635) MigrateDown() error {
	database.PostgresDB.Exec(`
		DELETE FROM _migrations WHERE name = 'LanguageSkillsAddDescriptionColumnMigration1684787635';
	`)
	database.PostgresDB.Exec(`
		ALTER TABLE languageskills DROP COLUMN description;
    `)
	return nil
}
