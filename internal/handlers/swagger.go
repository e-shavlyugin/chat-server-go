package handlers

import (
	"bytes"
	"chat_server_v2/middleware"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerHandler(engine *gin.RouterGroup) {
	// Serve the OpenAPI v3 YAML file
	engine.GET("/openapi.yaml", renderSwaggerYaml())

	// Serve the Swagger UI
	url := ginSwagger.URL("/openapi.yaml") // The URL pointing to your OpenAPI file
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

type SwaggerConfig struct {
	SwagHost   string
	SwagScheme string
}

func renderSwaggerYaml() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := SwaggerConfig{SwagHost: ctx.Request.Host}

		tmpl, err := template.ParseFiles("./docs/v1/openapi/openapi.yaml")
		if err != nil {
			middleware.AbortWithErrorObject(ctx, http.StatusInternalServerError, err)
			return
		}

		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, config)
		if err != nil {
			middleware.AbortWithErrorObject(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.Data(http.StatusOK, "text/plain; charset=utf-8", renderedTemplate.Bytes())
	}
}
