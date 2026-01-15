package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sprint-backlog/internal/handler"
	"sprint-backlog/internal/middleware"
	"sprint-backlog/internal/repository"
	"sprint-backlog/internal/service"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(projectRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Sprint Backlog API is running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes (auth)
		auth := api.Group("/auth")
		{
			auth.POST("/google/verify", authHandler.VerifyGoogleCode)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth - protected
			protected.GET("/auth/me", authHandler.GetCurrentUser)

			// Users
			users := protected.Group("/users")
			{
				users.GET("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				users.GET("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				users.GET("/:id/activities", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				users.PUT("/profile", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
			}

			// Projects
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.GetAll)
				projects.POST("", projectHandler.Create)
				projects.GET("/:id", projectHandler.GetByID)
				projects.PUT("/:id", projectHandler.Update)
				projects.DELETE("/:id", projectHandler.Delete)
			}

			// Backlog
			backlog := protected.Group("/backlog")
			{
				backlog.GET("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.POST("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.GET("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.PUT("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.DELETE("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.PATCH("/:id/status", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.PATCH("/:id/priority", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.POST("/:id/comments", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.POST("/:id/labels", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.DELETE("/:id/labels/:label", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				backlog.GET("/:id/history", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
			}

			// Sprints
			sprints := protected.Group("/sprints")
			{
				sprints.GET("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.POST("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.GET("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.PUT("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.DELETE("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.POST("/:id/start", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.POST("/:id/complete", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.POST("/:id/cancel", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.POST("/:id/items", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.DELETE("/:id/items/:itemId", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.GET("/history", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.GET("/:id/history", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				sprints.GET("/:id/history/report", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
			}

			// Board
			board := protected.Group("/board")
			{
				board.GET("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
				board.PATCH("/items/:id/move", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "Not implemented yet"})
				})
			}
		}
	}

	return r
}
