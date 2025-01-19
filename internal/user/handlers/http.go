package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/otp"
	"github.com/hoitek/Maja-Service/internal/_shared/security"
	lPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	oPorts "github.com/hoitek/Maja-Service/internal/otp/ports"
	rolePorts "github.com/hoitek/Maja-Service/internal/role/ports"
	nPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	"github.com/hoitek/Maja-Service/internal/user/config"
	"github.com/hoitek/Maja-Service/internal/user/models"
	"github.com/hoitek/Maja-Service/internal/user/ports"
	"log"
	"net/http"
	"strings"
)

type UserHandler struct {
	UserService          ports.UserService
	RoleService          rolePorts.RoleService
	LanguageSkillService lPorts.LanguageSkillService
	StaffService         nPorts.StaffService
	OTPService           oPorts.OTPService
}

func NewUserHandler(
	r *mux.Router,
	uService ports.UserService,
	rService rolePorts.RoleService,
	lService lPorts.LanguageSkillService,
	nService nPorts.StaffService,
	oService oPorts.OTPService,
) (UserHandler, error) {
	userHandler := UserHandler{
		RoleService:          rService,
		UserService:          uService,
		LanguageSkillService: lService,
		StaffService:         nService,
		OTPService:           oService,
	}
	if r == nil {
		return UserHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.UserConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.UserConfig.ApiVersion1).Subrouter()

	// Authenticated router for v1
	rv1Auth := rv1.PathPrefix("/").Subrouter()
	rv1Auth.Use(middlewares.OAuth2Middleware)

	// Authenticated routes
	rv1Auth.Handle("/users/accept-policy", userHandler.AcceptPolicy()).Methods(http.MethodPut)
	rv1Auth.Handle("/users/change-password", userHandler.ChangePassword()).Methods(http.MethodPut)
	rv1Auth.Handle("/users/reset-password-request", userHandler.ResetPasswordRequest()).Methods(http.MethodPut)

	// Unauthenticated routes
	rv1.Handle("/users", userHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/users", userHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/users", userHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/users/reset-password", userHandler.ResetPassword()).Methods(http.MethodPut)
	rv1.Handle("/users/{id}", userHandler.Update()).Methods(http.MethodPut)

	// Return handler
	return userHandler, nil
}

/*
* @apiTag: user
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: users
* @apiResponseRef: UsersQueryResponse
* @apiSummary: Query users
* @apiParametersRef: UsersQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: UsersQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: UsersQueryNotFoundResponse
 */
func (h *UserHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.UsersQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		users, err := h.UserService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(users)
	})
}

/*
 * @apiTag: user
 * @apiPath: /users
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: UsersCreateRequestBody
 * @apiResponseRef: UsersCreateResponse
 * @apiSummary: Create user
 * @apiDescription: Create user
 */
