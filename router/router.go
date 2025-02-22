package router

import (
	"inventoryapp/controllers"
	"inventoryapp/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("register", controllers.UserRegister)

		userRouter.POST("login", controllers.UserLogin)

		userRouter.POST("logout", controllers.UserLogout)

		userRouter.GET("profile", controllers.UserProfile)
	}

	productRouter := r.Group("/products")
	{
		productRouter.Use(middlewares.Authentication())
		productRouter.GET("/", controllers.GetProducts)
		productRouter.GET("/:productId", controllers.GetProducts)
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.PUT("/:productId", controllers.UpdateProduct)
		productRouter.DELETE("/:productId", controllers.DeleteProduct)
	}

	incomingItemRouter := r.Group("/incoming-items")
	{
		incomingItemRouter.Use(middlewares.Authentication())
		incomingItemRouter.GET("/", controllers.GetIncomingItems)
		incomingItemRouter.GET("/:incomingItemId", controllers.GetIncomingItems)
		incomingItemRouter.POST("/", controllers.CreateIncomingItem)
		incomingItemRouter.PUT("/:incomingItemId", controllers.UpdateIncomingItem)
		incomingItemRouter.PUT("/cancel/:incomingItemId", controllers.CancelIncomingItem)
	}

	outgoingItemRouter := r.Group("/outgoing-items")
	{
		outgoingItemRouter.Use(middlewares.Authentication())
		outgoingItemRouter.GET("/", controllers.GetOutgoingItems)
		outgoingItemRouter.GET("/:outgoingItemId", controllers.GetOutgoingItems)
		outgoingItemRouter.POST("/", controllers.CreateOutgoingItem)
		outgoingItemRouter.PUT("/:outgoingItemId", controllers.UpdateOutgoingItem)
		outgoingItemRouter.PUT("/cancel/:outgoingItemId", controllers.CancelOutgoingItem)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
