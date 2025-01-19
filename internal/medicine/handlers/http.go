package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/medicine/models"
	"github.com/hoitek/Maja-Service/internal/medicine/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/medicine/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type MedicineHandler struct {
	UserService     uPorts.UserService
	MedicineService ports.MedicineService
}

func NewMedicineHandler(r *mux.Router, s ports.MedicineService, u uPorts.UserService) (MedicineHandler, error) {
	medicineHandler := MedicineHandler{
		UserService:     u,
		MedicineService: s,
	}
	if r == nil {
		return MedicineHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.MedicineConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.MedicineConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/medicines", medicineHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/medicines", medicineHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/medicines", medicineHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/medicines/{id}", medicineHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/medicines/csv/download", medicineHandler.Download()).Methods(http.MethodGet)

	return medicineHandler, nil
}

/*
* @apiTag: medicine
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: medicines
* @apiResponseRef: MedicinesQueryResponse
* @apiSummary: Query medicines
* @apiParametersRef: MedicinesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: MedicinesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: MedicinesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *MedicineHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.MedicinesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query medicines
		medicines, err := h.MedicineService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(medicines)
	})
}

/*
* @apiTag: medicine
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: medicines/csv/download
* @apiResponseRef: MedicinesQueryResponse
* @apiSummary: Query medicines
* @apiParametersRef: MedicinesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: MedicinesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: MedicinesQueryNotFoundResponse
 */
func (h *MedicineHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.MedicinesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		medicines, err := h.MedicineService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(medicines)
	})
}

/*
 * @apiTag: medicine
 * @apiPath: /medicines
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: MedicinesCreateRequestBody
 * @apiResponseRef: MedicinesCreateResponse
 * @apiSummary: Create medicine
 * @apiDescription: Create medicine
 * @apiSecurity: apiKeySecurity
 */
func (h *MedicineHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.MedicinesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find medicine by name
		res, _ := h.MedicineService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorBadRequest(nil, "Medicine already exists")
		}

		// Find medicine by code
		res, _ = h.MedicineService.FindByCode(payload.Code)
		if res != nil {
			return response.ErrorBadRequest(nil, "Medicine already exists")
		}

		// Create medicine
		medicine, err := h.MedicineService.Create(payload)
		if err != nil {
			log.Printf("Failed to create medicine: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Failed to create medicine, please try again later")
		}

		return response.Created(medicine)
	})
}

/*
 * @apiTag: medicine
 * @apiPath: /medicines/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: MedicinesUpdateRequestParams
 * @apiRequestRef: MedicinesCreateRequestBody
 * @apiResponseRef: MedicinesCreateResponse
 * @apiSummary: Update medicine
 * @apiDescription: Update medicine
 * @apiSecurity: apiKeySecurity
 */
func (h *MedicineHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.MedicinesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.MedicinesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if medicine already exists
		res, _ := h.MedicineService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "medicine does not exist")
		}

		// Check if medicine already exists with the same name
		res, _ = h.MedicineService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "medicine already exists")
		}

		// Update medicine
		data, err := h.MedicineService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: medicine
 * @apiPath: /medicines
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: MedicinesDeleteRequestBody
 * @apiResponseRef: MedicinesDeleteResponse
 * @apiSummary: Delete medicine
 * @apiDescription: Delete medicine
 * @apiSecurity: apiKeySecurity
 */
func (h *MedicineHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.MedicinesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.MedicineService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
