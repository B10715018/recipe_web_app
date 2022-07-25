// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Contact: Brandon<brandon.wymer@rocketmail.com>
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipe_api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var err error
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

// populate mock data into inconsistent data
func init() {
	ctx := context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).
		Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection,
		redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).
		Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	// initialize the router
	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.DELETE("recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	}
	// run the web application, can specify port too here
	router.Run()
}
