package constants

// Environment is the environment in which the application is running
const (
	ENVIRONMENT_DEVELOPMENT = "development"
	ENVIRONMENT_PRODUCTION  = "production"
	ENVIRONMENT_TESTING     = "testing"
)

// OpenApiIgnored is a list of directories that are ignored when generating the OpenAPI documentation
var OpenApiIgnored = []string{
	"accessor",
	"database",
	"cmd",
	"config",
	"constants",
	"middlewares",
	"router",
	"utils",
	"scripts",
	"apigen",
	"apidocs",
}
