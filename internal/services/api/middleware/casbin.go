package middleware

import (
	"fmt"
	"os"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/rbac"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/fiber/v2"
)

var enforcer *casbin.Enforcer

// InitCasbin initializes the Casbin enforcer with embedded model and policy
func InitCasbin() error {
	// Load model from embedded file
	m, err := model.NewModelFromString(string(rbac.ModelConf))
	if err != nil {
		return fmt.Errorf("failed to load casbin model: %w", err)
	}

	// Create temp file for policy
	tmpFile, err := os.CreateTemp("", "policy-*.csv")
	if err != nil {
		return fmt.Errorf("failed to create temp policy file: %w", err)
	}
	// defer os.Remove(tmpFile.Name()) // Keep file for now, or manage lifecycle

	if _, err := tmpFile.Write(rbac.PolicyCSV); err != nil {
		return fmt.Errorf("failed to write policy to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp policy file: %w", err)
	}

	// Create adapter from temp file
	adapter := fileadapter.NewAdapter(tmpFile.Name())

	// Create enforcer
	enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		return fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load policy
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load casbin policy: %w", err)
	}

	return nil
}

// CasbinMiddleware enforces RBAC authorization
func CasbinMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by JWT middleware)
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found in token",
			})
		}

		// Get request path and method
		path := c.Path()
		method := c.Method()

		// Check permission
		allowed, err := enforcer.Enforce(role, path, method)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Access denied: %s role cannot %s %s", role, method, path),
			})
		}

		return c.Next()
	}
}
