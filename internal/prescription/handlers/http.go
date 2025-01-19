package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/prescription/constants"
	"github.com/hoitek/Maja-Service/internal/prescription/models"
	stPorts "github.com/hoitek/Maja-Service/internal/prescription/ports"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/prescription/config"
)

type PrescriptionHandler struct {
	PrescriptionService stPorts.PrescriptionService
	S3Service           s3Ports.S3Service
}

func NewPrescriptionHandler(
	r *mux.Router,
	s stPorts.PrescriptionService,
	s3Service s3Ports.S3Service,
) (PrescriptionHandler, error) {
	prescriptionHandler := PrescriptionHandler{
		PrescriptionService: s,
		S3Service:           s3Service,
	}
	if r == nil {
		return PrescriptionHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.PrescriptionConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.PrescriptionConfig.ApiVersion1).Subrouter()

	rv1.Handle("/prescriptions", prescriptionHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/prescriptions", prescriptionHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/prescriptions", prescriptionHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/prescriptions/{id}", prescriptionHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/prescriptions/csv/download", prescriptionHandler.Download()).Methods(http.MethodGet)
	return prescriptionHandler, nil
}

/*
* @apiTag: prescription
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: prescriptions
* @apiResponseRef: PrescriptionsQueryResponse
* @apiSummary: Query prescriptions
* @apiParametersRef: PrescriptionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PrescriptionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PrescriptionsQueryNotFoundResponse
 */
func (h *PrescriptionHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PrescriptionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		prescriptions, err := h.PrescriptionService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(prescriptions)
	})
}

/*
* @apiTag: prescription
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: prescriptions/csv/download
* @apiResponseRef: PrescriptionsQueryResponse
* @apiSummary: Query prescriptions
* @apiParametersRef: PrescriptionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PrescriptionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PrescriptionsQueryNotFoundResponse
 */
func (h *PrescriptionHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PrescriptionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		prescriptions, err := h.PrescriptionService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(prescriptions)
	})
}

/*
 * @apiTag: prescription
 * @apiPath: /prescriptions
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: PrescriptionsCreateRequestBody
 * @apiResponseRef: PrescriptionsCreateResponse
 * @apiSummary: Create prescription
 * @apiDescription: Create prescription
 */
func (h *PrescriptionHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PrescriptionsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Create prescription
		createdPrescription, err := h.PrescriptionService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.PRESCRIPTION_BUCKET_NAME, payload.Attachments, int64(createdPrescription.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedPrescription, err := h.PrescriptionService.UpdatePrescriptionAttachments(nil, uploadedFilesMetadata, int64(createdPrescription.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdPrescription = updatedPrescription
		}

		return response.Created(createdPrescription)
	})
}

/*
 * @apiTag: prescription
 * @apiPath: /prescriptions/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: PrescriptionsUpdateRequestParams
 * @apiRequestRef: PrescriptionsUpdateRequestBody
 * @apiResponseRef: PrescriptionsCreateResponse
 * @apiSummary: Update prescription
 * @apiDescription: Update prescription
 */
func (h *PrescriptionHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.PrescriptionsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.PrescriptionsUpdateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		prescription, err := h.PrescriptionService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		var uploadedFilesMetadata []*types.UploadMetadata
		if len(payload.Attachments) > 0 {
			uploadedFilesMetadata, _ = h.S3Service.UploadFiles(constants.PRESCRIPTION_BUCKET_NAME, payload.Attachments, int64(prescription.ID))
		}
		updatedPrescription, err := h.PrescriptionService.UpdatePrescriptionAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, int64(prescription.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		prescription = updatedPrescription

		return response.Success(prescription)
	})
}

/*
 * @apiTag: prescription
 * @apiPath: /prescriptions
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: PrescriptionsDeleteRequestBody
 * @apiResponseRef: PrescriptionsCreateResponse
 * @apiSummary: Delete prescription
 * @apiDescription: Delete prescription
 */
func (h *PrescriptionHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PrescriptionsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.PrescriptionService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
