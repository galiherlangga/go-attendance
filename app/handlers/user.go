package handlers

import (
	"net/http"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns access/refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login credentials"
// @Success      200   {object}  map[string]string     "Tokens returned"
// @Failure      400   {object}  map[string]string     "Invalid input"
// @Failure      401   {object}  map[string]string     "Unauthorized"
// @Router       /auth/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var input models.LoginRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(&input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set access token in HttpOnly cookie
	ctx.SetCookie("access_token", accessToken, 3600, "/", "", false, true) // HttpOnly=true

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}