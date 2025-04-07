package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing Authorization header",
			})
		}

		// Prepare request to Auth Service
		req, err := http.NewRequest("POST", "http://localhost:3000/auth", nil)
		if err != nil {
			log.Println("Failed to create request:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		req.Header.Set("Authorization", authHeader)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Auth service request failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Authentication service error",
			})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			log.Println("Auth failed:", string(body))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		// Parse user ID from response
		var authResp struct {
			UserID int `json:"user_id"`
		}
		if err := json.Unmarshal(body, &authResp); err != nil {
			log.Println("Failed to parse auth response:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Invalid auth response",
			})
		}

		// Store user ID in context
		c.Locals("user_id", int(authResp.UserID))

		return c.Next()
	}
}
