package middlewares

import (
	"net/http"

	"github.com/hoitek/Kit/response"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

func AuthMiddleware(userService uPorts.UserService, permissions []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := userService.GetUserFromContext(r.Context())
			if user == nil {
				errMsg := "Not Authenticated"
				errResp := response.ErrorUnAuthorized(errMsg)
				response.ErrorWithWriter(w, errResp, http.StatusForbidden)
				return
			}

			if len(permissions) > 0 {
				for _, permission := range permissions {
					if user.HasPermission(permission) {
						next.ServeHTTP(w, r)
						return
					}
				}
				errMsg := "You don't have permission to access this resource"
				errResp := response.ErrorForbidden(errMsg)
				response.ErrorWithWriter(w, errResp, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
