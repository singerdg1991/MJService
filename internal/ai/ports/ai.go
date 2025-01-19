package ports

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/hoitek/Maja-Service/internal/ai/models"
)

type AIService interface {
	Create(payload *models.AIsCreateRequestBody) (*genai.GenerateContentResponse, error)
}
