package handlers

import (
	"errors"
	"log"
	"net/http"

	sharedconstants "github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/ticket/constants"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/ticket/models"
	"github.com/hoitek/Maja-Service/internal/ticket/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/ticket/config"
)

type TicketHandler struct {
	TicketService ports.TicketService
	UserService   uPorts.UserService
}

func NewTicketHandler(r *mux.Router, s ports.TicketService, u uPorts.UserService) (TicketHandler, error) {
	ticketHandler := TicketHandler{
		TicketService: s,
		UserService:   u,
	}
	if r == nil {
		return TicketHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.TicketConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.TicketConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/tickets", ticketHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/tickets/{id}/message", ticketHandler.CreateMessage()).Methods(http.MethodPost)
	rAuth.Handle("/tickets", ticketHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/tickets/messages", ticketHandler.QueryMessages()).Methods(http.MethodGet)
	rAuth.Handle("/tickets", ticketHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/tickets/csv/download", ticketHandler.Download()).Methods(http.MethodGet)

	return ticketHandler, nil
}

/*
* @apiTag: ticket
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: tickets
* @apiResponseRef: TicketsQueryResponse
* @apiSummary: Query Tickets
* @apiParametersRef: TicketsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TicketsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TicketsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *TicketHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TicketsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query tickets
		tickets, err := h.TicketService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(tickets)
	})
}

/*
* @apiTag: ticket
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /tickets/messages
* @apiResponseRef: TicketsQueryMessagesResponse
* @apiSummary: Query Tickets
* @apiParametersRef: TicketsQueryMessagesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TicketsQueryMessagesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TicketsQueryMessagesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *TicketHandler) QueryMessages() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TicketsQueryMessagesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query tickets
		tickets, err := h.TicketService.QueryMessages(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(tickets)
	})
}

/*
* @apiTag: ticket
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: tickets/csv/download
* @apiResponseRef: TicketsQueryResponse
* @apiSummary: Query Tickets
* @apiParametersRef: TicketsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TicketsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TicketsQueryNotFoundResponse
 */
func (h *TicketHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TicketsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		tickets, err := h.TicketService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(tickets)
	})
}

/*
 * @apiTag: ticket
 * @apiPath: /tickets
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: TicketsCreateRequestBody
 * @apiResponseRef: TicketsCreateResponse
 * @apiSummary: Create ticket
 * @apiDescription: Create ticket
 * @apiSecurity: apiKeySecurity
 */
func (h *TicketHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TicketsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		var (
			senderType   string
			receiverType string
		)

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if len(authenticatedUser.Roles) == 0 {
			return response.ErrorBadRequest(nil, "Authenticated user has no role")
		}

		// Set sender type
		for _, role := range authenticatedUser.Roles {
			if role.Name == "Dispatcher" {
				senderType = constants.TICKET_SENDER_TYPE_DISPATCHER
				break
			}
		}
		if senderType == "" {
			if authenticatedUser.UserType == sharedconstants.USER_TYPE_STAFF {
				senderType = constants.TICKET_SENDER_TYPE_STAFF
			}
			if authenticatedUser.UserType == sharedconstants.USER_TYPE_CUSTOMER {
				senderType = constants.TICKET_SENDER_TYPE_CUSTOMER
			}
		}
		if senderType == "" {
			senderType = constants.TICKET_SENDER_TYPE_SYSTEM
		}

		// Find user by id if provided
		if payload.UserID != nil {
			user, err := h.UserService.FindByID(int(*payload.UserID))
			if err != nil {
				log.Printf("Error in finding user by id in ticket handler: %v", err)
				return response.ErrorInternalServerError(nil, "User not found")
			}
			if user == nil {
				return response.ErrorBadRequest(nil, "User not found")
			}
			if len(user.Roles) == 0 {
				return response.ErrorBadRequest(nil, "User has no role")
			}

			for _, role := range user.Roles {
				if role.Name == "Dispatcher" {
					receiverType = constants.TICKET_RECEIVER_TYPE_DISPATCHER
					break
				}
			}
			if receiverType == "" {
				if user.UserType == sharedconstants.USER_TYPE_STAFF {
					receiverType = constants.TICKET_RECEIVER_TYPE_STAFF
				}
				if user.UserType == sharedconstants.USER_TYPE_CUSTOMER {
					receiverType = constants.TICKET_RECEIVER_TYPE_CUSTOMER
				}
			}
		}
		if payload.DepartmentID != nil {
			receiverType = constants.TICKET_RECEIVER_TYPE_SYSTEM
		}
		if senderType == constants.TICKET_SENDER_TYPE_STAFF || senderType == constants.TICKET_SENDER_TYPE_CUSTOMER {
			receiverType = constants.TICKET_RECEIVER_TYPE_DISPATCHER
		}

		// Create ticket
		ticket, err := h.TicketService.Create(payload, authenticatedUser.ID, senderType, receiverType)
		if err != nil {
			log.Printf("Error in creating ticket in ticket handler: %v", err)
			return response.ErrorInternalServerError(nil, "Error in creating ticket, please try again later")
		}

		return response.Success(ticket)
	})
}

