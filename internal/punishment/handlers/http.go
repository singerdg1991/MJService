package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/punishment/models"
	"github.com/hoitek/Maja-Service/internal/punishment/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/punishment/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type PunishmentHandler struct {
	UserService       uPorts.UserService
	PunishmentService ports.PunishmentService
}

func NewPunishmentHandler(r *mux.Router, s ports.PunishmentService, u uPorts.UserService) (PunishmentHandler, error) {
	punishmentHandler := PunishmentHandler{
		UserService:       u,
		PunishmentService: s,
	}
	if r == nil {
		return PunishmentHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.PunishmentConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.PunishmentConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/punishments", punishmentHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/punishments", punishmentHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/punishments", punishmentHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/punishments/{id}", punishmentHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/punishments/csv/download", punishmentHandler.Download()).Methods(http.MethodGet)

	return punishmentHandler, nil
}

/*
* @apiTag: punishment
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: punishments
* @apiResponseRef: PunishmentsQueryResponse
* @apiSummary: Query punishments
* @apiParametersRef: PunishmentsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PunishmentsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PunishmentsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *PunishmentHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PunishmentsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query punishments
		punishments, err := h.PunishmentService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(punishments)
	})
}

/*
* @apiTag: punishment
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: punishments/csv/download
* @apiResponseRef: PunishmentsQueryResponse
* @apiSummary: Query punishments
* @apiParametersRef: PunishmentsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PunishmentsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PunishmentsQueryNotFoundResponse
 */
func (h *PunishmentHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PunishmentsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		punishments, err := h.PunishmentService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(punishments)
	})
}

/*
 * @apiTag: punishment
 * @apiPath: /punishments
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: PunishmentsCreateRequestBody
 * @apiResponseRef: PunishmentsCreateResponse
 * @apiSummary: Create punishment
 * @apiDescription: Create punishment
 * @apiSecurity: apiKeySecurity
 */
func (h *PunishmentHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PunishmentsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find punishment by name
		res, _ := h.PunishmentService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Punishment already exists")
		}

		// Create punishment
		punishment, err := h.PunishmentService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(punishment)
	})
}

/*
 * @apiTag: punishment
 * @apiPath: /punishments/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: PunishmentsUpdateRequestParams
 * @apiRequestRef: PunishmentsCreateRequestBody
 * @apiResponseRef: PunishmentsCreateResponse
 * @apiSummary: Update punishment
 * @apiDescription: Update punishment
 * @apiSecurity: apiKeySecurity
 */
func (h *PunishmentHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.PunishmentsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.PunishmentsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.PunishmentService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.PunishmentService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.PunishmentService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: punishment
 * @apiPath: /punishments
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: PunishmentsDeleteRequestBody
 * @apiResponseRef: PunishmentsDeleteResponse
 * @apiSummary: Delete punishment
 * @apiDescription: Delete punishment
 * @apiSecurity: apiKeySecurity
 */
func (h *PunishmentHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PunishmentsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.PunishmentService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
