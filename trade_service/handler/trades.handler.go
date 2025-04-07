package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"trade_service/database"
	"trade_service/model/entity"
	"trade_service/utils"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
)

func GoogleLogin(c *fiber.Ctx) error {

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != "randomstate" {
		return c.SendString("State validation failed")
	}

	code := c.Query("code")

	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return c.SendString("Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.SendString("User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.SendString("JSON Parsing Failed")
	}

	var userMap map[string]interface{}
	if err := json.Unmarshal(userData, &userMap); err != nil {
		return c.SendString("JSON Unmarshal Failed")
	}

	var user entity.User
	email := userMap["email"].(string)
	name := userMap["name"].(string)

	registered := database.DB.First(&user, "email = ?", email).Error
	if registered != nil {
		newUser, err := RegisterUser(name, email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to store user",
			})
		}
		user = *newUser
	}

	// Generate JWT
	claims := jwt.MapClaims{}
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	jwtToken, errGenerateToken := utils.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create JWT",
		})
	}

	return c.Redirect("http://localhost:3000/dummy_success_page.html?token=" + jwtToken)
}

func RegisterUser(name, email string) (*entity.User, error) {
	user := &entity.User{
		Name:  name,
		Email: email,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
