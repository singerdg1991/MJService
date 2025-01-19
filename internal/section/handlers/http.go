package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/section/models"
	"github.com/hoitek/Maja-Service/internal/section/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/section/config"
)

type SectionHandler struct {
	SectionService service.SectionService
}

func NewSectionHandler(r *mux.Router, s service.SectionService) (SectionHandler, error) {
	sectionHandler := SectionHandler{
		SectionService: s,
	}
	if r == nil {
		return SectionHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.SectionConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.SectionConfig.ApiVersion1).Subrouter()

	rv1.Handle("/sections", sectionHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/sections", sectionHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/sections", sectionHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/sections/{id}", sectionHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/sections/csv/download", sectionHandler.Download()).Methods(http.MethodGet)
	return sectionHandler, nil
}

/*
* @apiTag: section
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: sections
* @apiResponseRef: SectionsQueryResponse
* @apiSummary: Query sections
* @apiParametersRef: SectionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: SectionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: SectionsQueryNotFoundResponse
 */
func (h *SectionHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.SectionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		sections, err := h.SectionService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(sections)
	})
}

/*
* @apiTag: section
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: sections/csv/download
* @apiResponseRef: SectionsQueryResponse
* @apiSummary: Query sections
* @apiParametersRef: SectionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: SectionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: SectionsQueryNotFoundResponse
 */
func (h *SectionHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.SectionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		sections, err := h.SectionService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(sections)
	})
}

/*
 * @apiTag: section
 * @apiPath: /sections
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: SectionsCreateRequestBody
 * @apiResponseRef: SectionsCreateResponse
 * @apiSummary: Create section
 * @apiDescription: Create section
 */
func (h *SectionHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.SectionsCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		section, err := h.SectionService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(section)
	})
}

/*
 * @apiTag: section
 * @apiPath: /sections/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: SectionsUpdateRequestParams
 * @apiRequestRef: SectionsCreateRequestBody
 * @apiResponseRef: SectionsCreateResponse
 * @apiSummary: Update section
 * @apiDescription: Update section
 */
func (h *SectionHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.SectionsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.SectionsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.SectionService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: section
 * @apiPath: /sections
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: SectionsDeleteRequestBody
 * @apiResponseRef: SectionsCreateResponse
 * @apiSummary: Delete section
 * @apiDescription: Delete section
 */
func (h *SectionHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.SectionsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.SectionService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
