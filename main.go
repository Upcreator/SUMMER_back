package main

import (
	"github.com/Upcreator/SUMMER_back/internal/controllers"
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     initializers.AppConfig.FrontendUrl,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Mount("/api", micro)
	app.Static("/uploads", "./uploads")
	// News

	micro.Route("/news", func(router fiber.Router) {
		router.Get("/", controllers.FindNews)
		router.Get("/admin", middleware.AuthMiddleware, middleware.RoleRequired("admin"), controllers.FindAdminNews)
		router.Get("/:newsId", controllers.FindNewsById)
		router.Use(middleware.AuthMiddleware, middleware.RoleRequired("admin"))
		router.Post("/", controllers.CreateNews)
		router.Patch("/:newsId", controllers.UpdateNews)
		router.Delete("/:newsId", controllers.DeleteNews)
	})

	// Transition applications
	micro.Route("/transition_applications", func(router fiber.Router) {
		router.Use(middleware.AuthMiddleware)
		router.Get("/", controllers.FindTransitionApplications)
		router.Post("/", controllers.CreateTransitionApplication)
		router.Patch("/:transitionApplicationId", controllers.UpdateTransitionApplication)
		router.Get("/:transitionApplicationId", controllers.FindTransitionApplicationById)
		router.Delete("/:transitionApplicationId", controllers.DeleteTransitionApplication)
	})

	// Questions
	micro.Route("/questions", func(router fiber.Router) {
		router.Post("/", controllers.CreateQuestion)
		router.Get("/", controllers.FindQuestion)
	})
	micro.Route("/questions/:questionId", func(router fiber.Router) {
		router.Patch("/", controllers.UpdateQuestion)
		router.Get("/", controllers.FindQuestionById)
		router.Delete("/", controllers.DeleteQuestion)
	})

	// Users
	micro.Route("/users", func(router fiber.Router) {
		router.Use(middleware.AuthMiddleware, middleware.RoleRequired("admin"))
		router.Post("/", controllers.CreateUser)
		router.Get("/", controllers.FindUsers)
		router.Patch("/:userId", controllers.UpdateUser)
		router.Get("/:userId", controllers.FindUserById)
		router.Delete("/:userId", controllers.DeleteUser)
	})

	// Elections
	micro.Route("/elections", func(router fiber.Router) {
		router.Post("/", controllers.CreateElection)
		router.Get("/", controllers.FindElections)
	})
	micro.Route("/elections/:voteId", func(router fiber.Router) {
		router.Patch("/", controllers.UpdateElection)
		router.Get("/", controllers.FindElectionById)
		router.Delete("/", controllers.DeleteElection)
	})

	// Votes
	micro.Route("/votes", func(router fiber.Router) {
		router.Use(middleware.AuthMiddleware)
		router.Get("/", controllers.FindVotes)
		router.Post("/:voteId", controllers.UserVote)
		router.Use(middleware.RoleRequired("admin"))
		router.Get("/ended", controllers.FindEndedVotes)
		router.Get("/:voteId", controllers.FindVoteById)
		router.Get("/:voteId/results", controllers.VoteResults)
		router.Post("/", controllers.CreateVote)
		router.Delete("/:voteId", controllers.DeleteVote)
		router.Patch("/:voteId", controllers.UpdateVote)
		router.Post("/:voteId/end", controllers.EndVote)
	})
	// Healtchecker
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, and GORM",
		})
	})

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/register", controllers.RegisterUser)
		router.Post("/login", controllers.LoginUser)
		router.Use(middleware.AuthMiddleware)
		router.Get("/me", controllers.GetUser)
		router.Post("/logout", controllers.LogoutUser)
	})

	log.Fatal(app.Listen(":8000"))
}
