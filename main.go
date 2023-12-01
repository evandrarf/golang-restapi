package main

import (
	"net/http"
	"restful_api/app"
	"restful_api/controller"
	"restful_api/helper"
	"restful_api/middleware"
	"restful_api/repository"
	"restful_api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	validate := validator.New()
	db := app.NewDB()

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := controller.NewProductController(productService)

	router := app.NewRouter(productController)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}