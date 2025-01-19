package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/diagnose/models"
	stPorts "github.com/hoitek/Maja-Service/internal/diagnose/ports"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/diagnose/config"
)

type DiagnoseHandler struct {
	DiagnoseService stPorts.DiagnoseService
}

func NewDiagnoseHandler(r *mux.Router, s stPorts.DiagnoseService) (DiagnoseHandler, error) {
	diagnoseHandler := DiagnoseHandler{
		DiagnoseService: s,
	}
	if r == nil {
		return DiagnoseHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.DiagnoseConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.DiagnoseConfig.ApiVersion1).Subrouter()

	rv1.Handle("/diagnoses", diagnoseHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/diagnoses", diagnoseHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/diagnoses", diagnoseHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/diagnoses/{id}", diagnoseHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/diagnoses/csv/download", diagnoseHandler.Download()).Methods(http.MethodGet)
	return diagnoseHandler, nil
}

/*
* @apiTag: diagnose
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: diagnoses
* @apiResponseRef: DiagnosesQueryResponse
* @apiSummary: Query diagnoses
* @apiParametersRef: DiagnosesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: DiagnosesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: DiagnosesQueryNotFoundResponse
 */
func (h *DiagnoseHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.DiagnosesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		diagnoses, err := h.DiagnoseService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(diagnoses)
	})
}

/*
* @apiTag: diagnose
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: diagnoses/csv/download
* @apiResponseRef: DiagnosesQueryResponse
* @apiSummary: Query diagnoses
* @apiParametersRef: DiagnosesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: DiagnosesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: DiagnosesQueryNotFoundResponse
 */
func (h *DiagnoseHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.DiagnosesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		diagnoses, err := h.DiagnoseService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(diagnoses)
	})
}

/*
 * @apiTag: diagnose
 * @apiPath: /diagnoses
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: DiagnosesCreateRequestBody
 * @apiResponseRef: DiagnosesCreateResponse
 * @apiSummary: Create diagnose
 * @apiDescription: Create diagnose
 */
func (h *DiagnoseHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.DiagnosesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		diagnose, err := h.DiagnoseService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(diagnose)
	})
}

/*
 * @apiTag: diagnose
 * @apiPath: /diagnoses/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: DiagnosesUpdateRequestParams
 * @apiRequestRef: DiagnosesCreateRequestBody
 * @apiResponseRef: DiagnosesCreateResponse
 * @apiSummary: Update diagnose
 * @apiDescription: Update diagnose
 */
func (h *DiagnoseHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.DiagnosesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.DiagnosesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.DiagnoseService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: diagnose
 * @apiPath: /diagnoses
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: DiagnosesDeleteRequestBody
 * @apiResponseRef: DiagnosesDeleteResponse
 * @apiSummary: Delete diagnose
 * @apiDescription: Delete diagnose
 */
func (h *DiagnoseHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.DiagnosesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.DiagnoseService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
