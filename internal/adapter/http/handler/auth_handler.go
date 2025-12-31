package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/adapter/repository"
	"main.go/internal/usecase/auth"
	"main.go/pkg/logger"
)

type AuthHandler struct {
	loginUserUC       *auth.LoginUserUseCase
	logoutUserUC      *auth.LogoutUserUseCase
	ssoRefreshTokenUC *auth.RefreshTokenUseCase
	tokenMakerUC      *auth.TokenMakerUseCase
}

func NewAuthHandler(
	loginUserUC *auth.LoginUserUseCase,
	logoutUserUC *auth.LogoutUserUseCase,
	ssoRefreshTokenUC *auth.RefreshTokenUseCase,
	tokenMakerUC *auth.TokenMakerUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUserUC:       loginUserUC,
		logoutUserUC:      logoutUserUC,
		ssoRefreshTokenUC: ssoRefreshTokenUC,
		tokenMakerUC:      tokenMakerUC,
	}
}

// POST /v1/auth/sso/login
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req repository.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credentials := &repository.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.loginUserUC.Execute(c.Request.Context(), credentials)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged successfully!",
		"data":    response,
	})
}

// POST /v1/auth/sso/logout
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	session_id := c.Param("session_id")

	if session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id param is required"})
		return
	}

	err := h.logoutUserUC.Execute(c.Request.Context(), session_id)

	if err != nil {
		logger.Error("Cannot login the user.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged successfully",
	})
}

// POST /v1/auth/sso/refresh-token
func (h *AuthHandler) SSORefreshToken(c *gin.Context) {
	var req repository.SSORefreshAccessTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshed_session, err := h.ssoRefreshTokenUC.Execute(c.Request.Context(), req.RefreshToken)

	if err != nil {
		logger.Error("Error at refresh access token.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Access token refreshed successfully",
		"data":    refreshed_session,
	})
}
