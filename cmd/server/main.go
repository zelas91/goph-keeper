package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/models"
)

func main() {
	//fmt.Println("start ")
	//user := controller.NewUserHandler(service.NewUserService())
	//
	//router := chi.NewRouter()
	//router.Mount("/", user.InitRoutes())
	//
	//http.ListenAndServe(":9095", router)

	user := models.User{
		Email:    "zelas@gmail.com",
		Password: "asd",
		Login:    "asd",
		Card:     "4283183072528759",
	}
	v := validator.New()
	fmt.Println(v.Struct(user))
}
