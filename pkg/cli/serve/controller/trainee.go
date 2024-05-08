package serve

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type TraineeController struct{}

func (c *TraineeController) RegisterRoutes(group fiber.Router) {
	group.Put("/profile/:id", c.EditProfile)
	group.Get("/profile/:id", c.GetTraineeProfile)
	group.Post("/request/:id", c.CreateProgramRequest)
	group.Get("/:id", c.getTrainee)
	group.Get("/request/:id", c.GetRequest)
}

func (c *TraineeController) getTrainee(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainee, err := serve.GetTraineeById(uint(id))
	if err != nil {
		return err
	}
	return ctx.JSON(trainee)
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

// CreateProgramRequest creates a new program request
// @Summary Create program request
// @Description Creates a new program request with the provided data
// @Tags trainee
// @Accept json
// @Produce json
// @Param id path string true "Trainee ID"
// @Param request body dto.ProgramRequest true "Program request data"
// @Success 200 {object} dto.ProgramRequest "Created program request"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Invalid user ID or not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainee/request/{id} [post]
func (c *TraineeController) CreateProgramRequest(ctx *fiber.Ctx) error {
	var request dto.ProgramRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	request.TraineeID = uint(id)

	if err := utils.ValidateUser(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	requestModel, err := serve.CreateProgramRequest(request)
	if err != nil {
		return err
	}
	r := dto.ProgramRequest{
		ID:          requestModel.ID,
		TrainerID:   requestModel.TrainerID,
		TraineeID:   requestModel.TraineeID,
		Description: requestModel.Description,
		ActiveDays:  request.ActiveDays,
	}
	return ctx.JSON(r)
}

// GetRequest
// @Summary Get trainee request
// @Description Retrieves the request of a trainee by ID
// @Tags trainee
// @Accept json
// @Produce json
// @Param id path string true "Trainee ID"
// @Success 200 {object} dto.RequestsInTrainerPage "Trainee request"
// @Failure 404 {object} string "Trainee not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainee/request/{id} [get]
func (c *TraineeController) GetRequest(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainee, err := serve.GetTraineeById(uint(id))
	if err != nil {
		return err
	}

	r, err := serve.GetRequest(trainee.RequestID)
	if err != nil {
		return err
	}
	req := dto.RequestsInTrainerPage{
		Date:   r.Date,
		Price:  r.Price,
		Status: r.Status,
	}
	return ctx.JSON(req)
}
