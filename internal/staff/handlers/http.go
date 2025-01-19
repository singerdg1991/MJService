package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/response"
	sharedconstants "github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/security"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	aPorts "github.com/hoitek/Maja-Service/internal/ability/ports"
	cPorts "github.com/hoitek/Maja-Service/internal/contracttype/ports"
	lsPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	permPorts "github.com/hoitek/Maja-Service/internal/license/ports"
	nPorts "github.com/hoitek/Maja-Service/internal/notification/ports"
	pPorts "github.com/hoitek/Maja-Service/internal/paymenttype/ports"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	sPorts "github.com/hoitek/Maja-Service/internal/section/ports"
	stPorts "github.com/hoitek/Maja-Service/internal/shifttype/ports"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"github.com/hoitek/Maja-Service/internal/staff/constants"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"github.com/hoitek/Maja-Service/internal/staff/models"
	"github.com/hoitek/Maja-Service/internal/staff/ports"
	"github.com/hoitek/Maja-Service/internal/staff/service"
	ntPorts "github.com/hoitek/Maja-Service/internal/stafftype/ports"
	umodels "github.com/hoitek/Maja-Service/internal/user/models"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/utils"
)

type StaffHandler struct {
	StaffService         ports.StaffService
	UserService          uPorts.UserService
	RoleService          rPorts.RoleService
	SectionService       sPorts.SectionService
	AbilityService       aPorts.AbilityService
	PaymentTypeService   pPorts.PaymentTypeService
	ContractTypeService  cPorts.ContractTypeService
	ShiftTypeService     stPorts.ShiftTypeService
	LicenseService       permPorts.LicenseService
	StaffTypeService     ntPorts.StaffTypeService
	LanguageSkillService lsPorts.LanguageSkillService
	NotificationService  nPorts.NotificationService
	S3Service            s3Ports.S3Service
}

