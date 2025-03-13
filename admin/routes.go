package admin

import (
	"go-gin-e-comm/common"
	"go-gin-e-comm/products"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, productsRepo products.ProductRepository) {

	r.POST("/create-product", func(c *gin.Context) {
		product := new(products.Product)
		if err := common.ValidateRequest(c, product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := productsRepo.CreateProduct(c.Request.Context(), product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Product added"})
	})
}
