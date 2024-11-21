package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TODO struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Content   string             `json:"content"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	MONGODB_URL := os.Getenv("MONGODB_URL")

	clientoption := options.Client().ApplyURI(MONGODB_URL)
	client, err := mongo.Connect(context.Background(), clientoption)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MONGODB Connected Successfully")

	collection = client.Database("go_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos/create", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []TODO

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo TODO
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(TODO)

	if err := c.BodyParser(todo); err != nil {
		return nil
	}

	if todo.Content == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	res, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = res.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	ID := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Todo Id"})
	}

	filter := bson.M{"_id": objectId}
	value := bson.M{"$set": bson.M{"Completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, value)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteTodo(c *fiber.Ctx) error {
	ID := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Todo Id"})
	}
	filter := bson.M{"_id": objectId}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}
