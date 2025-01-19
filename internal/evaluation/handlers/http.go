package handlers

import (
	"errors"
	"net/http"

	staffPorts "github.com/hoitek/Maja-Service/internal/staff/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/evaluation/models"
	"github.com/hoitek/Maja-Service/internal/evaluation/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/evaluation/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type EvaluationHandler struct {
	EvaluationService ports.EvaluationService
	StaffService      staffPorts.StaffService
	UserService       uPorts.UserService
}

func NewEvaluationHandler(r *mux.Router, s ports.EvaluationService, ss staffPorts.StaffService, us uPorts.UserService) (EvaluationHandler, error) {
	evaluationHandler := EvaluationHandler{
		EvaluationService: s,
		StaffService:      ss,
		UserService:       us,
	}
	if r == nil {
		return EvaluationHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.EvaluationConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.EvaluationConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/evaluations", evaluationHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/evaluations", evaluationHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/evaluations", evaluationHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/evaluations/{id}", evaluationHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/evaluations/csv/download", evaluationHandler.Download()).Methods(http.MethodGet)

	return evaluationHandler, nil
}

/*
* @apiTag: evaluation
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: evaluations
* @apiResponseRef: EvaluationsQueryResponse
* @apiSummary: Query evaluations
* @apiParametersRef: EvaluationsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: EvaluationsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: EvaluationsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *EvaluationHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.EvaluationsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query evaluations
		evaluations, err := h.EvaluationService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(evaluations)
	})
}

/*
* @apiTag: evaluation
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: evaluations/csv/download
* @apiResponseRef: EvaluationsQueryResponse
* @apiSummary: Query evaluations
* @apiParametersRef: EvaluationsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: EvaluationsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: EvaluationsQueryNotFoundResponse
 */
func (h *EvaluationHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.EvaluationsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		evaluations, err := h.EvaluationService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(evaluations)
	})
}

/*
 * @apiTag: evaluation
 * @apiPath: /evaluations
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: EvaluationsCreateRequestBody
 * @apiResponseRef: EvaluationsCreateResponse
 * @apiSummary: Create evaluation
 * @apiDescription: Create evaluation
 * @apiSecurity: apiKeySecurity
 */
func (h *EvaluationHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EvaluationsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find evaluation by name
		//res, _ := h.EvaluationService.FindByTitle(payload.Title)
		//if res != nil {
		//	return response.ErrorInternalServerError(nil, "Evaluation already exists")
		//}

		// Find staff by id
		// TODO: Uncomment this when staff service is ready
		//staff, _ := h.StaffService.FindByID(payload.StaffID)
		//if staff == nil {
		//	return response.ErrorInternalServerError(nil, "Staff does not exist")
		//}

		// Create evaluation
		evaluation, err := h.EvaluationService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(evaluation)
	})
}

/*
 * @apiTag: evaluation
 * @apiPath: /evaluations/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: EvaluationsUpdateRequestParams
 * @apiRequestRef: EvaluationsCreateRequestBody
 * @apiResponseRef: EvaluationsCreateResponse
 * @apiSummary: Update evaluation
 * @apiDescription: Update evaluation
 * @apiSecurity: apiKeySecurity
 */
func (h *EvaluationHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.EvaluationsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.EvaluationsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Evaluation already exists
		res, _ := h.EvaluationService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Evaluation does not exist")
		}

		// Check if Evaluation already exists with the same title
		//res, _ = h.EvaluationService.FindByTitle(payload.Title)
		//log.Println(res, payload.Title, params.ID)
		//if res != nil && int64(res.ID) != int64(params.ID) {
		//	return response.ErrorBadRequest(nil, "Evaluation already exists")
		//}

		// Find staff by id
		// TODO: Uncomment this when staff service is ready
		//staff, _ := h.StaffService.FindByID(payload.StaffID)
		//if staff == nil {
		//	return response.ErrorInternalServerError(nil, "Staff does not exist")
		//}

		// Update Evaluation
		data, err := h.EvaluationService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: evaluation
 * @apiPath: /evaluations
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: EvaluationsDeleteRequestBody
 * @apiResponseRef: EvaluationsDeleteResponse
 * @apiSummary: Delete evaluation
 * @apiDescription: Delete evaluation
 * @apiSecurity: apiKeySecurity
 */
func (h *EvaluationHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EvaluationsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.EvaluationService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
