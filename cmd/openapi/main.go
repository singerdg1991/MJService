package main

import (
	"fmt"
	openengine "github.com/hoitek/OpenEngine"
	"log"
	"os"

	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/openapi"
	"github.com/hoitek/OpenEngine/engine"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	exportPath := GetExportPathFromArgs()
	internalFolderPath := GetInternalDirPathFromArgs()

	config.LoadDefault()

	errResponses := openapi.DataAndMessageSwaggerErrorResponses

	majaUrl := fmt.Sprintf("%s://%s:%d%s%s",
		config.AppConfig.Protocol, config.AppConfig.HostUri, config.AppConfig.Port, config.AppConfig.ApiPrefix, config.AppConfig.ApiVersion1)

	if config.AppConfig.Environment == "production" {
		majaUrl = fmt.Sprintf("https://%s%s%s", config.AppConfig.HostUri, config.AppConfig.ApiPrefix, config.AppConfig.ApiVersion1)
	}

	majaServer := engine.ApiServer{Url: majaUrl}

	swaggerUiConfig := engine.SwaggerUiConfig{
		ExportPath: exportPath,
		ServeURI:   "/apidocs",
		Title:      "Maja API",
	}

	_, err = openengine.NewPackage().
		AddServers(engine.ApiServers{majaServer}).
		AddIgnoredPaths(constants.OpenApiIgnored).
		AddIgnoredPaths([]string{exportPath}).
		ParseTags(internalFolderPath, []string{"static", "healthcheck", "protobuf", "welcome", "_shared", "s3", "_mock"}).
		ParseSchemas(rootPath).
		ParseEnums(rootPath).
		AddDefaultErrors(404, 400, 401, 500).
		AddErrorResponses(errResponses).
		ParsePaths(rootPath).
		AddSecuritySchemes(openapi.SecuritySchemes).
		ExportSwaggerUi(swaggerUiConfig).
		Generate(exportPath)

	if err != nil {
		log.Print("\n\n --> Generate Yaml Errors: ", err, "\n\n")
	}
}
