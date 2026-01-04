package api

import (
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/handlers"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(app *fiber.App, client *ent.Client) {
	// Handlers
	authHandler := handlers.NewAuthHandler(client)
	userHandler := handlers.NewUserHandler(client)
	blogHandler := handlers.NewBlogHandler(client)
	commentHandler := handlers.NewCommentHandler(client)

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API Group
	api := app.Group("/api/v1")

	// Public Routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Protected Routes (JWT + Casbin)
	// Apply JWT middleware to all routes that require authentication
	// Apply Casbin middleware to enforce RBAC
	protected := api.Group("/", middleware.JWTMiddleware(), middleware.CasbinMiddleware())

	// User Routes
	users := protected.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Patch("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Blog Routes
	blogs := protected.Group("/blogs")
	blogs.Get("/", blogHandler.GetBlogs)
	blogs.Post("/", blogHandler.CreateBlog)
	blogs.Get("/:id", blogHandler.GetBlog)
	blogs.Patch("/:id", blogHandler.UpdateBlog)
	blogs.Delete("/:id", blogHandler.DeleteBlog)

	// Comment Routes
	comments := protected.Group("/comments")
	comments.Get("/", commentHandler.GetComments)
	comments.Post("/", commentHandler.CreateComment)
	comments.Get("/:id", commentHandler.GetComment)
	comments.Patch("/:id", commentHandler.UpdateComment)
	comments.Delete("/:id", commentHandler.DeleteComment)
}
