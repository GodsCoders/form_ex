package main

import (
	"fmt"
	"form_ex/dbops"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rohanthewiz/logger"
)

func main() {
	// Variable types (simple types on the left | aggregate types on the right
	// int, float, string, bool   |   array, map, channel, structs, func def

	// Example of treating a function as a variable
	/*	var doYourThing  = func(input string) (output string) {
			// fmt.Println("**-> output", output)
			return
		}

		result := doYourThing("Mary")
		fmt.Println(result)
	*/

	// Create a new web server
	app := fiber.New()

	// Define some routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Fiber web server!")
	})

	app.Get("/sayhello/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name + "!")
	})

	app.Get("/show-persons", func(c *fiber.Ctx) error {
		persons, err := dbops.GetPersons()
		if err != nil {
			html(c, "<p>There was an error: "+err.Error(), true)
			return nil
		}

		const pre = `<body style="background-color: slategrey;font-weight: bold;">
		<table border=1 cellpadding=3>`
		const post = `</table></body>`

		var rows []string
		for _, person := range persons {
			row := fmt.Sprintf(
				`<tr><td style="color:#3adeda">Name: %s</td><td> Age: %d</td></tr>`,
				person.Name, person.Age)
			rows = append(rows, row)
		}

		html(c, pre+strings.Join(rows, "\n")+post, true)
		return nil
	})

	// New route to serve the HTML form
	app.Get("/form",
		func(c *fiber.Ctx) error {
			text := `
			<body style="background-color: #89a3bd;font-weight: bold;">
			<form action="/submit" method="post">
				<label for="name">Name:</label>
				<input type="text" id="name" name="name" required><br><br>

				<label for="age">Age:</label>
				<input type="number" id="age" name="age" required><br><br>

				<input type="submit" value="Submit">
			</form>
			</body>`
			html(c, text, true)
			html(c, "<p>The End</p>", false)
			return nil
		},
	)

	// New route to handle form submission
	app.Post("/submit", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		strAge := c.FormValue("age")

		age, err := strconv.Atoi(strAge)
		if err != nil {
			fmt.Println("error converting age", err)
		}

		// Save the Person's info to the DB
		err = dbops.SavePerson(name, age)
		if err != nil {
			logger.LogErr(err, "Error saving "+name)
		}

		// Render the results as HTML
		resp := fmt.Sprintf(`
			<body style="background-color: slategrey;font-weight: bold;">
			<table border=1 cellpadding=3>
				<tr><td style="color:#3adeda">Name: %s</td><td> Age: %d</td></tr>
			</table>
			</body>`, name, age)

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
