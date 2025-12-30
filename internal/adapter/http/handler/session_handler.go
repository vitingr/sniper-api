package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/domain/entity"
	"main.go/internal/usecase/session"
	"main.go/pkg/logger"
)

type SessionHandler struct {
	createSessionUC *session.CreateSessionUseCase
	getSessionUC    *session.GetSessionUseCase
	deleteSessionUC *session.DeleteSessionUseCase
	revokeSessionUC *session.RevokeSessionUseCase
}

type CreateSessionRequest struct {
	Email string `json:"email"`
}

func NewSessionHandler(
	createSessionUC *session.CreateSessionUseCase,
	getSessionUC *session.GetSessionUseCase,
	deleteSessionUC *session.DeleteSessionUseCase,
	revokeSessionUC *session.RevokeSessionUseCase,
) *SessionHandler {
	return &SessionHandler{
		createSessionUC: createSessionUC,
		getSessionUC:    getSessionUC,
		deleteSessionUC: deleteSessionUC,
		revokeSessionUC: revokeSessionUC,
	}
}

// POST /v1/sessions
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := &entity.Session{
		UserEmail: req.Email,
	}

	session, err := h.createSessionUC.Execute(c.Request.Context(), session)

	if err != nil {
		logger.Error("Cannot create a new session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session created successfully!",
		"data":    session,
	})
}
// GET /v1/sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	session_id := c.Param("id")

	if session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
	}

	session, err := h.getSessionUC.Execute(c.Request.Context(), session_id)
	if err != nil {
		logger.Error("Session not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session found successfully!",
		"data":    session,
	})
}

// POST /v1/sessions/revoke/:id
func (h *SessionHandler) RevokeSession(c *gin.Context) {
	session_id := c.Param("id")

	if session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
	}

	err := h.revokeSessionUC.Execute(c.Request.Context(), session_id)

	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Session has been revoked successfully!",
	})
}

// DELETE /v1/sessions/:id
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	session_id := c.Param("id")

	if session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
	}

	err := h.deleteSessionUC.Execute(c.Request.Context(), session_id)

	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Session deleted successfully!",
	})
}
