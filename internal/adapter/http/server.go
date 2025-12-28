package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewServer(handler *gin.Engine, port string) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
