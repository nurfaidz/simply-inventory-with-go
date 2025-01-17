package controllers

import (
	"inventoryapp/database"
	"inventoryapp/helpers"
	"inventoryapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	db := database.GetDB()
	products := []models.Products{}
	productId := c.Param("productId")

	if productId != "" {
		id, err := strconv.Atoi(productId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid Parameter",
			})

			return
		}

		result := db.Where("id = ?", id).Find(&products)
		count := result.RowsAffected
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid Parameter",
			})

			return
		}

		if count < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Data Not Found",
				"message": "Data Doesn't Exist",
			})

			return
		}

		c.JSON(http.StatusOK, products[0])
		return
	}

	err := db.Find(&products).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	Product := models.Products{}

	productId, _ := strconv.Atoi(c.Param("productId"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.ID = uint(productId)

	err := db.Model(&Product).Where("id = ?", productId).Updates(models.Products{Name: Product.Name, Stock: Product.Stock}).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, Product)
}

func DeleteProduct(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	Product := models.Products{}

	productId, _ := strconv.Atoi(c.Param("productId"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.ID = uint(productId)

	// err := db.Model(&Product).Where("id = ?", productId).Delete(&Product).Error
	err := db.Debug().Where("id = ?", productId).Delete(&Product).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, Product)
}

func CreateProduct(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	Product := models.Products{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	err := db.Debug().Create(&Product).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, Product)
}

func HelloProduct(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
