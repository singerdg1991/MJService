package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/keikkala/config"
	"github.com/hoitek/Maja-Service/internal/keikkala/models"
	"github.com/hoitek/Maja-Service/internal/keikkala/ports"
)

type KeikkalaHandler struct {
	KeikkalaService ports.KeikkalaService
	UserService     uPorts.UserService
}

func NewKeikkalaHandler(r *mux.Router, s ports.KeikkalaService, us uPorts.UserService) (KeikkalaHandler, error) {
	keikkalaHandler := KeikkalaHandler{
		KeikkalaService: s,
		UserService:     us,
	}
	if r == nil {
		return KeikkalaHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.KeikkalaConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.KeikkalaConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/keikkala", keikkalaHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/keikkala", keikkalaHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/keikkala", keikkalaHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/keikkala/shifts/statistics", keikkalaHandler.QueryShiftsStatistics()).Methods(http.MethodGet)

	return keikkalaHandler, nil
}

/*
* @apiTag: keikkala
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /keikkala
* @apiResponseRef: KeikkalasQueryResponse
* @apiSummary: Query keikkala
* @apiParametersRef: KeikkalasQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: KeikkalasQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: KeikkalasQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *KeikkalaHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.KeikkalasQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query keikkala
		keikkala, err := h.KeikkalaService.Query(queries)
		if err != nil {
			log.Printf("Error querying keikkala: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while querying keikkala, please try again later")
		}

		return response.Success(keikkala)
	})
}

/*
 * @apiTag: keikkala
 * @apiPath: /keikkala
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: KeikkalasCreateRequestBody
 * @apiResponseRef: KeikkalasCreateResponse
 * @apiSummary: Create keikkala
 * @apiDescription: Create keikkala
 * @apiSecurity: apiKeySecurity
 */
func (h *KeikkalaHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.KeikkalasCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Create keikkala shift
		keikkala, err := h.KeikkalaService.Create(payload)
		if err != nil {
			log.Printf("Error creating keikkala: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while creating keikkala, please try again later")
		}

		return response.Success(keikkala)
	})
}

/*
 * @apiTag: keikkala
 * @apiPath: /keikkala
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: KeikkalasDeleteRequestBody
 * @apiResponseRef: KeikkalasDeleteResponse
 * @apiSummary: Delete keikkala
 * @apiDescription: Delete keikkala
 * @apiSecurity: apiKeySecurity
 */
func (h *KeikkalaHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.KeikkalasDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		// Delete keikkala
		data, err := h.KeikkalaService.Delete(payload)
		if err != nil {
			log.Printf("Error deleting keikkala: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while deleting keikkala, please try again later")
		}

		return response.Success(data)
	})
}

/*
* @apiTag: keikkala
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /keikkala/shifts/statistics
* @apiResponseRef: KeikkalasQueryShiftStatisticsResponse
* @apiSummary: Query keikkala statistics
* @apiParametersRef: KeikkalasQueryShiftStatisticsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: KeikkalasQueryShiftStatisticsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: KeikkalasQueryShiftStatisticsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *KeikkalaHandler) QueryShiftsStatistics() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.KeikkalasQueryShiftStatisticsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query keikkala shift statistics
		keikkala, err := h.KeikkalaService.QueryShiftStatistics(queries)
		if err != nil {
			log.Printf("Error querying keikkala shift statistics: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while querying keikkala shift statistics, please try again later")
		}
		return response.Success(keikkala)
	})
}
