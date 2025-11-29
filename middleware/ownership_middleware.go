package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// OwnershipMiddleware checks if user is accessing their own resources
// Usage: app.Use(middleware.OwnershipMiddleware())
// This middleware extracts owner_id from URL param and compares with user_id in context
func OwnershipMiddleware(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "forbidden",
			"error":  "user_id not found in token",
		})
	}

	// Get the owner ID from URL params (assumes :owner_id or :user_id in route)
	ownerID := c.Params("owner_id")
	if ownerID == "" {
		ownerID = c.Params("user_id")
	}

	// If there's an ID in params, check ownership
	if ownerID != "" && ownerID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "forbidden",
			"error":  "you can only access your own resources",
		})
	}

	return c.Next()
}

// CheckAchievementOwnership checks if user is the owner of the achievement
// Call this inside handler before processing
func CheckAchievementOwnership(userID string, achievementOwnerID string) error {
	if userID != achievementOwnerID {
		return fmt.Errorf("you can only modify your own achievements")
	}
	return nil
}
