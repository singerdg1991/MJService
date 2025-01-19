package service

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/ai/constants"
	"github.com/hoitek/Maja-Service/internal/ai/models"
	"github.com/hoitek/Maja-Service/storage"
	"google.golang.org/api/option"
)

type AIService struct {
	MinIOStorage *storage.MinIO
}

func NewAIService(m *storage.MinIO) AIService {
	go minio.SetupMinIOStorage(constants.AI_BUCKET_NAME, m)
	return AIService{
		MinIOStorage: m,
	}
}

func (s *AIService) Create(payload *models.AIsCreateRequestBody) (*genai.GenerateContentResponse, error) {
	// Create a new client.
	var index = 0
	API_KEYS := []string{
		"AIzaSyBmJA27v8z2sHvV7Fr6x1A8qQSdkiPPZck",
		"AIzaSyBmJA27v8z2sHvV7Fr6x1A8qQSdkiPPZck",
		"AIzaSyBmJA27v8z2sHvV7Fr6x1A8qQSdkiPPZck",
		"AIzaSyBmJA27v8z2sHvV7Fr6x1A8qQSdkiPPZck",
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(API_KEYS[index]))
	if err != nil {
		log.Fatalf("genai.NewClient: %v\n", err)
	}
	defer client.Close()

	// Create a new model.
	model := client.GenerativeModel("gemini-1.5-flash")

	// Generate content.
	resp, err := model.GenerateContent(
		ctx,
		genai.Text(payload.Subject),
	)
	if err != nil {
		log.Fatalf("model.GenerateContent: %v\n", err)
	}

	// Increment the index.
	if index == len(API_KEYS)-1 {
		index = 0
	} else {
		index++
	}

	return resp, nil
}
