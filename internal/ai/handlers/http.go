package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/ai/models"
	"github.com/hoitek/Maja-Service/internal/ai/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/ai/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type AIHandler struct {
	AIService   ports.AIService
	UserService uPorts.UserService
}

func NewAIHandler(r *mux.Router, a ports.AIService, u uPorts.UserService) (AIHandler, error) {
	aiHandler := AIHandler{
		AIService:   a,
		UserService: u,
	}
	if r == nil {
		return AIHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.AIConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.AIConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rv1.Handle("/ai", aiHandler.Create()).Methods(http.MethodPost)

	return aiHandler, nil
}

/*
 * @apiTag: ai
 * @apiPath: /ai
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: AIsCreateRequestBody
 * @apiResponseRef: AIsCreateResponse
 * @apiSummary: Give response from AI
 * @apiDescription: Give response from AI
 * @apiSecurity: apiKeySecurity
 */
func (h *AIHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AIsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		resp, err := h.AIService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		ai := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

		return response.Success(ai)
	})
}
