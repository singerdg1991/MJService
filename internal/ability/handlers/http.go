package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/ability/models"
	"github.com/hoitek/Maja-Service/internal/ability/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/ability/config"
)

type AbilityHandler struct {
	AbilityService service.AbilityService
}

func NewAbilityHandler(r *mux.Router, s service.AbilityService) (AbilityHandler, error) {
	abilityHandler := AbilityHandler{
		AbilityService: s,
	}
	if r == nil {
		return AbilityHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.AbilityConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.AbilityConfig.ApiVersion1).Subrouter()

	rv1.Handle("/abilities", abilityHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/abilities", abilityHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/abilities", abilityHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/abilities/{name}", abilityHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/abilities/csv/download", abilityHandler.Download()).Methods(http.MethodGet)
	return abilityHandler, nil
}

/*
* @apiTag: ability
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: abilities
* @apiResponseRef: AbilitiesQueryResponse
* @apiSummary: Query abilities
* @apiParametersRef: AbilitiesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AbilitiesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AbilitiesQueryNotFoundResponse
 */
func (h *AbilityHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AbilitiesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		abilities, err := h.AbilityService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(abilities)
	})
}

/*
* @apiTag: ability
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: abilities/csv/download
* @apiResponseRef: AbilitiesQueryResponse
* @apiSummary: Query abilities
* @apiParametersRef: AbilitiesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AbilitiesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AbilitiesQueryNotFoundResponse
 */
func (h *AbilityHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AbilitiesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		abilities, err := h.AbilityService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(abilities)
	})
}

/*
 * @apiTag: ability
 * @apiPath: /abilities
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: AbilitiesCreateRequestBody
 * @apiResponseRef: AbilitiesCreateResponse
 * @apiSummary: Create ability
 * @apiDescription: Create ability
 */
func (h *AbilityHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AbilitiesCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		user, err := h.AbilityService.Create(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: ability
 * @apiPath: /abilities/{name}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: AbilitiesUpdateRequestParams
 * @apiRequestRef: AbilitiesCreateRequestBody
 * @apiResponseRef: AbilitiesCreateResponse
 * @apiSummary: Update ability
 * @apiDescription: Update ability
 */
func (h *AbilityHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.AbilitiesUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.AbilitiesCreateRequestBody{}
		errs = payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.AbilityService.Update(payload, params.Name)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: ability
 * @apiPath: /abilities
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: AbilitiesDeleteRequestBody
 * @apiResponseRef: AbilitiesCreateResponse
 * @apiSummary: Delete ability
 * @apiDescription: Delete ability
 */
func (h *AbilityHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AbilitiesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.AbilityService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
