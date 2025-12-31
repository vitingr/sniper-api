package handler

import (
	"net/http"
	"os"

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
	meUC              *auth.MeUseCase
}

func NewAuthHandler(
	loginUserUC *auth.LoginUserUseCase,
	logoutUserUC *auth.LogoutUserUseCase,
	ssoRefreshTokenUC *auth.RefreshTokenUseCase,
	tokenMakerUC *auth.TokenMakerUseCase,
	meUC *auth.MeUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUserUC:       loginUserUC,
		logoutUserUC:      logoutUserUC,
		ssoRefreshTokenUC: ssoRefreshTokenUC,
		tokenMakerUC:      tokenMakerUC,
		meUC:              meUC,
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

	isProduction := os.Getenv("APP_ENV") == "production"

	if isProduction {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie(
		"access_token",
		response.AccessToken,
		3600,
		"/",
		"",
		isProduction,
		true,
	)

	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		7*24*3600,
		"/",
		"",
		isProduction,
		true,
	)

	c.SetCookie(
		"session_id",
		response.SessionID,
		7*24*3600,
		"/",
		"",
		isProduction,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged successfully!",
		"data": gin.H{
			"user": response.User,
		},
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

	isProduction := os.Getenv("APP_ENV") == "production"

	if isProduction {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
	c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
	c.SetCookie("session_id", "", -1, "/", "", isProduction, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logout!",
	})
}

// POST /v1/auth/sso/refresh-token
func (h *AuthHandler) SSORefreshToken(c *gin.Context) {
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	refreshed_session, err := h.ssoRefreshTokenUC.Execute(c.Request.Context(), refresh_token)

	if err != nil {
		logger.Error("Error at refresh access token.")
		return
	}

	isProduction := os.Getenv("APP_ENV") == "production"

	if isProduction {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie(
		"access_token",
		refreshed_session.AccessToken,
		3600,
		"/",
		"",
		isProduction,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Access token refreshed successfully",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No active session."})
		return
	}

	user, err := h.meUC.Execute(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
