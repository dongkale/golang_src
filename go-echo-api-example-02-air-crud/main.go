package main

import (
	"log"
	"os"
	"simple-rest-go-echo/Config"
	"simple-rest-go-echo/Models"
	"simple-rest-go-echo/Routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "simple-rest-go-echo/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Wookiist Sample Swagger API
// @version 1.0
// @host localhost:3000
// @BasePath /api/v1
func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
         log.Fatalf("Error loading .env file: %v", err)
    }

    // Initialize Echo instance
    e := echo.New()

    // e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    //     StackSize: 1 << 10, // 1 KB        
    // }))

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.GET("/swagger/*", echoSwagger.WrapHandler)

    /// Connect to the database
    Config.DatabaseInit()
    defer Config.GetDB().DB()

    // Perform migrations using AutoMigrate
    db := Config.GetDB()
    err := db.AutoMigrate(&Models.Course{})
    if err != nil {
        panic(err)
    }
    
    // Set up Routes
    Routes.SetupRoutes(e)

    // Start the server
    serverPort := os.Getenv("PORT")
    e.Logger.Fatal(e.Start(":" + serverPort))
}