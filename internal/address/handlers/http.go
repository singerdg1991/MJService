package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/address/models"
	"github.com/hoitek/Maja-Service/internal/address/service"
	cPorts "github.com/hoitek/Maja-Service/internal/city/ports"
	customerPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	nPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/address/config"
)

type AddressHandler struct {
	AddressService  service.AddressService
	CityService     cPorts.CityService
	StaffService    nPorts.StaffService
	CustomerService customerPorts.CustomerService
	UserService     uPorts.UserService
}

func NewAddressHandler(r *mux.Router, s service.AddressService, c cPorts.CityService, n nPorts.StaffService, cs customerPorts.CustomerService, us uPorts.UserService) (AddressHandler, error) {
	addressHandler := AddressHandler{
		AddressService:  s,
		CityService:     c,
		StaffService:    n,
		CustomerService: cs,
		UserService:     us,
	}
	if r == nil {
		return AddressHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.AddressConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.AddressConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/addresses", addressHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/addresses", addressHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/addresses", addressHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/addresses/{id}", addressHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/addresses/csv/download", addressHandler.Download()).Methods(http.MethodGet)
	return addressHandler, nil
}

/*
* @apiTag: address
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: addresses
* @apiResponseRef: AddressesQueryResponse
* @apiSummary: Query addresses
* @apiParametersRef: AddressesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AddressesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AddressesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *AddressHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AddressesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		addresses, err := h.AddressService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(addresses)
	})
}

/*
* @apiTag: address
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: addresses/csv/download
* @apiResponseRef: AddressesQueryResponse
* @apiSummary: Query addresses
* @apiParametersRef: AddressesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: AddressesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: AddressesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *AddressHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.AddressesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		addresses, err := h.AddressService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(addresses)
	})
}

/*
 * @apiTag: address
 * @apiPath: /addresses
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: AddressesCreateRequestBody
 * @apiResponseRef: AddressesCreateResponse
 * @apiSummary: Create address
 * @apiDescription: Create address
* @apiSecurity: apiKeySecurity
*/
func (h *AddressHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AddressesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if city exists
		city, err := h.CityService.GetCityByID(int64(payload.City.ID))
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}
		if city == nil {
			return response.ErrorBadRequest(nil, "City not found")
		}

		// Check if staff or customer exists
		if payload.StaffID != nil {
			staff, err := h.StaffService.FindByID(*payload.StaffID)
			if err != nil {
				return response.ErrorBadRequest(nil, err.Error())
			}
			if staff == nil {
				return response.ErrorBadRequest(nil, "Staff not found")
			}
		}
		if payload.CustomerID != nil {
			customer, err := h.CustomerService.FindByID(*payload.CustomerID)
			if err != nil {
				return response.ErrorBadRequest(nil, err.Error())
			}
			if customer == nil {
				return response.ErrorBadRequest(nil, "Customer not found")
			}
		}

		// Create address
		address, err := h.AddressService.Create(payload)
		if err != nil {
			if strings.Contains(err.Error(), "addresses_name_unique") {
				return response.ErrorBadRequest(nil, "Address name already exists")
			}
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(address)
	})
}

/*
 * @apiTag: address
 * @apiPath: /addresses/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: AddressesUpdateRequestParams
 * @apiRequestRef: AddressesCreateRequestBody
 * @apiResponseRef: AddressesCreateResponse
 * @apiSummary: Update address
 * @apiDescription: Update address
* @apiSecurity: apiKeySecurity
*/
func (h *AddressHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.AddressesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.AddressesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if city exists
		city, err := h.CityService.GetCityByID(int64(payload.City.ID))
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}
		if city == nil {
			return response.ErrorBadRequest(nil, "City not found")
		}

		data, err := h.AddressService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: address
 * @apiPath: /addresses
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: AddressesDeleteRequestBody
 * @apiResponseRef: AddressesCreateResponse
 * @apiSummary: Delete address
 * @apiDescription: Delete address
 * @apiSecurity: apiKeySecurity
 */
func (h *AddressHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.AddressesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.AddressService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
