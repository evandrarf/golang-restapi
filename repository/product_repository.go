package repository

import (
	"context"
	"database/sql"
	"restful_api/model/domain"
)

type ProductRepository interface {
	Store(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	Delete(ctx context.Context, tx *sql.Tx, product domain.Product) 
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Product, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Product
}