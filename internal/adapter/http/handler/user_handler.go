package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	entity "main.go/internal/domain/entity"
	user "main.go/internal/usecase/users"
)

type UserHandler struct {
	createUserUC     *user.CreateUserUseCase
	getUserByEmailUC *user.GetUserByEmailUseCase
	listAllUsersUC   *user.ListAllUsersUseCase
	updateUserUC     *user.UpdateUserUseCase
	deleteUserUC     *user.DeleteUserUseCase
}

type CreateUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	ProfilePicture string `json:"profile_picture"`
}

type UpdateUserRequest struct {
	ID             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
}

func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	getUserByEmailUC *user.GetUserByEmailUseCase,
	listAllUsersUC *user.ListAllUsersUseCase,
	updateUserUC *user.UpdateUserUseCase,
	deleteUserUC *user.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC:     createUserUC,
		getUserByEmailUC: getUserByEmailUC,
		listAllUsersUC:   listAllUsersUC,
		updateUserUC:     updateUserUC,
		deleteUserUC:     deleteUserUC,
	}
}

// POST /v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		ProfilePicture: req.ProfilePicture,
	}

	if err := h.createUserUC.Execute(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully!",
		"data":    user,
	})
}

// GET /v1/users?exclude_email=test@example.com
func (h *UserHandler) ListAllUsers(c *gin.Context) {
	excludeEmail := c.Query("exclude_email")
	if excludeEmail == "" {
		excludeEmail = ""
	}

	users, err := h.listAllUsersUC.Execute(c.Request.Context(), excludeEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"count": len(users),
	})
}

// GET /v1/users/:email
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := h.getUserByEmailUC.Execute(c.Request.Context(), email)
	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// PATCH /v1/users
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		ID:             req.ID,
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		ProfilePicture: req.ProfilePicture,
	}

	if err := h.updateUserUC.Execute(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully!",
		"data":    user,
	})
}

// DELETE /v1/users
func (h *UserHandler) DeleteUser(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email param is required"})
		return
	}

	deleted_user, err := h.deleteUserUC.Execute(c.Request.Context(), email)

	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "User deleted successfully!",
		"data":    deleted_user,
	})
}
