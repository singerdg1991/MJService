package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/geartype/models"
	"github.com/hoitek/Maja-Service/internal/geartype/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/geartype/config"
)

type GearTypeHandler struct {
	GearTypeService service.GearTypeService
}

func NewGearTypeHandler(r *mux.Router, s service.GearTypeService) (GearTypeHandler, error) {
	geartypeHandler := GearTypeHandler{
		GearTypeService: s,
	}
	if r == nil {
		return GearTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.GearTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.GearTypeConfig.ApiVersion1).Subrouter()

	rv1.Handle("/geartypes", geartypeHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/geartypes", geartypeHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/geartypes", geartypeHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/geartypes/{name}", geartypeHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/geartypes/csv/download", geartypeHandler.Download()).Methods(http.MethodGet)
	return geartypeHandler, nil
}

/*
* @apiTag: geartype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: geartypes
* @apiResponseRef: GearTypesQueryResponse
* @apiSummary: Query geartypes
* @apiParametersRef: GearTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: GearTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: GearTypesQueryNotFoundResponse
 */
func (h *GearTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.GearTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		geartypes, err := h.GearTypeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(geartypes)
	})
}

/*
* @apiTag: geartype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: geartypes/csv/download
* @apiResponseRef: GearTypesQueryResponse
* @apiSummary: Query geartypes
* @apiParametersRef: GearTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: GearTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: GearTypesQueryNotFoundResponse
 */
func (h *GearTypeHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.GearTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		geartypes, err := h.GearTypeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(geartypes)
	})
}

/*
 * @apiTag: geartype
 * @apiPath: /geartypes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: GearTypesCreateRequestBody
 * @apiResponseRef: GearTypesCreateResponse
 * @apiSummary: Create geartype
 * @apiDescription: Create geartype
 */
func (h *GearTypeHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.GearTypesCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		user, err := h.GearTypeService.Create(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: geartype
 * @apiPath: /geartypes/{name}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: GearTypesUpdateRequestParams
 * @apiRequestRef: GearTypesCreateRequestBody
 * @apiResponseRef: GearTypesCreateResponse
 * @apiSummary: Update geartype
 * @apiDescription: Update geartype
 */
func (h *GearTypeHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.GearTypesUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.GearTypesCreateRequestBody{}
		errs = payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.GearTypeService.Update(payload, params.Name)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: geartype
 * @apiPath: /geartypes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: GearTypesDeleteRequestBody
 * @apiResponseRef: GearTypesCreateResponse
 * @apiSummary: Delete geartype
 * @apiDescription: Delete geartype
 */
func (h *GearTypeHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.GearTypesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.GearTypeService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
