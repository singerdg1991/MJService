package handlers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/push/models"
	"github.com/hoitek/Maja-Service/internal/push/ports"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/push/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type PushHandler struct {
	UserService uPorts.UserService
	PushService ports.PushService
}

func NewPushHandler(r *mux.Router, s ports.PushService, u uPorts.UserService) (PushHandler, error) {
	pushHandler := PushHandler{
		UserService: u,
		PushService: s,
	}
	if r == nil {
		return PushHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.PushConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.PushConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/pushes", pushHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/pushes", pushHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/pushes/send", pushHandler.Send()).Methods(http.MethodPost)

	return pushHandler, nil
}

/*
* @apiTag: push
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: pushes
* @apiResponseRef: PushesQueryResponse
* @apiSummary: Query pushes
* @apiParametersRef: PushesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PushesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PushesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *PushHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PushesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query pushes
		pushes, err := h.PushService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(pushes)
	})
}

/*
 * @apiTag: push
 * @apiPath: /pushes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: PushesCreateRequestBody
 * @apiResponseRef: PushesCreateResponse
 * @apiSummary: Create push
 * @apiDescription: Create push
 * @apiSecurity: apiKeySecurity
 */
func (h *PushHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PushesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find user by id
		user, err := h.UserService.FindByID(payload.UserID)
		if err != nil {
			log.Printf("Error in finding user in creating push: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}
		if user == nil {
			return response.ErrorNotFound(nil, "User not found")
		}

		// Create push
		push, err := h.PushService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(push)
	})
}

/*
 * @apiTag: push
 * @apiPath: /pushes/send
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: PushesSendRequestBody
 * @apiResponseRef: PushesSendResponse
 * @apiSummary: Send push
 * @apiDescription: Send push
 * @apiSecurity: apiKeySecurity
 */
func (h *PushHandler) Send() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PushesSendRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find user by id
		user, err := h.UserService.FindByID(payload.UserID)
		if err != nil {
			log.Printf("Error in finding user in creating push: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}
		if user == nil {
			return response.ErrorNotFound(nil, "User not found")
		}

		// Find push by userId
		push, err := h.PushService.FindByUserID(payload.UserID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if push == nil {
			return response.ErrorNotFound(nil, "Push not found")
		}

		// Get config and prepare data
		var (
			data            = []byte(fmt.Sprintf(`{"title":"%s","body":"%s"}`, payload.Title, payload.Body))
			vapidPublicKey  = config.PushConfig.VAPIDPublicKey
			vapidPrivateKey = config.PushConfig.VAPIDPrivateKey
		)

		// Send push
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		_, err = webpush.SendNotification(data, &webpush.Subscription{
			Endpoint: push.Endpoint,
			Keys: webpush.Keys{
				Auth:   push.KeysAuth,
				P256dh: push.KeysP256dh,
			},
		}, &webpush.Options{
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             30,
			HTTPClient:      httpClient,
		})
		if err != nil {
			log.Printf("Error in sending push: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}
		return response.Created(models.PushesSendResponseData{Status: "success"})
	})
}
