package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/company/models"
	"github.com/hoitek/Maja-Service/internal/company/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/company/config"
)

type CompanyHandler struct {
	CompanyService service.CompanyService
}

func NewCompanyHandler(r *mux.Router, s service.CompanyService) (CompanyHandler, error) {
	companyHandler := CompanyHandler{
		CompanyService: s,
	}
	if r == nil {
		return CompanyHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.CompanyConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.CompanyConfig.ApiVersion1).Subrouter()

	rv1.Handle("/companies", companyHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/companies", companyHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/companies", companyHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/companies/{name}", companyHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/companies/csv/download", companyHandler.Download()).Methods(http.MethodGet)
	return companyHandler, nil
}

/*
* @apiTag: company
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: companies
* @apiResponseRef: CompaniesQueryResponse
* @apiSummary: Query companies
* @apiParametersRef: CompaniesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CompaniesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CompaniesQueryNotFoundResponse
 */
func (h *CompanyHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CompaniesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		companies, err := h.CompanyService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(companies)
	})
}

/*
* @apiTag: company
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: companies/csv/download
* @apiResponseRef: CompaniesQueryResponse
* @apiSummary: Query companies
* @apiParametersRef: CompaniesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CompaniesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CompaniesQueryNotFoundResponse
 */
func (h *CompanyHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CompaniesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		companies, err := h.CompanyService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(companies)
	})
}

/*
 * @apiTag: company
 * @apiPath: /companies
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: CompaniesCreateRequestBody
 * @apiResponseRef: CompaniesCreateResponse
 * @apiSummary: Create company
 * @apiDescription: Create company
 */
func (h *CompanyHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CompaniesCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		user, err := h.CompanyService.Create(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: company
 * @apiPath: /companies/{name}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: CompaniesUpdateRequestParams
 * @apiRequestRef: CompaniesCreateRequestBody
 * @apiResponseRef: CompaniesCreateResponse
 * @apiSummary: Update company
 * @apiDescription: Update company
 */
func (h *CompanyHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.CompaniesUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.CompaniesCreateRequestBody{}
		errs = payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.CompanyService.Update(payload, params.Name)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: company
 * @apiPath: /companies
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: CompaniesDeleteRequestBody
 * @apiResponseRef: CompaniesCreateResponse
 * @apiSummary: Delete company
 * @apiDescription: Delete company
 */
func (h *CompanyHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CompaniesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.CompanyService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
