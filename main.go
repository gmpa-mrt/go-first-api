package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` // MongoDB format
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	if os.Getenv("ENV") != "production" {
		// Load environment file if not in production
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	// Set DB url and connect it
	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	/*
		app.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:5173",
			AllowHeaders: "Origin,Content-type,Accept",
		}))
	*/

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))
}

// Create a todo
func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	// Open a cursor to iterate over all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{}) // bson.M for options research
	if err != nil {
		return err
	}

	// Close the cursor at the end
	defer cursor.Close(context.Background())

	// Populate the todos slice with data from each document
	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.Status(200).JSON(todos)
}

// Create a todo
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	// Bind and parse the incoming JSON request body to the Todo struct
	if err := c.BodyParser(todo); err != nil {
		return err
	}

	// Validate to ensure the todo body is not empty
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	// Insert the new todo item into the database
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	// Assign the MongoDB-generated ID to the new todo item
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

// Update a todo
func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	// Convert the string ID parameter to a MongoDB ObjectID type
	objectID, err := primitive.ObjectIDFromHex(id)

	// Not found
	if err != nil {
		// Return a 400 Bad Request if the ID format is invalid
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	// Define the filter and update to mark the todo as completed
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	// Execute the update operation on the database
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

// Delete a todo
func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}

	// Execute the delete operation on the database
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
