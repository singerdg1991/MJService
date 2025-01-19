package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/city/models"
	"github.com/hoitek/Maja-Service/internal/city/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/city/config"
)

type CityHandler struct {
	CityService service.CityService
}

func NewCityHandler(r *mux.Router, s service.CityService) (CityHandler, error) {
	cityHandler := CityHandler{
		CityService: s,
	}
	if r == nil {
		return CityHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.CityConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.CityConfig.ApiVersion1).Subrouter()

	rv1.Handle("/cities", cityHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/cities", cityHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/cities", cityHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/cities/{name}", cityHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/cities/csv/download", cityHandler.Download()).Methods(http.MethodGet)
	return cityHandler, nil
}

/*
* @apiTag: city
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: cities
* @apiResponseRef: CitiesQueryResponse
* @apiSummary: Query cities
* @apiParametersRef: CitiesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CitiesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CitiesQueryNotFoundResponse
 */
func (h *CityHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CitiesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		cities, err := h.CityService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(cities)
	})
}

/*
* @apiTag: city
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: cities/csv/download
* @apiResponseRef: CitiesQueryResponse
* @apiSummary: Query cities
* @apiParametersRef: CitiesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CitiesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CitiesQueryNotFoundResponse
 */
func (h *CityHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CitiesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		cities, err := h.CityService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(cities)
	})
}

/*
 * @apiTag: city
 * @apiPath: /cities
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: CitiesCreateRequestBody
 * @apiResponseRef: CitiesCreateResponse
 * @apiSummary: Create city
 * @apiDescription: Create city
 */
func (h *CityHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CitiesCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		user, err := h.CityService.Create(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: city
 * @apiPath: /cities/{name}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: CitiesUpdateRequestParams
 * @apiRequestRef: CitiesCreateRequestBody
 * @apiResponseRef: CitiesCreateResponse
 * @apiSummary: Update city
 * @apiDescription: Update city
 */
func (h *CityHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.CitiesUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.CitiesCreateRequestBody{}
		errs = payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.CityService.Update(payload, params.Name)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: city
 * @apiPath: /cities
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: CitiesDeleteRequestBody
 * @apiResponseRef: CitiesCreateResponse
 * @apiSummary: Delete city
 * @apiDescription: Delete city
 */
func (h *CityHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CitiesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.CityService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
