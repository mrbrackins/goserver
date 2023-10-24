package main

import (
	"context"
	"fmt"
	"log"
	"os"

	database "github.com/mrbrackins/goserver/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/joho/godotenv"
)

func main() {
	// init app
	initApp()
	// if err != nil {
	// 	panic(err)
	// }

	//defer close db
	// defer database.CloseMongoDB()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// write a todo to the db

		collection := database.GetCollection("todos")
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

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ENV! " + os.Getenv("MONGODB_URI"))
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func initApp() error {
	// setup env
	err := loadENV()
	if err != nil {
		return err
	}

	// setup db
	err = database.StartMongoDB()
	if err != nil {
		return err
	}
	return nil
}

func loadENV() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}
	return nil
}
