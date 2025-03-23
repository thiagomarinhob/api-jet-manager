// internal/api/routes/router.go
package routes

import (
	"api-jet-manager/internal/api/handlers"
	"api-jet-manager/internal/api/middlewares"
	"api-jet-manager/internal/config"
	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/auth"
	"api-jet-manager/internal/infrastructure/database"
	repoImpl "api-jet-manager/internal/infrastructure/repositories"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *database.PostgresDB) *gin.Engine {
	// Configurar modo do Gin
	gin.SetMode(cfg.GinMode)
	router := gin.Default()

	// Middleware CORS
	router.Use(middlewares.CORSMiddleware())

	// Inicialização do serviço JWT
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration)

	// Repositórios
	userRepo := repoImpl.NewPostgresUserRepository(db)
	tableRepo := repoImpl.NewPostgresTableRepository(db)
	orderRepo := repoImpl.NewPostgresOrderRepository(db)
	financeRepo := repoImpl.NewPostgresFinanceRepository(db)
	productRepo := repoImpl.NewPostgresProductRepository(db)
	productCategoryRepo := repoImpl.NewPostgresProductCategoryRepository(db)
	restaurantRepo := repoImpl.NewPostgresRestaurantRepository(db)

	// Serviços
	authService := services.NewAuthService(userRepo, jwtService)
	tableService := services.NewTableService(tableRepo)
	orderService := services.NewOrderService(orderRepo, tableRepo, financeRepo, productRepo)
	financeService := services.NewFinanceService(financeRepo)
	productService := services.NewProductService(productRepo)
	productCategoryService := services.NewProductCategoryService(productCategoryRepo)
	restaurantService := services.NewRestaurantService(restaurantRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, restaurantService)
	tableHandler := handlers.NewTableHandler(tableService)
	orderHandler := handlers.NewOrderHandler(orderService, tableService)
	financeHandler := handlers.NewFinanceHandler(financeService)
	productHandler := handlers.NewProductHandler(productService, productCategoryService)
	productCategoryHandler := handlers.NewProductCategoryHandler(productCategoryService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, authService)

	// Rotas públicas
	router.POST("/v1/auth/login", authHandler.Login)
	router.POST("/v1/auth/register-superadmin", authHandler.RegisterSuperAdmin) // Rota para o primeiro superadmin
	router.POST("/v1/auth/register-admin", authHandler.Register)

	// Grupo de rotas autenticadas
	api := router.Group("/v1")
	api.Use(middlewares.AuthMiddleware(jwtService))

	// Rotas de perfil de usuário
	api.GET("/profile", authHandler.GetProfile)
	api.PUT("/profile", authHandler.UpdateProfile)

	// Rotas de gestão de restaurantes
	restaurantsApi := api.Group("/restaurants")
	restaurantsApi.GET("", restaurantHandler.List) // Com filtro para usuários normais

	// IMPORTANTE: Todas as rotas de restaurante usam o mesmo parâmetro ":restaurant_id"
	restaurantsApi.GET("/:restaurant_id", middlewares.RestaurantMiddleware(), restaurantHandler.GetByID)

	// Operações que exigem superadmin
	restaurantAdminApi := restaurantsApi.Group("/")
	restaurantAdminApi.Use(middlewares.SuperAdminMiddleware())
	restaurantAdminApi.POST("", restaurantHandler.Create)
	restaurantAdminApi.PUT("/:restaurant_id", restaurantHandler.Update)
	restaurantAdminApi.DELETE("/:restaurant_id", restaurantHandler.Delete)
	restaurantAdminApi.PATCH("/:restaurant_id/status", restaurantHandler.UpdateStatus)

	// Rotas de usuário (agrupadas por restaurante)
	restaurantsApi.POST("/:restaurant_id/users", middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		authHandler.Register)

	// Rotas de categorias (agrupadas por restaurante)
	restaurantsApi.POST("/:restaurant_id/categories", middlewares.RestaurantMiddleware(), productCategoryHandler.Create)
	restaurantsApi.GET("/:restaurant_id/categories", middlewares.RestaurantMiddleware(), productCategoryHandler.List)
	restaurantsApi.GET("/:restaurant_id/categories/active", middlewares.RestaurantMiddleware(), productCategoryHandler.ListActive)
	restaurantsApi.GET("/:restaurant_id/categories/:category_id", middlewares.RestaurantMiddleware(), productCategoryHandler.GetByID)
	restaurantsApi.PUT("/:restaurant_id/categories/:category_id", middlewares.RestaurantMiddleware(), productCategoryHandler.Update)
	restaurantsApi.DELETE("/:restaurant_id/categories/:category_id", middlewares.RestaurantMiddleware(), productCategoryHandler.Delete)
	restaurantsApi.PATCH("/:restaurant_id/categories/:category_id/status", middlewares.RestaurantMiddleware(), productCategoryHandler.UpdateStatus)

	// Rotas de mesas (agrupadas por restaurante)
	restaurantsApi.GET("/:restaurant_id/tables", middlewares.RestaurantMiddleware(), tableHandler.List)
	restaurantsApi.GET("/:restaurant_id/tables/:table_id", middlewares.RestaurantMiddleware(), tableHandler.GetByID)
	restaurantsApi.POST("/:restaurant_id/tables",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		tableHandler.Create)
	restaurantsApi.PUT("/:restaurant_id/tables/:table_id",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		tableHandler.Update)
	restaurantsApi.DELETE("/:restaurant_id/tables/:table_id",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin),
		tableHandler.Delete)
	restaurantsApi.PATCH("/:restaurant_id/tables/:table_id/status",
		middlewares.RestaurantMiddleware(),
		tableHandler.UpdateStatus)

	// Rotas de pedidos (agrupadas por restaurante)
	restaurantsApi.GET("/:restaurant_id/orders", middlewares.RestaurantMiddleware(), orderHandler.List)
	restaurantsApi.GET("/:restaurant_id/orders/:order_id", middlewares.RestaurantMiddleware(), orderHandler.GetByID)
	restaurantsApi.POST("/:restaurant_id/orders", middlewares.RestaurantMiddleware(), orderHandler.Create)
	restaurantsApi.PATCH("/:restaurant_id/orders/:order_id/status", middlewares.RestaurantMiddleware(), orderHandler.UpdateStatus)
	restaurantsApi.POST("/:restaurant_id/orders/:order_id/items", middlewares.RestaurantMiddleware(), orderHandler.AddItem)
	restaurantsApi.DELETE("/:restaurant_id/orders/:order_id/items/:item_id", middlewares.RestaurantMiddleware(), orderHandler.RemoveItem)

	restaurantsApi.GET("/:restaurant_id/delivery/today", orderHandler.FindTodayDeliveryOrders)
	restaurantsApi.GET("/:restaurant_id/delivery/by-date", orderHandler.FindDeliveryOrdersByDate)
	restaurantsApi.GET("/:restaurant_id/delivery/by-type-and-date", orderHandler.FindOrdersByDateAndType)
	restaurantsApi.GET("/:restaurant_id/delivery/by-date-range", orderHandler.FindOrdersByDateRangeAndType)

	// Rotas de produtos (agrupadas por restaurante)
	restaurantsApi.GET("/:restaurant_id/products", middlewares.RestaurantMiddleware(), productHandler.List)
	restaurantsApi.GET("/:restaurant_id/products/:product_id", middlewares.RestaurantMiddleware(), productHandler.GetByID)
	restaurantsApi.POST("/:restaurant_id/products",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		productHandler.Create)
	restaurantsApi.PUT("/:restaurant_id/products/:product_id",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		productHandler.Update)
	restaurantsApi.DELETE("/:restaurant_id/products/:product_id",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin),
		productHandler.Delete)
	restaurantsApi.PATCH("/:restaurant_id/products/:product_id/stock",
		middlewares.RestaurantMiddleware(),
		middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager),
		productHandler.UpdateStock)

	// Rotas de finanças (agrupadas por restaurante)
	financeApi := restaurantsApi.Group("/:restaurant_id/finance")
	financeApi.Use(middlewares.RestaurantMiddleware())
	financeApi.Use(middlewares.UserTypeMiddleware(models.UserTypeAdmin, models.UserTypeManager))

	financeApi.GET("/transactions", financeHandler.List)
	financeApi.GET("/transactions/:transaction_id", financeHandler.GetByID)
	financeApi.POST("/transactions", financeHandler.Create)
	financeApi.PUT("/transactions/:transaction_id", financeHandler.Update)
	financeApi.DELETE("/transactions/:transaction_id",
		middlewares.UserTypeMiddleware(models.UserTypeAdmin),
		financeHandler.Delete)
	financeApi.GET("/summary", financeHandler.GetSummary)

	return router
}
