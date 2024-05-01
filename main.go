package main

import (
	"gin-freemarket-app/controllers"
	"gin-freemarket-app/infra"
	"gin-freemarket-app/middlewares"
	"github.com/gin-contrib/cors"

	//"gin-freemarket-app/models"
	"gin-freemarket-app/repositories"
	"gin-freemarket-app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	//items := []models.Item{
	//	{ID: 1, Name: "Item 1", Price: 1000, Description: "Description 1", SoldOut: false},
	//	{ID: 2, Name: "Item 2", Price: 2000, Description: "Description 2", SoldOut: true},
	//	{ID: 3, Name: "Item 3", Price: 3000, Description: "Description 3", SoldOut: false},
	//}

	//itemRepository := repositories.NewItemMemoryRepository(items)
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	r := gin.Default()
	/// Reactなど異なるサーバーからのアクセス制限を設定できるミドルウェア
	r.Use(cors.Default())
	itemRouter := r.Group("/items")
	itemRouterWithAuth := r.Group("/items", middlewares.AuthMiddleware(authService))
	authRouter := r.Group("/auth")

	itemRouter.GET("", itemController.FindAll)
	itemRouterWithAuth.GET("/:id", itemController.FindById)
	itemRouterWithAuth.POST("", itemController.Create)
	itemRouterWithAuth.PUT("/:id", itemController.Update)
	itemRouterWithAuth.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)
	r.Run("localhost:8080")
}
