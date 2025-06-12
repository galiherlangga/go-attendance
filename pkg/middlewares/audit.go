package middleware

import (
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := utils.GetUserFromContext(ctx)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		requestID := uuid.New().String()
		// Store in context
		ctx.Set("user_id", userID)
		ctx.Set("request_id", requestID)
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Next()
	}
}
