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
	backlogRepo := repository.NewBacklogRepository(db)
	historyRepo := repository.NewItemHistoryRepository(db)
	sprintRepo := repository.NewSprintRepository(db)
	sprintHistoryRepo := repository.NewSprintHistoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	backlogService := service.NewBacklogService(backlogRepo, historyRepo)
	sprintService := service.NewSprintService(sprintRepo, sprintHistoryRepo, backlogRepo, historyRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	backlogHandler := handler.NewBacklogHandler(backlogService)
	sprintHandler := handler.NewSprintHandler(sprintService)

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
				backlog.GET("", backlogHandler.GetAll)
				backlog.POST("", backlogHandler.Create)
				backlog.GET("/:id", backlogHandler.GetByID)
				backlog.PUT("/:id", backlogHandler.Update)
				backlog.DELETE("/:id", backlogHandler.Delete)
				backlog.PATCH("/:id/status", backlogHandler.UpdateStatus)
				backlog.PATCH("/:id/priority", backlogHandler.UpdatePriority)
				backlog.POST("/:id/comments", backlogHandler.AddComment)
				backlog.POST("/:id/labels", backlogHandler.AddLabel)
				backlog.DELETE("/:id/labels/:label", backlogHandler.RemoveLabel)
				backlog.GET("/:id/history", backlogHandler.GetHistory)
			}

			// Sprints
			sprints := protected.Group("/sprints")
			{
				sprints.GET("", sprintHandler.GetAll)
				sprints.POST("", sprintHandler.Create)
				sprints.GET("/:id", sprintHandler.GetByID)
				sprints.PUT("/:id", sprintHandler.Update)
				sprints.DELETE("/:id", sprintHandler.Delete)
				sprints.POST("/:id/start", sprintHandler.Start)
				sprints.POST("/:id/complete", sprintHandler.Complete)
				sprints.POST("/:id/cancel", sprintHandler.Cancel)
				sprints.POST("/:id/items", sprintHandler.AddItem)
				sprints.DELETE("/:id/items/:itemId", sprintHandler.RemoveItem)
				sprints.GET("/:id/history", sprintHandler.GetHistory)
				sprints.GET("/:id/report", sprintHandler.GetReport)
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
