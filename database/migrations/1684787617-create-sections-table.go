package migrations

import (
	"context"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/internal/section/domain"
	"github.com/hoitek/Maja-Service/internal/section/models"
	"github.com/hoitek/Maja-Service/internal/section/ports"
	"github.com/hoitek/Maja-Service/internal/section/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SectionsMigration1684787617 struct {
}

func NewSectionsMigration1684787617() *SectionsMigration1684787617 {
	return &SectionsMigration1684787617{}
}

func (m *SectionsMigration1684787617) MigrateUp() error {
	database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS sections (
			id SERIAL PRIMARY KEY,
			parentId INTEGER DEFAULT NULL,
			name VARCHAR(255) NOT NULL,
            color VARCHAR(255) DEFAULT NULL,
            description VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			FOREIGN KEY (parentId) REFERENCES sections(id) ON DELETE CASCADE
		);
		ALTER TABLE sections ALTER COLUMN id SET DEFAULT nextval('sections_id_seq'::regclass);
		ALTER TABLE sections ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE sections ALTER COLUMN updated_at SET DEFAULT now();
    `)
	var (
		currentTime = time.Now()
	)
	sections := []*domain.Section{
		{
			ID:     1,
			Parent: nil,
			Children: []*domain.Section{
				{
					ID:        6,
					Parent:    nil,
					Children:  nil,
					Name:      "Section 6",
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
					DeletedAt: nil,
				},
				{
					ID:        7,
					Parent:    nil,
					Children:  nil,
					Name:      "Section 7",
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
					DeletedAt: nil,
				},
			},
			Name:      "Section 1",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		{
			ID:        2,
			Parent:    nil,
			Children:  nil,
			Name:      "Section 2",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		{
			ID:        3,
			Parent:    nil,
			Children:  nil,
			Name:      "Section 3",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		{
			ID:        4,
			Parent:    nil,
			Children:  nil,
			Name:      "Section 4",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		{
			ID:        5,
			Parent:    nil,
			Children:  nil,
			Name:      "Section 5",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
	}

	// Insert sections to PostgreSQL
	for _, section := range sections {
		_, err := database.PostgresDB.Exec(`
			INSERT INTO sections
				(id, parentId, name, color, description, created_at, updated_at, deleted_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT(id) DO NOTHING;
		`, section.ID, section.ParentID, section.Name, section.Color, section.Description, section.CreatedAt, section.UpdatedAt, section.DeletedAt)
		if err != nil {
			return err
		}
		for _, child := range section.Children {
			_, err := database.PostgresDB.Exec(`
				INSERT INTO sections
					(id, parentId, name, color, description, created_at, updated_at, deleted_at)
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT(id) DO NOTHING;
			`, child.ID, child.ParentID, child.Name, child.Color, child.Description, child.CreatedAt, child.UpdatedAt, child.DeletedAt)
			if err != nil {
				return err
			}
		}
	}

	// Set sequence value
	_, err := database.PostgresDB.Exec(`SELECT setval('sections_id_seq', (SELECT MAX(id) FROM sections));`)
	if err != nil {
		return err
	}

	// Insert sections to MongoDB
	for _, row := range sections {
		// Check if section already exists
		filter := bson.M{
			"id": row.ID,
		}
		count, err := database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewSection().TableName()).CountDocuments(context.Background(), filter)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		sectionRepositoryPostgresDB := exp.TerIf[ports.SectionRepositoryPostgresDB](
			config.AppConfig.Environment == constants.ENVIRONMENT_TESTING,
			repositories.NewSectionRepositoryStub(),
			repositories.NewSectionRepositoryPostgresDB(database.PostgresDB),
		)

		s, err := sectionRepositoryPostgresDB.Query(&models.SectionsQueryRequestParams{
			ID: int(row.ID),
		})
		if err != nil {
			return err
		}
		if len(s) == 0 {
			continue
		}

		data, err := s[0].ToJson()
		if err != nil {
			return err
		}

		// Convert data to BSON document
		var sectionData bson.M
		err = bson.UnmarshalExtJSON([]byte(data), true, &sectionData)
		if err != nil {
			return err
		}

		// Create section
		_, err = database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewSection().TableName()).InsertOne(context.Background(), sectionData, options.InsertOne())
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SectionsMigration1684787617) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'SectionsMigration1684787617';`)
	database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS sections CASCADE;
    `)
	database.MongoDB.Database(config.AppConfig.DatabaseMongoDBName).Collection(domain.NewSection().TableName()).Drop(context.Background())
	return nil
}
