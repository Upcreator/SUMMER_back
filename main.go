package main

import (
	"log"

	"github.com/Upcreator/SUMMER_back/internal/controllers"
	"github.com/Upcreator/SUMMER_back/internal/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	// News
	micro.Route("/news", func(router fiber.Router) {
		router.Post("/", controllers.CreateNews)
		router.Get("/", controllers.FindNews)
	})
	micro.Route("/news/:newsId", func(router fiber.Router) {
		router.Patch("", controllers.UpdateNews)
		router.Get("", controllers.FindNewsById)
		router.Delete("", controllers.DeleteNews)
	})

	// Transition applications
	micro.Route("/transition_applications", func(router fiber.Router) {
		router.Post("/", controllers.CreateTransitionApplication)
		router.Get("/", controllers.FindTransitionApplications)
	})
	micro.Route("/transition_applications/:transitionApplicationId", func(router fiber.Router) {
		router.Patch("/", controllers.UpdateTransitionApplication)
		router.Get("/", controllers.FindTransitionApplicationById)
		router.Delete("/", controllers.DeleteTransitionApplication)
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
		router.Post("/", controllers.CreateUser)
		router.Get("/", controllers.FindUsers)
	})
	micro.Route("/users/:userId", func(router fiber.Router) {
		router.Patch("/", controllers.UpdateUser)
		router.Get("/", controllers.FindUserById)
		router.Delete("/", controllers.DeleteUser)
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
		router.Post("/", controllers.CreateVote)
		router.Get("/", controllers.FindVotes)
	})
	micro.Route("/votes/:voteId", func(router fiber.Router) {
		router.Patch("/", controllers.UpdateVote)
		router.Get("/", controllers.FindVoteById)
		router.Delete("/", controllers.DeleteVote)
	})

	// Healtchecker
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, and GORM",
		})
	})

	log.Fatal(app.Listen(":8000"))
}
