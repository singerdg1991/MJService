package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/security"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	aPorts "github.com/hoitek/Maja-Service/internal/address/ports"
	"github.com/hoitek/Maja-Service/internal/customer/config"
	"github.com/hoitek/Maja-Service/internal/customer/constants"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/customer/models"
	"github.com/hoitek/Maja-Service/internal/customer/ports"
	lssPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	lPorts "github.com/hoitek/Maja-Service/internal/limitation/ports"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	sPorts "github.com/hoitek/Maja-Service/internal/section/ports"
	ssPorts "github.com/hoitek/Maja-Service/internal/service/ports"
	sgsPorts "github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	stPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type CustomerHandler struct {
	CustomerService      ports.CustomerService
	UserService          uPorts.UserService
	RoleService          rPorts.RoleService
	SectionService       sPorts.SectionService
	LimitationService    lPorts.LimitationService
	AddressService       aPorts.AddressService
	ServiceService       ssPorts.ServiceService
	StaffService         stPorts.StaffService
	ServiceGradeService  sgsPorts.ServiceGradeService
	LanguageSkillService lssPorts.LanguageSkillService
	S3Service            s3Ports.S3Service
}

func NewCustomerHandler(r *mux.Router,
	nService ports.CustomerService,
	uService uPorts.UserService,
	rService rPorts.RoleService,
	sService sPorts.SectionService,
	lService lPorts.LimitationService,
	aService aPorts.AddressService,
	ssService ssPorts.ServiceService,
	stService stPorts.StaffService,
	sgsService sgsPorts.ServiceGradeService,
	lssService lssPorts.LanguageSkillService,
	s3Service s3Ports.S3Service,
) (CustomerHandler, error) {
	customerHandler := CustomerHandler{
		CustomerService:      nService,
		UserService:          uService,
		RoleService:          rService,
		SectionService:       sService,
		LimitationService:    lService,
		AddressService:       aService,
		ServiceService:       ssService,
		StaffService:         stService,
		ServiceGradeService:  sgsService,
		LanguageSkillService: lssService,
		S3Service:            s3Service,
	}
	if r == nil {
		return CustomerHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.CustomerConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.CustomerConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(uService, []string{}))

	rv1.Handle("/customers/creditdetails", customerHandler.CreateCreditDetails()).Methods(http.MethodPost)
	rv1.Handle("/customers/creditdetails", customerHandler.QueryCreditDetails()).Methods(http.MethodGet)
	rv1.Handle("/customers/creditdetails/{id}/{customerid}", customerHandler.UpdateCreditDetails()).Methods(http.MethodPut)
	rv1.Handle("/customers/creditdetails", customerHandler.DeleteCreditDetails()).Methods(http.MethodDelete)
	rv1.Handle("/customers/absences", customerHandler.CreateAbsences()).Methods(http.MethodPost)
	rv1.Handle("/customers/absences", customerHandler.QueryAbsences()).Methods(http.MethodGet)
	rv1.Handle("/customers/absences/{id}", customerHandler.UpdateAbsence()).Methods(http.MethodPut)
	rv1.Handle("/customers/absences", customerHandler.DeleteAbsences()).Methods(http.MethodDelete)
	rv1.Handle("/customers/services", customerHandler.CreateServices()).Methods(http.MethodPost)
	rv1.Handle("/customers/services", customerHandler.QueryServices()).Methods(http.MethodGet)
	rv1.Handle("/customers/services/{id}/{customerid}", customerHandler.UpdateService()).Methods(http.MethodPut)
	rv1.Handle("/customers/services", customerHandler.DeleteServices()).Methods(http.MethodDelete)
	rv1.Handle("/customers/medicines", customerHandler.CreateMedicines()).Methods(http.MethodPost)
	rv1.Handle("/customers/medicines", customerHandler.QueryMedicines()).Methods(http.MethodGet)
	rv1.Handle("/customers/medicines/{id}/{customerid}", customerHandler.UpdateMedicine()).Methods(http.MethodPut)
	rv1.Handle("/customers/medicines", customerHandler.DeleteMedicines()).Methods(http.MethodDelete)
	rv1.Handle("/customers", customerHandler.Query()).Methods(http.MethodGet)

	rAuth.Handle("/customers/personalinfo/{id}/{customerid}", customerHandler.UpdatePersonalInfo()).Methods(http.MethodPut)
	rAuth.Handle("/customers/personalinfo", customerHandler.CreatePersonalInfo()).Methods(http.MethodPost)
	rAuth.Handle("/customers/additionalinfo/{userid}", customerHandler.UpdateAdditionalInfo()).Methods(http.MethodPut)
	rAuth.Handle("/customers/otherattachments", customerHandler.CreateOtherAttachments()).Methods(http.MethodPost)
	rAuth.Handle("/customers/otherattachments", customerHandler.QueryOtherAttachments()).Methods(http.MethodGet)
	rAuth.Handle("/customers/otherattachments/{id}", customerHandler.UpdateOtherAttachment()).Methods(http.MethodPut)
	rAuth.Handle("/customers/otherattachments", customerHandler.DeleteOtherAttachments()).Methods(http.MethodDelete)
	rAuth.Handle("/customers/relatives", customerHandler.CreateRelatives()).Methods(http.MethodPost)
	rAuth.Handle("/customers/relatives", customerHandler.QueryRelatives()).Methods(http.MethodGet)
	rAuth.Handle("/customers/relatives/{id}", customerHandler.UpdateRelative()).Methods(http.MethodPut)
	rAuth.Handle("/customers/relatives", customerHandler.DeleteRelatives()).Methods(http.MethodDelete)
	rAuth.Handle("/customers/logs/contractualmobilityrestriction", customerHandler.QueryContractualMobilityRestrictionLogs()).Methods(http.MethodGet)
	rAuth.Handle("/customers/logs/status", customerHandler.QueryStatusLogs()).Methods(http.MethodGet)
	return customerHandler, nil
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: customers
* @apiResponseRef: CustomersQueryResponse
* @apiSummary: Query customers
* @apiParametersRef: CustomersQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
 */
func (h *CustomerHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customers, err := h.CustomerService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customers)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/personalinfo
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreatePersonalInfoRequestBody
 * @apiResponseRef: CustomersCreatePersonalInfoResponse
 * @apiSummary: Create customer personal info
 * @apiDescription: Create customer personal info
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) CreatePersonalInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreatePersonalInfoRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Find mother lang ids based on ids
		motherLangIds, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.MotherLangIDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "MotherLangIDs are invalid")
		}
		languageSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(motherLangIds)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if len(languageSkills) != len(motherLangIds) {
			return response.ErrorBadRequest(nil, "MotherLangIDs are invalid")
		}
		payload.MotherLangIDsInt64 = motherLangIds

		// Check limitations
		if payload.Limitations != nil {
			for _, limitation := range payload.Limitations {
				_, err := h.LimitationService.FindByID(int64(limitation.LimitationID))
				if err != nil {
					return response.ErrorBadRequest(nil, fmt.Sprintf("limitation with id %d is not found", limitation.LimitationID))
				}
			}
		}

		// Create customer
		createdCustomer, err := h.CustomerService.CreatePersonalInfo(payload)
		if err != nil {
			log.Printf("Error in creating customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in creating customer")
		}

		// Create user for customer
		payload.Password = security.HashPassword(payload.Password)
		user, err := h.UserService.CreateUserForCustomer(&sharedmodels.CustomersCreatePersonalInfo{
			CustomerID:           createdCustomer.ID,
			FirstName:            payload.FirstName,
			LastName:             payload.LastName,
			Gender:               payload.Gender,
			DateOfBirth:          payload.DateOfBirth,
			NationalCode:         payload.NationalCode,
			Email:                payload.Email,
			PhoneNumber:          payload.PhoneNumber,
			Password:             payload.Password,
			ForcedChangePassword: payload.ForcedChangePassword,
		})
		if err != nil {
			log.Printf("Error in creating user for customer: %s", err.Error())
			if strings.Contains(err.Error(), "users_email_unique") {
				return response.ErrorBadRequest(nil, "email is already taken")
			}
			return response.ErrorInternalServerError(nil, "something went wrong in creating customer. contact support team for more information")
		}
		if user == nil {
			return response.ErrorInternalServerError(nil, "something went wrong in creating customer. contact support team for more information")
		}
		payload.UserID = int64(user.ID)

		// Update userID in customer
		updatedCustomer, err := h.CustomerService.UpdateUserInformation(createdCustomer.ID, payload)
		if err != nil {
			log.Printf("Error in updating customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "something went wrong in creating customer. contact support team for more information")
		}

		// Return response
		return response.Success(updatedCustomer)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/personalinfo/{id}/{customerid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdatePersonalInfoRequestParams
 * @apiRequestRef: CustomersCreatePersonalInfoRequestBody
 * @apiResponseRef: CustomersCreatePersonalInfoResponse
 * @apiSummary: Update customer personal info
 * @apiDescription: Update customer personal info
 */
func (h *CustomerHandler) UpdatePersonalInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.CustomersUpdatePersonalInfoRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersCreatePersonalInfoRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Find mother lang ids based on ids
		motherLangIds, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.MotherLangIDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "MotherLangIDs are invalid")
		}
		languageSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(motherLangIds)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if len(languageSkills) != len(motherLangIds) {
			return response.ErrorBadRequest(nil, "MotherLangIDs are invalid")
		}
		payload.MotherLangIDsInt64 = motherLangIds

		// Check limitations
		if payload.Limitations != nil {
			for _, limitation := range payload.Limitations {
				_, err := h.LimitationService.FindByID(int64(limitation.LimitationID))
				if err != nil {
					return response.ErrorBadRequest(nil, fmt.Sprintf("limitation with id %d is not found", limitation.LimitationID))
				}
			}
		}

		// Update customer
		updatedCustomer, err := h.CustomerService.UpdatePersonalInfo(int64(params.CustomerID), payload)
		if err != nil {
			log.Printf("Error in updating customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in updating customer")
		}

		// Get user for customer
		payload.UserID = *updatedCustomer.UserID
		user, err := h.UserService.UpdateUserForCustomer(payload.UserID, &sharedmodels.CustomersCreatePersonalInfo{
			CustomerID:   updatedCustomer.ID,
			FirstName:    payload.FirstName,
			LastName:     payload.LastName,
			Gender:       payload.Gender,
			NationalCode: payload.NationalCode,
			DateOfBirth:  payload.DateOfBirth,
			Email:        payload.Email,
			PhoneNumber:  payload.PhoneNumber,
		})
		if err != nil {
			log.Printf("Error in updating user for customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "something went wrong in updating customer. contact support team for more information")
		}
		if user == nil {
			return response.ErrorInternalServerError(nil, "something went wrong in updating customer. contact support team for more information")
		}

		// Update userID in customer
		updatedUserInformation, err := h.CustomerService.UpdateUserInformation(updatedCustomer.ID, payload)
		if err != nil {
			log.Printf("Error in updating customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "something went wrong in updating customer. contact support team for more information")
		}

		// Return response
		return response.Success(updatedUserInformation)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/additionalinfo/{userid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateAdditionalInfoRequestParams
 * @apiRequestRef: CustomersUpdateAdditionalInfoRequestBody
 * @apiResponseRef: CustomersCreatePersonalInfoResponse
 * @apiSummary: Update customer additional info
 * @apiDescription: Create customer additional info
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) UpdateAdditionalInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.CustomersUpdateAdditionalInfoRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersUpdateAdditionalInfoRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Find user by id
		user, err := h.UserService.FindByID(params.UserID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if user == nil {
			return response.ErrorBadRequest(nil, "user is not found")
		}
		if user.CustomerID == nil {
			return response.ErrorBadRequest(nil, "user is not a customer")
		}
		payload.UserID = int64(params.UserID)
		payload.CustomerID = int64(*user.CustomerID)

		// Find sections based on ids
		sectionIds, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.SectionIDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "SectionIDs is invalid")
		}
		payload.SectionIDsInt64 = sectionIds
		sections, err := h.SectionService.GetSectionsByIds(sectionIds)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if len(sections) != len(sectionIds) {
			return response.ErrorBadRequest(nil, "sections are invalid")
		}
		payload.SectionIDs = sectionIds

		// Find relative ids based on ids
		relativeIDsInt64, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.RelativeIDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "RelativeIDs is invalid")
		}
		payload.RelativeIDsInt64 = relativeIDsInt64

		// Find diagnose ids based on ids
		diagnoseIDsInt64, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.DiagnoseIDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "DiagnoseIDs is invalid")
		}
		payload.DiagnoseIDsInt64 = diagnoseIDsInt64

		// Update user
		updatedUser, err := h.UserService.UpdateUserAdditionalInfoForCustomer(&sharedmodels.UpdateUserAdditionalInfoForCustomer{
			UserID:     payload.UserID,
			CustomerID: payload.CustomerID,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if updatedUser == nil {
			return response.ErrorInternalServerError(nil, "something went wrong in updating user. contact support team for more information")
		}

		// Update customer
		updatedCustomer, err := h.CustomerService.UpdateAdditionalInfo(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(updatedCustomer)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/creditdetails
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateCreditDetailsRequestBody
 * @apiResponseRef: CustomersCreateCreditDetailsResponse
 * @apiSummary: Create customer credit details
 * @apiDescription: Create customer credit details
 */
func (h *CustomerHandler) CreateCreditDetails() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreateCreditDetailsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find user by id
		user, err := h.UserService.FindByID(int(payload.UserID))
		if err != nil {
			log.Printf("Error in finding user in creating credit details: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding user in creating credit details")
		}
		if user == nil {
			return response.ErrorBadRequest(nil, "user is not found")
		}
		if user.CustomerID == nil {
			return response.ErrorBadRequest(nil, "user is not a customer")
		}
		payload.CustomerID = int64(*user.CustomerID)

		// Find address by id
		address, err := h.AddressService.FindByID(payload.BillingAddressID)
		if err != nil {
			log.Printf("Error in finding address in creating credit details: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding address in creating credit details")
		}
		if address == nil {
			return response.ErrorBadRequest(nil, "address is not found")
		}
		if *address.CustomerID != *user.CustomerID {
			return response.ErrorBadRequest(nil, "address is not belong to customer")
		}

		// Create credit details
		creditDetails, err := h.CustomerService.CreateCreditDetails(payload)
		if err != nil {
			log.Printf("Error in creating credit details for customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in creating credit details for customer")
		}

		// Return response
		return response.Success(creditDetails)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/creditdetails
* @apiResponseRef: CustomersQueryCreditDetailsResponse
* @apiSummary: Query customer credit details
* @apiParametersRef: CustomersQueryCreditDetailsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryCreditDetailsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryCreditDetailsNotFoundResponse
 */
func (h *CustomerHandler) QueryCreditDetails() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryCreditDetailsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		creditDetails, err := h.CustomerService.QueryCreditDetails(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(creditDetails)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/creditdetails/{id}/{customerid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateCreditDetailsRequestParams
 * @apiRequestRef: CustomersUpdateCreditDetailsRequestBody
 * @apiResponseRef: CustomersCreateCreditDetailsResponse
 * @apiSummary: Update customer credit details
 * @apiDescription: Update customer credit details
 */
func (h *CustomerHandler) UpdateCreditDetails() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.CustomersUpdateCreditDetailsRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.CustomersUpdateCreditDetailsRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find credit details by id
		foundCreditDetails, err := h.CustomerService.FindCreditDetailsByIDAndCustomerID(params.ID, params.CustomerID)
		if err != nil {
			log.Printf("Error in finding credit details in updating credit details: %s", err.Error())
			return response.ErrorBadRequest(nil, "credit details is not found")
		}
		if foundCreditDetails == nil {
			return response.ErrorBadRequest(nil, "credit details is not found")
		}
		payload.ID = int64(params.ID)

		// Find customer by id
		customer, err := h.CustomerService.FindByID(params.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer in updating credit details: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding customer in updating credit details")
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "customer is not found")
		}
		payload.CustomerID = int64(params.CustomerID)

		// Find address by id
		address, err := h.AddressService.FindByID(payload.BillingAddressID)
		if err != nil {
			log.Printf("Error in finding address in creating credit details: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding address in updating credit details")
		}
		if address == nil {
			return response.ErrorBadRequest(nil, "address is not found")
		}
		if address.CustomerID == nil || (*address.CustomerID != uint(params.CustomerID)) {
			return response.ErrorBadRequest(nil, "address is not belong to customer")
		}

		// Update credit details
		creditDetails, err := h.CustomerService.UpdateCustomerCreditDetails(payload)
		if err != nil {
			log.Printf("Error in creating credit details for customer: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in creating credit details for customer")
		}

		return response.Success(creditDetails)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/creditdetails
 * @apiMethod: DELETE
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersDeleteCreditDetailsRequestBody
 * @apiResponseRef: CustomersDeleteCreditDetailsResponse
 * @apiSummary: Delete customer credit details
 * @apiDescription: Delete customer credit details
 */
func (h *CustomerHandler) DeleteCreditDetails() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteCreditDetailsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteCustomerCreditDetails(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/absences
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateAbsencesRequestBody
 * @apiResponseRef: CustomersCreateAbsencesResponse
 * @apiSummary: Create customer absences
 * @apiDescription: Create customer absences
 */
func (h *CustomerHandler) CreateAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreateAbsencesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer based on id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer in creating absences: %s", err.Error())
			return response.ErrorInternalServerError(nil, "customer not found")
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Create customer absence
		createdCustomerAbsence, err := h.CustomerService.CreateAbsences(customer, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, createdCustomerAbsence.ID)
		if len(uploadedFilesMetadata) > 0 {
			updatedAbsence, err := h.CustomerService.UpdateAbsenceAttachments(nil, uploadedFilesMetadata, createdCustomerAbsence.ID)
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdCustomerAbsence = updatedAbsence
		}

		// Return response
		return response.Success(createdCustomerAbsence)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/absences
* @apiResponseRef: CustomersQueryAbsencesResponse
* @apiSummary: Query customer absenses
* @apiParametersRef: CustomersQueryAbsencesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
 */
func (h *CustomerHandler) QueryAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryAbsencesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customerAbsences, err := h.CustomerService.QueryAbsences(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customerAbsences)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/absences/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateAbsenceRequestParams
 * @apiRequestRef: CustomersUpdateAbsenceRequestBody
 * @apiResponseRef: CustomersUpdateAbsenceResponse
 * @apiSummary: Update customer absence
 * @apiDescription: Update customer absence
 */
func (h *CustomerHandler) UpdateAbsence() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CustomersUpdateAbsenceRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersUpdateAbsenceRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer absence based on id
		customerAbsence, err := h.CustomerService.FindCustomerAbsenceByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if customerAbsence == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update customer absence
		updatedCustomerAbsence, err := h.CustomerService.UpdateAbsence(customerAbsence, payload)
		if err != nil {
			log.Printf("Error in updating customer absence: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in updating customer absence, please try again later")
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, updatedCustomerAbsence.ID)
		if len(uploadedFilesMetadata) > 0 {
			updatedAbsence, err := h.CustomerService.UpdateAbsenceAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, updatedCustomerAbsence.ID)
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedCustomerAbsence = &domain.CustomerAbsence{
				ID:          updatedAbsence.ID,
				CustomerID:  updatedAbsence.CustomerID,
				StartDate:   updatedAbsence.StartDate,
				EndDate:     updatedAbsence.EndDate,
				Reason:      updatedAbsence.Reason,
				Attachments: updatedAbsence.Attachments,
			}
		}

		// Return response
		return response.Success(updatedCustomerAbsence)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/absences
 * @apiMethod: DELETE
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersDeleteAbsencesRequestBody
 * @apiResponseRef: CustomersDeleteAbsencesResponse
 * @apiSummary: Delete customer absences
 * @apiDescription: Delete customer absences
 */
func (h *CustomerHandler) DeleteAbsences() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteAbsencesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteAbsences(payload)
		if err != nil {
			log.Printf("Error in deleting customer absences: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in deleting customer absences, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/services
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateServicesRequestBody
 * @apiResponseRef: CustomersCreateServicesResponse
 * @apiSummary: Create customer services
 * @apiDescription: Create customer services
 */
func (h *CustomerHandler) CreateServices() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreateServicesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer based on id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer in creating services: %s", err.Error())
			return response.ErrorInternalServerError(nil, "customer not found")
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Find service based on id
		service, err := h.ServiceService.FindByID(int64(payload.ServiceID))
		if err != nil {
			log.Printf("Error in finding service in creating services: %s", err.Error())
			return response.ErrorInternalServerError(nil, "service not found")
		}
		if service == nil {
			return response.ErrorBadRequest(nil, "service not found")
		}

		// Find staff based on id
		//staffWish, err := h.StaffService.FindByID(payload.NurseWishID)
		//if err != nil {
		//	log.Printf("Error in finding staff in creating services: %s", err.Error())
		//	return response.ErrorInternalServerError(nil, "staffWish not found")
		//}
		//if staffWish == nil {
		//	return response.ErrorBadRequest(nil, "staffWish not found")
		//}

		// Create customer service
		createdCustomerService, err := h.CustomerService.CreateServices(customer, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Return response
		return response.Success(createdCustomerService)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/services
* @apiResponseRef: CustomersQueryServicesResponse
* @apiSummary: Query customer services
* @apiParametersRef: CustomersQueryServicesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryServicesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryServicesNotFoundResponse
 */
func (h *CustomerHandler) QueryServices() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryServicesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		services, err := h.CustomerService.QueryServices(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(services)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/services/{id}/{customerid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateServiceRequestParams
 * @apiRequestRef: CustomersCreateServicesRequestBody
 * @apiResponseRef: CustomersCreateServicesResponse
 * @apiSummary: Update customer service
 * @apiDescription: Update customer service
 */
func (h *CustomerHandler) UpdateService() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CustomersUpdateServiceRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersCreateServicesRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer service based on id
		customerService, err := h.CustomerService.FindCustomerServiceByIDAndCustomerID(params.ID, params.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer service: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding customer service, please try again later")
		}
		if customerService == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Find service based on id
		service, err := h.ServiceService.FindByID(int64(payload.ServiceID))
		if err != nil {
			log.Printf("Error in finding service in creating services: %s", err.Error())
			return response.ErrorInternalServerError(nil, "service not found")
		}
		if service == nil {
			return response.ErrorBadRequest(nil, "service not found")
		}

		// Find staff based on id
		//staffWish, err := h.StaffService.FindByID(payload.NurseWishID)
		//if err != nil {
		//	log.Printf("Error in finding staff in creating services: %s", err.Error())
		//	return response.ErrorInternalServerError(nil, "staffWish not found")
		//}
		//if staffWish == nil {
		//	return response.ErrorBadRequest(nil, "staffWish not found")
		//}

		// Update customer service
		updatedCustomerService, err := h.CustomerService.UpdateService(customerService, payload)
		if err != nil {
			log.Printf("Error in updating customer service: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in updating customer service, please try again later")
		}

		// Return response
		return response.Success(updatedCustomerService)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/services
 * @apiMethod: DELETE
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersDeleteServicesRequestBody
 * @apiResponseRef: CustomersDeleteServicesResponse
 * @apiSummary: Delete customer services
 * @apiDescription: Delete customer services
 */
func (h *CustomerHandler) DeleteServices() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteServicesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteServices(payload)
		if err != nil {
			log.Printf("Error in deleting customer services: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in deleting customer services, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/medicines
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateMedicinesRequestBody
 * @apiResponseRef: CustomersCreateMedicinesResponse
 * @apiSummary: Create customer services
 * @apiDescription: Create customer services
 */
func (h *CustomerHandler) CreateMedicines() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreateMedicinesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer based on id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer in creating medicines: %s", err.Error())
			return response.ErrorInternalServerError(nil, "customer not found")
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Create customer medicine
		createdCustomerMedicine, err := h.CustomerService.CreateMedicines(customer, payload)
		if err != nil {
			if strings.Contains(err.Error(), "fk_customermedicines_prescriptionid") {
				return response.ErrorBadRequest(nil, "Prescription not found")
			}
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, int64(createdCustomerMedicine.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedMedicine, err := h.CustomerService.UpdateMedicineAttachments(nil, uploadedFilesMetadata, int64(createdCustomerMedicine.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdCustomerMedicine = updatedMedicine
		}

		// Return response
		return response.Success(createdCustomerMedicine)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/medicines
* @apiResponseRef: CustomersQueryMedicinesResponse
* @apiSummary: Query customer medicines
* @apiParametersRef: CustomersQueryMedicinesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryMedicinesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryMedicinesNotFoundResponse
 */
func (h *CustomerHandler) QueryMedicines() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryMedicinesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		medicines, err := h.CustomerService.QueryMedicines(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(medicines)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/medicines/{id}/{customerid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateMedicineRequestParams
 * @apiRequestRef: CustomersUpdateMedicinesRequestBody
 * @apiResponseRef: CustomersCreateMedicinesResponse
 * @apiSummary: Update customer medicine
 * @apiDescription: Update customer medicine
 */
func (h *CustomerHandler) UpdateMedicine() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CustomersUpdateMedicineRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersUpdateMedicinesRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer medicine based on id
		customerMedicine, err := h.CustomerService.FindCustomerMedicineByIDAndCustomerID(params.ID, params.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer medicine: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in finding customer medicine, please try again later")
		}
		if customerMedicine == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update customer medicine
		updatedCustomerMedicine, err := h.CustomerService.UpdateMedicine(customerMedicine, payload)
		if err != nil {
			log.Printf("Error in updating customer medicine: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in updating customer medicine, please try again later")
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, int64(updatedCustomerMedicine.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedMedicine, err := h.CustomerService.UpdateMedicineAttachments(payload.PreviousAttachmentsMetadata, uploadedFilesMetadata, int64(updatedCustomerMedicine.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedCustomerMedicine = updatedMedicine
		}

		// Return response
		return response.Success(updatedCustomerMedicine)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/medicines
 * @apiMethod: DELETE
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersDeleteMedicinesRequestBody
 * @apiResponseRef: CustomersDeleteMedicinesResponse
 * @apiSummary: Delete customer medicines
 * @apiDescription: Delete customer medicines
 */
func (h *CustomerHandler) DeleteMedicines() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteMedicinesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteMedicines(payload)
		if err != nil {
			log.Printf("Error in deleting customer medicines: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in deleting customer medicines, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/otherattachments
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateOtherAttachmentsRequestBody
 * @apiResponseRef: CustomersCreateOtherAttachmentsResponse
 * @apiSummary: Create customer other attachments
 * @apiDescription: Create other attachments
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) CreateOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CustomersCreateOtherAttachmentsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer based on id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "customer not found")
		}

		// Create customer other attachment
		createdCustomerOtherAttachments, err := h.CustomerService.CreateOtherAttachments(customer, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, int64(createdCustomerOtherAttachments.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedOtherAttachment, err := h.CustomerService.UpdateCustomerOtherAttachments(uploadedFilesMetadata, int64(createdCustomerOtherAttachments.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			createdCustomerOtherAttachments = updatedOtherAttachment
		}

		// Return response
		return response.Success(createdCustomerOtherAttachments)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/otherattachments
* @apiResponseRef: CustomersQueryOtherAttachmentsResponse
* @apiSummary: Query customer absenses
* @apiParametersRef: CustomersQueryOtherAttachmentsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) QueryOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryOtherAttachmentsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customerOtherAttachments, err := h.CustomerService.QueryOtherAttachments(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customerOtherAttachments)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/otherattachments/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateOtherAttachmentRequestParams
 * @apiRequestRef: CustomersUpdateOtherAttachmentRequestBody
 * @apiResponseRef: CustomersUpdateOtherAttachmentResponse
 * @apiSummary: Update customer other attachment
 * @apiDescription: Update customer other attachment
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) UpdateOtherAttachment() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CustomersUpdateOtherAttachmentRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersUpdateOtherAttachmentRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer other attachment based on id
		customerOtherAttachment, err := h.CustomerService.FindCustomerOtherAttachmentByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if customerOtherAttachment == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update customer other attachment
		updatedCustomerOtherAttachment, err := h.CustomerService.UpdateCustomerOtherAttachment(customerOtherAttachment, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CUSTOMER_BUCKET_NAME, payload.Attachments, int64(updatedCustomerOtherAttachment.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedCustomerOtherAttachments, err := h.CustomerService.UpdateCustomerOtherAttachments(uploadedFilesMetadata, int64(updatedCustomerOtherAttachment.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			updatedCustomerOtherAttachment = updatedCustomerOtherAttachments
		}

		// Return response
		return response.Success(updatedCustomerOtherAttachment)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/otherattachments
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: CustomersDeleteOtherAttachmentsRequestBody
 * @apiResponseRef: CustomersDeleteOtherAttachmentsResponse
 * @apiSummary: Delete customer other attachments
 * @apiDescription: Delete customer other attachments
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) DeleteOtherAttachments() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteOtherAttachmentsRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteCustomerOtherAttachments(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/relatives
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersCreateRelativesRequestBody
 * @apiResponseRef: CustomersCreateRelativesResponse
 * @apiSummary: Create customer relatives
 * @apiDescription: Create customer relatives
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) CreateRelatives() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CustomersCreateRelativesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer based on id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil {
			log.Printf("Error in finding customer in creating relatives: %s", err.Error())
			return response.ErrorInternalServerError(nil, "customer not found")
		}
		if customer == nil {
			return response.ErrorBadRequest(nil, "staff not found")
		}

		// Create customer relative
		createdCustomerRelative, err := h.CustomerService.CreateRelatives(customer, payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Return response
		return response.Success(createdCustomerRelative)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/relatives
* @apiResponseRef: CustomersQueryRelativesResponse
* @apiSummary: Query customer relatives
* @apiParametersRef: CustomersQueryRelativesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) QueryRelatives() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryRelativesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customerRelatives, err := h.CustomerService.QueryRelatives(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customerRelatives)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/relatives/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CustomersUpdateRelativeRequestParams
 * @apiRequestRef: CustomersCreateRelativesRequestBody
 * @apiResponseRef: CustomersCreateRelativesResponse
 * @apiSummary: Update customer relative
 * @apiDescription: Update customer relative
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) UpdateRelative() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CustomersUpdateRelativeRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CustomersCreateRelativesRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find customer relative based on id
		customerRelative, err := h.CustomerService.FindCustomerRelativeByID(params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if customerRelative == nil {
			return response.ErrorBadRequest(nil, "id is invalid")
		}

		// Update customer relative
		updatedCustomerRelative, err := h.CustomerService.UpdateRelative(customerRelative, payload)
		if err != nil {
			log.Printf("Error in updating customer relative: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in updating customer relative, please try again later")
		}

		// Return response
		return response.Success(updatedCustomerRelative)
	})
}

/*
 * @apiTag: customer
 * @apiPath: /customers/relatives
 * @apiMethod: DELETE
 * @apiStatusCode: 200
 * @apiRequestRef: CustomersDeleteRelativesRequestBody
 * @apiResponseRef: CustomersDeleteRelativesResponse
 * @apiSummary: Delete customer relatives
 * @apiDescription: Delete customer relatives
 * @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) DeleteRelatives() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CustomersDeleteRelativesRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := sharedutils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CustomerService.DeleteRelatives(payload)
		if err != nil {
			log.Printf("Error in deleting customer relatives: %s", err.Error())
			return response.ErrorInternalServerError(nil, "Error in deleting customer relatives, please try again later")
		}

		return response.Success(data)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/logs/contractualmobilityrestriction
* @apiResponseRef: CustomersQueryContractualMobilityRestrictionLogsResponse
* @apiSummary: Query customer contractual mobility restriction
* @apiParametersRef: CustomersQueryContractualMobilityRestrictionLogsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) QueryContractualMobilityRestrictionLogs() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryContractualMobilityRestrictionLogsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customerContractualMobilityRestrictionLogs, err := h.CustomerService.QueryContractualMobilityRestrictionLogs(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customerContractualMobilityRestrictionLogs)
	})
}

/*
* @apiTag: customer
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /customers/logs/status
* @apiResponseRef: CustomersQueryStatusLogsResponse
* @apiSummary: Query customer status
* @apiParametersRef: CustomersQueryStatusLogsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CustomersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CustomersQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CustomerHandler) QueryStatusLogs() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CustomersQueryStatusLogsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		customerStatusLogs, err := h.CustomerService.QueryStatusLogs(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(customerStatusLogs)
	})
}
