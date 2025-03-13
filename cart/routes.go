package cart

import (
	"go-gin-e-comm/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupRoutes(r *gin.RouterGroup, cartRepo CartRepository) {

	r.GET("/", func(c *gin.Context) {
		userID := c.Value("userID").(primitive.ObjectID)
		cart, err := cartRepo.GetCart(c.Request.Context(), userID)
		if err != nil {
			if err == ErrCartNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cart)
	})

	r.POST("/add-product", func(c *gin.Context) {
		var item CartItem
		if err := common.ValidateRequest(c, &item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.Value("userID").(primitive.ObjectID)
		err := cartRepo.AddProductToCart(c.Request.Context(), userID, item.ProductID, item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Product added to cart"})
	})

	r.POST("/remove-product", func(c *gin.Context) {
		var item CartItem
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.Value("userID").(primitive.ObjectID)
		prodObjectID, _ := primitive.ObjectIDFromHex(item.ProductID.Hex())
		err := cartRepo.RemoveProductFromCart(c.Request.Context(), userID, prodObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
	})

	r.GET("/checkout", func(c *gin.Context) {
		userID := c.Value("userID").(primitive.ObjectID)
		order, err := cartRepo.Checkout(c.Request.Context(), userID)
		if err != nil {
			if err == ErrProductInCartNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

}
