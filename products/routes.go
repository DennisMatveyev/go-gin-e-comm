package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupRoutes(r *gin.RouterGroup, productsRepo ProductRepository) {

	r.GET("/", func(c *gin.Context) {
		var pagination Pagination
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")
		pageInt, _ := strconv.Atoi(page)
		pageSizeInt, _ := strconv.Atoi(pageSize)
		pagination.Page = pageInt
		pagination.PageSize = pageSizeInt
		products, paging, err := productsRepo.GetProducts(c.Request.Context(), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products, "paging": paging})
	})

	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, _ := primitive.ObjectIDFromHex(id)
		product, err := productsRepo.GetProductByID(objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	r.GET("/search", func(c *gin.Context) {
		var params SearchParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products, err := productsRepo.SearchProducts(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
			"query":    params.Query,
			"filters": gin.H{
				"min_price": params.MinPrice,
				"max_price": params.MaxPrice,
			},
		})
	})
}