func NewStaffHandler(r *mux.Router,
	nService service.StaffService,
	uService uPorts.UserService,
	rService rPorts.RoleService,
	sService sPorts.SectionService,
	aService aPorts.AbilityService,
	pService pPorts.PaymentTypeService,
	cService cPorts.ContractTypeService,
	stService stPorts.ShiftTypeService,
	permService permPorts.LicenseService,
	ntService ntPorts.StaffTypeService,
	lsService lsPorts.LanguageSkillService,
	nfService nPorts.NotificationService,
	s3Service s3Ports.S3Service,
) (StaffHandler, error) {
	staffHandler := StaffHandler{
		StaffService:         &nService,
		UserService:          uService,
		RoleService:          rService,
		SectionService:       sService,
		AbilityService:       aService,
		PaymentTypeService:   pService,
		ContractTypeService:  cService,
		ShiftTypeService:     stService,
		LicenseService:       permService,
		StaffTypeService:     ntService,
		LanguageSkillService: lsService,
		NotificationService:  nfService,
		S3Service:            s3Service,
	}
	if r == nil {
		return StaffHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.StaffConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.StaffConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(uService, []string{}))

	rAuth.Handle("/staffs/personalinfo", staffHandler.CreatePersonalInfo()).Methods(http.MethodPost)
	rAuth.Handle("/staffs/personalinfo/{userId}", staffHandler.UpdatePersonalInfo()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/contract/{userId}", staffHandler.CreateOrUpdateContract()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/licenses", staffHandler.CreateLicenses()).Methods(http.MethodPost)
	rAuth.Handle("/staffs/absences", staffHandler.CreateAbsences()).Methods(http.MethodPost)
	rAuth.Handle("/staffs/absences", staffHandler.QueryAbsences()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/absences/{id}", staffHandler.UpdateAbsence()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/licenses", staffHandler.DeleteLicenses()).Methods(http.MethodDelete)
	rAuth.Handle("/staffs/licenses", staffHandler.QueryLicenses()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/licenses/{id}", staffHandler.UpdateLicense()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/absences", staffHandler.DeleteAbsences()).Methods(http.MethodDelete)
	rAuth.Handle("/staffs", staffHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/staffs", staffHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/staffs/csv/download", staffHandler.Download()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/profile", staffHandler.QueryProfile()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/otherattachments", staffHandler.CreateOtherAttachments()).Methods(http.MethodPost)
	rAuth.Handle("/staffs/otherattachments", staffHandler.QueryOtherAttachments()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/otherattachments/{id}", staffHandler.UpdateOtherAttachment()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/otherattachments", staffHandler.DeleteOtherAttachments()).Methods(http.MethodDelete)
	rAuth.Handle("/staffs/libraries", staffHandler.CreateLibraries()).Methods(http.MethodPost)
	rAuth.Handle("/staffs/libraries", staffHandler.QueryLibraries()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/libraries/{id}", staffHandler.UpdateLibrary()).Methods(http.MethodPut)
	rAuth.Handle("/staffs/libraries", staffHandler.DeleteLibraries()).Methods(http.MethodDelete)
	rAuth.Handle("/staffs/chats", staffHandler.QueryChats()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/chats/messages", staffHandler.QueryChatMessages()).Methods(http.MethodGet)
	rAuth.Handle("/staffs/chats/messages", staffHandler.CreateChatMessage()).Methods(http.MethodPost)

	return staffHandler, nil
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: staffs
* @apiResponseRef: StaffsQueryResponse
* @apiSummary: Query staffs
* @apiParametersRef: StaffsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffs, err := h.StaffService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffs)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/licenses
* @apiResponseRef: StaffsQueryLicensesResponse
* @apiSummary: Query staff licenses
* @apiParametersRef: StaffsQueryLicensesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryLicenses() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryLicensesRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffs, err := h.StaffService.QueryLicenses(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffs)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: staffs/csv/download
* @apiResponseRef: StaffsQueryResponse
* @apiSummary: Query staffs
* @apiParametersRef: StaffsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffs, err := h.StaffService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Generate csv file
		csvDownloadFilePath := "/reports/staffs.csv"
		csvFilePath := fmt.Sprintf("public%s", csvDownloadFilePath)
		_, err = h.StaffService.ExportToCsvAndSave(staffs.Items.([]*domain.Staff), csvFilePath)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffs)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/personalinfo
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: StaffsCreatePersonalInfoRequestBody
 * @apiResponseRef: StaffsCreatePersonalInfoResponse
 * @apiSummary: Create staff personal info
 * @apiDescription: Create staff personal info
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreatePersonalInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.StaffsCreatePersonalInfoRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if language skills exists
		langSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(payload.LanguageSkillsInt64)
		if err != nil {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		if len(langSkills) != len(payload.LanguageSkillsInt64) {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		payload.LanguageSkillIds = payload.LanguageSkillsInt64

		var userLanguageSkills []umodels.UsersCreateRequestBodyLanguageSkill
		for _, languageSkill := range langSkills {
			userLanguageSkills = append(userLanguageSkills, umodels.UsersCreateRequestBodyLanguageSkill{
				ID:   languageSkill.ID,
				Name: languageSkill.Name,
			})
		}

		// Generate registration number
		registrationNumber, err := h.StaffService.GenerateRegistrationNumber()
		if err != nil {
			log.Printf("Error in generating registration number: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		// Create user
		payload.Password = security.HashPassword(payload.Password)
		user, err := h.UserService.Create(&umodels.UsersCreateRequestBody{
			UserType:             sharedconstants.USER_TYPE_STAFF,
			LanguageSkillIds:     payload.LanguageSkillsInt64,
			FirstName:            payload.FirstName,
			LastName:             payload.LastName,
			Username:             payload.Username,
			Password:             payload.Password,
			Gender:               payload.Gender,
			WorkPhoneNumber:      payload.WorkPhoneNumber,
			AccountNumber:        payload.AccountNumber,
			LanguageSkills:       userLanguageSkills,
			RegistrationNumber:   registrationNumber,
			ForcedChangePassword: payload.ForcedChangePassword,
			Email:                payload.Email,
			Phone:                payload.Phone,
			NationalCode:         payload.NationalCode,
			BirthDate:            payload.BirthDate,
			AvatarUrl:            payload.AvatarUrl,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Create staff
		staff, err := h.StaffService.CreateStaff(int64(user.ID), payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		user.StaffID = &staff.ID
		user.VehicleTypes = staff.VehicleTypes
		user.VehicleLicenseTypes = staff.VehicleLicenseTypes
		user.Limitations = staff.Limitations

		// Return response
		return response.Success(user)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/personalinfo/{userId}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: StaffsUpdatePersonalInfoRequestParams
 * @apiRequestRef: StaffsUpdatePersonalInfoRequestBody
 * @apiResponseRef: StaffsCreatePersonalInfoResponse
 * @apiSummary: Update staff personal info
 * @apiDescription: Update staff personal info
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) UpdatePersonalInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsUpdatePersonalInfoRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsUpdatePersonalInfoRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if language skills exists
		langSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(payload.LanguageSkillsInt64)
		if err != nil {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		if len(langSkills) != len(payload.LanguageSkillsInt64) {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		payload.LanguageSkillIds = payload.LanguageSkillsInt64

		var userLanguageSkills []umodels.UsersUpdateRequestBodyLanguageSkill
		for _, languageSkill := range langSkills {
			userLanguageSkills = append(userLanguageSkills, umodels.UsersUpdateRequestBodyLanguageSkill{
				ID:   languageSkill.ID,
				Name: languageSkill.Name,
			})
		}

		// Generate registration number
		registrationNumber, err := h.StaffService.GenerateRegistrationNumber()
		if err != nil {
			log.Printf("Error in generating registration number: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		// Update user
		user, err := h.UserService.Update(&umodels.UsersUpdateRequestBody{
			UserType:           sharedconstants.USER_TYPE_STAFF,
			LanguageSkillIds:   payload.LanguageSkillsInt64,
			FirstName:          payload.FirstName,
			LastName:           payload.LastName,
			Username:           payload.Username,
			Gender:             payload.Gender,
			WorkPhoneNumber:    payload.WorkPhoneNumber,
			AccountNumber:      payload.AccountNumber,
			LanguageSkills:     userLanguageSkills,
			RegistrationNumber: registrationNumber,
			Email:              payload.Email,
			Phone:              payload.Phone,
			NationalCode:       payload.NationalCode,
			BirthDate:          payload.BirthDate,
			AvatarUrl:          payload.AvatarUrl,
		}, params.UserID)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				if strings.Contains(err.Error(), "users_username_unique") {
					return response.ErrorBadRequest(nil, "username is duplicate")
				}
				if strings.Contains(err.Error(), "users_email_unique") {
					return response.ErrorBadRequest(nil, "email is duplicate")
				}
			}
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Update staff
		staff, err := h.StaffService.UpdateStaff(int64(user.ID), payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		user.StaffID = &staff.ID
		user.VehicleTypes = staff.VehicleTypes
		user.VehicleLicenseTypes = staff.VehicleLicenseTypes
		user.Limitations = staff.Limitations

		// Return response
		return response.Success(user)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/contract/{userId}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: StaffsCreateOrUpdateContractRequestParams
 * @apiRequestRef: StaffsCreateOrUpdateContractRequestBody
 * @apiResponseRef: StaffsCreateOrUpdateContractResponse
 * @apiSummary: Create or Update staff contract
 * @apiDescription: Create or Update staff contract
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateOrUpdateContract() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsCreateOrUpdateContractRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsCreateOrUpdateContractRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find user based on id
		user, _ := h.UserService.FindByID(params.UserID)
		if user == nil {
			return response.ErrorBadRequest(nil, "user not found")
		}
		payload.UserID = int(user.ID)

		// Generate unique organization number
		organizationNumber, err := h.StaffService.GenerateStaffOrganizationNumber()
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		payload.OrganizationNumber = organizationNumber

		// Check if paymentType exists
		_, err = h.PaymentTypeService.GetPaymentTypeByID(payload.PaymentTypeAsMetadata.ID)
		if err != nil {
			return response.ErrorBadRequest(nil, "paymentType not found")
		}

		// Find sections based on ids
		var sectionIds []int64
		for _, section := range payload.SectionsAsMetadata {
			sectionIds = append(sectionIds, section.ID)
		}
		sections, err := h.SectionService.GetSectionsByIds(sectionIds)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if len(sections) != len(sectionIds) {
			return response.ErrorBadRequest(nil, "sections are invalid")
		}
		payload.SectionIDs = sectionIds

		// Find contract types based on ids
		var staffContractIds []int64
		for _, contractType := range payload.ContractTypesAsMetadata {
			staffContractIds = append(staffContractIds, int64(contractType.ID))
		}
		contractTypes, err := h.ContractTypeService.GetContractTypesByIds(staffContractIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "contract types are invalid")
		}
		if len(contractTypes) != len(staffContractIds) {
			return response.ErrorBadRequest(nil, "contract types are invalid")
		}
		payload.ContractTypeIDs = staffContractIds

		// Find shift types based on ids
		var shiftTypeIds []int64
		for _, shiftType := range payload.ShiftTypesAsMetadata {
			shiftTypeIds = append(shiftTypeIds, int64(shiftType.ID))
		}
		shiftTypes, err := h.ShiftTypeService.GetShiftTypesByIds(shiftTypeIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "shift types are invalid")
		}
		if len(shiftTypes) != len(shiftTypeIds) {
			return response.ErrorBadRequest(nil, "shift types are invalid")
		}
		payload.ShiftTypeIDs = shiftTypeIds

		// Find staff types based on ids
		var staffTypeIds []int64
		for _, staffType := range payload.StaffTypesAsMetadata {
			staffTypeIds = append(staffTypeIds, int64(staffType.ID))
		}
		staffTypes, err := h.StaffTypeService.GetStaffTypesByIds(staffTypeIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "staff types are invalid")
		}
		if len(staffTypes) != len(staffTypeIds) {
			return response.ErrorBadRequest(nil, "staff types are invalid")
		}
		payload.StaffTypeIDs = staffTypeIds

		// Find roles based on ids
		var roleIds []int64
		for _, role := range payload.RolesAsMetadata {
			roleIds = append(roleIds, int64(role.ID))
		}
		roles, err := h.RoleService.GetRolesByIds(roleIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "roles are invalid")
		}
		if len(roles) != len(roleIds) {
			return response.ErrorBadRequest(nil, "roles are invalid")
		}
		payload.RoleIDs = roleIds

		// Validate joinedAt
		if payload.JoinedAt != "" {
			_, err = utils.TryParseToDateTime(payload.JoinedAt)
			if err != nil {
				return response.ErrorBadRequest(nil, "joinedAt is invalid")
			}
		}

		// Validate contractStartedAt
		if payload.ContractStartedAt != "" {
			_, err = utils.TryParseToDateTime(payload.ContractStartedAt)
			if err != nil {
				return response.ErrorBadRequest(nil, "contractStartedAt is invalid")
			}
		}

		// Validate contractExpiresAt
		if payload.ContractExpiresAt != "" {
			_, err = utils.TryParseToDateTime(payload.ContractExpiresAt)
			if err != nil {
				return response.ErrorBadRequest(nil, "contractExpiresAt is invalid")
			}
		}

		// Validate trialTime
		if payload.TrialTime != "" {
			_, err = utils.TryParseToDateTime(payload.TrialTime)
			if err != nil {
				return response.ErrorBadRequest(nil, "trialTime is invalid")
			}
		}

		// Update staff contract
		updatedStaff, err := h.StaffService.CreateOrUpdateContract(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Check if staff is not dispatcher remove attachments
		isDispatcher := updatedStaff.IsSubcontractor
		if !isDispatcher {
			payload.Attachments = nil
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(updatedStaff.ID))
		if len(uploadedFilesMetadata) > 0 {
			updated, err := h.StaffService.UpdateStaffAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, int64(updatedStaff.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedStaff = updated
		}

		// Return response
		return response.Success(updatedStaff)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/licenses
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsCreateLicensesRequestBody
 * @apiResponseRef: StaffsCreateLicensesResponse
 * @apiSummary: Create staff licenses
 * @apiDescription: Create staff licenses
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateLicenses() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.StaffsCreateLicensesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate expireDate
		if _, err := utils.TryParseToDateTime(payload.ExpireDate); err != nil {
			return response.ErrorBadRequest(map[string]interface{}{
				"expire_date": []string{"Expire date is invalid"},
			}, "Your request data is invalid")
		}

		// Find license based on id
		license, err := h.LicenseService.FindByID(int64(payload.License.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if license == nil {
			return response.ErrorBadRequest(map[string]interface{}{
				"license": []string{"License not found"},
			}, "Your request data is invalid")
		}

		// Find staff based on id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staff == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Check staff has license or not
		hasLicense, err := h.StaffService.HasLicense(staff.ID, uint(payload.License.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if hasLicense {
			return response.ErrorBadRequest(nil, "staff already has this license")
		}

		// Create staff licenses
		createdStaffLicense, err := h.StaffService.CreateLicenses(staff, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(createdStaffLicense.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedLicense, err := h.StaffService.UpdateLicenseAttachments(nil, uploadedFilesMetadata, int64(createdStaffLicense.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdStaffLicense = updatedLicense
		}

		// Return response
		return response.Success(createdStaffLicense)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/absences
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsCreateAbsencesRequestBody
 * @apiResponseRef: StaffsCreateAbsencesResponse
 * @apiSummary: Create staff absences
 * @apiDescription: Create staff absences
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.StaffsCreateAbsencesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate expireDate
		if _, err := utils.TryParseToDateTime(payload.StartDate); err != nil {
			return response.ErrorBadRequest(map[string]interface{}{
				"start_date": []string{"startDate is invalid"},
			}, "Your request data is invalid")
		}

		// Validate endDate
		if payload.EndDate != "" {
			if _, err := utils.TryParseToDateTime(payload.EndDate); err != nil {
				return response.ErrorBadRequest(map[string]interface{}{
					"end_date": []string{"endDate is invalid"},
				}, "Your request data is invalid")
			}
		}

		// Find staff based on id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staff == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Check staff has more than 5 absences or not
		absences, err := h.StaffService.QueryAbsences(&models.StaffsQueryAbsencesRequestParams{
			StaffID: int(staff.ID),
			Filters: models.StaffsQueryAbsencesFilterType{
				Status: filters.FilterValue[string]{
					Op:    operators.EQUALS,
					Value: constants.ABSENCE_STATUS_PENDING,
				},
			},
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if absences.TotalRows >= 5 {
			return response.ErrorBadRequest(nil, "staff can't have more than 5 absences")
		}

		// Create staff absence
		createdStaffAbsence, err := h.StaffService.CreateAbsences(staff, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		isDispatcher := h.UserService.IsDispatcher(int64(user.ID))
		if isDispatcher {
			var (
				title       = "Staff absence"
				description = fmt.Sprintf("You have new absence from %s to %s", payload.StartDate, payload.EndDate)
				status      = sharedconstants.NOTIFICATION_STATUS_ACCEPTED
			)
			_, err := h.NotificationService.CreateNotification(int64(user.ID), title, description, status, map[string]uint{
				"absenceId": createdStaffAbsence.ID,
			})
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
		} else {
			var (
				title       = "Staff absence"
				description = fmt.Sprintf("You have new absence request from %s %s", user.FirstName, user.LastName)
				status      = sharedconstants.NOTIFICATION_STATUS_PENDING
			)
			_, err := h.NotificationService.CreateSystemNotification(int64(user.ID), title, description, status, map[string]uint{
				"absenceId": createdStaffAbsence.ID,
			})
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(createdStaffAbsence.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedAbsence, err := h.StaffService.UpdateAbsenceAttachments(nil, uploadedFilesMetadata, int64(createdStaffAbsence.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdStaffAbsence = updatedAbsence
		}

		// Return response
		return response.Success(createdStaffAbsence)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsDeleteRequestBody
 * @apiResponseRef: StaffsDeleteResponse
 * @apiSummary: Delete staff
 * @apiDescription: Delete staff
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.StaffsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.StaffService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/licenses
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsDeleteLicensesRequestBody
 * @apiResponseRef: StaffsDeleteLicensesResponse
 * @apiSummary: Delete staff licenses
 * @apiDescription: Delete staff licenses
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) DeleteLicenses() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.StaffsDeleteLicensesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.StaffService.DeleteLicenses(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/licenses/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: StaffsUpdateLicenseRequestParams
 * @apiRequestRef: StaffsUpdateLicenseRequestBody
 * @apiResponseRef: StaffsUpdateLicenseResponse
 * @apiSummary: Update staff license
 * @apiDescription: Update staff license
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) UpdateLicense() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsUpdateLicenseRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsUpdateLicenseRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate expireDate
		if _, err := utils.TryParseToDateTime(payload.ExpireDate); err != nil {
			return response.ErrorBadRequest(nil, "expireDate is invalid")
		}

		// Find license based on id
		license, err := h.LicenseService.FindByID(int64(payload.License.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if license == nil {
			return response.ErrorBadRequest(nil, "license not found")
		}

		// Find staff license based on id
		staffLicense, err := h.StaffService.FindStaffLicenseByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staffLicense == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Check staff has license or not
		hasLicense, err := h.StaffService.HasLicenseExcept(staffLicense.StaffID, uint(payload.License.ID), staffLicense.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if hasLicense {
			return response.ErrorBadRequest(nil, "staff already has this license")
		}

		// Update staff license
		updatedStaffLicense, err := h.StaffService.UpdateLicense(staffLicense, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(updatedStaffLicense.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedLicense, err := h.StaffService.UpdateLicenseAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, int64(updatedStaffLicense.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedStaffLicense = updatedLicense
		}

		// Return response
		return response.Success(updatedStaffLicense)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/absences
* @apiResponseRef: StaffsQueryAbsencesResponse
* @apiSummary: Query staff absenses
* @apiParametersRef: StaffsQueryAbsencesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryAbsencesRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffs, err := h.StaffService.QueryAbsences(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffs)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/absences/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: StaffsUpdateAbsenceRequestParams
 * @apiRequestRef: StaffsUpdateAbsenceRequestBody
 * @apiResponseRef: StaffsUpdateAbsenceResponse
 * @apiSummary: Update staff absence
 * @apiDescription: Update staff absence
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) UpdateAbsence() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsUpdateAbsenceRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsUpdateAbsenceRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate startDate
		if _, err := utils.TryParseToDateTime(payload.StartDate); err != nil {
			return response.ErrorBadRequest(map[string]interface{}{
				"start_date": []string{"startDate is invalid"},
			}, "Your request data is invalid")
		}

		// Validate endDate
		if payload.EndDate != "" {
			if _, err := utils.TryParseToDateTime(payload.EndDate); err != nil {
				return response.ErrorBadRequest(map[string]interface{}{
					"end_date": []string{"endDate is invalid"},
				}, "Your request data is invalid")
			}
		}

		// Find staff license based on id
		staffAbsence, err := h.StaffService.FindStaffAbsenceByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staffAbsence == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update staff absence
		updatedStaffAbsence, err := h.StaffService.UpdateAbsence(staffAbsence, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(updatedStaffAbsence.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedAbsence, err := h.StaffService.UpdateAbsenceAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, int64(updatedStaffAbsence.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedStaffAbsence = &domain.StaffAbsenceRes{
				ID:          updatedAbsence.ID,
				StaffID:     updatedAbsence.StaffID,
				StartDate:   updatedAbsence.StartDate,
				EndDate:     updatedAbsence.EndDate,
				Reason:      updatedAbsence.Reason,
				Attachments: updatedAbsence.Attachments,
				Status:      updatedAbsence.Status,
				StatusBy:    updatedAbsence.StatusBy,
				StatusAt:    updatedAbsence.StatusAt,
			}
		}

		// Return response
		return response.Success(updatedStaffAbsence)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/absences
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsDeleteAbsencesRequestBody
 * @apiResponseRef: StaffsDeleteAbsencesResponse
 * @apiSummary: Delete staff absences
 * @apiDescription: Delete staff absences
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) DeleteAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.StaffsDeleteAbsencesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.StaffService.DeleteAbsences(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/profile
* @apiResponseRef: StaffsQueryProfileResponse
* @apiSummary: Query staff profile
* @apiParametersRef: StaffsQueryProfileRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryProfile() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryProfileRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		result, err := h.StaffService.QueryProfile(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(result)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/otherattachments
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: StaffsCreateOtherAttachmentsRequestBody
 * @apiResponseRef: StaffsCreateOtherAttachmentsResponse
 * @apiSummary: Create staff other attachments
 * @apiDescription: Create other attachments
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.StaffsCreateOtherAttachmentsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff based on id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staff == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Create staff other attachment
		createdStaffOtherAttachments, err := h.StaffService.CreateOtherAttachments(staff, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(createdStaffOtherAttachments.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedOtherAttachment, err := h.StaffService.UpdateStaffOtherAttachments(uploadedFilesMetadata, int64(createdStaffOtherAttachments.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdStaffOtherAttachments = updatedOtherAttachment
		}

		// Return response
		return response.Success(createdStaffOtherAttachments)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/otherattachments
* @apiResponseRef: StaffsQueryOtherAttachmentsResponse
* @apiSummary: Query staff absenses
* @apiParametersRef: StaffsQueryOtherAttachmentsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryOtherAttachmentsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffOtherAttachments, err := h.StaffService.QueryOtherAttachments(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffOtherAttachments)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/otherattachments/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: StaffsUpdateOtherAttachmentRequestParams
 * @apiRequestRef: StaffsUpdateOtherAttachmentRequestBody
 * @apiResponseRef: StaffsUpdateOtherAttachmentResponse
 * @apiSummary: Update staff other attachment
 * @apiDescription: Update staff other attachment
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) UpdateOtherAttachment() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsUpdateOtherAttachmentRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsUpdateOtherAttachmentRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff other attachment based on id
		staffOtherAttachment, err := h.StaffService.FindStaffOtherAttachmentByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staffOtherAttachment == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update staff other attachment
		updatedStaffOtherAttachment, err := h.StaffService.UpdateStaffOtherAttachment(staffOtherAttachment, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(updatedStaffOtherAttachment.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedStaffOtherAttachments, err := h.StaffService.UpdateStaffOtherAttachments(uploadedFilesMetadata, int64(updatedStaffOtherAttachment.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedStaffOtherAttachment = updatedStaffOtherAttachments
		}

		// Return response
		return response.Success(updatedStaffOtherAttachment)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/otherattachments
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsDeleteOtherAttachmentsRequestBody
 * @apiResponseRef: StaffsDeleteOtherAttachmentsResponse
 * @apiSummary: Delete staff other attachments
 * @apiDescription: Delete staff other attachments
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) DeleteOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.StaffsDeleteOtherAttachmentsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.StaffService.DeleteStaffOtherAttachments(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/libraries
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: StaffsCreateLibrariesRequestBody
 * @apiResponseRef: StaffsCreateLibrariesResponse
 * @apiSummary: Create staff librarys
 * @apiDescription: Create librarys
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateLibraries() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.StaffsCreateLibrariesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff based on id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staff == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Create staff library
		createdStaffLibraries, err := h.StaffService.CreateLibraries(staff, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(createdStaffLibraries.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedOtherAttachment, err := h.StaffService.UpdateStaffLibraries(uploadedFilesMetadata, int64(createdStaffLibraries.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdStaffLibraries = updatedOtherAttachment
		}

		// Return response
		return response.Success(createdStaffLibraries)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/libraries
* @apiResponseRef: StaffsQueryLibrariesResponse
* @apiSummary: Query staff libraries
* @apiParametersRef: StaffsQueryLibrariesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryLibraries() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryLibrariesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffLibraries, err := h.StaffService.QueryLibraries(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffLibraries)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/libraries/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: StaffsUpdateLibraryRequestParams
 * @apiRequestRef: StaffsUpdateLibraryRequestBody
 * @apiResponseRef: StaffsUpdateLibraryResponse
 * @apiSummary: Update staff library
 * @apiDescription: Update staff library
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) UpdateLibrary() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.StaffsUpdateLibraryRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.StaffsUpdateLibraryRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff library based on id
		staffLibrary, err := h.StaffService.FindStaffLibraryByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if staffLibrary == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update staff library
		updatedStaffLibrary, err := h.StaffService.UpdateStaffLibrary(staffLibrary, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(updatedStaffLibrary.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedStaffLibraries, err := h.StaffService.UpdateStaffLibraries(uploadedFilesMetadata, int64(updatedStaffLibrary.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedStaffLibrary = updatedStaffLibraries
		}

		// Return response
		return response.Success(updatedStaffLibrary)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/libraries
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: StaffsDeleteLibrariesRequestBody
 * @apiResponseRef: StaffsDeleteLibrariesResponse
 * @apiSummary: Delete staff librarys
 * @apiDescription: Delete staff librarys
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) DeleteLibraries() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.StaffsDeleteLibrariesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.StaffService.DeleteStaffLibraries(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/chats
* @apiResponseRef: StaffsQueryChatsResponse
* @apiSummary: Query staff chats
* @apiParametersRef: StaffsQueryChatsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryChatsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryChatsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryChats() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryChatsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query chats
		chats, err := h.StaffService.QueryChats(queries)
		if err != nil {
			log.Printf("StaffHandler.QueryChats: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(chats)
	})
}

/*
* @apiTag: staff
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffs/chats/messages
* @apiResponseRef: StaffsQueryChatMessagesResponse
* @apiSummary: Query staff chat messages
* @apiParametersRef: StaffsQueryChatMessagesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffsQueryChatMessagesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffsQueryChatMessagesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) QueryChatMessages() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffsQueryChatMessagesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query chat messages
		chatMessages, err := h.StaffService.QueryChatMessages(queries)
		if err != nil {
			log.Printf("StaffHandler.QueryChatMessages: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(chatMessages)
	})
}

/*
 * @apiTag: staff
 * @apiPath: /staffs/chats/messages
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: StaffsCreateChatMessageRequestBody
 * @apiResponseRef: StaffsCreateChatMessageResponse
 * @apiSummary: Create chat message
 * @apiDescription: Create chat message
 * @apiSecurity: apiKeySecurity
 */
func (h *StaffHandler) CreateChatMessage() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.StaffsCreateChatMessageRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Create chat message
		chatMessage, err := h.StaffService.CreateChatMessage(payload)
		if err != nil {
			log.Printf("StaffHandler.CreateChatMessage: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.STAFF_BUCKET_NAME, payload.Attachments, int64(chatMessage.ID))
		if len(uploadedFilesMetadata) > 0 {
			updated, err := h.StaffService.UpdateChatMessageAttachments(nil, uploadedFilesMetadata, int64(chatMessage.ID))
			if err != nil {
				log.Printf("StaffHandler.CreateChatMessage: %s\n", err.Error())
				return response.ErrorInternalServerError(nil, err.Error())
			}
			chatMessage = updated
		}

		// Return created chat message
		return response.Success(chatMessage)
	})
}
