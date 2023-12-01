package service

import (
	"context"
	"restful_api/model/web"
)

type ProductService interface {
	Create(ctx context.Context, request web.ProductCreateRequest) web.ProductResponse
	Update(ctx context.Context, request web.ProductUpdateRequest) web.ProductResponse
	Delete(ctx context.Context, id int)
	FindById(ctx context.Context, id int) web.ProductResponse
	FindAll(ctx context.Context) []web.ProductResponse
}