package route

import "github.com/hoitek/Kit/response"

type SecureRoute struct {
	Path        string
	Method      string
	Handler     response.Handler
	Permissions []string
}
