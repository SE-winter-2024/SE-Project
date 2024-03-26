package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type TraineeController struct{}

func (c *TraineeController) RegisterRoutes(group fiber.Router) {
	group.Put("/profile/:id", c.EditProfile)
	group.Get("/profile/:id", c.GetTraineeProfile)
}

// EditProfile
// @Summary Edit trainee profile
// @Description Updates the profile information of a trainee by UserID
// @Tags trainee
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param trainer body dto.TraineeEdit true "Trainee profile data"
// @Success 200 {object} dto.TraineeResponse "Updated trainee profile"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Trainee not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainee/profile/{id} [put]
func (c *TraineeController) EditProfile(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	var trainee dto.TraineeEdit
	if err := ctx.BodyParser(&trainee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(trainee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	newTrainee, err := serve.EditTraineeProfile(id, trainee)
	if err != nil {
		return err
	}
	return ctx.JSON(newTrainee)
}

// GetTraineeProfile
// @Summary Get trainee profile
// @Description Retrieves the profile information of a trainee by ID
// @Tags trainee
// @Accept json
// @Produce json
// @Param id path string true "Trainee ID"
// @Success 200 {object} dto.TraineeResponse "Trainee profile information"
// @Failure 404 {object} string "Trainee not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainee/profile/{id} [get]
func (c *TraineeController) GetTraineeProfile(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	traineeModel, err := serve.GetTraineeProfile(uint(id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	profileCard := dto.TraineeProfileCard{
		UserName: traineeModel.UserName,
		Email:    traineeModel.User.Email,
		Status:   traineeModel.Status,
		Wallet:   traineeModel.User.Wallet,
		Contact:  traineeModel.Contact,
		Language: traineeModel.Language,
		Country:  traineeModel.Country,
	}

	traineeDto := dto.TraineeResponse{
		TraineeProfileCard: profileCard,
		SportExperience:    traineeModel.Sports,
		HealthProblems:     traineeModel.MedicalHistory,
	}
	return ctx.JSON(traineeDto)
}
