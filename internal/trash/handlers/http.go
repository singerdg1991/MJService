package handlers

import (
	"errors"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/trash"
	"github.com/hoitek/Maja-Service/internal/trash/models"
	tPorts "github.com/hoitek/Maja-Service/internal/trash/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/trash/config"
)

type TrashHandler struct {
	TrashService tPorts.TrashService
	UserService  uPorts.UserService
}

func NewTrashHandler(r *mux.Router, t tPorts.TrashService, u uPorts.UserService) (TrashHandler, error) {
	trashHandler := TrashHandler{
		TrashService: t,
		UserService:  u,
	}
	if r == nil {
		return TrashHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.TrashConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.TrashConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/trashes/models", trashHandler.Models()).Methods(http.MethodGet)
	rAuth.Handle("/trashes", trashHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/trashes", trashHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/trashes", trashHandler.Delete()).Methods(http.MethodDelete)

	return trashHandler, nil
}

/*
* @apiTag: trash
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: trashes
* @apiResponseRef: TrashesQueryResponse
* @apiSummary: Query trashes models
* @apiParametersRef: TrashesModelsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TrashesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TrashesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *TrashHandler) Models() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		return response.Success(trash.TrashModels)
	})
}

/*
* @apiTag: trash
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: trashes
* @apiResponseRef: TrashesQueryResponse
* @apiSummary: Query trashes
* @apiParametersRef: TrashesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TrashesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TrashesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *TrashHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TrashesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		trashes, err := h.TrashService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(trashes)
	})
}

/*
 * @apiTag: trash
 * @apiPath: /trashes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: TrashesCreateRequestBody
 * @apiResponseRef: TrashesCreateResponse
 * @apiSummary: Create trash
 * @apiDescription: Create trash
 * @apiSecurity: apiKeySecurity
 */
func (h *TrashHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TrashesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}

		// Set created by
		payload.CreatedBy = models.TrashesCreateRequestBodyCreatedBy{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}

		// Create trash
		data, err := h.TrashService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Return created trash
		return response.Created(data)
	})
}

/*
 * @apiTag: trash
 * @apiPath: /trashes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: TrashesDeleteRequestBody
 * @apiResponseRef: TrashesCreateResponse
 * @apiSummary: Delete trash
 * @apiDescription: Delete trash
 * @apiSecurity: apiKeySecurity
 */
func (h *TrashHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TrashesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.TrashService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
