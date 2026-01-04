package api

import (
	"fmt"
	"log"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
	"github.com/kardianos/service"
)

// APIService implements service.Interface
type APIService struct {
	server *Server
	exit   chan struct{}
}

// NewAPIService creates a new APIService
func NewAPIService() *APIService {
	return &APIService{
		exit: make(chan struct{}),
	}
}

// Start starts the service
func (s *APIService) Start(svc service.Service) error {
	// Start in a separate goroutine so service manager doesn't block
	go s.run()
	return nil
}

// run is the actual execution logic
func (s *APIService) run() {
	// Load config
	if err := config.Load(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	// Create server
	s.server = NewServer()

	// Start server in goroutine
	go func() {
		if err := s.server.Start(); err != nil {
			log.Printf("Server start failed: %v", err)
		}
	}()

	// Block until exit
	<-s.exit
}

// Stop stops the service
func (s *APIService) Stop(svc service.Service) error {
	// Signal exit
	close(s.exit)

	// Graceful shutdown
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// ManageService handles service lifecycle commands
func ManageService(action string) error {
	// Load config to get service name (optional, but good practice)
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	svcConfig := &service.Config{
		Name:        "CRUDSolution_WebService",
		DisplayName: "CRUD Solution Web Service",
		Description: "High-performance Go Fiber Web Service",
		Arguments:   []string{"services", "api", "serve"}, // Arguments to run the service
		Option: service.KeyValue{
			"UserService": true,
		},
	}

	prg := NewAPIService()
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	if len(action) > 0 {
		return service.Control(s, action)
	}

	// If no action, run directly (for 'serve' command)
	return prg.runDirect()
}

// runDirect runs the service in foreground (blocking)
func (s *APIService) runDirect() error {
	// Load config
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	s.server = NewServer()
	log.Printf("Starting API server on %s:%d...", config.AppConfig.API.Host, config.AppConfig.API.Port)
	return s.server.Start()
}
