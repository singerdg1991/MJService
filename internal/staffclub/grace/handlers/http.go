package handlers

import (
	"errors"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	rPorts "github.com/hoitek/Maja-Service/internal/reward/ports"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type GraceHandler struct {
	UserService   uPorts.UserService
	GraceService  ports.GraceService
	RewardService rPorts.RewardService
}

func NewGraceHandler(r *mux.Router, s ports.GraceService, rs rPorts.RewardService, us uPorts.UserService) (GraceHandler, error) {
	graceHandler := GraceHandler{
		UserService:   us,
		GraceService:  s,
		RewardService: rs,
	}
	if r == nil {
		return GraceHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.GraceConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.GraceConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/staffclub/graces", graceHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/staffclub/graces", graceHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/staffclub/graces", graceHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/staffclub/graces/{id}", graceHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/staffclub/graces/csv/download", graceHandler.Download()).Methods(http.MethodGet)

	return graceHandler, nil
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/graces
* @apiResponseRef: GracesQueryResponse
* @apiSummary: Query graces
* @apiParametersRef: GracesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: GracesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: GracesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *GraceHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.GracesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query graces
		graces, err := h.GraceService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(graces)
	})
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/graces/csv/download
* @apiResponseRef: GracesQueryResponse
* @apiSummary: Query graces
* @apiParametersRef: GracesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: GracesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: GracesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *GraceHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.GracesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		graces, err := h.GraceService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(graces)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/graces
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: GracesCreateRequestBody
 * @apiResponseRef: GracesCreateResponse
 * @apiSummary: Create grace
 * @apiDescription: Create grace
 * @apiSecurity: apiKeySecurity
 */
func (h *GraceHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.GracesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Grace already exists with the same grace number
		graceByNumber, _ := h.GraceService.FindByGraceNumber(int64(payload.GraceNumber))
		if graceByNumber != nil {
			return response.ErrorBadRequest(nil, "Grace with this grace number already exists")
		}

		// Find reward
		reward, _ := h.RewardService.FindByID(int64(payload.RewardID))
		if reward == nil {
			return response.ErrorInternalServerError(nil, "Reward not found")
		}
		payload.Reward = &domain.GraceReward{
			ID:          reward.ID,
			Name:        reward.Name,
			Description: *reward.Description,
		}

		// Create grace
		grace, err := h.GraceService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(grace)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/graces/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: GracesUpdateRequestParams
 * @apiRequestRef: GracesCreateRequestBody
 * @apiResponseRef: GracesCreateResponse
 * @apiSummary: Update grace
 * @apiDescription: Update grace
 * @apiSecurity: apiKeySecurity
 */
func (h *GraceHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.GracesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.GracesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Grace already exists
		grace, _ := h.GraceService.FindByID(int64(params.ID))
		if grace == nil {
			return response.ErrorBadRequest(nil, "Grace does not exist")
		}

		// Check if Grace already exists with the same grace number
		graceByNumber, _ := h.GraceService.FindByGraceNumber(int64(grace.GraceNumber))
		if graceByNumber != nil && int64(graceByNumber.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Grace with this grace number already exists")
		}

		// Find reward
		reward, _ := h.RewardService.FindByID(int64(payload.RewardID))
		if reward == nil {
			return response.ErrorInternalServerError(nil, "Reward not found")
		}
		payload.Reward = &domain.GraceReward{
			ID:          reward.ID,
			Name:        reward.Name,
			Description: *reward.Description,
		}

		// Update Grace
		data, err := h.GraceService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/graces
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: GracesDeleteRequestBody
 * @apiResponseRef: GracesDeleteResponse
 * @apiSummary: Delete grace
 * @apiDescription: Delete grace
 * @apiSecurity: apiKeySecurity
 */
func (h *GraceHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.GracesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.GraceService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
