package repository

import (
	"context"
	"database/sql"
	"errors"
	"restful_api/helper"
	"restful_api/model/domain"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (p *ProductRepositoryImpl) Store(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	sql := "INSERT INTO products(name, description) VALUES(?,?)"

	res, err := tx.ExecContext(ctx, sql, product.Name, product.Description)
	helper.PanicIfError(err)

	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	product.Id = int(id)
	return product
}

func (p *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	sql := "UPDATE products SET name=?, description=? WHERE id=?"

	_, err := tx.ExecContext(ctx, sql, product.Name, product.Description, product.Id)
	helper.PanicIfError(err)

	return product
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, product domain.Product) {
	sql := "DELETE FROM products WHERE id=?"

	_, err := tx.ExecContext(ctx, sql, product.Id)
	helper.PanicIfError(err)
}

func (p *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Product, error) {
	sql := "SELECT id, name, description, created_at, updated_at FROM products WHERE id=?"

	rows, err := tx.QueryContext(ctx, sql, id)
	helper.PanicIfError(err)

	defer rows.Close()

	product := domain.Product{}
	if !rows.Next() {
		return product, errors.New("product not found")
	}

	err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	helper.PanicIfError(err)

	return product, nil
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Product {
	sql := "SELECT id, name, description, created_at, updated_at FROM products"

	rows, err := tx.QueryContext(ctx, sql)
	helper.PanicIfError(err)

	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product := domain.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
		helper.PanicIfError(err)

		products = append(products, product)
	}

	return products
}