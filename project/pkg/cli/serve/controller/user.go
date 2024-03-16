package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils/authService"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct{}

func (c *UserController) RegisterRoutes(group fiber.Router) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the user group!")
	})
	group.Get("/login", c.LogIn)

	group.Post("/sign-up", c.SignUp)
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (c *UserController) LogIn(ctx *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	user, err := GetUserByEmail(loginRequest.Email)
	if err != nil {
		return err
	}
	fmt.Println("user: ", user)
	if user == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid email or password"})
	}

	token, err := authService.JwtGenerator(*user)
	fmt.Println(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}

	return ctx.JSON(user)
}

func (c *UserController) SignUp(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	if err := user.HashPassword(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	if err := database.DB.Create(&user); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user", "error": err})
	}

	token, err := authService.JwtGenerator(user)
	fmt.Println(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}
