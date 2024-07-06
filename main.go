package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Fiber web server!")
	})

	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name + "!")
	})

	// New route to serve the HTML form
	app.Get("/form", func(c *fiber.Ctx) error {

		text := `<form action="/submit" method="post">
				<label for="name">Name:</label>
				<input type="text" id="name" name="name" required><br><br>
				<label for="age">Age:</label>
				<input type="number" id="age" name="age" required><br><br>
				<input type="submit" value="Submit">
			</form>`
		html(c, text, true)
		html(c, "<p>The End</p>", false)
		return nil
	})

	// New route to handle form submission
	app.Post("/submit", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		age := c.FormValue("age")

		resp := fmt.Sprintf(`<table border=1><tr><td style="color:#f0e020">Name: %s</td><td> Age: %s</td></tr></table>`, name, age)
		html(c, resp, true)
		return nil
	})

	// Start the server
	log.Println("Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func html(c *fiber.Ctx, text string, setHeader bool) {
	if setHeader {
		c.Set("Content-Type", "text/html; charset=utf-8")
	}

	_, _ = c.WriteString(text)
}
