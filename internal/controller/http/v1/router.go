// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "bhs/docs"
	"bhs/internal/usecase"
	"bhs/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go BHS test project
// @description Go BHS test project
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(
	handler *gin.Engine,
	l logger.Interface,
	a usecase.Auth,
	as usecase.Assets) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, a, l)
		newAssetsRoutes(h, as, a, l)
	}
}