func (h *UserHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.UsersCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Sanitize phone and work phone number
		payload.Phone = strings.ReplaceAll(payload.Phone, " ", "")
		payload.WorkPhoneNumber = strings.ReplaceAll(payload.WorkPhoneNumber, " ", "")

		log.Printf("%#v payload\n", payload)

		// Check if language skills exists
		languageSkillIds := make([]int64, 0)
		for _, languageSkill := range payload.LanguageSkills {
			langSkillId := int64(languageSkill.ID)
			languageSkillIds = append(languageSkillIds, langSkillId)
		}
		langSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(languageSkillIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		if len(langSkills) != len(languageSkillIds) {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		payload.LanguageSkillIds = languageSkillIds

		// Create user
		user, err := h.UserService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Create staff
		var staffId *uint
		if payload.UserType == constants.USER_TYPE_STAFF {
			staff, err := h.StaffService.CreateEmptyStaffForUser(int(user.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			staffId = &staff.ID
			user.StaffID = staffId
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: user
 * @apiPath: /users/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: UsersUpdateRequestParams
 * @apiRequestRef: UsersUpdateRequestBody
 * @apiResponseRef: UsersCreateResponse
 * @apiSummary: Update user
 * @apiDescription: Update user
 */
func (h *UserHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.UsersUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.UsersUpdateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if language skills exists
		languageSkillIds := make([]int64, 0)
		for _, languageSkill := range payload.LanguageSkills {
			langSkillId := int64(languageSkill.ID)
			languageSkillIds = append(languageSkillIds, langSkillId)
		}
		langSkills, err := h.LanguageSkillService.GetLanguageSkillsByIds(languageSkillIds)
		if err != nil {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		if len(langSkills) != len(languageSkillIds) {
			return response.ErrorBadRequest(nil, "language skills are invalid")
		}
		payload.LanguageSkillIds = languageSkillIds

		// Update user
		payload.Password = security.HashPassword(payload.Password)
		data, err := h.UserService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: user
 * @apiPath: /users
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: UsersDeleteRequestBody
 * @apiResponseRef: UsersCreateResponse
 * @apiSummary: Delete user
 * @apiDescription: Delete user
 */
func (h *UserHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.UsersDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.UserService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: user
 * @apiPath: /users/accept-policy
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiResponseRef: UsersAcceptPolicyResponse
 * @apiSummary: Accept policy
 * @apiDescription: Accept policy
 * @apiSecurity: apiKeySecurity
 */
func (h *UserHandler) AcceptPolicy() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorBadRequest(nil, "token is invalid")
		}

		// Update user
		_, err := h.UserService.UpdateAcceptPolicy(true, int(user.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Return response
		return response.Success(&models.UsersAcceptPolicyResponse{
			Accepted: true,
		})
	})
}

/*
 * @apiTag: user
 * @apiPath: /users/change-password
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiRequestRef: UsersChangePasswordRequestBody
 * @apiResponseRef: UsersChangePasswordResponse
 * @apiSummary: Change Password
 * @apiDescription: Change Password
 * @apiSecurity: apiKeySecurity
 */
func (h *UserHandler) ChangePassword() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.UsersChangePasswordRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorBadRequest(nil, "token is invalid")
		}

		// Check if password is correct
		if err := user.ValidatePassword(payload.CurrentPassword); err != nil {
			return response.ErrorBadRequest(nil, "current password is incorrect")
		}

		// Update password
		password := security.HashPassword(payload.NewPassword)
		updatedUser, err := h.UserService.UpdatePassword(password, int(user.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(&models.UsersChangePasswordResponse{
			User: *updatedUser,
		})
	})
}

/*
 * @apiTag: user
 * @apiPath: /users/reset-password
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiRequestRef: UsersResetPasswordRequestBody
 * @apiResponseRef: UsersResetPasswordResponse
 * @apiSummary: Reset Password
 * @apiDescription: Reset Password
 */
func (h *UserHandler) ResetPassword() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.UsersResetPasswordRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		user, err := h.UserService.FindByEmail(payload.Email)
		if user == nil || err != nil {
			return response.ErrorBadRequest(nil, "user not found")
		}

		// Verify exchange code
		err = h.OTPService.VerifyExchangeCode(payload.ExchangeCode, otp.ReasonResetPassword, int64(user.ID), *user.WorkPhoneNumber, user.Email)
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}

		// Update password
		password := security.HashPassword(payload.NewPassword)
		updatedUser, err := h.UserService.UpdatePassword(password, int(user.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(&models.UsersResetPasswordResponse{
			User: *updatedUser,
		})
	})
}

/*
 * @apiTag: user
 * @apiPath: /users/reset-password-request
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiRequestRef: UsersResetPasswordRequestRequestBody
 * @apiResponseRef: UsersResetPasswordRequestResponse
 * @apiSummary: Reset Password
 * @apiDescription: Reset Password
 * @apiSecurity: apiKeySecurity
 */
func (h *UserHandler) ResetPasswordRequest() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.UsersResetPasswordRequestRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get user from context
		user, err := h.UserService.FindByID(payload.UserID)
		if user == nil || err != nil {
			return response.ErrorBadRequest(nil, "user not found")
		}

		return response.Success(&models.UsersResetPasswordRequestResponseData{
			Message: fmt.Sprintf("New password has been sent to email %s. Please check email.", user.Email),
		})
	})
}
