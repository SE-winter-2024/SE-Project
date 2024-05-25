package serve

import (
	"os"

	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller"

	_ "bitbucket.org/dyfrag-internal/mass-media-core/docs"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/configs"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title						SE Project
// @version						1.0
// @description					Software Engineering Course Project
// @contact.name				Mahdieh Moghiseh
// @contact.email				mahdiehmoghiseh81@gmail.com
// @securityDefinitions.apikey	AdminID
// @in							header
// @name						X-Admin-Id
// @securityDefinitions.apikey	UserID
// @in							header
// @name						X-User-Id
// @securityDefinitions.apikey
// @in							header
// @name						api-key
// @BasePath					/v1
// @externalDocs.description	OpenAPI
func main() {
	app := fiber.New()
	initialization()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	user := app.Group("/user")
	var userController serve.UserController
	userController.RegisterRoutes(user)

	trainee := app.Group("/trainee")
	var traineeController serve.TraineeController
	traineeController.RegisterRoutes(trainee)

	trainer := app.Group("/trainer")
	var trainerController serve.TrainerController
	trainerController.RegisterRoutes(trainer)

	admin := app.Group("/admin")
	var adminController serve.AdminController
	adminController.RegisterRoutes(admin)

	port := os.Getenv("SERVER_PORT")
	connection := ":" + port
	app.Listen(connection)
}

func initialization() {
	configs.SetUpConfigs()
	database.SetUpDB()
}

func New() *cobra.Command {

	return &cobra.Command{
		Use:   "serve",
		Short: "runs http server",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
}
