package middlewares

import (
	"net/http"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	AdminRoleID     = 2
	AssistantRoleID = 3
)

// getUserInfoFromPath はクエリパラメータからユーザー情報を取得します。
func getUserInfoFromPath(ctx *gin.Context) (uid uint, cid uint, err error) {
	uidStr, cidStr := ctx.Query("uid"), ctx.Query("cid")
	uidUint, uidErr := strconv.ParseUint(uidStr, 10, 32)
	cidUint, cidErr := strconv.ParseUint(cidStr, 10, 32)
	if uidErr != nil || cidErr != nil {
		return 0, 0, err
	}

	return uint(uidUint), uint(cidUint), nil
}

// ClassUserRoleMiddleware は指定された権限を持っているかどうかを確認するミドルウェアです。
func ClassUserRoleMiddleware(roleService services.ClassUserService, requiredRoleID int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, cid, err := getUserInfoFromPath(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(constants.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid user or class ID"})
			return
		}

		roleID, err := roleService.GetRole(uid, cid)
		if err != nil {
			ctx.AbortWithStatusJSON(constants.StatusUnauthorized, gin.H{"error": "Unauthorized: role ID check failed"})
			return
		}

		if roleID != requiredRoleID {
			ctx.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"error": "Forbidden: insufficient privileges"})
			return
		}

		ctx.Next()
	}
}

// AdminMiddleware は管理者権限を持っているかどうかを確認するミドルウェアです。
func AdminMiddleware(roleService services.ClassUserService) gin.HandlerFunc {
	return ClassUserRoleMiddleware(roleService, AdminRoleID)
}

// AssistantMiddleware はアシスタント権限を持っているかどうかを確認するミドルウェアです。
func AssistantMiddleware(roleService services.ClassUserService) gin.HandlerFunc {
	return ClassUserRoleMiddleware(roleService, AssistantRoleID)
}

func AuthMiddleware(authenticate func(token string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !authenticate(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenAuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}

		tokenString := header[len(BearerSchema):]
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		c.Set("userID", userID)

		c.Next()
	}
}
