package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok":      true,
			"message": "Hello World",
		})
	})

	log.Fatal(app.Listen(":8000"))
}
