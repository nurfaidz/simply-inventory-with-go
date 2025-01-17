package controllers

import (
	"inventoryapp/database"
	"inventoryapp/helpers"
	"inventoryapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOutgoingItems(c *gin.Context) {
	db := database.GetDB()

	outgoingItems := []models.OutgoingItems{}
	outgoingItemID := c.Param("outgoing_item_id")

	if outgoingItemID != "" {
		id, err := strconv.Atoi(outgoingItemID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid Parameter",
			})

			return
		}

		result := db.Debug().Where("id = ?", id).Preload("Products").Preload("Users").Find(&outgoingItems)
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

		c.JSON(http.StatusOK, outgoingItems[0])

		return
	}

	err := db.Find(&outgoingItems).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, outgoingItems)
}

func CreateOutgoingItem(c *gin.Context) {
	db := database.GetDB()
	OutgoingItem := models.OutgoingItems{}

	if err := c.ShouldBindBodyWithJSON(&OutgoingItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	if err := db.Debug().Create(&OutgoingItem).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	// above code is for reducing stock of product from outgoing item
	Product := models.Products{}
	if err := db.Debug().Where("id = ?", OutgoingItem.ProductID).First(&Product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product Not Found",
		})

		return
	}

	// and then update the stock of product
	Product.Stock -= OutgoingItem.Qty

	if err := db.Debug().Save(&Product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	// preload product and user for response
	if err := db.Debug().Preload("Products").Preload("Users").Find(&OutgoingItem, OutgoingItem.ID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, OutgoingItem)
}

func UpdateOutgoingItem(c *gin.Context) {
	db := database.GetDB()

	OutgoingItem := models.OutgoingItems{}
	outgoingItemId, _ := strconv.Atoi(c.Param("outgoingItemId"))

	if err := c.ShouldBindJSON(&OutgoingItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	OutgoingItem.ID = uint(outgoingItemId)

	tx := db.Begin()

	previousOutgoingItem := models.OutgoingItems{}
	if err := tx.Debug().Where("id = ?", outgoingItemId).First(&previousOutgoingItem).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Outgoing Item Not Found",
		})

		return
	}

	previousQty := previousOutgoingItem.Qty

	if err := tx.Debug().Model(&previousOutgoingItem).Updates(models.OutgoingItems{
		Qty:        OutgoingItem.Qty,
		OutgoingAt: OutgoingItem.OutgoingAt,
		UserID:     OutgoingItem.UserID,
		ProductID:  OutgoingItem.ProductID,
	}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	diff := previousQty - OutgoingItem.Qty

	// get the associated product
	Product := models.Products{}
	if err := tx.Debug().Where("id = ?", OutgoingItem.ProductID).First(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product Not Found",
		})

		return
	}

	// adjust the product stock
	newStock := Product.Stock + diff

	if newStock < 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Stock of Product Can't be Negative",
		})

		return
	}
	Product.Stock = newStock

	// update stock
	if err := tx.Debug().Save(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	tx.Commit()

	// preload product and user for response
	if err := db.Debug().Preload("Products").Preload("Users").Find(&OutgoingItem, OutgoingItem.ID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, OutgoingItem)
}

func DeleteOutgoingItem(c *gin.Context) {
	db := database.GetDB()

	contentType := helpers.GetContentType(c)

	OutgoingItem := models.OutgoingItems{}
	outgoingItemId, _ := strconv.Atoi(c.Param("outgoingItemId"))

	if contentType == appJSON {
		c.ShouldBindJSON(&OutgoingItem)
	} else {
		c.ShouldBind(&OutgoingItem)
	}

	OutgoingItem.ID = uint(outgoingItemId)

	err := db.Debug().Where("id = ?", outgoingItemId).Delete(&OutgoingItem).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, OutgoingItem)
}

func HelloOutgoingItem(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
