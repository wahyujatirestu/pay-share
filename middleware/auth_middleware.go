package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/utils/service"
)


type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JWTService
}

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func NewAuthMiddleware(jwtService service.JWTService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		var authHeader authHeader
		if err := ctx.ShouldBindHeader(&authHeader); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return 
		}

		token := strings.Replace(authHeader.Authorization, "Bearer ", "", 1)
		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return 
		}

		userId, err := uuid.Parse(tokenClaim.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return 
		}

		ctx.Set("user", model.User{
			ID: userId,
			Role: tokenClaim.Role,
		})
		validRole := false

		for _, role := range roles{
			if role == tokenClaim.Role{
				validRole = true
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		}

		ctx.Next()
	}
}