package auth

import (
	"errors"
	"go-gin-e-comm/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticationMiddleware(userRepo user.UserRepository, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuthHeader.Error()})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		if _, err := userRepo.FindUserByID(c.Request.Context(), objectID); err != nil {
			var status int
			if err == user.ErrUserNotFound {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", objectID)
		c.Next()
	}
}

func AdminMiddleware(userRepo user.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Value("userID").(primitive.ObjectID)
		userDB, err := userRepo.FindUserByID(c.Request.Context(), userID)
		var status int
		if err != nil {
			if err == user.ErrUserNotFound {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !userDB.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": ErrAdminRequired.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
