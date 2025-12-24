package http

import (
	"github.com/gin-gonic/gin"
	"main.go/internal/di"
	"main.go/internal/infra/http/middleware"
)

type Router struct {
	engine    *gin.Engine
	container *di.Container
}

func NewRouter(container *di.Container) *Router {
	return &Router{
		engine:    gin.New(),
		container: container,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(gin.Recovery())
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())
	r.engine.Use(middleware.RateLimitMiddleware())

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.engine.Group("/api/v1")
	{
		companies := v1.Group("/companies")
		{
			companies.GET("", r.container.CompanyHandler.GetCompanies)
			companies.GET("/:id", r.container.CompanyHandler.GetCompany)
		}
	}

	return r.engine
}
