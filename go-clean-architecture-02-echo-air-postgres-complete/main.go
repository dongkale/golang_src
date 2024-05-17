package main

import (
	"go-echo-rest-api/controller"
	"go-echo-rest-api/db"
	"go-echo-rest-api/repository"
	"go-echo-rest-api/router"
	"go-echo-rest-api/usecase"
	"go-echo-rest-api/validator"
)

func main() {	
	db := db.NewDB()

	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)

	e := router.NewRouter(userController, taskController)
	
	e.Logger.Fatal(e.Start(":8080"))
}
