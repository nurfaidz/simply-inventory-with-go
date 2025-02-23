package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"inventoryapp/database"
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

		result := db.Where("id = ?", id).Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).Preload("Users").Find(&incomingItems)
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

	if err := db.Debug().Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Users").Find(&incomingItems).Error; err != nil {
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

	// add status success to incoming item
	IncomingItem.Status = "succeed"

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
		// ProductID:  IncomingItem.ProductID,
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
	if err := tx.Debug().Where("id = ?", previousIncomingItem.ProductID).First(&Product).Error; err != nil {
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

func CancelIncomingItem(c *gin.Context) {
	db := database.GetDB()

	incomingItemId, err := strconv.Atoi(c.Param("incomingItemId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid Parameter",
		})

		return
	}

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

	if err := tx.Debug().Model(&previousIncomingItem).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	diff := previousQty

	Product := models.Products{}
	if err := tx.Debug().Where("id = ?", previousIncomingItem.ProductID).First(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Product Not Found",
		})

		return
	}

	newStock := Product.Stock - diff

	if newStock < 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Stock of Product Can't be negative",
		})

		return
	}

	Product.Stock = newStock

	if err := tx.Debug().Save(&Product).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	tx.Commit()

	IncomingItem := models.IncomingItems{}
	if err := db.Debug().First(&IncomingItem, incomingItemId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	if err := db.Debug().Preload("Products").Preload("Users").First(&IncomingItem, incomingItemId).Error; err != nil {
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
