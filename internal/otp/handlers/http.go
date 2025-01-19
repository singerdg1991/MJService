package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hoitek/Kit/response"
	jwtlib "github.com/hoitek/Maja-Service/internal/_shared/jwt-lib"
	"github.com/hoitek/Maja-Service/internal/_shared/otp"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/otp/config"
	"github.com/hoitek/Maja-Service/internal/otp/domain"
	"github.com/hoitek/Maja-Service/internal/otp/models"
	oPorts "github.com/hoitek/Maja-Service/internal/otp/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"github.com/gorilla/mux"
)

type OTPHandler struct {
	OTPService  oPorts.OTPService
	UserService uPorts.UserService
}

func NewOTPHandler(r *mux.Router, o oPorts.OTPService, u uPorts.UserService) (OTPHandler, error) {
	otpHandler := OTPHandler{
		OTPService:  o,
		UserService: u,
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.OTPConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.OTPConfig.ApiVersion1).Subrouter()

	rv1.Handle("/otp/send", otpHandler.Send()).Methods(http.MethodPost)
	rv1.Handle("/otp/verify", otpHandler.Verify()).Methods(http.MethodPost)
	rv1.Handle("/otp/reasons", otpHandler.Reasons()).Methods(http.MethodGet)
	rv1.Handle("/otp/cooldown/remain", otpHandler.CoolDownRemain()).Methods(http.MethodGet)
	rv1.Handle("/otp/settings", otpHandler.Settings()).Methods(http.MethodGet)

	return otpHandler, nil
}

/*
 * @apiTag: otp
 * @apiPath: /otp/send
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: OTPSendRequest
 * @apiResponseRef: OTPSendResponse
 * @apiSummary: OTP Send
 * @apiDescription: OTP Send
 */
func (h *OTPHandler) Send() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.OTPSendRequest{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		var userInterface interface{}
		if payload.Type == otp.TypePhone {
			userInterface, _ = h.UserService.FindByWorkPhoneNumber(payload.Username)
		} else if payload.Type == otp.TypeEmail {
			userInterface, _ = h.UserService.FindByEmail(payload.Username)
		}
		if userInterface == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}
		user := h.UserService.AssertToUserDomain(userInterface)
		if user == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}

		log.Println("user", user)

		if payload.Reason == otp.ReasonLogin {
			// Validate password
			if err := user.ValidatePassword(payload.Password); err != nil {
				return response.ErrorBadRequest(nil, err.Error())
			}
		}

		// Get Ip and UserAgent
		ip, err := utils.GetIP(r)
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}
		userAgent := r.UserAgent()

		// Generate exchange code
		exchangeCode, err := uuid.NewRandom()
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}

		// Calculate cool down
		var coolDownValues = []int{60, 120, 300, 3600}
		var coolDown = coolDownValues[0]
		userOtps, err := h.OTPService.Query(fmt.Sprintf("userId = '%d' AND isUsed = false AND isVerified = false", user.ID))
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}
		count := len(userOtps)
		if count > 0 {
			// Get last otp
			lastOtp := userOtps[count-1]

			// Calculate cool down
			expireTime := int(time.Until(lastOtp.ExpiredAt).Seconds())
			if expireTime > 0 {
				if count < len(coolDownValues)-1 {
					coolDown = coolDownValues[count]
				} else {
					coolDown = coolDownValues[len(coolDownValues)-1]
				}
			} else {
				coolDown = coolDownValues[0]
			}
		}

		// Define variables
		var (
			currentTime = time.Now()
			code        = h.OTPService.GenerateOTPCode(config.OTPConfig.OTPCodeLength)
			otpType     = payload.Type
			expiredAt   = time.Now().Add(time.Second * time.Duration(config.OTPConfig.OTPCodeExpirationSeconds))
		)

		// Send email via tikka
		// _, err = tikka.Default.SendEmail(&protobuf.EmailRequest{
		// 	EmailEntry: &protobuf.Email{
		// 		Recipient: user.Email,
		// 		Subject:   "Login Code",
		// 		Body:      fmt.Sprintf("Hello dear %s,\n Use this code to login to Hoivalani application: %s", user.FirstName, code),
		// 	},
		// })
		// if err != nil {
		// 	log.Printf("Error while sending email via tika in otp: %s\n", err.Error())
		// 	return response.ErrorInternalServerError(nil, "Something went wrong with sending email. Please contact administrator and try again.")
		// }

		// Create OTP object
		otp := &domain.OTP{
			UserID:       int64(user.ID),
			Username:     payload.Username,
			Password:     payload.Password,
			Code:         code,
			ExchangeCode: exchangeCode.String(),
			Ip:           ip,
			UserAgent:    userAgent,
			Type:         otpType,
			IsUsed:       false,
			IsVerified:   false,
			Reason:       payload.Reason,
			ExpiredAt:    expiredAt,
			CreatedAt:    currentTime,
			UpdatedAt:    currentTime,
		}

		// Send OTP
		err = h.OTPService.Send(otp)
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}

		// Generate token
		token, err := jwtlib.Encrypt(jwtlib.JwtPayload{
			ID:       otp.ID,
			Username: user.Email,
			Data:     nil,
		}, 0)
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}

		log.Println("code: ", code)

		// Return response
		return response.Success(&models.OTPSendResponse{
			Token:    token,
			CoolDown: coolDown,
			Code:     code,
		})
	})
}

/*
 * @apiTag: otp
 * @apiPath: /otp/verify
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: OTPVerifyRequest
 * @apiResponseRef: OTPVerifyResponse
 * @apiSummary: OTP Verify
 * @apiDescription: OTP Verify
 */
func (h *OTPHandler) Verify() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.OTPVerifyRequest{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Verify OTP
		otp, err := h.OTPService.Verify(payload)
		if err != nil {
			return response.ErrorBadRequest(nil, err.Error())
		}

		// Return response
		return response.Success(&models.OTPVerifyResponse{
			ExchangeCode: otp.ExchangeCode,
		})
	})
}

/*
* @apiTag: otp
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /otp/reasons
* @apiResponseRef: OTPReasonsResponse
* @apiSummary: Query reasons for OTP
* @apiErrorStatusCodes: 400, 500, 404
 */
func (h *OTPHandler) Reasons() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Create map for reasons
		res := map[string]string{}

		// Set reasons in map
		res["Login"] = otp.ReasonLogin
		res["ResetPassword"] = otp.ReasonResetPassword
		res["Register"] = otp.ReasonRegister
		res["ChangePassword"] = otp.ReasonChangePassword

		// Return response
		return response.Success(&models.OTPReasonsResponse{
			Reasons: res,
		})
	})
}

/*
* @apiTag: otp
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /otp/settings
* @apiResponseRef: OTPSettingsResponse
* @apiSummary: Query settings for OTP
* @apiErrorStatusCodes: 400, 500, 404
 */
func (h *OTPHandler) Settings() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		var (
			OTPCodeLength = config.OTPConfig.OTPCodeLength
			OTPEnable     = config.OTPConfig.OTPEnable
		)
		return response.Success(&models.OTPSettingsResponse{
			Length: OTPCodeLength,
			Enable: OTPEnable,
		})
	})
}

func (h *OTPHandler) CoolDownRemain() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		panic("impement me")
	})
}
