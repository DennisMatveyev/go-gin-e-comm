package admin

import (
	"go-gin-e-comm/common"
	"go-gin-e-comm/products"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, productsRepo products.ProductRepository) {

	r.POST("/products", func(c *gin.Context) {
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

	r.GET("/products", func(c *gin.Context) {
		var pagination products.Pagination
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
}
