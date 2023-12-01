package helper

import (
	"restful_api/model/domain"
	"restful_api/model/web"
)

func ToProductResponse(product domain.Product) web.ProductResponse {
	return web.ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ToProductsResponse(products []domain.Product) []web.ProductResponse {
	var productsResponse []web.ProductResponse
	for _, product := range products {
		productsResponse = append(productsResponse, ToProductResponse(product))
	}
	return productsResponse
}