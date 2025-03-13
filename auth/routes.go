package auth

import (
	"go-gin-e-comm/common"
	"go-gin-e-comm/user"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, userRepo user.UserRepository, log *slog.Logger, jwtSecret string) {
	authService := NewAuthService(userRepo, jwtSecret, log)

	r.POST("/signup", func(c *gin.Context) {
		user := new(UserAuth)
		if err := common.ValidateRequest(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := authService.Register(c.Request.Context(), user.Email, user.Password); err != nil {
			switch err {
			case common.ErrDatabase, ErrCreateUser, ErrHashPassword:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully"})
	})

	r.POST("/login", func(c *gin.Context) {
		user := new(UserAuth)
		if err := common.ValidateRequest(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := authService.Login(c.Request.Context(), user.Email, user.Password)
		if err != nil {
			switch err {
			case ErrGenerateToken, common.ErrDatabase:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
