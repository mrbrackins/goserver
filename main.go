package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// write a todo to the db

		collection := client.Database("test").Collection("movies")
		sampleDoc := bson.M{"name": "sample todo"}
		nDoc, err := collection.InsertOne(context.TODO(), sampleDoc)

		fmt.Println(nDoc)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error inserting Todo")
		}

		// send down info about the todo
		return c.JSON(nDoc)

		//return c.SendString("Hello, World!")
	})

	app.Get("/check", func(c *fiber.Ctx) error {
		collection := client.Database("test").Collection("todos")
		name := "sample todo"
		var result bson.M
		err = collection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the name %s\n", name)
			return c.SendString("Not Found!")
		}
		if err != nil {
			panic(err)
		}
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
		return c.SendString("Found!" + string(jsonData))
	})

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ENV! " + os.Getenv("MONGODB_URI"))
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}
