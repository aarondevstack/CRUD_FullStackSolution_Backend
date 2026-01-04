package api

import (
	"fmt"
	"log"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/go-sql-driver/mysql"
)

// Server represents the API server
type Server struct {
	App    *fiber.App
	Client *ent.Client
}

// @title CRUD Solution API
// @version 1.0
// @description RESTful API for CRUD operations in Golang + Fiber
// @host localhost:8888
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func NewServer() *Server {
	// Initialize Ent Client
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name,
	)

	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: config.AppConfig.API.Name,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Initialize Casbin
	if err := middleware.InitCasbin(); err != nil {
		log.Fatalf("failed to initialize casbin: %v", err)
	}

	// Create Server instance
	server := &Server{
		App:    app,
		Client: client,
	}

	// Register Routes
	RegisterRoutes(app, client)

	// Conditionally register dev routes (Swagger)
	RegisterDevRoutes(app)

	return server
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", config.AppConfig.API.Host, config.AppConfig.API.Port)
	return s.App.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	if s.Client != nil {
		s.Client.Close()
	}
	return s.App.Shutdown()
}
