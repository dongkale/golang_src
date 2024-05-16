package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joisandresky/go-echo-mysql-boilerplate/database"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/injector"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	LoadConfig()
	port := viper.GetString("app.port")

	dbPool := database.ConnectMYSQL().(*database.DatabaseProviderConnection)
	defer dbPool.Db.Close()

	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))
	server.Use((middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		},
	)))

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to MS-GB User Service")
	})

	server.Any("/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"code":    http.StatusNotFound,
			"message": "Route is Not Exist",
		})
	})

	injection := injector.NewInjector(*dbPool, server)
	injection.InjectModules()

	log.Println(fmt.Printf("%v is Running at port %v ... ", viper.GetString("app.name"), port))

	// Start Server
	go func() {
		if err := server.Start(":" + port); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatal("shutting down the Polaris GBMarket SERVICE ENDPOINT!")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}

func LoadConfig() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Fatal Error to Read Config \n", err)
		os.Exit(1)
	}
}
