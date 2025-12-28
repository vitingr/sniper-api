package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(startTime)

		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if statusCode >= 500 {
			log.Error().
				Int("status", statusCode).
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Dur("latency", latency).
				Str("error", errorMessage).
				Msg("Server error")
		} else if statusCode >= 400 {
			log.Warn().
				Int("status", statusCode).
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Dur("latency", latency).
				Msg("Client error")
		} else {
			log.Info().
				Int("status", statusCode).
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Dur("latency", latency).
				Msg("Request completed")
		}
	}
}

func CustomLogger(logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 500 {
			logger.Error("request completed",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency,
			)
		} else if statusCode >= 400 {
			logger.Warn("request completed",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency,
			)
		} else {
			logger.Info("request completed",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency,
			)
		}
	}
}
