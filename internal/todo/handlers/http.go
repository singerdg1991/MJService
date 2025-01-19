package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/todo/domain"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/todo/models"
	"github.com/hoitek/Maja-Service/internal/todo/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/todo/config"
)

type TodoHandler struct {
	TodoService ports.TodoService
	UserService uPorts.UserService
}

func NewTodoHandler(r *mux.Router, s ports.TodoService, u uPorts.UserService) (TodoHandler, error) {
	todoHandler := TodoHandler{
		TodoService: s,
		UserService: u,
	}
	if r == nil {
		return TodoHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.TodoConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.TodoConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/todos", todoHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/todos", todoHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/todos", todoHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/todos/status/{id}/{userid}", todoHandler.UpdateStatus()).Methods(http.MethodPut)
	rAuth.Handle("/todos/{id}", todoHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/todos/csv/download", todoHandler.Download()).Methods(http.MethodGet)

	return todoHandler, nil
}

/*
* @apiTag: todo
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: todos
* @apiResponseRef: TodosQueryResponse
* @apiSummary: Query todos
* @apiParametersRef: TodosQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TodosQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TodosQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *TodoHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TodosQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query todos
		todos, err := h.TodoService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(todos)
	})
}

/*
* @apiTag: todo
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: todos/csv/download
* @apiResponseRef: TodosQueryResponse
* @apiSummary: Query todos
* @apiParametersRef: TodosQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: TodosQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: TodosQueryNotFoundResponse
 */
func (h *TodoHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.TodosQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		todos, err := h.TodoService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(todos)
	})
}

/*
 * @apiTag: todo
 * @apiPath: /todos
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: TodosCreateRequestBody
 * @apiResponseRef: TodosCreateResponse
 * @apiSummary: Create todo
 * @apiDescription: Create todo
 * @apiSecurity: apiKeySecurity
 */
func (h *TodoHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TodosCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = &domain.TodoUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Find user by UserID
		user, err := h.UserService.FindByID(int(payload.UserID))
		if user == nil {
			log.Printf("Create todo: User not found, err: %v", err.Error())
			return response.ErrorBadRequest(nil, "User not found")
		}
		payload.User = &domain.TodoUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			AvatarUrl: user.AvatarUrl,
		}

		// Create todos
		todo, err := h.TodoService.Create(payload)
		if err != nil {
			log.Printf("Create todo: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Create todo failed, please try again later")
		}

		return response.Created(todo)
	})
}

/*
 * @apiTag: todo
 * @apiPath: /todos/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: TodosUpdateRequestParams
 * @apiRequestRef: TodosCreateRequestBody
 * @apiResponseRef: TodosCreateResponse
 * @apiSummary: Update todo
 * @apiDescription: Update todo
 * @apiSecurity: apiKeySecurity
 */
func (h *TodoHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.TodosUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.TodosCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = &domain.TodoUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Check if todos already exists
		res, _ := h.TodoService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "todo does not exist")
		}

		// Update todos
		data, err := h.TodoService.Update(payload, int64(params.ID))
		if err != nil {
			log.Printf("Update todo: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Update todo failed, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: todo
 * @apiPath: /todos/status/{id}/{userid}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: TodosUpdateStatusRequestParams
 * @apiRequestRef: TodosUpdateStatusRequestBody
 * @apiResponseRef: TodosCreateResponse
 * @apiSummary: Update todo
 * @apiDescription: Update todo
 * @apiSecurity: apiKeySecurity
 */
func (h *TodoHandler) UpdateStatus() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.TodosUpdateStatusRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.TodosUpdateStatusRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if todos already exists
		res, _ := h.TodoService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "todo does not exist")
		}
		if res.UserID != uint(params.UserID) {
			return response.ErrorBadRequest(nil, "You are not authorized to update this todo, because you are not the owner of this todo")
		}

		// Update todos
		data, err := h.TodoService.UpdateStatus(payload, int64(params.ID))
		if err != nil {
			log.Printf("Update todo status: %v", err.Error())
			return response.ErrorInternalServerError(nil, "Update todo status failed, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: todo
 * @apiPath: /todos
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: TodosDeleteRequestBody
 * @apiResponseRef: TodosDeleteResponse
 * @apiSummary: Delete todo
 * @apiDescription: Delete todo
 * @apiSecurity: apiKeySecurity
 */
func (h *TodoHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.TodosDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.TodoService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
