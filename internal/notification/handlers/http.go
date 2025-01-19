package handlers

import (
	"errors"
	"net/http"

	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/notification/models"
	"github.com/hoitek/Maja-Service/internal/notification/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/notification/config"
)

type NotificationHandler struct {
	NotificationService ports.NotificationService
	UserService         uPorts.UserService
}

func NewNotificationHandler(r *mux.Router, s ports.NotificationService, u uPorts.UserService) (NotificationHandler, error) {
	notificationHandler := NotificationHandler{
		NotificationService: s,
		UserService:         u,
	}
	if r == nil {
		return NotificationHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.NotificationConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.NotificationConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/notifications", notificationHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/notifications", notificationHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/notifications/{id}", notificationHandler.UpdateStatus()).Methods(http.MethodPut)

	return notificationHandler, nil
}

/*
* @apiTag: notification
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: notifications
* @apiResponseRef: NotificationsQueryResponse
* @apiSummary: Query notifications
* @apiParametersRef: NotificationsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: NotificationsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: NotificationsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *NotificationHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.NotificationsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorUnAuthorized(nil, "You are not authorized to access this resource")
		}

		// Query notifications
		notifications, err := h.NotificationService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(notifications)
	})
}

/*
 * @apiTag: notification
 * @apiPath: /notifications
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: NotificationsDeleteRequestBody
 * @apiResponseRef: NotificationsDeleteResponse
 * @apiSummary: Delete notification
 * @apiDescription: Delete notification
 * @apiSecurity: apiKeySecurity
 */
func (h *NotificationHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.NotificationsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.NotificationService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: notification
 * @apiPath: /notifications/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: NotificationsUpdateStatusRequestParams
 * @apiRequestRef: NotificationsUpdateStatusRequestBody
 * @apiResponseRef: NotificationsUpdateStatusResponse
 * @apiSummary: Update notification status
 * @apiDescription: Update notification status
 * @apiSecurity: apiKeySecurity
 */
func (h *NotificationHandler) UpdateStatus() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.NotificationsUpdateStatusRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.NotificationsUpdateStatusRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if notification already exists
		notification, _ := h.NotificationService.FindByID(int64(params.ID))
		if notification == nil {
			return response.ErrorBadRequest(nil, "notification does not exist")
		}
		if notification.Status == nil {
			return response.ErrorBadRequest(nil, "The type of the notification is not request")
		}

		// Update notification
		data, err := h.NotificationService.UpdateStatus(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
