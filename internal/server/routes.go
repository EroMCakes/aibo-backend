package server

import (
	"aibo/internal/database"
	"aibo/internal/handlers"
	"aibo/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the routes for the server.
//
// It creates an instance of AuthService and assigns it to handle the "/register" and "/login" routes.
//
// It creates a route group "/premium" that requires authentication and a premium subscription.
// The premium routes are not implemented yet.
//
// The "/profile" route is accessible only if the user is authenticated.
func SetupRoutes(router *gin.Engine, db database.Service) {

	router.GET("/health", handlers.DBHealthHandler(db))
	router.POST("/migrate", handlers.HandleMigrate(db))

	authHandler := handlers.NewAuthService(db.GetDB())
	cbRepo := handlers.NewCatBudService(db.GetDB())

	// setupRoutes sets up the routes for the server.
	//
	// It creates an instance of AuthService and assigns it to handle the "/register" and "/login" routes.
	//
	// It creates a route group "/premium" that requires authentication and a premium subscription.
	// The premium routes are not implemented yet.
	//
	// The "/profile" route is accessible only if the user is authenticated.

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/update-profile", authHandler.UpdateProfile)
		protected.PUT("/update-password", authHandler.UpdatePassword)
		protected.POST("/logout", authHandler.Logout)

		catbuds := protected.Group("/catbud")
		{
			catbuds.GET("/:aiboId", cbRepo.GetCatBuds)
			catbuds.POST("/", cbRepo.CreateCatBuds)

		}
	}

	aiborepo := authHandler.AiboRepository

	// Premium routes
	premium := router.Group("/premium")
	premium.Use(middlewares.AuthMiddleware(), middlewares.PremiumMiddleware(aiborepo))
	{
		// TO-DO Add premium-only routes here
	}
}
