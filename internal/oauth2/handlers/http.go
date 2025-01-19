package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	jwtlib "github.com/hoitek/Maja-Service/internal/_shared/jwt-lib"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/oauth2"
	"github.com/hoitek/Maja-Service/internal/_shared/otp"
	"github.com/hoitek/Maja-Service/internal/oauth2/config"
	"github.com/hoitek/Maja-Service/internal/oauth2/models"
	oPorts "github.com/hoitek/Maja-Service/internal/otp/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type OAuth2Handler struct {
	OTPService  oPorts.OTPService
	UserService uPorts.UserService
}

func NewOAuth2Handler(r *mux.Router, o oPorts.OTPService, u uPorts.UserService) (OAuth2Handler, error) {
	oAuth2Handler := OAuth2Handler{
		OTPService:  o,
		UserService: u,
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.OAuth2Config.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.OAuth2Config.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	// Routes without auth
	rv1.Handle("/oauth2/auth", oAuth2Handler.Auth()).Methods(http.MethodPost)

	// Routes with auth
	rAuth.Handle("/oauth2/refresh", oAuth2Handler.Refresh()).Methods(http.MethodPost)
	rAuth.Handle("/oauth2/userinfo", oAuth2Handler.UserInfo()).Methods(http.MethodGet)

	return oAuth2Handler, nil
}

/*
 * @apiTag: oauth2
 * @apiPath: /oauth2/auth
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: OAuth2AuthRequestBody
 * @apiResponseRef: OAuth2AuthResponse
 * @apiSummary: Authenticate
 * @apiDescription: Authenticate
 */
func (h *OAuth2Handler) Auth() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.OAuth2AuthRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}
		var userInterface interface{}
		if payload.Type == oauth2.OAuth2TypeEmail {
			u, err := h.UserService.FindByEmail(payload.Username)
			if err != nil {
				return response.ErrorBadRequest(err, "User not found")
			}
			userInterface = u
		} else if payload.Type == oauth2.OAuth2TypeUsername {
			u, err := h.UserService.FindByUsername(payload.Username)
			if err != nil {
				return response.ErrorBadRequest(err, "User not found")
			}
			userInterface = u
		}
		if userInterface == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}
		user := h.UserService.AssertToUserDomain(userInterface)
		if user == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}

		// Validate password
		if err := user.ValidatePassword(payload.Password); err != nil {
			return response.ErrorBadRequest(nil, "User not found")
		}

		// Validate exchange code if OTP is enabled
		if config.OAuth2Config.OTPEnable {
			otps, err := h.OTPService.Query(fmt.Sprintf("userId = '%d' AND exchangeCode = '%s' AND isVerified = '%t' AND isUsed = '%t'", user.ID, payload.ExchangeCode, true, true))
			if err != nil {
				return response.ErrorBadRequest(err, "Exchange code is not valid")
			}
			if len(otps) == 0 {
				return response.ErrorBadRequest(err, "Exchange code is not valid")
			}
			otpIdentity := otps[0]
			if otpIdentity.Type == otp.TypeEmail {
				if otpIdentity.Username != user.Email || otpIdentity.Password != payload.Password {
					return response.ErrorBadRequest(err, "Exchange code is not valid")
				}
			} else if otpIdentity.Type == otp.TypePhone {
				if otpIdentity.Username != user.Phone || otpIdentity.Password != payload.Password {
					return response.ErrorBadRequest(err, "Exchange code is not valid")
				}
			} else {
				return response.ErrorBadRequest(err, "Exchange code is not valid")
			}
			// TODO: Check expired time of otp
		}

		// Generate the access token and refresh token
		accessToken, refreshToken, err := jwtlib.GenerateTokenPair(jwtlib.JwtPayload{
			ID:       int64(user.ID),
			Username: user.Username,
			Data: map[string]interface{}{
				"info": user,
			},
		})
		if err != nil {
			return response.ErrorBadRequest(err, "Failed to generate token")
		}

		return response.Success(models.OAuth2AuthResponse{
			OAuth2: models.OAuth2Response{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				User:         *user,
			},
		})
	})
}

/*
 * @apiTag: oauth2
 * @apiPath: /oauth2/refresh
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: OAuth2RefreshRequestBody
 * @apiResponseRef: OAuth2RefreshResponse
 * @apiSummary: Refresh Token
 * @apiDescription: Refresh Token
 * @apiSecurity: apiKeySecurity
 */
func (h *OAuth2Handler) Refresh() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.OAuth2RefreshRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		refreshToken := payload.RefreshToken
		parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.OAuth2Config.JwtSigningKey), nil
		})
		if err != nil {
			return response.ErrorBadRequest(err, "Refresh token is not valid1")
		}
		if parsedToken == nil || !parsedToken.Valid {
			return response.ErrorBadRequest(nil, "Refresh token is not valid2")
		}
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return response.ErrorBadRequest(nil, "Refresh token is not valid3")
		}
		payloadData, ok := claims["sub"].(map[string]interface{})
		if !ok {
			return response.ErrorBadRequest(nil, "Refresh token is not valid4")
		}
		accessToken, ok := payloadData["accessToken"]
		if !ok {
			return response.ErrorBadRequest(nil, "Refresh token is not valid5")
		}
		info, ok := payloadData["info"]
		if !ok {
			return response.ErrorBadRequest(nil, "Refresh token is not valid6")
		}
		userInfo, ok := info.(map[string]interface{})
		if !ok {
			return response.ErrorBadRequest(nil, "Refresh token is not valid7")
		}

		// Get token from context
		token := h.UserService.GetTokenFromContext(r.Context())
		if token == "" || token != accessToken {
			return response.ErrorBadRequest(nil, "Refresh token is not valid9")
		}

		userId := int64(userInfo["id"].(float64))

		// Delete accessToken from payloadData if exists
		delete(payloadData, "accessToken")

		// Generate the access token and refresh token
		newAccessToken, newRefreshToken, err := jwtlib.GenerateTokenPair(jwtlib.JwtPayload{
			ID:       userId,
			Username: fmt.Sprintf("%s", userInfo["username"]),
			Data:     payloadData,
		})
		if err != nil {
			return response.ErrorBadRequest(err, "Failed to generate token")
		}

		return response.Success(&models.OAuth2RefreshResponse{
			OAuth2: models.OAuth2Response{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			},
		})
	})
}

/*
* @apiTag: oauth2
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /oauth2/userinfo
* @apiResponseRef: OAuth2UserInfoResponse
* @apiSummary: Query users
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: OAuth2UserInfoNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: OAuth2UserInfoNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *OAuth2Handler) UserInfo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorBadRequest(nil, "User not found")
		}
		return response.Success(&models.OAuth2UserInfoResponse{
			User: *user,
		})
	})
}
