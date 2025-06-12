package middleware

import (
	"fmt"

	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil {
			// If cookie not found, check Authorization header
			authHeader := ctx.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				ctx.JSON(401, gin.H{"error": "Unauthorized – token missing"})
				ctx.Abort()
				return
			}
		}

		// Parse and validate the JWT token here
		userID, err := utils.ParseJWT(token)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// Store user ID in context for later use
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

func IsAdminMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil {
			// If cookie not found, check Authorization header
			authHeader := ctx.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				ctx.JSON(401, gin.H{"error": "Unauthorized – token missing"})
				ctx.Abort()
				return
			}
		}

		// Parse and validate the JWT token here
		userID, err := utils.ParseJWT(token)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// Fetch user from DB
		user, err := userRepo.FindByID(userID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch user"})
			ctx.Abort()
			return
		}

		// Check role
		if user.Role.Name != "admin" {
			ctx.JSON(403, gin.H{"error": "Forbidden – Admins only"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
