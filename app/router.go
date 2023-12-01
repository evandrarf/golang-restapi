package app

import (
	"restful_api/controller"
	"restful_api/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(productController controller.ProductController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/products", productController.FindAll)
	router.GET("/api/products/:id", productController.FindById)
	router.POST("/api/products", productController.Create)
	router.PUT("/api/products/:id", productController.Update)
	router.DELETE("/api/products/:id", productController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}