package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/languageskill/models"
	"github.com/hoitek/Maja-Service/internal/languageskill/ports"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/languageskill/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type LanguageSkillHandler struct {
	UserService          uPorts.UserService
	LanguageSkillService ports.LanguageSkillService
}

func NewLanguageSkillHandler(r *mux.Router, s ports.LanguageSkillService, u uPorts.UserService) (LanguageSkillHandler, error) {
	languageskillHandler := LanguageSkillHandler{
		UserService:          u,
		LanguageSkillService: s,
	}
	if r == nil {
		return LanguageSkillHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.LanguageSkillConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.LanguageSkillConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/languageskills",
			Method:      http.MethodPost,
			Handler:     languageskillHandler.Create(),
			Permissions: []string{permissions.LANGUAGE_SKILLS_CREATE_NEW_LANGUAGE_SKILL},
		},
		{
			Path:        "/languageskills",
			Method:      http.MethodGet,
			Handler:     languageskillHandler.Query(),
			Permissions: []string{permissions.LANGUAGE_SKILLS_VIEW_ALL_LANGUAGE_SKILLS},
		},
		{
			Path:        "/languageskills",
			Method:      http.MethodDelete,
			Handler:     languageskillHandler.Delete(),
			Permissions: []string{permissions.LANGUAGE_SKILLS_CREATE_NEW_LANGUAGE_SKILL},
		},
		{
			Path:        "/languageskills/{id}",
			Method:      http.MethodPut,
			Handler:     languageskillHandler.Update(),
			Permissions: []string{permissions.LANGUAGE_SKILLS_CREATE_NEW_LANGUAGE_SKILL},
		},
		{
			Path:        "/languageskills/csv/download",
			Method:      http.MethodGet,
			Handler:     languageskillHandler.Download(),
			Permissions: []string{permissions.LANGUAGE_SKILLS_VIEW_ALL_LANGUAGE_SKILLS},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return languageskillHandler, nil
}

/*
* @apiTag: languageskill
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: languageskills
* @apiResponseRef: LanguageSkillsQueryResponse
* @apiSummary: Query languageskills
* @apiParametersRef: LanguageSkillsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LanguageSkillsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LanguageSkillsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *LanguageSkillHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LanguageSkillsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query languageSkills
		languageskills, err := h.LanguageSkillService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(languageskills)
	})
}

/*
* @apiTag: languageskill
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: languageskills/csv/download
* @apiResponseRef: LanguageSkillsQueryResponse
* @apiSummary: Query languageskills
* @apiParametersRef: LanguageSkillsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LanguageSkillsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LanguageSkillsQueryNotFoundResponse
 */
func (h *LanguageSkillHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LanguageSkillsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		languageSkills, err := h.LanguageSkillService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(languageSkills)
	})
}

/*
 * @apiTag: languageskill
 * @apiPath: /languageskills
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: LanguageSkillsCreateRequestBody
 * @apiResponseRef: LanguageSkillsCreateResponse
 * @apiSummary: Create languageskill
 * @apiDescription: Create languageskill
 * @apiSecurity: apiKeySecurity
 */
func (h *LanguageSkillHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LanguageSkillsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find languageSkill by name
		res, _ := h.LanguageSkillService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "LanguageSkill already exists")
		}

		// Create languageSkill
		languageSkill, err := h.LanguageSkillService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(languageSkill)
	})
}

/*
 * @apiTag: languageskill
 * @apiPath: /languageskills/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: LanguageSkillsUpdateRequestParams
 * @apiRequestRef: LanguageSkillsCreateRequestBody
 * @apiResponseRef: LanguageSkillsCreateResponse
 * @apiSummary: Update languageskill
 * @apiDescription: Update languageskill
 * @apiSecurity: apiKeySecurity
 */
func (h *LanguageSkillHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.LanguageSkillsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.LanguageSkillsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.LanguageSkillService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.LanguageSkillService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.LanguageSkillService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: languageskill
 * @apiPath: /languageskills
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: LanguageSkillsDeleteRequestBody
 * @apiResponseRef: LanguageSkillsDeleteResponse
 * @apiSummary: Delete languageskill
 * @apiDescription: Delete languageskill
 * @apiSecurity: apiKeySecurity
 */
func (h *LanguageSkillHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LanguageSkillsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.LanguageSkillService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