/*
 * @apiTag: ticket
 * @apiPath: /tickets/{id}/message
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: TicketsCreateRequestBody
 * @apiResponseRef: TicketsCreateMessageResponse
 * @apiSummary: Create ticket message
 * @apiDescription: Create ticket message
 * @apiSecurity: apiKeySecurity
 */
func (h *TicketHandler) CreateMessage() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate params
		p := mux.Vars(r)
		params := &models.TicketsCreateMessageRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate body
		payload := &models.TicketsCreateMessageRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find ticket by id
		ticket, err := h.TicketService.FindByID(int64(params.ID))
		if err != nil {
			log.Printf("Error in finding ticket by id in ticket handler: %v", err)
			return response.ErrorInternalServerError(nil, "Error in finding ticket, please try again later")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if len(authenticatedUser.Roles) == 0 {
			return response.ErrorBadRequest(nil, "Authenticated user has no role")
		}

		var (
			senderType   string
			receiverType string
			senderID     = authenticatedUser.ID
			recipientId  *int64
		)

		// Set sender type
		for _, role := range authenticatedUser.Roles {
			if role.Name == "Dispatcher" {
				senderType = constants.TICKET_SENDER_TYPE_DISPATCHER
				break
			}
		}
		if senderType == "" {
			if authenticatedUser.UserType == sharedconstants.USER_TYPE_STAFF {
				senderType = constants.TICKET_SENDER_TYPE_STAFF
			}
			if authenticatedUser.UserType == sharedconstants.USER_TYPE_CUSTOMER {
				senderType = constants.TICKET_SENDER_TYPE_CUSTOMER
			}
		}
		if senderType == "" {
			senderType = constants.TICKET_SENDER_TYPE_SYSTEM
		}

		// Set receiver type
		if ticket.SenderType == senderType {
			rid := int64(ticket.UserID)
			receiverType = ticket.RecipientType
			recipientId = &rid
		} else if ticket.RecipientType == senderType {
			rid := int64(ticket.CreatedByID)
			receiverType = ticket.SenderType
			recipientId = &rid
		} else {
			return response.ErrorBadRequest(nil, "You are not allowed to send message to this ticket")
		}

		// Create ticket message
		ticketMessage, err := h.TicketService.CreateMessage(int64(ticket.ID), payload, int64(senderID), recipientId, senderType, receiverType)
		if err != nil {
			log.Printf("Error in creating ticket in ticket handler: %v", err)
			return response.ErrorInternalServerError(nil, "Error in creating ticket, please try again later")
		}

		return response.Success(ticketMessage)
	})
}

/*
 * @apiTag: ticket
 * @apiPath: /tickets
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: TicketsDeleteRequestBody
 * @apiResponseRef: TicketsDeleteResponse
 * @apiSummary: Delete ticket
 * @apiDescription: Delete ticket
 * @apiSecurity: apiKeySecurity
 */
func (h *TicketHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TicketsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.TicketService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
