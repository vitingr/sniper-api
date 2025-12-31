package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		isProduction := os.Getenv("APP_ENV") == "production"

		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			"http://localhost:5174": true,
			"http://localhost:3000": true,
		}

		if isProduction {
			// TODO: add the real domain here
			allowedOrigins = map[string]bool{
				"https://sniper-api.com":     true,
			}
		}
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-platform")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
