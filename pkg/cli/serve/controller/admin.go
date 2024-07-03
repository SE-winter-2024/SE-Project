package serve

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AdminController struct{}

func (c *AdminController) RegisterRoutes(group fiber.Router) {
	group.Get("/users", c.GetUsers)
	group.Post("/sport", c.AddSport)
	group.Get("/reports", c.GetReports)
}

// GetUsers
// @Summary Get users
// @Description get all users
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} []dto.User "User information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/users [get]
func (c *AdminController) GetUsers(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	users := serve.GetUsers()
	return ctx.JSON(users)
}

// AddSport
// @Summary Add sport
// @Description add sport
// @Tags admin
// @Accept json
// @Produce json
// @Param Sport body dto.Sport true "Sport information"
// @Success 200 {object} dto.Sport "Sport information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/sport-activity [post]
func (c *AdminController) AddSport(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	sport := dto.Sport{}
	if err := ctx.BodyParser(&sport); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	s, err := serve.AddSport(sport)
	if err != nil {
		return err
	}
	return ctx.JSON(s)
}

func (c *AdminController) GetReports(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	reports := serve.GetReports()
	return ctx.JSON(reports)
}

// BlockUser
// @Summary Block user
// @Description block user
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param report_id query uint true "Report ID"
// @Success 200 {object} dto.Response "Report information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/report/block [put]
func (c *AdminController) BlockUser(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	reportID := ctx.Params("report_id")
	report, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		return err
	}
	blocked, err := serve.BlockUser(uint(report))
	if err != nil {
		return err
	}
	fmt.Println(blocked)
	res := dto.Response{
		Message: "True",
		Success: true,
	}
	return ctx.JSON(res)
}
