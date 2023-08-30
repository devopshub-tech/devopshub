// internal/infrastructure/http/router.go
package http

import (
	"net/http"

	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HealthHandlers []func() bool

type GinRouter struct {
	router         *gin.Engine
	healthHandlers HealthHandlers
}

// NewGinRouter creates a new instance of GinRouter.
func NewGinRouter(mode string) *GinRouter {
	router := gin.New()

	switch mode {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	logger := logging.NewLogger()

	router.Use(logger.GinLogger())
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	return &GinRouter{
		router: router,
	}
}

// Engine returns the underlying gin.Engine instance.
func (r *GinRouter) Engine() *gin.Engine {
	return r.router
}

// AddReadinessCheck adds a readiness check handler.
func (r *GinRouter) AddReadinessCheck(handler func() bool) {
	r.healthHandlers = append(r.healthHandlers, handler)
}

// SetupHealthRoutes sets up health check routes.
func (r *GinRouter) SetupHealthRoutes() {
	r.router.GET("/health/liveness", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	r.router.GET("/health/readiness", func(ctx *gin.Context) {
		for _, handler := range r.healthHandlers {
			if !handler() {
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "NOT_READY"})
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "READY"})
	})
}
