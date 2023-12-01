package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"restful_api/app"
	"restful_api/controller"
	"restful_api/helper"
	"restful_api/middleware"
	"restful_api/model/domain"
	"restful_api/repository"
	"restful_api/service"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setUpTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golang_restful_api_test?parseTime=true")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := controller.NewProductController(productService)

	router := app.NewRouter(productController)

	return middleware.NewAuthMiddleware(router)
}

func truncateProducts(db *sql.DB) {
	db.Exec("TRUNCATE TABLE products")
}

func TestCreateProductSuccess(t *testing.T) {
	db := setUpTestDB()
	truncateProducts(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":"Product 1","description": "Description 1"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusCreated, int(responseBody["code"].(float64)))
	assert.Equal(t, "CREATED", responseBody["status"])
	assert.Equal(t, 1, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Product 1", responseBody["data"].(map[string]interface{})["name"])	
	assert.Equal(t, "Description 1", responseBody["data"].(map[string]interface{})["description"])
}

func TestCreateProductFailed(t *testing.T) {
	db := setUpTestDB()
	truncateProducts(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":"","description": "Description 1"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateProductSuccess(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Update",
		Description: "Description Update",
	})
	tx.Commit()
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":"Product 1","description": "Description 1"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Product 1", responseBody["data"].(map[string]interface{})["name"])	
	assert.Equal(t, "Description 1", responseBody["data"].(map[string]interface{})["description"])
}

func TestUpdateProductFailed(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Update",
		Description: "Description Update",
	})
	tx.Commit()
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":"","description": ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestDeleteProductSuccess(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Delete",
		Description: "Description Delete",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteProductFailed(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Delete",
		Description: "Description Delete",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id + 1), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestFindAllProductsSuccess(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product1 := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Update 1",
		Description: "Description Update 1",
	})
	product2 := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product Update 2",
		Description: "Description Update 2",
	})
	products := []domain.Product{product1, product2}

	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	data := responseBody["data"].([]interface{})
	assert.Equal(t, len(products), len(data))

	for index, product := range products {
		assert.Equal(t, product.Id, int(data[index].(map[string]interface{})["id"].(float64)))
		assert.Equal(t, product.Name, data[index].(map[string]interface{})["name"])
		assert.Equal(t, product.Description, data[index].(map[string]interface{})["description"])
	}
}

func TestFindByIdProductSuccess(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product 1",
		Description: "Description 1",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Product 1", responseBody["data"].(map[string]interface{})["name"])	
	assert.Equal(t, "Description 1", responseBody["data"].(map[string]interface{})["description"])
}

func TestFindByIdProductFailed(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product 1",
		Description: "Description 1",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id + 1), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestUnauthorized(t *testing.T){
	db := setUpTestDB()
	truncateProducts(db)
	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Store(context.Background(), tx, domain.Product{
		Name: "Product 1",
		Description: "Description 1",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products/" + strconv.Itoa(product.Id ), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "WRONG_SECRET")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	res := recorder.Result()
	
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}