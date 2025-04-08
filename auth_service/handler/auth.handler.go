package handler

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/model/entity"
	"auth_service/model/request"
	"auth_service/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
		newUser, err := RegisterUserGoogle(name, email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to store user",
			})
		}
		user = *newUser
	}

	// Generate JWT
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
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

func LoginAuth(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)

	if err := ctx.BodyParser(loginRequest); err != nil {
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(loginRequest)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	var user entity.User

	err := database.DB.First(&user, "email = ?", loginRequest.Email).Error
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	isValid := utils.CheckPasswordHash(loginRequest.Password, user.Password)

	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong pass",
		})
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	jwtToken, errGenerateToken := utils.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to create JWT",
		})
	}

	return ctx.JSON(fiber.Map{
		"access_token": jwtToken,
	})
}

func RegisterAuth(ctx *fiber.Ctx) error {
	user := new(request.UserCreateRequest)
	if err := ctx.BodyParser(user); err != nil {
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(user)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	var userEmail entity.User
	result := database.DB.First(&userEmail, "email = ?", user.Email)
	if result.Error == nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email already registered",
		})
	}

	newUser := entity.User{
		Email:    user.Email,
		Password: user.Password,
	}

	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	newUser.Password = hashedPassword

	errCreateUser := database.DB.Create(&newUser).Error
	if errCreateUser != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newUser,
	})
}

func RegisterUserGoogle(name, email string) (*entity.User, error) {
	user := &entity.User{
		Name:  name,
		Email: email,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func AuthorizeUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.DecodeToken(tokenStr)
	if err != nil {
		if err.Error() == "token is expired" {
			return c.Status(401).JSON(fiber.Map{
				"message": "Token has expired",
			})
		}
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to decode JWT",
		})
	}

	userID := claims["id"].(float64)

	return c.JSON(fiber.Map{"user_id": userID})
}
