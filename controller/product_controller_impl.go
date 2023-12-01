package controller

import (
	"net/http"
	"restful_api/helper"
	"restful_api/model/web"
	"restful_api/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
	}
}

func (controller *ProductControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productCreateRequest := web.ProductCreateRequest{}
	helper.ReadFromRequestBody(request, &productCreateRequest)

	categoryResponse := controller.ProductService.Create(request.Context(), productCreateRequest)
	webResponse := web.WebResponse{
		Code:   201,
		Status: "CREATED",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productUpdateRequest := web.ProductUpdateRequest{}
	helper.ReadFromRequestBody(request, &productUpdateRequest)

	productId, err := strconv.Atoi(params.ByName("id"))
	helper.PanicIfError(err)
	productUpdateRequest.Id = productId

	categoryResponse := controller.ProductService.Update(request.Context(), productUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId, err := strconv.Atoi(params.ByName("id"))
	helper.PanicIfError(err)

	controller.ProductService.Delete(request.Context(), productId)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId, err := strconv.Atoi(params.ByName("id"))
	helper.PanicIfError(err)

	categoryResponse := controller.ProductService.FindById(request.Context(), productId)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryResponses := controller.ProductService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
