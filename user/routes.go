package user

import (
	"go-gin-e-comm/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupRoutes(r *gin.RouterGroup, userRepo UserRepository) {

	r.POST("/user_data", func(c *gin.Context) {
		user := new(User)
		if err := common.ValidateRequest(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.Value("userID").(primitive.ObjectID)
		if err := userRepo.CreateUserData(c.Request.Context(), user, userID); err != nil {
			switch err {
			case ErrUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User data updated successfully"})
	})

}
