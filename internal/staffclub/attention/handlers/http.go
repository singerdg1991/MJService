package handlers

import (
	"errors"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	rPorts "github.com/hoitek/Maja-Service/internal/punishment/ports"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type AttentionHandler struct {
	UserService       uPorts.UserService
	AttentionService  ports.AttentionService
	PunishmentService rPorts.PunishmentService
}

func NewAttentionHandler(r *mux.Router, s ports.AttentionService, rs rPorts.PunishmentService, us uPorts.UserService) (AttentionHandler, error) {
	attentionHandler := AttentionHandler{
		UserService:       us,
		AttentionService:  s,
		PunishmentService: rs,
	}
	if r == nil {
		return AttentionHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.AttentionConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.AttentionConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/staffclub/attentions", attentionHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/staffclub/attentions", attentionHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/staffclub/attentions", attentionHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/staffclub/attentions/{id}", attentionHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/staffclub/attentions/csv/download", attentionHandler.Download()).Methods(http.MethodGet)

	return attentionHandler, nil
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/attentions
* @apiResponseRef: AttentionsQueryResponse
* @apiSummary: Query attentions
* @apiParametersRef: AttentionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AttentionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AttentionsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *AttentionHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AttentionsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query attentions
		attentions, err := h.AttentionService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(attentions)
	})
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/attentions/csv/download
* @apiResponseRef: AttentionsQueryResponse
* @apiSummary: Query attentions
* @apiParametersRef: AttentionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AttentionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AttentionsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *AttentionHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AttentionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		attentions, err := h.AttentionService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(attentions)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/attentions
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: AttentionsCreateRequestBody
 * @apiResponseRef: AttentionsCreateResponse
 * @apiSummary: Create attention
 * @apiDescription: Create attention
 * @apiSecurity: apiKeySecurity
 */
func (h *AttentionHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AttentionsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Attention already exists with the same attention number
		attentionByNumber, _ := h.AttentionService.FindByAttentionNumber(int64(payload.AttentionNumber))
		if attentionByNumber != nil {
			return response.ErrorBadRequest(nil, "Attention with this attention number already exists")
		}

		// Find punishment
		punishment, _ := h.PunishmentService.FindByID(int64(payload.PunishmentID))
		if punishment == nil {
			return response.ErrorInternalServerError(nil, "Punishment not found")
		}
		payload.Punishment = &domain.AttentionPunishment{
			ID:          punishment.ID,
			Name:        punishment.Name,
			Description: *punishment.Description,
		}

		// Create attention
		attention, err := h.AttentionService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(attention)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/attentions/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: AttentionsUpdateRequestParams
 * @apiRequestRef: AttentionsCreateRequestBody
 * @apiResponseRef: AttentionsCreateResponse
 * @apiSummary: Update attention
 * @apiDescription: Update attention
 * @apiSecurity: apiKeySecurity
 */
func (h *AttentionHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.AttentionsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.AttentionsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Attention already exists
		attention, _ := h.AttentionService.FindByID(int64(params.ID))
		if attention == nil {
			return response.ErrorBadRequest(nil, "Attention does not exist")
		}

		// Check if Attention already exists with the same attention number
		attentionByNumber, _ := h.AttentionService.FindByAttentionNumber(int64(attention.AttentionNumber))
		if attentionByNumber != nil && int64(attentionByNumber.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Attention with this attention number already exists")
		}

		// Find punishment
		punishment, _ := h.PunishmentService.FindByID(int64(payload.PunishmentID))
		if punishment == nil {
			return response.ErrorInternalServerError(nil, "Punishment not found")
		}
		payload.Punishment = &domain.AttentionPunishment{
			ID:          punishment.ID,
			Name:        punishment.Name,
			Description: *punishment.Description,
		}

		// Update Attention
		data, err := h.AttentionService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/attentions
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: AttentionsDeleteRequestBody
 * @apiResponseRef: AttentionsDeleteResponse
 * @apiSummary: Delete attention
 * @apiDescription: Delete attention
 * @apiSecurity: apiKeySecurity
 */
func (h *AttentionHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AttentionsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.AttentionService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
