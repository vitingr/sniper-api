package http

import (
	"github.com/gin-gonic/gin"
	"main.go/internal/adapter/http/middleware"
	"main.go/internal/di"
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

	v1 := r.engine.Group("/v1")
	{
		companies := v1.Group("/companies")
		{
			companies.GET("", r.container.CompanyHandler.GetCompanies)
			companies.GET("/:id", r.container.CompanyHandler.GetCompany)
		}
		users := v1.Group("/users")
		{
			users.GET("", r.container.UserHandler.ListAllUsers)
			users.GET("/:email", r.container.UserHandler.GetUserByEmail)
			users.POST("", r.container.UserHandler.CreateUser)
			users.PATCH("", r.container.UserHandler.UpdateUser)
			users.DELETE("/:email", r.container.UserHandler.DeleteUser)
		}
		auth := v1.Group("/auth")
		{
			auth.GET("/me", r.container.AuthHandler.Me)
			auth.POST("/sso/login", r.container.AuthHandler.LoginUser)
			auth.POST("/sso/logout", r.container.AuthHandler.LogoutUser)
			auth.POST("/sso/refresh-token", r.container.AuthHandler.SSORefreshToken)
		}
		sessions := v1.Group("/sessions")
		{
			sessions.GET("/:id", r.container.SessionHandler.GetSession)
			sessions.POST("", r.container.SessionHandler.CreateSession)
			sessions.DELETE("/:id", r.container.SessionHandler.DeleteSession)
			sessions.POST("/revoke/:id", r.container.SessionHandler.RevokeSession)
		}
	}

	return r.engine
}
