package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/archive/constants"
	"github.com/hoitek/Maja-Service/internal/archive/models"
	"github.com/hoitek/Maja-Service/internal/archive/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/archive/config"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
)

type ArchiveHandler struct {
	ArchiveService ports.ArchiveService
	UserService    uPorts.UserService
	S3Service      s3Ports.S3Service
}

func NewArchiveHandler(r *mux.Router, s ports.ArchiveService, u uPorts.UserService, s3 s3Ports.S3Service) (ArchiveHandler, error) {
	archiveHandler := ArchiveHandler{
		ArchiveService: s,
		UserService:    u,
		S3Service:      s3,
	}
	if r == nil {
		return ArchiveHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ArchiveConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ArchiveConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/archives",
			Method:      http.MethodPost,
			Handler:     archiveHandler.Create(),
			Permissions: []string{permissions.ARCHIVES_CREATE_NEW_ARCHIVE},
		},
		{
			Path:        "/archives",
			Method:      http.MethodGet,
			Handler:     archiveHandler.Query(),
			Permissions: []string{permissions.ARCHIVES_VIEW_ALL_ARCHIVES},
		},
		{
			Path:        "/archives",
			Method:      http.MethodDelete,
			Handler:     archiveHandler.Delete(),
			Permissions: []string{permissions.ARCHIVES_CREATE_NEW_ARCHIVE},
		},
		{
			Path:        "/archives/{id}",
			Method:      http.MethodPut,
			Handler:     archiveHandler.Update(),
			Permissions: []string{permissions.ARCHIVES_CREATE_NEW_ARCHIVE},
		},
		{
			Path:        "/archives/csv/download",
			Method:      http.MethodGet,
			Handler:     archiveHandler.Download(),
			Permissions: []string{permissions.ARCHIVES_VIEW_ALL_ARCHIVES},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return archiveHandler, nil
}

/*
* @apiTag: archive
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: archives
* @apiResponseRef: ArchivesQueryResponse
* @apiSummary: Query archives
* @apiParametersRef: ArchivesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ArchivesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ArchivesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ArchiveHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ArchivesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query archives
		archives, err := h.ArchiveService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(archives)
	})
}

/*
* @apiTag: archive
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: archives/csv/download
* @apiResponseRef: ArchivesQueryResponse
* @apiSummary: Query archives
* @apiParametersRef: ArchivesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ArchivesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ArchivesQueryNotFoundResponse
 */
func (h *ArchiveHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ArchivesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		archives, err := h.ArchiveService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(archives)
	})
}

/*
 * @apiTag: archive
 * @apiPath: /archives
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: ArchivesCreateRequestBody
 * @apiResponseRef: ArchivesCreateResponse
 * @apiSummary: Create archive
 * @apiDescription: Create archive
 */
func (h *ArchiveHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ArchivesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if user exists or not
		user, err := h.UserService.FindByID(int(payload.UserID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if user == nil {
			return response.ErrorNotFound(nil, "User not found")
		}

		// Check date format is valid 2021-01-01 with regexp `\d{4}-\d{2}-\d{2}`
		match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, payload.Date)
		if !match {
			return response.ErrorBadRequest(nil, "Date format should be YYYY-MM-DD")
		}

		// Check time format is valid 00:00:00 with regexp `\d{2}:\d{2}:\d{2}`
		match, _ = regexp.MatchString(`\d{2}:\d{2}:\d{2}`, payload.Time)
		if !match {
			return response.ErrorBadRequest(nil, "Time format should be HH:MM:SS")
		}

		// Prepare payload DateTime from date and time with the format YYYY-MM-DDTHH:MM:SS to time.Time
		payloadDateTime, err := time.Parse("2006-01-02T15:04:05", payload.Date+"T"+payload.Time)
		if err != nil {
			return response.ErrorBadRequest(nil, "Date and time format should be YYYY-MM-DDTHH:MM:SS")
		}
		payload.DateTime = payloadDateTime

		// Create archive
		archive, err := h.ArchiveService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.ARCHIVE_BUCKET_NAME, payload.Attachments, int64(archive.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedArchive, err := h.ArchiveService.UpdateAttachments(uploadedFilesMetadata, int64(archive.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			archive = updatedArchive
		}

		return response.Success(archive)
	})
}

/*
 * @apiTag: archive
 * @apiPath: /archives/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: ArchivesUpdateRequestParams
 * @apiRequestRef: ArchivesCreateRequestBody
 * @apiResponseRef: ArchivesCreateResponse
 * @apiSummary: Update archive
 * @apiDescription: Update archive
 * @apiSecurity: apiKeySecurity
 */
func (h *ArchiveHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ArchivesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ArchivesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if user exists or not
		user, err := h.UserService.FindByID(int(payload.UserID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if user == nil {
			return response.ErrorNotFound(nil, "User not found")
		}

		// Check if Archive already exists
		res, _ := h.ArchiveService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Archive does not exist")
		}

		// Update Archive
		data, err := h.ArchiveService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: archive
 * @apiPath: /archives
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ArchivesDeleteRequestBody
 * @apiResponseRef: ArchivesDeleteResponse
 * @apiSummary: Delete archive
 * @apiDescription: Delete archive
 * @apiSecurity: apiKeySecurity
 */
func (h *ArchiveHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ArchivesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ArchiveService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
