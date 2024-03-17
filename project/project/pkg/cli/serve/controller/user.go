package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils/authService"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

func (c *UserController) LogIn(ctx *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	user, err := serve.GetUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Cannot get user", "error": err})
	}
	token, err := authService.JwtGenerator(user)
	fmt.Println(token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}

	return ctx.JSON(user)
}

func (c *UserController) SignUp(ctx *fiber.Ctx) error {
	var user dto.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	if err := user.HashPassword(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}
	userModel, err := serve.CreateUser(user)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "cannot create user", "error": err})
	}
	_, err = authService.JwtGenerator(userModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}

	return ctx.JSON(userModel)
}
