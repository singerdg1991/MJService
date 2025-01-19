package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/reward/models"
	"github.com/hoitek/Maja-Service/internal/reward/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/reward/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type RewardHandler struct {
	UserService   uPorts.UserService
	RewardService ports.RewardService
}

func NewRewardHandler(r *mux.Router, s ports.RewardService, u uPorts.UserService) (RewardHandler, error) {
	rewardHandler := RewardHandler{
		UserService:   u,
		RewardService: s,
	}
	if r == nil {
		return RewardHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.RewardConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.RewardConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/rewards", rewardHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/rewards", rewardHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/rewards", rewardHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/rewards/{id}", rewardHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/rewards/csv/download", rewardHandler.Download()).Methods(http.MethodGet)

	return rewardHandler, nil
}

/*
* @apiTag: reward
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: rewards
* @apiResponseRef: RewardsQueryResponse
* @apiSummary: Query rewards
* @apiParametersRef: RewardsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: RewardsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: RewardsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *RewardHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.RewardsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query rewards
		rewards, err := h.RewardService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(rewards)
	})
}

/*
* @apiTag: reward
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: rewards/csv/download
* @apiResponseRef: RewardsQueryResponse
* @apiSummary: Query rewards
* @apiParametersRef: RewardsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: RewardsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: RewardsQueryNotFoundResponse
 */
func (h *RewardHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.RewardsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		rewards, err := h.RewardService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(rewards)
	})
}

/*
 * @apiTag: reward
 * @apiPath: /rewards
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: RewardsCreateRequestBody
 * @apiResponseRef: RewardsCreateResponse
 * @apiSummary: Create reward
 * @apiDescription: Create reward
 * @apiSecurity: apiKeySecurity
 */
func (h *RewardHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.RewardsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find reward by name
		res, _ := h.RewardService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Reward already exists")
		}

		// Create reward
		reward, err := h.RewardService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(reward)
	})
}

/*
 * @apiTag: reward
 * @apiPath: /rewards/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: RewardsUpdateRequestParams
 * @apiRequestRef: RewardsCreateRequestBody
 * @apiResponseRef: RewardsCreateResponse
 * @apiSummary: Update reward
 * @apiDescription: Update reward
 * @apiSecurity: apiKeySecurity
 */
func (h *RewardHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.RewardsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.RewardsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.RewardService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.RewardService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.RewardService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: reward
 * @apiPath: /rewards
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: RewardsDeleteRequestBody
 * @apiResponseRef: RewardsDeleteResponse
 * @apiSummary: Delete reward
 * @apiDescription: Delete reward
 * @apiSecurity: apiKeySecurity
 */
func (h *RewardHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.RewardsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.RewardService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
