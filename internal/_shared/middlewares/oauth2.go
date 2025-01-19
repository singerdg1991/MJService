package middlewares

import (
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

// GetUserService returns a new instance of the user service
func GetUserService(pDB *sql.DB, mDB *mongo.Client) uPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[uPorts.UserRepositoryPostgresDB](
		config.AppConfig.Environment == constants.ENVIRONMENT_TESTING,
		uRepositories.NewUserRepositoryStub(),
		uRepositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := uRepositories.NewUserRepositoryMongoDB(mDB)

	// user service and inject the user repository database and grpc
	userService := uService.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, storage.MinIOStorage)
	return userService
}

func OAuth2Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithUser, err := func(r *http.Request) (*context.Context, response.Response) {
			userService := GetUserService(database.PostgresDB, database.MongoDB)

			// Get authorization from header
			authorization := r.Header.Get("authorization")
			if authorization == "" {
				return nil, response.ErrorUnAuthorized("")
			}

			// Extract token from bearer token
			tokenSlice := strings.Split(authorization, " ")
			if len(tokenSlice) < 2 {
				return nil, response.ErrorUnAuthorized("")
			}
			if strings.ToLower(tokenSlice[0]) != "bearer" {
				return nil, response.ErrorUnAuthorized("Your token is not a bearer token")
			}
			token := tokenSlice[1]

			// Parse token via jwt
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JwtSigningKey), nil
			})
			if err != nil {
				return nil, response.ErrorUnAuthorized("You're token is not valid!")
			}

			// Check token validity
			if parsedToken == nil || !parsedToken.Valid {
				return nil, response.ErrorUnAuthorized("You're token is not valid!")
			}

			// Get claims from parsed token
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				return nil, response.ErrorUnAuthorized("You're token is not valid!")
			}

			// Get user id from claims
			id, ok := claims["jti"].(float64)
			if !ok {
				return nil, response.ErrorUnAuthorized("You're token is not valid!")
			}
			// Find user by id
			var userId = int(id)
			res, err := userService.Query(&models.UsersQueryRequestParams{
				ID: userId,
			})
			if err != nil {
				return nil, response.ErrorUnAuthorized("Not Authenticated")
			}
			items := res.Items
			if items == nil {
				return nil, response.ErrorUnAuthorized("Not Authenticated")
			}
			users, ok := items.([]*domain.User)
			if !ok {
				return nil, response.ErrorUnAuthorized("Not Authenticated")
			}
			if len(users) == 0 {
				return nil, response.ErrorUnAuthorized("Not Authenticated")
			}

			// Bind user to context
			ctx := userService.BindUserToContext(r.Context(), users[0])

			// Bind token to context
			ctx = userService.BindTokenToContext(ctx, token)

			return &ctx, nil
		}(r)

		if err != nil {
			errorResponse := err.(response.ErrorResponse)
			if err := response.ErrorWithWriter(w, errorResponse, errorResponse.StatusCode); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			next.ServeHTTP(w, r.WithContext(*ctxWithUser))
		}
	})
}
