package controllers

import (
	"gorm.io/gorm"
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
	productId, err := strconv.Atoi(c.Param("productId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid Parameter",
		})

		return
	}

	product := models.Products{}

	var incomingCount, outgoingCount int64
	db.Model(&models.IncomingItems{}).Where("product_id = ?", productId).Count(&incomingCount)
	db.Model(&models.OutgoingItems{}).Where("product_id = ?", productId).Count(&outgoingCount)

	if incomingCount > 0 || outgoingCount > 0 {
		err := db.Debug().Where("id = ?", productId).First(&product).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		err = db.Model(&product).Update("deleted_at", gorm.Expr("NOW()")).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully archived product",
			"product": product,
		})
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted product",
		"product": product,
	})
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
