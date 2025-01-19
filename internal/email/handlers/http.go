package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/email/config"
	"github.com/hoitek/Maja-Service/internal/email/models"
	"github.com/hoitek/Maja-Service/internal/email/ports"
)

type EmailHandler struct {
	EmailService ports.EmailService
	UserService  uPorts.UserService
}

func NewEmailHandler(r *mux.Router, s ports.EmailService, us uPorts.UserService) (EmailHandler, error) {
	emailHandler := EmailHandler{
		EmailService: s,
		UserService:  us,
	}
	if r == nil {
		return EmailHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.EmailConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.EmailConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/emails", emailHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/emails", emailHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/emails", emailHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/emails/category/{id}", emailHandler.UpdateCategory()).Methods(http.MethodPut)
	rAuth.Handle("/emails/star/{id}", emailHandler.UpdateStar()).Methods(http.MethodPut)

	return emailHandler, nil
}

/*
* @apiTag: email
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /emails
* @apiResponseRef: EmailsQueryResponse
* @apiSummary: Query emails
* @apiParametersRef: EmailsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: EmailsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: EmailsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *EmailHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.EmailsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query emails
		emails, err := h.EmailService.Query(queries)
		if err != nil {
			log.Printf("Error querying emails: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while querying emails, please try again later")
		}

		return response.Success(emails)
	})
}

/*
 * @apiTag: email
 * @apiPath: /emails
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: EmailsCreateRequestBody
 * @apiResponseRef: EmailsCreateResponse
 * @apiSummary: Create email
 * @apiDescription: Create email
 * @apiSecurity: apiKeySecurity
 */
func (h *EmailHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EmailsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Create email
		email, err := h.EmailService.Create(payload)
		if err != nil {
			log.Printf("Error creating email: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while creating email, please try again later")
		}

		return response.Success(email)
	})
}

/*
 * @apiTag: email
 * @apiPath: /emails
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: EmailsDeleteRequestBody
 * @apiResponseRef: EmailsDeleteResponse
 * @apiSummary: Delete email
 * @apiDescription: Delete email
 * @apiSecurity: apiKeySecurity
 */
func (h *EmailHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EmailsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		// Delete email
		data, err := h.EmailService.Delete(payload)
		if err != nil {
			log.Printf("Error deleting email: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while deleting email, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: email
 * @apiPath: /emails/category/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: EmailsUpdateCategoryRequestParams
 * @apiRequestRef: EmailsUpdateCategoryRequestBody
 * @apiResponseRef: EmailsCreateResponse
 * @apiSummary: Update email category
 * @apiDescription: Update email category
 * @apiSecurity: apiKeySecurity
 */
func (h *EmailHandler) UpdateCategory() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.EmailsUpdateCategoryRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find email by id
		email, _ := h.EmailService.FindByID(int64(params.ID))
		if email == nil {
			return response.ErrorNotFound(nil, "Email not found")
		}

		// Validate request body
		payload := &models.EmailsUpdateCategoryRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Update email category
		data, err := h.EmailService.UpdateCategory(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error updating email category: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while updating email category, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: email
 * @apiPath: /emails/star/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: EmailsUpdateStarRequestParams
 * @apiRequestRef: EmailsUpdateStarRequestBody
 * @apiResponseRef: EmailsCreateResponse
 * @apiSummary: Update email star
 * @apiDescription: Update email star
 * @apiSecurity: apiKeySecurity
 */
func (h *EmailHandler) UpdateStar() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.EmailsUpdateStarRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find email by id
		email, _ := h.EmailService.FindByID(int64(params.ID))
		if email == nil {
			return response.ErrorNotFound(nil, "Email not found")
		}

		// Validate request body
		payload := &models.EmailsUpdateStarRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Update email star
		data, err := h.EmailService.UpdateStar(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error updating email star: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong while updating email star, please try again later")
		}

		return response.Success(data)
	})
}
