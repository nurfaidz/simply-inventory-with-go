package controllers

import (
	"fmt"
	"inventoryapp/database"
	"inventoryapp/helpers"
	"inventoryapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIncomingItems(c *gin.Context) {
	db := database.GetDB()

	incomingItems := []models.IncomingItems{}
	incomingItemID := c.Param("incomingItemId")

	if incomingItemID != "" {
		id, err := strconv.Atoi(incomingItemID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid Parameter",
			})

			return
		}

		result := db.Where("id = ?", id).Preload("Products").Preload("Users").Find(&incomingItems)
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
				"message": "Data Doesnt Exist",
			})

			return
		}

		c.JSON(http.StatusOK, incomingItems[0])
		return
	}

	if err := db.Debug().Preload("Products").Preload("Users").Find(&incomingItems).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, incomingItems)
}

func CreateIncomingItem(c *gin.Context) {
	db := database.GetDB()

	IncomingItem := models.IncomingItems{}

	if err := c.ShouldBindJSON(&IncomingItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	fmt.Printf("IncomingItem: %+v\n", IncomingItem)

	if err := db.Debug().Create(&IncomingItem).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	// above code is adding stock to product from incoming item
	Product := models.Products{}
	if err := db.Debug().Where("id = ?", IncomingItem.ProductID).First(&Product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product Not Found",
		})

		return
	}

	// and then update stock
	Product.Stock += IncomingItem.Qty
	if err := db.Debug().Save(&Product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to update stock",
		})

		return
	}

	// preload product and user for response
	if err := db.Debug().Preload("Products").Preload("Users").First(&IncomingItem, IncomingItem.ID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, IncomingItem)
}

func UpdateIncomingItem(c *gin.Context) {
	db := database.GetDB()
	// contentType := helpers.GetContentType(c)

	IncomingItem := models.IncomingItems{}
	incomingItemId, _ := strconv.Atoi(c.Param("incomingItemId"))

	if err := c.ShouldBindJSON(&IncomingItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	IncomingItem.ID = uint(incomingItemId)

	// start transaction
	tx := db.Begin()

	previousIncomingItem := models.IncomingItems{}
	if err := tx.Debug().Where("id = ?", incomingItemId).First(&previousIncomingItem).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Incoming Item Not Found",
		})

		return
	}

	previousQty := previousIncomingItem.Qty

	if err := tx.Debug().Model(&previousIncomingItem).Updates(models.IncomingItems{
		Qty:        IncomingItem.Qty,
		IncomingAt: IncomingItem.IncomingAt,
		UserID:     IncomingItem.UserID,
		ProductID:  IncomingItem.ProductID,
	}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	diff := IncomingItem.Qty - previousQty

	// get the associated product
	Product := models.Products{}
	if err := tx.Debug().Where("id = ?", IncomingItem.ProductID).First(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product Not Found",
		})

		return
	}

	// adjust the product stock
	newStock := Product.Stock + diff

	fmt.Printf("IncomingItem qty: %d, Previous stock: %d, Diff: %d, New stock: %d, Product Stock: %d\n", IncomingItem.Qty, previousIncomingItem.Qty, diff, newStock, Product.Stock)

	if newStock < 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Stock cannot be negative",
		})

		return
	}
	Product.Stock = newStock

	// update stock
	if err := tx.Debug().Save(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to update stock",
		})

		return
	}

	tx.Commit()

	// preload product and user for response
	if err := db.Debug().Preload("Products").Preload("Users").First(&IncomingItem, IncomingItem.ID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, IncomingItem)
}

func DeleteIncomingItem(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	IncomingItem := models.IncomingItems{}

	incomingItemId, _ := strconv.Atoi(c.Param("incomingItemId"))

	if contentType == appJSON {
		c.ShouldBindJSON(&IncomingItem)
	} else {
		c.ShouldBind(&IncomingItem)
	}

	IncomingItem.ID = uint(incomingItemId)

	err := db.Debug().Where("id = ?", incomingItemId).Delete(&IncomingItem).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, IncomingItem)
}

func HelloIncomingItem(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
