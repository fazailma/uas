package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RBACMiddleware checks if user has required permission
// Usage: app.Use(middleware.RBACMiddleware("permission:action"))
func RBACMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Extract JWT from header (already done by AuthMiddleware)
		// Step 2: Validate token (already done by AuthMiddleware)

		// Step 3 & 4: Check if user has required permission
		permissionsInterface := c.Locals("permissions")
		if permissionsInterface == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "forbidden",
				"error":  "no permissions found in token",
			})
		}

		// Type assertion to handle permissions as []string or []interface{}
		var permissions []string
		switch v := permissionsInterface.(type) {
		case []string:
			permissions = v
		case []interface{}:
			for _, p := range v {
				if pStr, ok := p.(string); ok {
					permissions = append(permissions, pStr)
				}
			}
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "forbidden",
				"error":  "invalid permissions format",
			})
		}

		// Check if user has required permission
		hasPermission := hasPermission(permissions, requiredPermission)
		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "forbidden",
				"error":  fmt.Sprintf("missing required permission: %s", requiredPermission),
			})
		}

		// Step 5: Allow request
		return c.Next()
	}
}

// hasPermission checks if user has specific permission
// Supports wildcard matching: e.g., "achievement:*" matches "achievement:create", "achievement:read"
func hasPermission(userPermissions []string, requiredPermission string) bool {
	for _, perm := range userPermissions {
		// Exact match
		if perm == requiredPermission {
			return true
		}

		// Wildcard match: "achievement:*" matches any "achievement:X"
		if strings.HasSuffix(perm, ":*") {
			prefix := strings.TrimSuffix(perm, ":*")
			if strings.HasPrefix(requiredPermission, prefix+":") {
				return true
			}
		}

		// Allow all permissions (super admin)
		if perm == "*" || perm == "*:*" {
			return true
		}
	}
	return false
}
