package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

type TrainerController struct{}

func (c *TrainerController) RegisterRoutes(group fiber.Router) {
	group.Get("/profile/", c.GetTrainerProfile)
	group.Put("/profile/", c.EditProfile)
	group.Get("/trainees/", c.GetTrainees)
	//group.Get("/requests/", c.GetRequests)
	group.Get("/requests/", c.GetAllRequests)
	group.Put("/request/set-price", c.SetPrice)
	group.Post("/program", c.CreateTrainingProgram)
	group.Put("/program/sport-activity", c.AddSportActivity)
}

// EditProfile
// @Summary Edit trainer profile
// @Description Updates the profile information of a trainer by UserID
// @Tags trainer
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the user"
// @Param trainer body dto.TrainerEdit true "Trainer profile data"
// @Success 200 {object} dto.TrainerResponse "Updated trainer profile"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/profile/{id} [put]
func (c *TrainerController) EditProfile(ctx *fiber.Ctx) error {
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
	var trainer dto.UserEditTraineeOrTrainer
	if err := ctx.BodyParser(&trainer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(trainer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	trainerModel, err := serve.GetTrainerByUserID(userID)
	if err != nil {
		return err
	}
	newTrainer, err := serve.EditTrainerProfile(uint64(trainerModel.ID), trainer)
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
// @Param X-User-ID header string true "ID of the user"
// @Success 200 {object} dto.TrainerResponse "Trainer profile information"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/profile/{id} [get]
func (c *TrainerController) GetTrainerProfile(ctx *fiber.Ctx) error {
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
	trainerModel, err := serve.GetTrainerByUserID(userID)
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
// @Param Authorization header string true "JWT token"
// @Success 200 {object} []dto.TraineeInTrainerPage "Trainer trainees"
// @Failure 400 {object} string "Bad Request: User ID header missing or invalid token"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/trainees/ [get]
func (c *TrainerController) GetTrainees(ctx *fiber.Ctx) error {
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
	trainerModel, err := serve.GetTrainerByUserID(userID)
	fmt.Println(trainerModel)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var trainees []dto.TraineeInTrainerPage
	for _, traineeID := range trainerModel.TraineeIDs {
		trainee, err := serve.GetTraineeById(uint(traineeID))
		if err != nil {
			return err
		}
		traineeInfo := dto.TraineeInTrainerPage{
			Name: trainee.User.FirstName + " " + trainee.User.LastName,
		}
		trainees = append(trainees, traineeInfo)
	}

	return ctx.JSON(trainees)
}

// GetAllRequests
// @Summary Get requests
// @Description get requests of a trainer by ID
// @Tags trainer
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} []dto.RequestsInTrainerPage "Trainer requests"
// @Failure 404 {object} string "Trainer not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /trainer/requests/ [get]
func (c *TrainerController) GetAllRequests(ctx *fiber.Ctx) error {
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
	trainerModel, err := serve.GetTrainerByUserID(userID)
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

// SetPrice updates the price
// @Summary Set price for a request
// @Description Trainer sets the price for a training request
// @Tags trainer
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the user"
// @Param TrainerSetPrice body dto.TrainerSetPrice true "Trainer Set Price Data"
// @Success 200 {object} dto.ProgramRequestSetPrice
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal server error"
// @Router /trainer/price [put]
func (c *TrainerController) SetPrice(ctx *fiber.Ctx) error {
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

	trainerModel, err := serve.GetTrainerById(userID)
	if err != nil {
		return err
	}
	var setPrice dto.TrainerSetPrice
	if err := ctx.BodyParser(&setPrice); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	setPrice.TrainerID = trainerModel.ID
	if err := utils.ValidateUser(setPrice); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	req, err := serve.SetPrice(setPrice)
	if err != nil {
		return err
	}
	res := dto.ProgramRequestSetPrice{
		ID:          req.ID,
		TrainerID:   req.TrainerID,
		TraineeID:   req.TraineeID,
		Price:       req.Price,
		Description: req.Description,
		Status:      req.Status,
	}
	return ctx.JSON(res)
}

// CreateTrainingProgram
// @Summary creates a program
// @Description create a training program by trainer
// @Tags trainer
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the user"
// @Param TrainingProgram body dto.TrainingProgram true "Trainer Create Program data"
// @Success 200 {object} dto.Response
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal server error"
// @Router /trainer/program [post]
func (c *TrainerController) CreateTrainingProgram(ctx *fiber.Ctx) error {
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

	trainerModel, err := serve.GetTrainerById(userID)
	if err != nil {
		return err
	}

	var program dto.TrainingProgram
	if err := ctx.BodyParser(&program); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	program.TrainerID = trainerModel.ID
	if err := utils.ValidateUser(program); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	p, err := serve.CreateTrainingProgram(program)
	if err != nil {
		return err
	}
	res := dto.Response{
		Message: "Training program created",
		Success: true,
		ID:      p.ID,
	}
	return ctx.JSON(res)
}

// AddSportActivity
// @Summary add sport activity
// @Description add sport activity to program by trainer
// @Tags trainer
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the user"
// @Param SportActivity body dto.AddSportActivity true "Add Sport Activity data"
// @Success 200 {object} dto.Response
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal server error"
// @Router /trainer/program/sport-activity [put]
func (c *TrainerController) AddSportActivity(ctx *fiber.Ctx) error {
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

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	var activity dto.AddSportActivity
	if err := ctx.BodyParser(&activity); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(activity); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	s, err := serve.AddSportActivity(activity)
	if err != nil {
		return err
	}
	res := dto.Response{
		Message: "Sport Activity Added successfully",
		Success: true,
		ID:      s.ID,
	}
	return ctx.JSON(res)
}
