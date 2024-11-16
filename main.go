package main

import (
	"fmt"
	"log"
	"os"

	// "github.com/gofiber/fiber"  //outdated
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type TODO struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Content   string `json:"content"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []TODO{}

	// Get all Todos
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Create Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &TODO{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Content == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo Content is empty!"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		fmt.Print(todos)
		return c.Status(201).JSON(todo)
	})

	// Update Todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		fmt.Print(todos)
		fmt.Print(id)

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		fmt.Print(todos)
		return c.Status(404).JSON(fiber.Map{"error": "Todo Not Found!"})
	})

	// Delete Todos
	app.Delete("/api/todos/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		fmt.Print(todos)
		fmt.Print(id)

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		fmt.Print(todos)
		return c.Status(404).JSON(fiber.Map{"error": "Todo Not Found!"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
