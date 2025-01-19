package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/config"
)

type HolidayHandler struct {
	HolidayService ports.HolidayService
	UserService    uPorts.UserService
}

func NewHolidayHandler(r *mux.Router, s ports.HolidayService, u uPorts.UserService) (HolidayHandler, error) {
	holidayHandler := HolidayHandler{
		HolidayService: s,
		UserService:    u,
	}
	if r == nil {
		return HolidayHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.HolidayConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.HolidayConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/staffclub/holidays", holidayHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/staffclub/holidays", holidayHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/staffclub/holidays", holidayHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/staffclub/holidays/status/{id}", holidayHandler.UpdateStatus()).Methods(http.MethodPut)
	rAuth.Handle("/staffclub/holidays/{id}", holidayHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/staffclub/holidays/csv/download", holidayHandler.Download()).Methods(http.MethodGet)

	return holidayHandler, nil
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/holidays
* @apiResponseRef: HolidaysQueryResponse
* @apiSummary: Query holidays
* @apiParametersRef: HolidaysQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: HolidaysQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: HolidaysQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.HolidaysQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query holidays
		holidays, err := h.HolidayService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(holidays)
	})
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/holidays/csv/download
* @apiResponseRef: HolidaysQueryResponse
* @apiSummary: Query holidays
* @apiParametersRef: HolidaysQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: HolidaysQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: HolidaysQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.HolidaysQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		holidays, err := h.HolidayService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(holidays)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/holidays
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: HolidaysCreateRequestBody
 * @apiResponseRef: HolidaysCreateResponse
 * @apiSummary: Create holiday
 * @apiDescription: Create holiday
 * @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.HolidaysCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user by userId
		if payload.UserID != nil {
			user, err := h.UserService.FindByID(*payload.UserID)
			if err != nil {
				return response.ErrorBadRequest(nil, "User not found")
			}
			payload.User = &sharedmodels.AuthenticatedUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				AvatarUrl: user.AvatarUrl,
			}
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Create holiday
		holiday, err := h.HolidayService.Create(payload)
		if err != nil {
			log.Printf("Error creating holiday: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Created(holiday)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/holidays/status/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: HolidaysUpdateStatusRequestParams
 * @apiRequestRef: HolidaysUpdateStatusRequestBody
 * @apiResponseRef: HolidaysCreateResponse
 * @apiSummary: Update holiday status
 * @apiDescription: Update holiday status
 * @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) UpdateStatus() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.HolidaysUpdateStatusRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.HolidaysUpdateStatusRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Check if Holiday already exists
		holiday, _ := h.HolidayService.FindByID(int64(params.ID))
		if holiday == nil {
			return response.ErrorBadRequest(nil, "Holiday does not exist")
		}

		// Update Holiday
		data, err := h.HolidayService.UpdateStatus(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error updating holiday status: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/holidays/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: HolidaysUpdateRequestParams
 * @apiRequestRef: HolidaysCreateRequestBody
 * @apiResponseRef: HolidaysCreateResponse
 * @apiSummary: Update holiday
 * @apiDescription: Update holiday
 * @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.HolidaysUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.HolidaysCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user by userId
		if payload.UserID != nil {
			user, err := h.UserService.FindByID(*payload.UserID)
			if err != nil {
				return response.ErrorBadRequest(nil, "User not found")
			}
			payload.User = &sharedmodels.AuthenticatedUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				AvatarUrl: user.AvatarUrl,
			}
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Check if Holiday already exists
		holiday, _ := h.HolidayService.FindByID(int64(params.ID))
		if holiday == nil {
			return response.ErrorBadRequest(nil, "Holiday does not exist")
		}

		// Update Holiday
		data, err := h.HolidayService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/holidays
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: HolidaysDeleteRequestBody
 * @apiResponseRef: HolidaysDeleteResponse
 * @apiSummary: Delete holiday
 * @apiDescription: Delete holiday
 * @apiSecurity: apiKeySecurity
 */
func (h *HolidayHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.HolidaysDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.HolidayService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
