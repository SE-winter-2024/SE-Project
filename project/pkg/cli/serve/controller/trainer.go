package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type TrainerController struct{}

func (c *TrainerController) RegisterRoutes(group fiber.Router) {
	group.Get("/profile/:id", c.GetTrainerProfile)
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
	id := ctx.Params("id")
	trainerModel, err := serve.GetTrainerProfile(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	profileCard := dto.TrainerProfileCard{
		UserName:        trainerModel.UserName,
		Email:           trainerModel.Email,
		Status:          trainerModel.Status,
		Role:            trainerModel.Role,
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
