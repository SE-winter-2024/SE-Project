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

type TrainerController struct{}

func (c *TrainerController) RegisterRoutes(group fiber.Router) {
	group.Get("/profile/:id", c.GetTrainerProfile)
	group.Put("/profile/:id", c.EditProfile)
	group.Get("/trainees/:id", c.GetTrainees)
	group.Get("/:id", c.getTrainer)
	group.Get("/requests/:id", c.GetRequests)
	group.Put("/request/set-price", c.SetPrice)
}

func (c *TrainerController) getTrainer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainer, err := serve.GetTrainerById(uint(id))
	if err != nil {
		return err
	}
	return ctx.JSON(trainer)
}

// EditProfile
// @Summary Edit trainer profile
// @Description Updates the profile information of a trainer by UserID
// @Tags trainer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param trainer body dto.TrainerEdit true "Trainer profile data"
// @Success 200 {object} dto.TrainerResponse "Updated trainer profile"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/profile/{id} [put]
func (c *TrainerController) EditProfile(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	var trainer dto.TrainerEdit
	if err := ctx.BodyParser(&trainer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(trainer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	newTrainer, err := serve.EditTrainerProfile(id, trainer)
	if err != nil {
		return err
	}
	return ctx.JSON(newTrainer)
}

// GetTrainerProfile
// @Summary Get trainer profile
// @Description Retrieves the profile information of a trainer by ID
// @Tags trainer
// @Accept json
// @Produce json
// @Param id path string true "Trainer ID"
// @Success 200 {object} dto.TrainerResponse "Trainer profile information"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/profile/{id} [get]
func (c *TrainerController) GetTrainerProfile(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainerModel, err := serve.GetTrainerProfile(uint(id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	profileCard := dto.TrainerProfileCard{
		UserName:        trainerModel.UserName,
		Email:           trainerModel.User.Email,
		Status:          trainerModel.Status,
		CoachExperience: trainerModel.CoachExperience,
		Contact:         trainerModel.Contact,
		Language:        trainerModel.Language,
		Country:         trainerModel.Country,
	}

	trainerDto := dto.TrainerResponse{
		TrainerProfileCard: profileCard,
		Sports:             trainerModel.Sport,
		Achievements:       trainerModel.Achievements,
		Education:          trainerModel.Education,
	}
	return ctx.JSON(trainerDto)
}

// GetTrainees
// @Summary Get trainees
// @Description get trainees of a trainer by ID
// @Tags trainer
// @Accept json
// @Produce json
// @Param id path string true "Trainer ID"
// @Success 200 {object} []dto.TraineeInTrainerPage "Trainer trainees"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/profile/{id} [get]
func (c *TrainerController) GetTrainees(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainerModel, err := serve.GetTrainerById(uint(id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	var trainees []dto.TraineeInTrainerPage
	for _, i := range trainerModel.TraineeIDs {
		t, err := serve.GetTraineeById(uint(i))
		if err != nil {
			return err
		}
		t1 := dto.TraineeInTrainerPage{Name: t.User.FirstName + " " + t.User.LastName}
		trainees = append(trainees, t1)
	}
	return ctx.JSON(trainees)
}

// GetRequests
// @Summary Get requests
// @Description get requests of a trainer by ID
// @Tags trainer
// @Accept json
// @Produce json
// @Param id path string true "Trainer ID"
// @Success 200 {object} []dto.RequestsInTrainerPage "Trainer requests"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/requests/{id} [get]
func (c *TrainerController) GetRequests(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	trainerModel, err := serve.GetTrainerById(uint(id))
	if err != nil {
		return err
	}
	var reqs []dto.RequestsInTrainerPage
	requests, err := serve.GetRequests(trainerModel)
	if err != nil {
		return err
	}
	for _, r := range requests {
		r1 := dto.RequestsInTrainerPage{
			TraineeName: r.TraineeName,
			Date:        r.Date,
			Price:       r.Price,
			Status:      r.Status,
		}
		reqs = append(reqs, r1)
	}
	return ctx.JSON(reqs)
}

func (c *TrainerController) SetPrice(ctx *fiber.Ctx) error {
	var setPrice dto.TrainerSetPrice
	if err := ctx.BodyParser(&setPrice); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(setPrice); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	req, err := serve.SetPrice(setPrice)
	if err != nil {
		return err
	}
	return ctx.JSON(req)
}
